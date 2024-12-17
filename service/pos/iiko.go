package pos

import (
	"context"
	"fmt"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	iikoConf "github.com/kwaaka-team/orders-core/pkg/iiko/clients"
	IIKOClient "github.com/kwaaka-team/orders-core/pkg/iiko/clients/http"
	iikoModels "github.com/kwaaka-team/orders-core/pkg/iiko/models"
	"github.com/pkg/errors"
)

type iikoService struct {
	*BasePosService

	iikoClient              iikoConf.IIKO
	organizationID          string
	terminalID              string
	transportToFrontTimeout int
}

func newIikoService(bps *BasePosService, baseURL, organizationId, terminalID string, apiLogin string, transportToFrontTimeout int, CustomDomain string) (*iikoService, error) {
	if bps == nil {
		return nil, errors.Wrap(constructorError, "iikoService constructor error")
	}

	if baseURL == "" {
		return nil, errors.Wrap(constructorError, "iiko baseURL is empty")
	}

	if apiLogin == "" {
		return nil, errors.Wrap(constructorError, "iiko apiLogin is empty")
	}

	if transportToFrontTimeout < 0 {
		return nil, errors.Wrap(constructorError, "iiko iikoTransportToFrontTimeout invalid value")
	}

	if organizationId == "" {
		return nil, errors.Wrap(constructorError, "iiko organizationId is empty")
	}

	if terminalID == "" {
		return nil, errors.Wrap(constructorError, "iiko terminalId is empty")
	}
	if CustomDomain != "" {
		baseURL = CustomDomain
	}
	iikoClient, err := IIKOClient.New(&iikoConf.Config{
		Protocol: "http",
		BaseURL:  baseURL,
		ApiLogin: apiLogin,
	})
	if err != nil {
		return nil, err
	}

	return &iikoService{bps, iikoClient, organizationId, terminalID, transportToFrontTimeout}, nil
}

func (iikoSvc *iikoService) GetBalanceLimit(ctx context.Context, store coreStoreModels.Store) int {
	return store.IikoCloud.StopListBalanceLimit
}

func (iikoSvc *iikoService) IsStopListByBalance(ctx context.Context, store coreStoreModels.Store) bool {
	return store.IikoCloud.StopListByBalance
}

func (iikoSvc *iikoService) GetMenu(ctx context.Context, store coreStoreModels.Store, systemMenuInDb coreMenuModels.Menu) (coreMenuModels.Menu, error) {
	if !store.IikoCloud.IsExternalMenu {
		rsp, err := iikoSvc.iikoClient.GetMenu(ctx, store.IikoCloud.OrganizationID)
		if err != nil {
			return coreMenuModels.Menu{}, err
		}

		products, err := iikoSvc.existProducts(ctx, systemMenuInDb.Products)
		if err != nil {
			return coreMenuModels.Menu{}, err
		}
		var (
			combos     iikoModels.GetCombosResponse
			comboExist map[string]coreMenuModels.Combo
		)

		if store.IikoCloud.HasCombo {
			combos, err = iikoSvc.iikoClient.GetCombos(ctx, iikoModels.GetCombosRequest{
				OrganizationID: store.IikoCloud.OrganizationID,
			})
			if err != nil {
				return coreMenuModels.Menu{}, err
			}

			comboExist, err = iikoSvc.existCombos(ctx, systemMenuInDb.Combos)
			if err != nil {
				return coreMenuModels.Menu{}, err
			}
		}

		return iikoSvc.menuFromClient(rsp, store.Settings, products, combos, comboExist), nil
	}

	rsp, err := iikoSvc.iikoClient.GetExternalMenu(ctx, store.IikoCloud.OrganizationID, store.IikoCloud.ExternalMenuID, store.IikoCloud.PriceCategory)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	existProducts, err := iikoSvc.getExistProducts(ctx, systemMenuInDb.Products)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	var collection coreMenuModels.MenuCollection
	if len(systemMenuInDb.Collections) > 0 {
		collection = systemMenuInDb.Collections[0]
	}

	return iikoSvc.externalMenuFromClient(rsp, store.Settings, existProducts, collection, store.IikoCloud.IgnoreExternalMenuProductsWithZeroNullPrice), err
}

func (iikoSvc *iikoService) GetStopList(ctx context.Context) (coreMenuModels.StopListItems, error) {

	isExists, err := iikoSvc.isTerminalIDExistsInPos(ctx)
	if err != nil {
		return coreMenuModels.StopListItems{}, err
	}
	if !isExists {
		return coreMenuModels.StopListItems{}, fmt.Errorf("terminalID %s not found in pos", iikoSvc.terminalID)
	}

	resp, err := iikoSvc.iikoClient.GetStopList(ctx, iikoModels.StopListRequest{
		Organizations: []string{iikoSvc.organizationID},
	})

	if err != nil {
		return coreMenuModels.StopListItems{}, err
	}

	terminalItems, err := resp.Item(iikoSvc.terminalID)
	if err != nil {
		return coreMenuModels.StopListItems{}, err
	}

	stopListItems := make(coreMenuModels.StopListItems, 0, len(terminalItems.Items))
	for _, item := range terminalItems.Items {
		stopListItems = append(stopListItems, coreMenuModels.StopListItem{
			ProductID: item.ProductID,
			Balance:   item.Balance,
		})
	}

	return stopListItems, nil
}

