package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/core/glovo/models"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	"github.com/kwaaka-team/orders-core/pkg/glovo/clients"
	dto2 "github.com/kwaaka-team/orders-core/pkg/glovo/clients/dto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

const storeAddressHeader = "Glovo-Store-Address-External-Id"

// [documentation]: https://api-docs.glovoapp.com/partners/index.html

type Client struct {
	ApiKey             string
	restyClient        *resty.Client
	StoreID            string
	BaseUrl            string
	Username, Password string
}

func NewClient(cfg *clients.Config) (*Client, error) {

	if cfg.BaseURL == "" {
		return nil, errors.New("base URL could not be empty")
	}

	client := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetHeaders(map[string]string{
			contentTypeHeader: jsonType,
			acceptHeader:      jsonType,
			authHeader:        cfg.ApiKey,
		}).
		SetRetryCount(retriesNumber).
		SetRetryWaitTime(retriesWaitTime)

	cl := Client{
		restyClient: client,
		ApiKey:      cfg.ApiKey,
		StoreID:     cfg.StoreID,
		BaseUrl:     cfg.BaseURL,
	}

	return &cl, nil
}

func (cl Client) UpdateStopListByProducts(ctx context.Context, storeId string, products []menuModels.Product, isAvailable bool) (string, error) {
	glovoProducts := cl.toStopListProducts(products, isAvailable)
	request := dto2.BulkUpdateRequest{
		Products: glovoProducts,
	}

	trID, err := cl.BulkUpdate(ctx, storeId, request)
	if err != nil {
		return "", err
	}
	return trID, nil
}

func (cl Client) UpdateStopListByProductsBulk(ctx context.Context, storeId string, products []menuModels.Product) (string, error) {
	glovoProducts := cl.toStopListProductsBulk(products)
	request := dto2.BulkUpdateRequest{
		Products: glovoProducts,
	}

	trID, err := cl.BulkUpdate(ctx, storeId, request)
	if err != nil {
		return "", err
	}
	return trID, nil
}

func (cl Client) toStopListProducts(products []menuModels.Product, isAvailable bool) []dto2.Product {
	result := make([]dto2.Product, 0, len(products))

	for _, current := range products {
		product := dto2.Product{
			ID:        current.ExtID,
			Available: &isAvailable,
		}
		result = append(result, product)
	}

	return result
}

func (cl Client) toStopListProductsBulk(products []menuModels.Product) []dto2.Product {
	result := make([]dto2.Product, 0, len(products))

	for _, cur := range products {
		current := cur
		product := dto2.Product{
			ID:        current.ExtID,
			Available: &current.IsAvailable,
		}
		result = append(result, product)
	}

	return result
}

func (cl Client) UpdateStopListByAttributesBulk(ctx context.Context, storeId string, attributes []menuModels.Attribute) (string, error) {
	glovoAttributes := cl.toStopListAttributesBulk(attributes)
	request := dto2.BulkUpdateRequest{
		Attributes: glovoAttributes,
	}

	trID, err := cl.BulkUpdate(ctx, storeId, request)
	if err != nil {
		return "", err
	}
	return trID, nil
}

func (cl Client) toStopListAttributesBulk(attributes []menuModels.Attribute) []dto2.Attribute {
	result := make([]dto2.Attribute, 0, len(attributes))

	for i := range attributes {
		current := attributes[i]
		attribute := dto2.Attribute{
			ID:          current.ExtID,
			Available:   current.IsAvailable,
			PriceImpact: current.Price,
		}
		result = append(result, attribute)
	}

	return result
}

func (cl Client) BulkUpdate(ctx context.Context, storeId string, request dto2.BulkUpdateRequest) (string, error) {

	path := fmt.Sprintf("/webhook/stores/%s/menu/updates", storeId)

	var (
		response dto2.BulkUpdateResponse
		errResp  dto2.CustomError
	)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&request).
		SetError(&errResp).
		SetResult(&response).
		Post(path)
	if err != nil {
		return "", err
	}

	utils.Beautify("request Glovo bulkUpdate", request)

	log.Info().Msgf("response Glovo %+v", resp.String())

	if resp.IsError() {
		return "", fmt.Errorf("glovo cli: %s, %s", errResp.Msg, resp.Error())
	}

	return response.TransactionID, nil
}

func (cl Client) UploadMenu(ctx context.Context, req dto2.UploadMenuRequest) (dto2.UploadMenuResponse, error) {

	var (
		response dto2.UploadMenuResponse
		errResp  dto2.CustomError
	)

	path := fmt.Sprintf("/webhook/stores/%s/menu", req.StoreId)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&req).
		SetError(&errResp).
		SetResult(&response).
		Post(path)

	if err != nil {
		return dto2.UploadMenuResponse{}, err
	}

	if resp.IsError() {
		return dto2.UploadMenuResponse{}, fmt.Errorf("glovo cli: %s", errResp.Msg)
	}

	return response, nil
}

