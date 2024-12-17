package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/talabat/models"
	"github.com/pkg/errors"
)

func (c *Client) SubmitCatalog(ctx context.Context, req models.SubmitCatalogRequest) (models.SubmitCatalogResponse, error) {
	path := fmt.Sprintf("/v2/chains/%s/catalog", req.ChainCode)

	var errResponse models.AuthMWErrorResponse
	var res models.SubmitCatalogResponse

	resp, err := c.cli.R().
		SetContext(ctx).
		SetError(&errResponse).
		SetResult(&res).
		SetHeader("Content-Type", "application/json").
		SetBody(req).
		Put(path)

	if err != nil {
		return models.SubmitCatalogResponse{}, err
	}

	if resp.IsError() {
		if resp.IsError() {
			return models.SubmitCatalogResponse{}, errors.New(fmt.Sprintf("code: %s, message: %s", errResponse.Code, errResponse.Message))
		}
	}

	return res, nil
}
