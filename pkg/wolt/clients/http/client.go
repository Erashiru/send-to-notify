package http

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	models2 "github.com/kwaaka-team/orders-core/core/wolt/models"
	"github.com/kwaaka-team/orders-core/core/wolt/models_v2"
	"github.com/kwaaka-team/orders-core/pkg/wolt/clients"
	dto2 "github.com/kwaaka-team/orders-core/pkg/wolt/clients/dto"
	"github.com/pkg/errors"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type Client struct {
	ApiKey             string
	restyClient        *resty.Client
	StoreID            string
	BaseUrl            string
	Username, Password string
}

var timeout = map[int]time.Duration{
	0: 1 * time.Second,
	1: 5 * time.Second,
	2: 10 * time.Second,
	3: 20 * time.Second,
	4: 20 * time.Second,
	5: 20 * time.Second,
	6: 30 * time.Second,
	7: 30 * time.Second,
	8: 30 * time.Second,
}

var (
	err500                   = errors.New("Wolt response is Internal server error")
	err409                   = errors.New("already accepted order")
	AcceptedWoltOrderStatus  = "production"
	ReadyWoltOrderStatus     = "ready"
	DeliveredWoltOrderStatus = "delivered"
)

func NewClient(cfg *clients.Config) (clients.Wolt, error) {

	if cfg.Username == "" || cfg.Password == "" {
		return nil, errors.New("username or password is not provided")
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = baseURL
	}

	client := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetHeaders(map[string]string{
			contentTypeHeader: jsonType,
			acceptHeader:      jsonType,
			WOLT_API_KEY:      cfg.ApiKey,
		}).
		SetBasicAuth(cfg.Username, cfg.Password). // think...
		SetRetryCount(retriesNumber).
		SetRetryWaitTime(retriesWaitTime)

	cl := Client{
		restyClient: client,
		ApiKey:      cfg.ApiKey,
		StoreID:     cfg.StoreID,
		BaseUrl:     cfg.BaseURL,
		Username:    cfg.Username,
		Password:    cfg.Password,
	}

	return cl, nil
}

func (cl Client) DeliveredOrder(ctx context.Context, orderID string) error {
	path := fmt.Sprintf("/orders/%s/delivered", orderID)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		Put(path)
	if err != nil {
		return err
	}

	log.Info().Msgf("%+v request DeliveredOrder, URL: %s", resp.Request.Body, resp.Request.URL)
	log.Info().Msgf("%+v response DeliveredOrder", resp.Body())

	if resp.IsError() {
		return errors.New(resp.Status() + " " + string(resp.Body()))
	}

	return nil
}

func (cl Client) GetOrder(ctx context.Context, orderID string) (models2.Order, error) {
	path := fmt.Sprintf("/orders/%s", orderID)

	order := models2.Order{}

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetResult(&order).
		Get(path)
	if err != nil {
		return models2.Order{}, err
	}

	log.Info().Msgf("%s request GetOrder", resp.Request.URL)
	log.Info().Msgf("%+v response GetOrder", string(resp.Body()))

	if resp.IsError() {
		return models2.Order{}, fmt.Errorf("response status: %s, message: %s", resp.Status(), string(resp.Body()))
	}

	return order, nil
}

func (cl Client) GetOrderByV2(ctx context.Context, orderID string) (models_v2.Order, error) {
	path := fmt.Sprintf("/orders/%s", orderID)

	order := models_v2.Order{}

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetHeader("Content-Type", "application/vnd.wolt.order+json;version=2beta1").
		SetResult(&order).
		Get(path)
	if err != nil {
		return models_v2.Order{}, err
	}

	log.Info().Msgf("%s request GetOrder v2", resp.Request.URL)
	log.Info().Msgf("%+v response GetOrder v2", string(resp.Body()))

	if resp.IsError() {
		return models_v2.Order{}, fmt.Errorf("response status: %s, message: %s", resp.Status(), string(resp.Body()))
	}

	return order, nil
}