// AcceptOrder - adjusted pick up time in UTC
func (cl Client) AcceptOrder(ctx context.Context, storeId, orderId string, adjustedPickUpTime time.Time) (dto2.Response, error) {
	var response dto2.Response

	path := fmt.Sprintf("/api/v0/integrations/orders/%s/accept", orderId)

	request := dto2.AcceptOrderRequest{
		CommittedPreparationTime: adjustedPickUpTime,
	}

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetHeader(storeAddressHeader, storeId).
		SetBody(&request).
		SetResult(&response).
		Put(path)
	if err != nil {
		return dto2.Response{}, err
	}

	if resp.IsError() {
		return dto2.Response{}, fmt.Errorf("accept glovo order status %w - status: %s, body: %s", ErrBadRequest, resp.Status(), string(resp.Body()))
	}

	return response, nil
}

// MarkOrderAsReady as ready for pickup
func (cl Client) MarkOrderAsReady(ctx context.Context, storeId, orderId string) (dto2.Response, error) {
	var response dto2.Response

	path := fmt.Sprintf("/api/v0/integrations/orders/%s/ready_for_pickup", orderId)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetHeader(storeAddressHeader, storeId).
		SetResult(&response).
		Put(path)
	if err != nil {
		return dto2.Response{}, err
	}

	if resp.IsError() {
		return dto2.Response{}, fmt.Errorf("mark glovo order as ready for pickup status %w - status: %s, body: %s", ErrBadRequest, resp.Status(), string(resp.Body()))
	}

	return response, nil
}

// MarkOrderAsOutForDelivery as out for delivery
func (cl Client) MarkOrderAsOutForDelivery(ctx context.Context, storeId, orderId string) (dto2.Response, error) {
	var response dto2.Response

	path := fmt.Sprintf("/api/v0/integrations/orders/%s/out_for_delivery", orderId)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetHeader(storeAddressHeader, storeId).
		SetResult(&response).
		Put(path)
	if err != nil {
		return dto2.Response{}, err
	}

	if resp.IsError() {
		return dto2.Response{}, fmt.Errorf("mark glovo order as out for delivery status %w - status: %s, body: %s", ErrBadRequest, resp.Status(), string(resp.Body()))
	}

	return response, nil
}

// MarkOrderAsCustomerPickedUp as ready for pickup
func (cl Client) MarkOrderAsCustomerPickedUp(ctx context.Context, storeId, orderId string) (dto2.Response, error) {
	var response dto2.Response

	path := fmt.Sprintf("/api/v0/integrations/orders/%s/customer_picked_up", orderId)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetHeader(storeAddressHeader, storeId).
		SetResult(&response).
		Put(path)
	if err != nil {
		return dto2.Response{}, err
	}

	if resp.IsError() {
		return dto2.Response{}, fmt.Errorf("mark glovo order as customer picked up status %w - status: %s, body: %s", ErrBadRequest, resp.Status(), string(resp.Body()))
	}

	return response, nil
}

func (cl Client) UpdateOrderStatus(ctx context.Context, request dto2.OrderUpdateRequest) (dto2.OrderUpdateResponse, error) {

	var response dto2.OrderUpdateResponse

	path := fmt.Sprintf("/webhook/stores/%s/orders/%v/status", request.StoreID, request.ID)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&request).
		SetResult(&response).
		Put(path)
	if err != nil {
		return dto2.OrderUpdateResponse{}, err
	}

	if resp.IsError() {
		return dto2.OrderUpdateResponse{}, fmt.Errorf("update glovo order status %w - status: %s, body: %s", ErrBadRequest, resp.Status(), string(resp.Body()))
	}

	return response, nil
}

func (cl Client) ModifyOrderProduct(ctx context.Context, request models.ModifyOrderProductRequest) (*models.Order, error) {

	response := models.Order{}

	path := fmt.Sprintf("/webhook/stores/%s/orders/%v/replace_products", cl.StoreID, request.ID)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&request).
		SetResult(&response).
		Post(path)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, errors.New(resp.Status() + " " + string(resp.Body()))
	}

	return &response, nil
}

