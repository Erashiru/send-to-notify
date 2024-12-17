package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_modifier_groups"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_modifier_schema_details"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_modifier_schemas"
	"net/http"
	"strings"
)

func (rkeeper *client) GetMenuModifierSchemas(ctx context.Context) (get_modifier_schemas.RK7QueryResult, error) {
	path := "/rk7api/v0/xmlinterface.xml"

	var (
		result  get_modifier_schemas.RK7QueryResult
		request = `
			<RK7Query>
			<RK7CMD CMD="GetRefData" RefName="MODISCHEMES" onlyActive="1"/>
			</RK7Query>`
	)

	response, err := rkeeper.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(strings.NewReader(request)).
		SetResult(&result).
		Post(path)
	if err != nil {
		return get_modifier_schemas.RK7QueryResult{}, err
	}

	if response.IsError() || response.StatusCode() >= http.StatusBadRequest {
		return get_modifier_schemas.RK7QueryResult{}, fmt.Errorf("get menu modifier schemas response error: %v", response.Error())
	}

	return result, nil
}

func (rkeeper *client) GetMenuModifierSchemaDetails(ctx context.Context) (get_modifier_schema_details.RK7QueryResult, error) {
	path := "/rk7api/v0/xmlinterface.xml"

	var (
		result  get_modifier_schema_details.RK7QueryResult
		request = `
			<RK7Query>
			<RK7CMD CMD="GetRefData" RefName="MODISCHEMEDETAILS" onlyActive="1"/>
			</RK7Query>`
	)

	response, err := rkeeper.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(strings.NewReader(request)).
		SetResult(&result).
		Post(path)
	if err != nil {
		return get_modifier_schema_details.RK7QueryResult{}, err
	}

	if response.IsError() || response.StatusCode() >= http.StatusBadRequest {
		return get_modifier_schema_details.RK7QueryResult{}, fmt.Errorf("get menu modifier schema details response error: %v", response.Error())
	}

	return result, nil
}

func (rkeeper *client) GetMenuModifierGroups(ctx context.Context) (get_modifier_groups.RK7QueryResult, error) {
	path := "/rk7api/v0/xmlinterface.xml"

	var (
		result  get_modifier_groups.RK7QueryResult
		request = `
			<RK7Query>
			<RK7CMD CMD="GetRefData" RefName="MODIGROUPS" onlyActive="1"/>
			</RK7Query>
		`
	)

	response, err := rkeeper.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(strings.NewReader(request)).
		SetResult(&result).
		Post(path)
	if err != nil {
		return get_modifier_groups.RK7QueryResult{}, err
	}

	if response.IsError() || response.StatusCode() >= http.StatusBadRequest {
		return get_modifier_groups.RK7QueryResult{}, fmt.Errorf("get menu modifier groups response error: %v", response.Error())
	}

	return result, nil
}
