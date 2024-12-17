package http

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	"github.com/kwaaka-team/orders-core/pkg/starterapp/clients"
	"github.com/kwaaka-team/orders-core/pkg/starterapp/clients/dto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Client struct {
	restyCli *resty.Client
	BaseURL  string
	ApyKey   string
}

func NewClient(cfg *clients.Config) (*Client, error) {
	if cfg.BaseURL == "" {
		return nil, errors.New("base URL could not be empty")
	}
	if cfg.ApiKey == "" {
		return nil, errors.New("apy key could not be empty")
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
		restyCli: client,
		BaseURL:  cfg.BaseURL,
	}

	return &cl, nil
}

func (c *Client) CreateCategories(ctx context.Context, req []dto.CategoryRequest) (dto.CreateMenuResponse, error) {
	path := "/api/categories"

	var (
		response dto.CreateMenuResponse
		errResp  dto.ErrorResponse
	)

	utils.Beautify("create categories request body: %s", req)

	resp, err := c.restyCli.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetBody(&req).
		SetResult(&response).
		Post(path)
	if err != nil {
		return dto.CreateMenuResponse{}, err
	}

	log.Info().Msgf("create categories request path: %s, body: %+v", resp.Request.URL, resp.Request.Body)

	if resp.IsError() {
		return dto.CreateMenuResponse{}, fmt.Errorf("starter app client, create categories error:  %+v, %+v", errResp, resp.Error())
	}

	log.Info().Msgf("create categories response status: %d, body: %+v", resp.RawResponse.StatusCode, string(resp.Body()))

	return response, nil
}

func (c *Client) UpdateCategories(ctx context.Context, req []dto.CategoryRequest) error {
	path := "/api/categories"

	var errResp dto.ErrorResponse

	resp, err := c.restyCli.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetBody(&req).
		Put(path)
	if err != nil {
		return err
	}

	log.Info().Msgf("update categories request path: %s, body: %+v", resp.Request.URL, resp.Request.Body)

	if resp.IsError() {
		return fmt.Errorf("starter app client, update categories error: %+v, %+v", errResp, resp.Error())
	}

	log.Info().Msgf("update categories response status: %d, body: %+v", resp.RawResponse.StatusCode, string(resp.Body()))

	return nil
}

func (c *Client) CreateModifierGroups(ctx context.Context, req []dto.ModifierGroupRequest) (dto.CreateMenuResponse, error) {
	path := "/api/modifier_groups"

	var (
		response dto.CreateMenuResponse
		errResp  dto.ErrorResponse
	)

	resp, err := c.restyCli.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetBody(&req).
		SetResult(&response).
		Post(path)
	if err != nil {
		return dto.CreateMenuResponse{}, err
	}

	log.Info().Msgf("create modidfier groups request path: %s, body: %+v", resp.Request.URL, resp.Request.Body)

	if resp.IsError() {
		return dto.CreateMenuResponse{}, fmt.Errorf("starter app client, create modifier groups error: %+v, %+v", errResp, resp.Error())
	}

	log.Info().Msgf("create modifier groups response status: %d, body: %+v", resp.RawResponse.StatusCode, string(resp.Body()))

	return response, nil
}

func (c *Client) UpdateModifierGroups(ctx context.Context, req []dto.ModifierGroupRequest) error {
	path := "/api/modifier_groups"

	var errResp dto.ErrorResponse

	resp, err := c.restyCli.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetBody(&req).
		Put(path)
	if err != nil {
		return err
	}

	log.Info().Msgf("udpate modidfier groups request path: %s, body: %+v", resp.Request.URL, resp.Request.Body)

	if resp.IsError() {
		return fmt.Errorf("starter app client, update modifier groups error: %+v, %+v", errResp, resp.Error())
	}

	log.Info().Msgf("udpate modifier groups response status: %d, body: %+v", resp.RawResponse.StatusCode, string(resp.Body()))

	return nil
}