func (iikoSvc *iikoService) isTerminalIDExistsInPos(ctx context.Context) (bool, error) {
	terminalGroupsResponse, err := iikoSvc.iikoClient.GetTerminalGroups(ctx, iikoSvc.organizationID)
	if err != nil {
		return false, err
	}
	for _, terminalGroup := range terminalGroupsResponse.TerminalGroups {
		for _, terminal := range terminalGroup.Items {
			if terminal.ID == iikoSvc.terminalID {
				return true, nil
			}
		}
	}
	return false, nil
}

func (iikoSvc *iikoService) IsAliveStatus(ctx context.Context, store coreStoreModels.Store) (bool, error) {

	result, err := iikoSvc.iikoClient.IsAlive(ctx, iikoModels.IsAliveRequest{
		OrganizationIds:  []string{store.IikoCloud.OrganizationID},
		TerminalGroupIds: []string{store.IikoCloud.TerminalID},
	})
	if err != nil {
		return false, err
	}

	if len(result.IsAliveStatus) > 0 {
		return result.IsAliveStatus[0].IsAlive, nil
	}

	return false, errors.New("iiko is alive status is empty")
}

func (iikoSvc *iikoService) MapPosStatusToSystemStatus(posStatus, currentSystemStatus string) (models.PosStatus, error) {

	switch posStatus {
	case "WAIT_SENDING":
		return models.WAIT_SENDING, nil
	case "PAYMENT_NEW":
		return models.PAYMENT_NEW, nil
	case "PAYMENT_IN_PROGRESS":
		return models.PAYMENT_IN_PROGRESS, nil
	case "PAYMENT_SUCCESS":
		return models.PAYMENT_SUCCESS, nil
	case "PAYMENT_CANCELLED":
		return models.PAYMENT_CANCELLED, nil
	case "PAYMENT_WAITING":
		return models.PAYMENT_DELETED, nil
	case "PAYMENT_DELETED":
		return models.PAYMENT_WAITING, nil
	case "CookingStarted":
		return models.COOKING_STARTED, nil
	case "WaitCooking", "ReadyForCooking", "Unconfirmed":
		return models.ACCEPTED, nil
	case "Waiting":
		return models.READY_FOR_PICKUP, nil
	case "CookingCompleted":
		return models.COOKING_COMPLETE, nil
	case "Delivered":
		return models.PICKED_UP_BY_CUSTOMER, nil
	case "OnWay":
		return models.OUT_FOR_DELIVERY, nil
	case "Closed":
		return models.CLOSED, nil
	case "NEW":
		return models.NEW, nil
	case "Cancelled":
		return models.CANCELLED_BY_POS_SYSTEM, nil
	case "Error":
		return models.FAILED, nil
	}

	return 0, models.StatusIsNotExist
}

func (iikoSvc *iikoService) AwakeTerminal(ctx context.Context, store coreStoreModels.Store) error {

	result, err := iikoSvc.iikoClient.AwakeTerminal(ctx, iikoModels.IsAliveRequest{
		OrganizationIds:  []string{store.IikoCloud.OrganizationID},
		TerminalGroupIds: []string{store.IikoCloud.TerminalID},
	})
	if err != nil {
		return err
	}

	if len(result.FailedProcessed) > 0 {
		return errors.Errorf("awake failed for terminals with id: %s", result.FailedProcessed)
	}

	return nil
}

func (iikoSvc *iikoService) CancelOrder(ctx context.Context, order models.Order, store coreStoreModels.Store) error {
	return nil
}

func (iikoSvc *iikoService) GetSeqNumber(ctx context.Context) (string, error) {
	return "", nil
}

func (iikoService *iikoService) SortStoplistItemsByIsIgnored(ctx context.Context, menu coreMenuModels.Menu, items coreMenuModels.StopListItems) (coreMenuModels.StopListItems, error) {
	itemsMap := make(map[string]coreMenuModels.StopListItem)
	for _, item := range items {
		itemsMap[item.ProductID] = item
	}

	res := coreMenuModels.StopListItems{}

	for _, product := range menu.Products {
		if item, found := itemsMap[product.ExtID]; found && !product.IsIgnored {
			res = append(res, item)
		}
	}

	for _, att := range menu.Attributes {
		if item, found := itemsMap[att.ExtID]; found && !att.IsIgnored {
			res = append(res, item)
		}
	}

	return res, nil
}

func (iikoService *iikoService) CloseOrder(ctx context.Context, posOrderId string) error {
	if err := iikoService.iikoClient.CloseOrder(ctx, posOrderId, iikoService.organizationID); err != nil {
		return err
	}

	return nil
}
