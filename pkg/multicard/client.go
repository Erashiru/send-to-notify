package multicard

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/pkg/multicard/dto"
	"github.com/redis/go-redis/v9"
)

const AuthKey = "multicard_token"

type Client interface {
	CreatePaymentInvoice(ctx context.Context, req dto.CreatePaymentInvoiceRequest) (*dto.CreatePaymentInvoiceResponse, error)
	ReturnFunds(ctx context.Context, uuid string) (*dto.Payment, *dto.ErrResponse, error)
}

type ServiceImpl struct {
	cli     *resty.Client
	storeId int
}

func NewClient(baseUrl, applicationId, secret, token string, storeId int) (Client, error) {
	cli, err := auth(context.TODO(), baseUrl, applicationId, secret, token)
	if err != nil {
		return nil, err
	}
	return &ServiceImpl{cli: cli, storeId: storeId}, nil
}

func GetAuthKey(ctx context.Context, redis *redis.Client) (string, error) {
	cmd := redis.Get(ctx, AuthKey)
	if cmd.Err() != nil {
		return "", cmd.Err()
	}

	return cmd.Val(), nil
}

func auth(ctx context.Context, baseUrl, applicationId, secret, token string) (*resty.Client, error) {
	const op = "multicard.auth"

	var (
		url            = "/auth"
		resp           = dto.AuthResponse{}
		multicardToken = ""
		cli            = resty.New().SetBaseURL(baseUrl)
		authReq        = dto.AuthRequest{
			ApplicationId: applicationId,
			Secret:        secret,
		}
	)

	switch token == "" {
	case true:
		r, err := cli.R().
			SetBody(authReq).
			SetContext(ctx).
			SetResult(&resp).
			Post(url)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if r.IsError() {
			return nil, fmt.Errorf("%s: %s", op, "could not authorize to multicard")
		}

		multicardToken = resp.Token
	default:
		multicardToken = token
	}

	cli.SetAuthToken(multicardToken)

	return cli, nil
}

// CreatePaymentInvoice TODO: check successions by status from resp
func (s *ServiceImpl) CreatePaymentInvoice(ctx context.Context, req dto.CreatePaymentInvoiceRequest) (*dto.CreatePaymentInvoiceResponse, error) {
	const op = "multicard.createPaymentInvoice"

	var (
		url = "/payment/invoice"
		res = dto.CreatePaymentInvoiceResponse{}
	)

	req.StoreId = s.storeId

	resp, err := s.cli.R().
		SetBody(&req).
		SetContext(ctx).
		SetResult(&res).
		Post(url)
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("%s: %s", op, "create payment invoice received fail status code")
	}

	return &res, nil
}

// ReturnFunds TODO: check successions by status from resp
func (s *ServiceImpl) ReturnFunds(ctx context.Context, uuid string) (*dto.Payment, *dto.ErrResponse, error) {
	const op = "multicard.returnFunds"

	var (
		url         = fmt.Sprintf("/payment/%s", uuid)
		errResponse dto.ErrResponse
		result      = struct {
			Success bool        `json:"success"`
			Data    dto.Payment `json:"data"`
		}{}
	)

	resp, err := s.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetError(&errResponse).
		SetResult(&result).
		Delete(url)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", op, err)
	}

	if resp.IsError() {
		return nil, &errResponse, fmt.Errorf("%s: %s", op, fmt.Sprintf("could not return funds: %s", errResponse.Error.Details))
	}

	return &result.Data, &errResponse, nil
}