func (c *Client) CreateModifiers(ctx context.Context, req []dto.ModifiersRequest) (dto.CreateMenuResponse, error) {
	path := "/api/modifiers"

	var (
		response dto.CreateMenuResponse
		errResp  dto.ErrorResponse
	)

	resp, err := c.restyCli.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetBody(&req).
		SetResult(&response).
		Post(path)
	if err != nil {
		return dto.CreateMenuResponse{}, err
	}

	log.Info().Msgf("create modifiers request path: %s, body: %+v", resp.Request.URL, resp.Request.Body)

	if resp.IsError() {
		return dto.CreateMenuResponse{}, fmt.Errorf("starter app client, create modifiers error: %+v, %+v", errResp, resp.Error())
	}

	log.Info().Msgf("create modifiers response status: %d, body: %+v", resp.RawResponse.StatusCode, string(resp.Body()))

	return response, nil
}

func (c *Client) UpdateModifiers(ctx context.Context, req []dto.ModifiersRequest) error {
	path := "/api/modifiers"

	var errResp dto.ErrorResponse

	resp, err := c.restyCli.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetBody(&req).
		Put(path)
	if err != nil {
		return err
	}

	log.Info().Msgf("update modifiers request path: %s, body: %+v", resp.Request.URL, resp.Request.Body)

	if resp.IsError() {
		return fmt.Errorf("starter app client, update modifiers error: %+v, %+v", errResp, resp.Error())
	}

	log.Info().Msgf("update modifiers response status: %d, body: %+v", resp.RawResponse.StatusCode, string(resp.Body()))

	return nil
}

func (c *Client) CreateMeals(ctx context.Context, req []dto.MealRequest) (dto.CreateMenuResponse, error) {
	path := "/api/meals"

	var (
		response dto.CreateMenuResponse
		errResp  dto.ErrorResponse
	)
	resp, err := c.restyCli.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetBody(&req).
		SetResult(&response).
		Post(path)
	if err != nil {
		return dto.CreateMenuResponse{}, err
	}

	utils.Beautify("Create meals req: ", req)

	log.Info().Msgf("create meals request path: %s, body: %+v", resp.Request.URL, resp.Request.Body)

	if resp.IsError() {
		return dto.CreateMenuResponse{}, fmt.Errorf("starter app client, create meals error: %+v, %+v", errResp, resp.Error())
	}

	log.Info().Msgf("create meals response status: %d, body: %+v", resp.RawResponse.StatusCode, string(resp.Body()))

	return response, nil
}

func (c *Client) UpdateMeals(ctx context.Context, req []dto.MealRequest) error {
	path := "/api/meals"

	var errResp dto.ErrorResponse

	resp, err := c.restyCli.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetBody(&req).
		Put(path)
	if err != nil {
		return err
	}

	log.Info().Msgf("update meals request path: %s, body: %+v", resp.Request.URL, resp.Request.Body)

	if resp.IsError() {
		return fmt.Errorf("starter app client, update meals error: %+v, %+v", errResp, resp.Error())
	}

	log.Info().Msgf("update meals response status: %d, body: %+v", resp.RawResponse.StatusCode, string(resp.Body()))

	return nil
}

func (c *Client) CreateMealOffers(ctx context.Context, req []dto.MealOfferRequest, shopID int) (dto.CreateMenuResponse, error) {
	path := fmt.Sprintf("/api/shop/%d/meals", shopID)

	var (
		response dto.CreateMenuResponse
		errResp  dto.ErrorResponse
	)

	utils.Beautify("CreateMealOffers body: ", req)

	resp, err := c.restyCli.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetBody(&req).
		SetResult(&response).
		Post(path)
	if err != nil {
		return dto.CreateMenuResponse{}, err
	}

	log.Info().Msgf("create meal offers request path: %s, body: %+v", resp.Request.URL, resp.Request.Body)

	if resp.IsError() {
		log.Info().Msgf("create meal offers error: %+v, %+v", errResp, *resp)
		return dto.CreateMenuResponse{}, fmt.Errorf("starter app client, create meal offers error: %+v, %+v", errResp, resp.Error())
	}

	log.Info().Msgf("create meal offers response status: %d, body: %+v", resp.RawResponse.StatusCode, string(resp.Body()))

	return response, nil
}

