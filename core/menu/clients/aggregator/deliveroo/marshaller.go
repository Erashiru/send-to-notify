package deliveroo

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
	deliverooModels "github.com/kwaaka-team/orders-core/pkg/deliveroo/clients/dto"
)

const item string = "ITEM"

func toDeliverooMenu(menu models.Menu) deliverooModels.MenuData {
	deliverooMenu := deliverooModels.MenuData{
		Categories: toCategories(menu.Sections, menu.Products),
		Modifiers:  toModifiers(menu.AttributesGroups),
		Items:      toItems(menu.Products),
		Mealtimes:  nil, //need to define
	}
	return deliverooMenu
}

func toCategories(sections []models.Section, products []models.Product) []deliverooModels.Category {

	newCategories := make([]deliverooModels.Category, 0, len(sections))
	for _, section := range sections {
		newCategory := deliverooModels.Category{
			ID: section.ExtID,
			Name: deliverooModels.Name{
				EN: section.Name,
			},
			Description: deliverooModels.Description{
				EN: section.Description[0].Value,
			},
			ItemIDs: nil, //need to define how attach item ids
		}
		newCategories = append(newCategories, newCategory)
	}
	return newCategories
}

func toModifiers(attributes []models.AttributeGroup) []deliverooModels.Modifier {
	newModifiers := make([]deliverooModels.Modifier, 0, len(attributes))

	for _, attribute := range attributes {
		newModifier := deliverooModels.Modifier{
			ID: attribute.ExtID,
			Name: deliverooModels.Name{
				EN: attribute.Name,
			},
			ItemIDs:      attribute.Attributes,
			MaxSelection: attribute.Max,
			MinSelection: attribute.Min,
		}
		newModifiers = append(newModifiers, newModifier)
	}
	return newModifiers
}

func toItems(products []models.Product) []deliverooModels.MenuItem {
	newItems := make([]deliverooModels.MenuItem, 0, len(products))

	for _, product := range products {
		newItem := deliverooModels.MenuItem{
			ID:   product.ExtID,
			Type: item,
		}
		if len(product.Name) > 0 {
			newItem.Name = deliverooModels.Name{
				EN: product.Name[0].Value,
			}
		}
		if len(product.Description) > 0 {
			newItem.Description = deliverooModels.Description{
				EN: product.Description[0].Value,
			}
		}
		if len(product.ImageURLs) > 0 {
			newItem.Image = deliverooModels.ItemImage{
				URL: product.ImageURLs[0],
			}
		}
		if len(product.Price) > 0 {
			newItem.PriceInfo = deliverooModels.PriceInfo{
				Price: int(product.Price[0].Value),
			}
		}

		newItem.ModifierIDs = append(newItem.ModifierIDs, product.AttributesGroups...)

		newItems = append(newItems, newItem)
	}
	return newItems
}