// ModifyProduct to update price & availability product in GLOVO.
func (cl Client) ModifyProduct(ctx context.Context, req dto2.ProductModifyRequest) (dto2.ProductModifyResponse, error) {

	path := fmt.Sprintf("/webhook/stores/%s/products/%s", req.StoreID, req.ID)

	var (
		res     dto2.ProductModifyResponse
		errResp dto2.CustomError
	)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&res).
		SetError(&errResp).
		Patch(path)

	if err != nil {
		return dto2.ProductModifyResponse{}, err
	}

	if resp.StatusCode() == 400 {
		return dto2.ProductModifyResponse{}, fmt.Errorf("product id %s from store id %s - %w", req.ID, req.StoreID, ErrNotExist)
	}

	if resp.IsError() {
		return dto2.ProductModifyResponse{}, errResp.Error()
	}

	return res, nil
}

func (cl Client) OpenStore(ctx context.Context, req dto2.StoreManageRequest) error {
	path := fmt.Sprintf("/webhook/stores/%s/closing", req.StoreID)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		Delete(path)

	if err != nil {
		return err
	}

	if resp.StatusCode() != 204 {
		return fmt.Errorf("something went wrong while closing: store_id - %s, err - %s", req.StoreID, string(resp.Body()))
	}

	if resp.IsError() {
		return err
	}

	return nil
}

func (cl Client) CloseStore(ctx context.Context, req dto2.StoreManageRequest) error {

	path := fmt.Sprintf("/webhook/stores/%s/closing", req.StoreID)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		SetBody(req).
		Put(path)

	if err != nil {
		return err
	}

	if resp.StatusCode() != 204 {
		return fmt.Errorf("something went wrong while closing: store_id - %s, err - %s", req.StoreID, string(resp.Body()))
	}

	if resp.IsError() {
		return err
	}
	return nil
}

func (cl Client) StoreStatus(ctx context.Context, storeID string) (dto2.StoreStatusResponse, error) {

	status, err := cl.GetStoreStatus(ctx, storeID)
	if err != nil {
		return dto2.StoreStatusResponse{}, err
	}

	schedule, err := cl.GetStoreSchedule(ctx, storeID)
	if err != nil {
		return dto2.StoreStatusResponse{}, err
	}

	location, err := time.LoadLocation(schedule.Timezone)
	if err != nil {
		return dto2.StoreStatusResponse{}, err
	}

	var Response dto2.StoreStatusResponse

	currentTime := time.Now().In(location)

	if status.Until == "" && checkSchedule(currentTime, schedule) {
		Response.IsActive = true
	}
	Response.StoreID = storeID

	return Response, nil
}

func (cl Client) GetStoreStatus(ctx context.Context, storeID string) (models.StoreStatus, error) {
	path := fmt.Sprintf("/webhook/stores/%s/closing", storeID)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		Get(path)

	if err != nil {
		return models.StoreStatus{}, err
	}

	if resp.StatusCode() != 200 {
		return models.StoreStatus{}, fmt.Errorf("something went wrong while checking status: store_id - %s, err - %s", storeID, string(resp.Body()))
	}

	var status models.StoreStatus
	if err := json.Unmarshal(resp.Body(), &status); err != nil {
		return models.StoreStatus{}, err
	}
	return status, nil
}

func (cl Client) GetStoreSchedule(ctx context.Context, storeID string) (dto2.StoreScheduleResponse, error) {
	path := fmt.Sprintf("/webhook/stores/%s/schedule", storeID)
	resp, err := cl.restyClient.R().
		SetContext(ctx).
		Get(path)

	if err != nil {
		return dto2.StoreScheduleResponse{}, err
	}

	if resp.StatusCode() != 200 {
		return dto2.StoreScheduleResponse{}, fmt.Errorf("something went wrong while getting store schedule: store_id - %s, err - %s", storeID, string(resp.Body()))
	}
	var response dto2.StoreScheduleResponse

	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return dto2.StoreScheduleResponse{}, err
	}

	return response, nil
}

// ModifyAttribute to update price & availability attribute in GLOVO.
func (cl Client) ModifyAttribute(ctx context.Context, req dto2.AttributeModifyRequest) (dto2.AttributeModifyResponse, error) {

	path := fmt.Sprintf("/%s/attributes/%s", req.StoreID, req.ID)

	var (
		res     dto2.AttributeModifyResponse
		errResp dto2.CustomError
	)

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&res).
		SetError(&errResp).
		Patch(path)

	if err != nil {
		return dto2.AttributeModifyResponse{}, err
	}

	if resp.StatusCode() == 400 {
		return dto2.AttributeModifyResponse{}, fmt.Errorf("attribute id %s from store id %s - %w", req.ID, req.StoreID, ErrNotExist)
	}

	if resp.IsError() {
		return dto2.AttributeModifyResponse{}, errResp.Error()
	}

	return res, nil
}