func (c *Client) UpdateMealOffers(ctx context.Context, req []dto.MealOfferRequest, shopID int) error {
	path := fmt.Sprintf("/api/shop/%d/meals", shopID)

	var errResp dto.ErrorResponse

	utils.Beautify("Update meals offers req: ", req)

	resp, err := c.restyCli.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetBody(&req).
		Put(path)
	if err != nil {
		return err
	}

	log.Info().Msgf("update meal offers request path: %s, body: %+v", resp.Request.URL, resp.Request.Body)

	if resp.IsError() {
		return fmt.Errorf("starter app client, update meal offers error: %+v, %+v", errResp, resp.Error())
	}

	log.Info().Msgf("update meal offers response status: %d, body: %+v", resp.RawResponse.StatusCode, string(resp.Body()))

	return nil
}

func (c *Client) CreateModifierOffers(ctx context.Context, req []dto.ModifierOfferRequest) (dto.CreateMenuResponse, error) {
	path := "/api/modifier_offer"

	var (
		response dto.CreateMenuResponse
		errResp  dto.ErrorResponse
	)

	utils.Beautify("CreateModifierOffers body: ", req)

	resp, err := c.restyCli.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetBody(&req).
		SetResult(&response).
		Post(path)
	if err != nil {
		return dto.CreateMenuResponse{}, err
	}

	log.Info().Msgf("create modifier offers request path: %s, body: %+v", resp.Request.URL, resp.Request.Body)

	if resp.IsError() {
		return dto.CreateMenuResponse{}, fmt.Errorf("starter app client, create modifier offers error: %+v, %+v", errResp, resp.Error())
	}

	log.Info().Msgf("create modifier offers response status: %d, body: %+v", resp.RawResponse.StatusCode, string(resp.Body()))

	return response, nil
}

func (c *Client) UpdateModifierOffers(ctx context.Context, req []dto.ModifierOfferRequest) error {
	path := "/api/modifier_offer"

	var errResp dto.ErrorResponse

	utils.Beautify("Update modifier offers req: ", req)

	resp, err := c.restyCli.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetBody(&req).
		Put(path)
	if err != nil {
		return err
	}

	log.Info().Msgf("update modifier offers request path: %s, body: %+v", resp.Request.URL, resp.Request.Body)

	if resp.IsError() {
		return fmt.Errorf("starter app client, update modifier offers error: %+v, %+v", errResp, resp.Error())
	}

	log.Info().Msgf("update modifier offers response status: %d, body: %+v", resp.RawResponse.StatusCode, string(resp.Body()))

	return nil
}

func (c *Client) SendOrderErrorNotification(ctx context.Context, req dto.SendOrderErrorNotificationRequest, orderID string) error {
	path := fmt.Sprintf("/api/order/%s/error", orderID)

	var errResp dto.ErrorResponse

	resp, err := c.restyCli.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetBody(&req).
		Post(path)
	if err != nil {
		return err
	}

	log.Info().Msgf("send order error notification request path: %s, body: %+v", resp.Request.URL, resp.Request.Body)

	if resp.IsError() {
		return fmt.Errorf("starter app client, send order error notification error: %+v, %+v", errResp, resp.Error())
	}

	log.Info().Msgf("send order error notification response status: %d, body: %+v", resp.RawResponse.StatusCode, string(resp.Body()))

	return nil
}

func (c *Client) ChangeOrderStatus(ctx context.Context, req dto.ChangeOrderStatusRequest, orderID string) error {
	path := fmt.Sprintf("/api/order/%s/status", orderID)

	var errResp dto.ErrorResponse

	resp, err := c.restyCli.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResp).
		SetBody(&req).
		Patch(path)
	if err != nil {
		return err
	}

	log.Info().Msgf("change order status request path: %s, body: %+v", resp.Request.URL, resp.Request.Body)

	if resp.IsError() {
		return fmt.Errorf("starter app client, change order status error: %+v, %+v", errResp, resp.Error())
	}

	log.Info().Msgf("change order status response status: %d, body: %+v", resp.RawResponse.StatusCode, string(resp.Body()))

	return nil
}