func (cl Client) AcceptSelfDeliveryOrder(ctx context.Context, order dto2.AcceptSelfDeliveryOrderOrderRequest) error {

	orderInAggregator, err := cl.GetOrder(ctx, order.ID)
	if err != nil {
		log.Error().Msgf("Error while get order ID: %s", order.ID)
		return err
	}

	if orderInAggregator.OrderStatus == AcceptedWoltOrderStatus {
		log.Info().Msgf("order already in PRODUCTION state")
		return nil
	}

	for i := 0; i < len(timeout); i++ {
		err := cl.acceptSelfDeliveryOrder(ctx, order)
		if err == nil {
			return nil
		}
		if errors.Is(err, err500) {
			log.Error().Msgf("Error while retrying: %s, retry number: %v, order ID: %s", err.Error(), i+1, order.ID)
			time.Sleep(timeout[i])
			continue
		}
		if errors.Is(err, err409) {
			log.Err(err).Msg("ignore 409 error")
			return nil
		}
		return err
	}

	return nil
}

func (cl Client) acceptSelfDeliveryOrder(ctx context.Context, order dto2.AcceptSelfDeliveryOrderOrderRequest) error {
	var err error

	path := fmt.Sprintf("/orders/%s/self-delivery/accept", order.ID)

	if order.DeliveryTime != nil {
		pickUpTime := order.DeliveryTime.Format(time.RFC3339)

		val, err := time.Parse(time.RFC3339, pickUpTime)
		if err != nil {
			return err
		}
		order.DeliveryTime = &val
	}

	resp, err := cl.restyClient.R().
		SetBody(order).
		Put(path)

	log.Info().Msgf("%+v request AcceptSelfDeliveryOrder, URL: %s\n", resp.Request.Body, resp.Request.URL)
	log.Info().Msgf("%+v response AcceptSelfDeliveryOrder\n", resp.Body())

	if err != nil {
		log.Info().Err(err)
		return err
	}

	if !resp.IsError() {
		return nil
	}

	if resp.StatusCode() == http.StatusInternalServerError {
		return errors.Wrap(err500, "Wolt Resp")
	}

	if resp.StatusCode() == http.StatusConflict {
		return errors.Wrap(err409, fmt.Sprintf("Wolt Order %s already accepted", order.ID))
	}

	return errors.New(resp.Status() + " " + string(resp.Body()))

}

func (cl Client) AcceptOrder(ctx context.Context, order dto2.AcceptOrderRequest) error {

	orderInAggregator, err := cl.GetOrder(ctx, order.ID)
	if err != nil {
		log.Error().Msgf("Error while get order ID: %s", order.ID)
		return err
	}

	if orderInAggregator.OrderStatus == AcceptedWoltOrderStatus {
		log.Info().Msgf("order already in PRODUCTION state")
		return nil
	}

	for i := 0; i < len(timeout); i++ {
		err := cl.acceptOrder(ctx, order)
		if err == nil {
			return nil
		}
		if errors.Is(err, err500) {
			log.Error().Msgf("Error while retrying: %s, retry number: %v, order ID: %s", err.Error(), i+1, order.ID)
			time.Sleep(timeout[i])
			continue
		}
		if errors.Is(err, err409) {
			log.Err(err).Msg("ignore 409 error")
			return nil
		}
		return err
	}

	return nil
}

func (cl Client) acceptOrder(ctx context.Context, order dto2.AcceptOrderRequest) error {
	var err error

	path := fmt.Sprintf("/orders/%s/accept", order.ID)

	if order.PickupTime != nil {
		pickUpTime := order.PickupTime.Format(time.RFC3339)

		val, err := time.Parse(time.RFC3339, pickUpTime)
		if err != nil {
			return err
		}
		order.PickupTime = &val
	}

	resp, err := cl.restyClient.R().
		SetBody(order).
		Put(path)

	log.Info().Msgf("%+v request AcceptOrder, URL: %s\n", resp.Request.Body, resp.Request.URL)
	log.Info().Msgf("%+v response AcceptOrder\n", resp.Body())

	if err != nil {
		log.Info().Err(err)
		return err
	}

	if !resp.IsError() {
		return nil
	}

	if resp.StatusCode() == http.StatusInternalServerError {
		return errors.Wrap(err500, "Wolt Resp")
	}

	if resp.StatusCode() == http.StatusConflict {
		return errors.Wrap(err409, fmt.Sprintf("Wolt Order %s already accepted", order.ID))
	}

	return errors.New(resp.Status() + " " + string(resp.Body()))

}