func (cl Client) VerifyMenu(ctx context.Context, storeId, trxId string) (dto2.UploadMenuResponse, error) {
	path := fmt.Sprintf("/webhook/stores/%s/menu/%s", storeId, trxId)

	var res dto2.UploadMenuResponse

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		SetResult(&res).
		Get(path)

	if err != nil {
		return dto2.UploadMenuResponse{}, err
	}

	if resp.StatusCode() == 400 || resp.IsError() {
		return dto2.UploadMenuResponse{}, fmt.Errorf("transaction_id %s error: %w", trxId, ErrInvalid)
	}

	return res, nil
}

func (cl Client) ValidateMenu(ctx context.Context, req dto2.ValidateMenuRequest) (dto2.ValidateMenuResponse, error) {
	path := "/paris/menu/validate"

	var res dto2.ValidateMenuResponse

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		SetResult(&res).
		SetError(&res).
		SetBody(req).
		Post(path)

	if err != nil {
		return dto2.ValidateMenuResponse{}, err
	}

	if resp.StatusCode() == 400 {
		return res, fmt.Errorf("validate menu: %w - %s", ErrBadRequest, string(resp.Body()))
	}

	if resp.IsError() {
		return dto2.ValidateMenuResponse{}, fmt.Errorf("validate menu: %w - %s", ErrBadRequest, string(resp.Body()))
	}

	return res, nil
}

func checkSchedule(currentTime time.Time, schedule dto2.StoreScheduleResponse) bool {
	currentDayOfWeek := int(currentTime.Weekday())
	if currentDayOfWeek == 0 {
		currentDayOfWeek = 7
	}

	for _, daySchedule := range schedule.Schedule {
		if daySchedule.DayOfWeek == currentDayOfWeek {
			for _, timeSlot := range daySchedule.TimeSlots {
				openingTime, _ := time.Parse("15:04", timeSlot.Opening)
				closingTime, _ := time.Parse("15:04", timeSlot.Closing)
				currentTimeParsed, _ := time.Parse("15:04", currentTime.Format("15:04"))
				if closingTime.Before(openingTime) {
					if currentTimeParsed.After(openingTime) || currentTimeParsed.Before(closingTime) {
						return true
					}
				}
				if currentTimeParsed.After(openingTime) && currentTimeParsed.Before(closingTime) {
					return true
				}

			}
		}
	}

	return false
}

func (cl Client) CreateBusyMode(ctx context.Context, storeId string, additionalPreparationTimeInMinutes int) error {

	path := fmt.Sprintf("/webhook/stores/%s/busy_mode", storeId)

	var errResp dto2.CustomError

	type BusyModeRequest struct {
		AdditionalPreparationTimeInMinutes int `json:"additionalPreparationTimeInMinutes"`
	}

	request := BusyModeRequest{
		AdditionalPreparationTimeInMinutes: additionalPreparationTimeInMinutes,
	}

	type BusyModeResponse struct {
		AdditionalPreparationTimeInMinutes int       `json:"additionalPreparationTimeInMinutes"`
		EndingAt                           time.Time `json:"endingAt"`
	}
	var response BusyModeResponse

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(&request).
		SetError(&errResp).
		SetResult(&response).
		Put(path)
	if err != nil {
		return err
	}

	utils.Beautify("request Glovo CreateBusyMode", request)
	log.Info().Msgf("response Glovo %+v", resp.String())

	if resp.IsError() {
		return fmt.Errorf("glovo cli: %s, %s", errResp.Msg, resp.Error())
	}

	return nil
}

func (cl Client) GetBusyMode(ctx context.Context, storeId string) (int, error) {

	path := fmt.Sprintf("/webhook/stores/%s/busy_mode", storeId)

	var errResp dto2.CustomError

	type BusyModeResponse struct {
		AdditionalPreparationTimeInMinutes int       `json:"additionalPreparationTimeInMinutes"`
		EndingAt                           time.Time `json:"endingAt"`
	}
	var response BusyModeResponse

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetResult(&response).
		Get(path)
	if err != nil {
		return 0, err
	}

	log.Info().Msgf("response Glovo %+v", resp.String())

	if resp.IsError() {
		return 0, fmt.Errorf("glovo cli: %s", errResp.Msg)
	}

	return response.AdditionalPreparationTimeInMinutes, nil
}

func (cl Client) DeleteBusyMode(ctx context.Context, storeId string) error {

	path := fmt.Sprintf("/webhook/stores/%s/busy_mode", storeId)

	var errResp dto2.CustomError

	resp, err := cl.restyClient.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		Delete(path)
	if err != nil {
		return err
	}

	log.Info().Msgf("response Glovo %+v", resp.String())

	if resp.IsError() {
		return fmt.Errorf("glovo cli: %s", errResp.Msg)
	}

	return nil
}
