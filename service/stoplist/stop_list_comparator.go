package stoplist

import (
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/pkg/errors"
	"strings"
)

type idExtractor interface {
	getProductID(menuModels.Product) (string, bool)
	getAttributeID(menuModels.Attribute) (string, bool)
}

type posMenuIDExtractor struct {
}

type aggregatorMenuIDExtractor struct {
	hasVirtualStore bool
	storeID         string
}

func (s posMenuIDExtractor) getProductID(product menuModels.Product) (string, bool) {
	id := product.ExtID

	if product.ProductID != "" && product.ProductID != product.ExtID {
		id = product.ProductID
	}
	return id, true
}

func (s posMenuIDExtractor) getAttributeID(attribute menuModels.Attribute) (string, bool) {
	return attribute.ExtID, true
}

func (s aggregatorMenuIDExtractor) getProductID(product menuModels.Product) (string, bool) {
	id := product.ExtID

	if product.PosID != "" && product.PosID != product.ExtID {
		id = product.PosID
	}

	if !s.hasVirtualStore {
		return id, true
	}

	arr := strings.Split(id, "_")
	if arr[0] != s.storeID {
		return "", false
	}
	if len(arr) >= 2 {
		id = arr[1]
	}

	return id, true
}

func (s aggregatorMenuIDExtractor) getAttributeID(attribute menuModels.Attribute) (string, bool) {
	id := attribute.ExtID

	if attribute.PosID != "" && attribute.PosID != attribute.ExtID {
		id = attribute.PosID
	}

	if !s.hasVirtualStore {
		return id, true
	}

	arr := strings.Split(id, "_")
	if arr[0] != s.storeID {
		return "", false
	}
	if len(arr) >= 2 {
		id = arr[1]
	}
	return id, true
}

type stopListMenuComparator struct {
	menu            *menuModels.Menu
	stopListFromPos map[string]menuModels.StopListItem
	isByBalance     bool
	balanceLimit    float64
	idExtractor     idExtractor
}

func newStopListMenuComparator(
	menu *menuModels.Menu, stopListItems []menuModels.StopListItem, isByBalance bool, balanceLimit float64,
	idExtractor idExtractor,
) (*stopListMenuComparator, error) {
	if stopListItems == nil {
		return nil, errors.New("stopListItems is nil")
	}
	if menu == nil {
		return nil, errors.New("menu is nil")
	}
	if isByBalance && balanceLimit < 0 {
		return nil, errors.New("balance error")
	}

	existStopLists := make(map[string]menuModels.StopListItem, len(stopListItems))
	for _, item := range stopListItems {
		existStopLists[item.ProductID] = item
	}

	comparator := stopListMenuComparator{
		menu:            menu,
		stopListFromPos: existStopLists,
		isByBalance:     isByBalance,
		balanceLimit:    balanceLimit,
		idExtractor:     idExtractor,
	}

	return &comparator, nil
}

func (s stopListMenuComparator) process() ([]menuModels.Product, []menuModels.Attribute) {
	products := s.processProducts()
	attributes := s.processAttributes()
	return products, attributes
}

func (s stopListMenuComparator) processProducts() (products []menuModels.Product) {

	for _, product := range s.menu.Products {
		productID, _ := s.idExtractor.getProductID(product)

		if product.IsDisabled || product.DisabledByValidation {
			continue
		}

		if product.IsDeleted {
			product.IsAvailable = false
			product.Balance = s.stopListFromPos[productID].Balance
			products = append(products, product)
			continue
		}

		productID, isProductIDBelongsToStore := s.idExtractor.getProductID(product)
		if !isProductIDBelongsToStore {
			continue
		}

		targetAvailability := s.isProductAvailable(productID)
		product.IsAvailable = targetAvailability
		product.Balance = s.stopListFromPos[productID].Balance
		products = append(products, product)
	}
	return
}

func (s stopListMenuComparator) processAttributes() (attributes []menuModels.Attribute) {

	for _, attribute := range s.menu.Attributes {

		if attribute.IsDisabled {
			continue
		}

		if attribute.IsDeleted {
			attribute.IsAvailable = false
			attributes = append(attributes, attribute)
			continue
		}

		attributeID, isAttributeIDBelongsToStore := s.idExtractor.getAttributeID(attribute)
		if !isAttributeIDBelongsToStore {
			continue
		}

		targetAvailability := s.isAttributeAvailable(attributeID)
		attribute.IsAvailable = targetAvailability
		attributes = append(attributes, attribute)
	}
	return
}

func (s stopListMenuComparator) isProductAvailable(productID string) bool {
	stopListItem, isOnStop := s.getStopListItem(productID)
	if !isOnStop {
		return true
	}
	if s.isByBalance {
		return stopListItem.Balance > s.balanceLimit
	}
	return false
}

func (s stopListMenuComparator) isAttributeAvailable(attributeID string) bool {
	stopListItem, isOnStop := s.getStopListItem(attributeID)
	if !isOnStop {
		return true
	}
	if s.isByBalance {
		return stopListItem.Balance > s.balanceLimit
	}
	return false
}

func (s stopListMenuComparator) getStopListItem(id string) (*menuModels.StopListItem, bool) {
	if item, ok := s.stopListFromPos[id]; ok {
		return &item, true
	}
	return nil, false
}
