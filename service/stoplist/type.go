package stoplist

import (
	"context"
	"errors"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/service/menu"
)

type stopListType interface {
	filterProducts(products []menuModels.Product) []menuModels.Product
	updateDisabledStatusByProductIDsInDatabase(ctx context.Context, menuID string, products []menuModels.Product) error

	filterAttributes(attributes []menuModels.Attribute) []menuModels.Attribute
	updateDisabledStatusByAttributeIDsInDatabase(ctx context.Context, menuID string, attributes []menuModels.Attribute) error
}

type cronStopList struct {
	menuService *menu.Service
}

func newCronStopList(menuService *menu.Service) (*cronStopList, error) {
	if menuService == nil {
		return nil, errors.New("menuService is nil")
	}
	return &cronStopList{
		menuService: menuService,
	}, nil
}

func (s *cronStopList) filterProducts(products []menuModels.Product) []menuModels.Product {
	return products
}

func (s *cronStopList) updateDisabledStatusByProductIDsInDatabase(ctx context.Context, menuID string, products []menuModels.Product) error {
	var (
		productIdsWithDisabledFalse = make([]string, 0, len(products))
		productIdsWithDisabledTrue  = make([]string, 0, len(products))
	)

	for _, product := range products {
		if product.IsAvailable {
			productIdsWithDisabledFalse = append(productIdsWithDisabledFalse, product.ExtID)
		} else {
			productIdsWithDisabledTrue = append(productIdsWithDisabledTrue, product.ExtID)
		}
	}

	if len(productIdsWithDisabledTrue) != 0 {
		if err := s.menuService.UpdateProductsDisabledStatus(ctx, menuID, productIdsWithDisabledTrue, true); err != nil {
			return err
		}
	}

	if len(productIdsWithDisabledFalse) != 0 {
		if err := s.menuService.UpdateProductsDisabledStatus(ctx, menuID, productIdsWithDisabledFalse, false); err != nil {
			return err
		}
	}
	return nil
}

func (s *cronStopList) filterAttributes(attributes []menuModels.Attribute) []menuModels.Attribute {
	return attributes
}

func (s *cronStopList) updateDisabledStatusByAttributeIDsInDatabase(ctx context.Context, menuID string, attributes []menuModels.Attribute) error {
	var (
		attributeIdsWithDisabledFalse = make([]string, 0, len(attributes))
		attributeIdsWithDisabledTrue  = make([]string, 0, len(attributes))
	)

	for _, attribute := range attributes {
		if attribute.IsAvailable {
			attributeIdsWithDisabledFalse = append(attributeIdsWithDisabledFalse, attribute.ExtID)
		} else {
			attributeIdsWithDisabledTrue = append(attributeIdsWithDisabledTrue, attribute.ExtID)
		}
	}

	if len(attributeIdsWithDisabledTrue) != 0 {
		if err := s.menuService.UpdateAttributesDisabledStatus(ctx, menuID, attributeIdsWithDisabledTrue, true); err != nil {
			return err
		}
	}

	if len(attributeIdsWithDisabledFalse) != 0 {
		if err := s.menuService.UpdateAttributesDisabledStatus(ctx, menuID, attributeIdsWithDisabledFalse, false); err != nil {
			return err
		}
	}

	return nil
}

type webhookStopList struct {
}

func newWebhookStopList() (*webhookStopList, error) {
	return &webhookStopList{}, nil
}

func (s *webhookStopList) filterProducts(products []menuModels.Product) []menuModels.Product {
	result := make([]menuModels.Product, 0)
	for i := range products {
		product := products[i]
		if product.IsDisabled {
			continue
		}
		result = append(result, product)
	}
	return result
}

func (s *webhookStopList) updateDisabledStatusByAttributeIDsInDatabase(ctx context.Context, menuID string, attributes []menuModels.Attribute) error {
	return nil
}

func (s *webhookStopList) filterAttributes(attributes []menuModels.Attribute) []menuModels.Attribute {
	result := make([]menuModels.Attribute, 0)
	for i := range attributes {
		attribute := attributes[i]
		if attribute.IsDisabled {
			continue
		}
		result = append(result, attribute)
	}
	return result
}

func (s *webhookStopList) updateDisabledStatusByProductIDsInDatabase(ctx context.Context, menuID string, products []menuModels.Product) error {
	return nil
}

type validateStopList struct{}

func newValidateStopList() *validateStopList {
	return &validateStopList{}
}

func (s *validateStopList) filterProducts(products []menuModels.Product) []menuModels.Product {
	return products
}
func (s *validateStopList) updateDisabledStatusByAttributeIDsInDatabase(ctx context.Context, menuID string, attributes []menuModels.Attribute) error {
	return nil
}

func (s *validateStopList) filterAttributes(attributes []menuModels.Attribute) []menuModels.Attribute {
	return attributes
}
func (s *validateStopList) updateDisabledStatusByProductIDsInDatabase(ctx context.Context, menuID string, products []menuModels.Product) error {
	return nil
}