func (cl Client) RejectOrder(ctx context.Context, order dto2.RejectOrderRequest) error {

	path := fmt.Sprintf("/orders/%s/reject", order.ID)

	resp, err := cl.restyClient.R().
		SetBody(order).
		Put(path)
	if err != nil {
		return err
	}

	log.Info().Msgf("%+v request RejectOrder, URL: %s\n", resp.Request.Body, resp.Request.URL)
	log.Info().Msgf("%+v response RejectOrder\n", string(resp.Body()))

	if resp.IsError() {
		return errors.New(resp.Status() + " " + string(resp.Body()))
	}
	return nil
}

func (cl Client) MarkOrder(ctx context.Context, orderID string) error {

	path := fmt.Sprintf("/orders/%s/%s", orderID, "ready")
	resp, err := cl.restyClient.R().
		Put(path)

	log.Info().Msgf("%+v request MarkOrder, URL: %s\n", resp.Request.Body, resp.Request.URL)
	log.Info().Msgf("%+v response MarkOrder\n", resp.Body())

	if err != nil {
		log.Info().Err(err)
		return err
	}
	if resp.IsError() {
		return errors.New(resp.Status() + " " + string(resp.Body()))
	}
	return nil
}

func (cl Client) ConfirmPreOrder(ctx context.Context, orderID string) error {
	path := fmt.Sprintf("/orders/%s/%s", orderID, "confirm-preorder")
	resp, err := cl.restyClient.R().
		Put(path)

	log.Info().Msgf("%+v request ConfirmPreOrder, URL: %s\n", resp.Request.Body, resp.Request.URL)
	log.Info().Msgf("%+v response ConfirmPreOrder\n", resp.Body())

	if err != nil {
		log.Info().Err(err)
		return err
	}
	if resp.IsError() {
		return errors.New(resp.Status() + " " + string(resp.Body()))
	}
	return nil
}

func (cl Client) ManageStore(ctx context.Context, storeStatus dto2.IsStoreOpen) error {
	path := fmt.Sprintf("venues/%s/online", storeStatus.VenueId)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		SetBody(dto2.IsStoreOpenRequest{
			Status: storeStatus.AvailableStore,
		}).
		Patch(path)

	if err != nil {
		return err
	}

	log.Info().Msgf("request Wolt store status(open/close) %+v \n", resp.Request.Body)

	log.Info().Msgf("response Wolt %+v", resp.String())

	if resp.IsError() {
		return fmt.Errorf(resp.String())
	}

	return nil
}

func (cl Client) GetStoreStatus(ctx context.Context, venueId string) (dto2.StoreStatusResponse, error) {
	path := fmt.Sprintf("/venues/%s/status", venueId)

	var response dto2.StoreStatusResponse

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		SetResult(&response).
		Get(path)

	if err != nil {
		return dto2.StoreStatusResponse{}, err
	}

	utils.Beautify("get store status request body", resp.Request.Body)

	utils.Beautify("get store status response body", resp.String())

	if resp.IsError() {
		return dto2.StoreStatusResponse{}, fmt.Errorf("%v", resp.Error())
	}

	return response, nil
}

func (cl Client) UploadMenu(ctx context.Context, menu dto2.Menu, storeID string) error {
	path := fmt.Sprintf("/v1/restaurants/%s/menu", storeID)

	resp, err := cl.restyClient.R().
		SetBody(menu).
		Post(path)

	log.Info().Msgf("%+v request Menu \n", resp.Request.Body)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return errors.New(resp.String() + ", " + resp.Status())
	}

	return nil
}

func (cl Client) UpdateStopListByProducts(ctx context.Context, storeId string, products []menuModels.Product, isAvailable bool) error {
	woltProducts := cl.toStopListProducts(products, isAvailable)
	request := dto2.UpdateProducts{
		Product: woltProducts,
	}

	if _, err := cl.BulkUpdate(ctx, storeId, request); err != nil {
		return err
	}
	return nil
}

func (cl Client) UpdateStopListByProductsBulk(ctx context.Context, storeId string, products []menuModels.Product) error {
	woltProducts := cl.toStopListProductsBulk(products)
	request := dto2.UpdateProducts{
		Product: woltProducts,
	}

	if _, err := cl.BulkUpdate(ctx, storeId, request); err != nil {
		return err
	}
	return nil
}

