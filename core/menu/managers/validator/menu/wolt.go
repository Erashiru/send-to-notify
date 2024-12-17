package menu

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/custom"
)

type wolt struct{}

func newWolt() *wolt {
	return &wolt{}
}

func (w *wolt) Validate(ctx context.Context, menu models.Menu) error {
	var errs custom.Error
	errs.Append(
		w.validateProducts(menu.Products),
	)
	return errs
}

func (w *wolt) validateProducts(products models.Products) error {
	var errs custom.Error
	// map for check section duplicate ids
	productExist := make(map[string]struct{}, len(products))

	for _, product := range products {

		if !product.IsDeleted {
			errs.Append(w.validateProduct(product))

			if _, ok := productExist[product.ExtID]; ok {
				errs.Append(fmt.Errorf("aggregator has products with same ids %s", product.ExtID))
			}
			productExist[product.ExtID] = struct{}{}
		}

	}

	return errs.ErrorOrNil()
}

func (w *wolt) validateProduct(product models.Product) error {

	var errs custom.Error

	if err := w.validatePrice(product); err != nil {
		errs.Append(err)
	}
	if err := w.validateName(product); err != nil {
		errs.Append(err)
	}

	return errs.ErrorOrNil()
}

func (w *wolt) validateName(product models.Product) error {

	if product.Name == nil || len(product.Name) == 0 || len(product.Name) != 0 && product.Name[0].Value == "" { //check loop?
		return fmt.Errorf("aggregator has product %s with empty Name", product.ExtID)
	}
	if err := w.validateLanguage(product.Name); err != nil {
		return fmt.Errorf("aggregator has product %s  with %w Name", product.ExtID, err)
	}
	return nil
}

func (w *wolt) validatePrice(product models.Product) error {
	if product.Price == nil || len(product.Price) == 0 {
		return fmt.Errorf("aggregator has product %s with empty Price", product.ExtID)
	}
	if err := w.validateCurrency(product.Price); err != nil {
		return fmt.Errorf("aggregator has product %s %w", product.ExtID, err)
	}
	return nil
}

func (w *wolt) validateCurrency(prices []models.Price) error {
	//compare all products - have same currency ?
	for _, price := range prices {
		if price.CurrencyCode == "" {
			return fmt.Errorf(" empty currency_code ")
		}
	}
	return nil
}

func (w *wolt) validateLanguage(items []models.LanguageDescription) error {
	//compare same lng?
	for _, item := range items {
		if item.LanguageCode == "" {
			return fmt.Errorf(" empty language ")
		}
	}
	return nil
}
