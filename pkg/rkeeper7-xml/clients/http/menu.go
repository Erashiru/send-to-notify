package http

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_menu_by_category"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models/get_menu_modifiers_response"
	"net/http"
	"strings"
)

func (rkeeper *client) GetMenuItems(ctx context.Context) (models.MenuRK7QueryResult, error) {
	path := "/rk7api/v0/xmlinterface.xml"

	var (
		result  models.MenuRK7QueryResult
		request = `
			<RK7Query>
			<RK7CMD CMD="GetRefData" RefName="MENUITEMS" onlyActive="1"/>
			</RK7Query>`
	)

	response, err := rkeeper.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(strings.NewReader(request)).
		SetResult(&result).
		Post(path)
	if err != nil {
		return models.MenuRK7QueryResult{}, err
	}

	if response.IsError() || response.StatusCode() >= http.StatusBadRequest {
		return models.MenuRK7QueryResult{}, fmt.Errorf("get menu items response: status code:%d; error: %v", response.StatusCode(), response.Error())
	}

	return result, nil
}

func (rkeeper *client) GetMenuModifiers(ctx context.Context) (get_menu_modifiers_response.RK7QueryResult, error) {
	path := "/rk7api/v0/xmlinterface.xml"

	var (
		result  get_menu_modifiers_response.RK7QueryResult
		request = `
			<RK7Query>
			<RK7CMD CMD="GetRefData" RefName="MODIFIERS" onlyActive="1"/>
			</RK7Query>`
	)

	response, err := rkeeper.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(strings.NewReader(request)).
		SetResult(&result).
		Post(path)
	if err != nil {
		return get_menu_modifiers_response.RK7QueryResult{}, err
	}

	if response.IsError() || response.StatusCode() >= http.StatusBadRequest {
		return get_menu_modifiers_response.RK7QueryResult{}, fmt.Errorf("get menu modifiers response error: %v", response.Error())
	}

	return result, nil
}

func (rkeeper *client) GetMenuByCategory(ctx context.Context) (get_menu_by_category.RK7QueryResult, error) {
	path := "/rk7api/v0/xmlinterface.xml"
	var (
		result  get_menu_by_category.RK7QueryResult
		request = fmt.Sprintf(`
			<RK7Query >
    <RK7Command CMD="GetRefData" 
						 RefName="ClassificatorGroups" 
						 WithChildItems="%d" 
						 OnlyActive="1" 
						 RefItemIdent="%d"  	 PropMask="%s">
 </RK7Command>
    <RK7Command 
						 CMD="GetRefData" 
						 RefName="MENUITEMS" 
						 OnlyActive="true" 
						 WithChildItems="%d" 
						 WithMacroProp="true" 
						 PropMask="%s">
       %s
    </RK7Command>
</RK7Query>
		`, rkeeper.childItems, rkeeper.classificatorItemIdent, rkeeper.classificatorPropMask, rkeeper.childItems, rkeeper.menuItemsPropMask, rkeeper.propFilter)
	)

	response, err := rkeeper.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(strings.NewReader(request)).
		SetResult(&result).
		Post(path)
	if err != nil {
		return get_menu_by_category.RK7QueryResult{}, err
	}

	if response.IsError() || response.StatusCode() >= http.StatusBadRequest {
		return get_menu_by_category.RK7QueryResult{}, fmt.Errorf("get menu by category response error: %v", response.Error())
	}

	return result, nil
}

func (rkeeper *client) GetOrderMenu(ctx context.Context) (models.OrderMenuRK7QueryResult, error) {
	path := "/rk7api/v0/xmlinterface.xml"

	var (
		result  models.OrderMenuRK7QueryResult
		request = fmt.Sprintf(`
			<RK7Query>
			<RK7CMD CMD="GetOrderMenu" >
			<Station code="%s"/>
			</RK7CMD>
			</RK7Query>
		`, rkeeper.stationCode)
	)

	response, err := rkeeper.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(strings.NewReader(request)).
		SetResult(&result).
		Post(path)
	if err != nil {
		return models.OrderMenuRK7QueryResult{}, err
	}

	if response.IsError() || response.StatusCode() >= http.StatusBadRequest {
		return models.OrderMenuRK7QueryResult{}, fmt.Errorf("get order menu response error: %v", response.Error())
	}

	return result, nil
}