func (cl Client) toStopListProducts(products []menuModels.Product, isAvailable bool) []dto2.UpdateProduct {
	result := make([]dto2.UpdateProduct, 0, len(products))

	for _, current := range products {
		product := dto2.UpdateProduct{
			ExtID:       current.ExtID,
			IsAvailable: &isAvailable,
		}
		result = append(result, product)
	}

	return result
}

func (cl Client) toStopListProductsBulk(products []menuModels.Product) []dto2.UpdateProduct {
	result := make([]dto2.UpdateProduct, 0, len(products))

	for _, cur := range products {
		current := cur
		product := dto2.UpdateProduct{
			ExtID:       current.ExtID,
			IsAvailable: &current.IsAvailable,
		}
		result = append(result, product)
	}

	return result
}

func (cl Client) BulkUpdate(ctx context.Context, storeID string, products dto2.UpdateProducts) (string, error) {
	products = getUniqueProducts(products)

	path := fmt.Sprintf("/venues/%s/items", storeID)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		SetBody(products).
		Patch(path)

	if err != nil {
		return "", err
	}

	log.Info().Msgf("request Wolt bulkUpdate %+v \n", resp.Request.Body)

	log.Info().Msgf("response Wolt %+v", resp.String())

	if resp.IsError() {
		log.Err(errors.New("BulkUpdate Wolt")).Msgf("response status: %s", resp.Status())
		return resp.Status(), fmt.Errorf(resp.String())
	}

	return resp.String(), nil
}

func (cl Client) BulkAttribute(ctx context.Context, storeID string, attributes dto2.UpdateAttributes) (string, error) {

	attributes = getUniqueAttributes(attributes)

	path := fmt.Sprintf("/venues/%s/options/values", storeID)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		SetBody(attributes).
		Patch(path)

	if err != nil {
		return "", err
	}

	log.Info().Msgf("request Wolt bulkAttributes %+v \n", resp.Request.Body)

	log.Info().Msgf("response Wolt %+v", resp.String())

	if resp.IsError() {
		log.Err(errors.New("BulkAttribute Wolt")).Msgf("response status: %s", resp.Status())
		return resp.Status(), fmt.Errorf(resp.String())
	}

	return resp.String(), nil
}

func (cl Client) GetMenu(ctx context.Context, storeID string) (dto2.Menu, error) {
	path := fmt.Sprintf("v2/venues/%s/menu", storeID)

	var woltMenu dto2.Menu

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		SetResult(&woltMenu).
		Get(path)

	if err != nil {
		return dto2.Menu{}, err
	}

	log.Info().Msgf("request Wolt menu by API %s \n", resp.Request.URL)
	log.Info().Msgf("response Wolt menu by API %+v", resp.String())

	if resp.IsError() {
		return dto2.Menu{}, fmt.Errorf(resp.String())
	}

	return woltMenu, nil
}

func (cl Client) UpdateMenuItemInventory(ctx context.Context, storeID string, woltInventory dto2.WoltInventory) error {
	path := fmt.Sprintf("/venues/%s/items/inventory", storeID)

	log.Info().Msgf("request (Wolt - UpdateMenuItemInventory) %+v", woltInventory)
	log.Info().Msgf("request path (Wolt - UpdateMenuItemInventory) %s", cl.restyClient.BaseURL+path)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		SetBody(&woltInventory).
		Patch(path)

	if err != nil {
		log.Err(err).Msgf("resty client (Wolt - UpdateMenuItemInventory) error: %+v", resp.Request.Body)
		return err
	}
	if resp.IsError() {
		log.Trace().Msgf("(Wolt - UpdateMenuItemInventory) request body: %+v", resp.Request.Body)
		log.Err(err).Msgf("response (Wolt - UpdateMenuItemInventory), code: %d, error: %+v", resp.StatusCode(), resp.Error())
		if resp.StatusCode() == http.StatusInternalServerError {
			log.Err(err500).Msg("")
		} else if resp.StatusCode() == http.StatusBadRequest {
			log.Err(errors.New("wolt gave bad request")).Msg("")
		}
		return fmt.Errorf("response error with status: %s", resp.Status())
	}

	log.Info().Msgf("(UpdateMenuItemInventory) for store_id: %s; response: %s", storeID, resp.String())

	return nil
}
