package menu

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/custom"
)

type glovo struct{}

func newGlovo() *glovo {
	return &glovo{}
}

func (g *glovo) Validate(ctx context.Context, menu models.Menu) error {
	var errs custom.Error

	errs.Append(
		g.validateProducts(menu.Products),
		g.validateSections(menu.Sections),
		g.validateAttributes(menu.Attributes, menu.AttributesGroups),
	)
	return errs
}

func (g *glovo) validateSections(sections models.Sections) error {
	var errs custom.Error

	// map for check section duplicate ids
	sectionExist := make(map[string]struct{}, len(sections))

	for _, section := range sections {
		if section.Name == "" {
			errs.Append(fmt.Errorf("aggregator has sections with empty nameid %s", section.ExtID))
		}

		if _, ok := sectionExist[section.ExtID]; ok {
			errs.Append(fmt.Errorf("aggregator has sections with same id %s", section.ExtID))
		}

		sectionExist[section.ExtID] = struct{}{}
	}

	return errs.ErrorOrNil()
}

func (g *glovo) validateProducts(products models.Products) error {
	var errs custom.Error
	// map for check section duplicate ids
	productExist := make(map[string]struct{}, len(products))

	for _, product := range products {

		errs.Append(g.validateProduct(product))

		if _, ok := productExist[product.ExtID]; ok {
			errs.Append(fmt.Errorf("aggregator has products with same ids %s", product.ExtID))
		}

		productExist[product.ExtID] = struct{}{}
	}

	return errs.ErrorOrNil()
}

func (g *glovo) validateProduct(product models.Product) error {
	var errs custom.Error

	if len(product.Name) == 0 || len(product.Name) != 0 && product.Name[0].Value == "" {
		errs.Append(fmt.Errorf("aggregator has product %s with empty name", product.ExtID))
	}
	return errs.ErrorOrNil()
}

func (g *glovo) validateAttributes(attributes models.Attributes, attributeGroups models.AttributeGroups) error {
	m := make(map[string]int)
	var errs custom.Error
	for _, attr := range attributes {
		m[attr.ExtID]++
		if m[attr.ExtID] > 1 {
			errs.Append(fmt.Errorf("attributes with id %s, name %s dublicated", attr.ExtID, attr.Name))
		}
	}

	return errs.ErrorOrNil()
}
