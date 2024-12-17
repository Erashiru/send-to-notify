package talabat

import (
	"context"
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	talabatModels "github.com/kwaaka-team/orders-core/pkg/talabat/models"
	"github.com/pkg/errors"
	"sort"
	"strconv"
	"time"
)

func setItemBranchesAvailability(itemBranchAvailabilityMap map[string][]talabatModels.BranchAvailability, brandBranches []string, branchIDs []string, itemID string, itemAvailability bool, itemPrice float64) map[string][]talabatModels.BranchAvailability {
	if _, ok := itemBranchAvailabilityMap[itemID]; !ok {
		for _, branch := range brandBranches {
			itemBranchAvailabilityMap[itemID] = append(itemBranchAvailabilityMap[itemID], talabatModels.BranchAvailability{
				BranchID: branch,
				Status:   false,
				Price:    0,
			})
		}
	}
	for i, ba := range itemBranchAvailabilityMap[itemID] {
		for _, b := range branchIDs {
			if ba.BranchID == b {
				itemBranchAvailabilityMap[itemID][i].Status = itemAvailability
				itemBranchAvailabilityMap[itemID][i].Price = itemPrice
			}
		}
	}
	return itemBranchAvailabilityMap
}

func (m mnm) constructTalabatNewMenu(ctx context.Context, menu models.Menu) (talabatModels.Catalog, error) {
	items := make(map[string]talabatModels.CatalogItem, len(menu.Products))

	attributesMap := make(map[string]models.Attribute, len(menu.Attributes))
	for _, attribute := range menu.Attributes {
		if !attribute.IsAvailable || attribute.IsDeleted {
			continue
		}
		attributesMap[attribute.ExtID] = attribute
	}

	attributeGroupsMap := make(map[string]models.AttributeGroup, len(menu.AttributesGroups))
	for _, ag := range menu.AttributesGroups {
		attributeGroupsMap[ag.ExtID] = ag
	}

	items[menu.ID] = talabatModels.CatalogItem{
		Id:       menu.ID,
		Type:     "Menu",
		MenuType: "DELIVERY",
		Title: &talabatModels.Title{
			Default: menu.Name,
		},
		Products: make(map[string]talabatModels.SubItem, len(menu.Products)),
	}

	for _, section := range menu.Sections {
		if _, ok := items[section.ExtID]; ok {
			continue
		}
		items[section.ExtID] = talabatModels.CatalogItem{
			Id:   section.ExtID,
			Type: "Category",
			Title: &talabatModels.Title{
				Default: section.Name,
			}}
	}

	for _, product := range menu.Products {
		if !isValidProduct(product) {
			continue
		}

		val, ok := items[product.Section]
		if !ok || val.Type != "Category" {
			continue
		}

		productImages := make(map[string]talabatModels.SubItem)
		for _, image := range product.ImageURLs {
			id := uuid.New().String()
			productImages[id] = talabatModels.SubItem{
				Id:   id,
				Type: "Image",
			}
			items[id] = talabatModels.CatalogItem{
				Id:   id,
				Type: "Image",
				URL:  image,
			}
		}

		productToppings := make(map[string]talabatModels.SubItem)
		productsAllAttributes := make(map[string]talabatModels.SubItem)

		for _, ag := range product.AttributesGroups {
			agVal, exist := attributeGroupsMap[ag]
			if !exist {
				continue
			}

			talabatAttributes := make(map[string]talabatModels.SubItem)

			for _, attributeID := range agVal.Attributes {
				attribute, ok := attributesMap[attributeID]
				if !ok {
					continue
				}

				talabatAttributes[attribute.ExtID] = talabatModels.SubItem{
					Id:    attributeID,
					Price: strconv.FormatFloat(attribute.Price, 'f', -1, 64),
					Type:  "Product",
				}

				items[attribute.ExtID] = talabatModels.CatalogItem{
					Id:   attributeID,
					Type: "Product",
				}

				productsAllAttributes[attribute.ExtID] = talabatModels.SubItem{
					Id:   attributeID,
					Type: "Product",
					Title: &talabatModels.Title{
						Default: attribute.Name,
					},
				}
			}

			items[agVal.ExtID] = talabatModels.CatalogItem{
				Id:   agVal.ExtID,
				Type: "Topping",
				Title: &talabatModels.Title{
					Default: agVal.Name,
				},
				Quantity: &talabatModels.Quantity{
					Max: agVal.Max,
					Min: agVal.Min,
				},
				Products: talabatAttributes,
			}

			productToppings[ag] = talabatModels.SubItem{
				Id:   ag,
				Type: "Topping",
			}
		}

		items[product.ExtID] = talabatModels.CatalogItem{
			Id:   product.ExtID,
			Type: "Product",
			Title: &talabatModels.Title{
				Default: product.Name[0].Value,
			},
			IsActive: pointerOfBool(product.IsAvailable),
			Price:    strconv.FormatFloat(product.Price[0].Value, 'f', -1, 64),
			Description: &talabatModels.Title{
				Default: product.Description[0].Value,
			},
			Images:   productImages,
			Toppings: productToppings,
		}

		sectionProducts := make(map[string]talabatModels.SubItem, len(val.Products))
		for k, v := range val.Products {
			sectionProducts[k] = v
		}
		for k, v := range productsAllAttributes {
			sectionProducts[k] = talabatModels.SubItem{
				Id:   v.Id,
				Type: v.Type,
			}
		}

		sectionProducts[product.ExtID] = talabatModels.SubItem{
			Id:   product.ExtID,
			Type: "Product",
		}

		items[product.Section] = talabatModels.CatalogItem{
			Id:   val.Id,
			Type: "Category",
			Title: &talabatModels.Title{
				Default: val.Title.Default,
			},
			Products: sectionProducts}

		menuProducts := make(map[string]talabatModels.SubItem, len(items[menu.ID].Products))
		for k, v := range items[menu.ID].Products {
			menuProducts[k] = v
		}
		menuProducts[product.ExtID] = talabatModels.SubItem{
			Id:   product.ExtID,
			Type: "Product",
		}

		items[menu.ID] = talabatModels.CatalogItem{
			Id:       menu.ID,
			Type:     "Menu",
			MenuType: "DELIVERY",
			Title: &talabatModels.Title{
				Default: menu.Name,
			},
			Products: menuProducts,
		}
	}

	return talabatModels.Catalog{
		Items: items,
	}, nil
}

func (m mnm) constructTalabatMenu(ctx context.Context, store storeModels.Store) (talabatModels.Menu, error) {
	stores, err := m.storeRepo.ListStoresByTalabatRestautantID(ctx, store.Talabat.RestaurantID)
	if err != nil {
		return talabatModels.Menu{}, err
	}

	branches := make([]string, 0, len(stores))
	for _, store := range stores {
		branches = append(branches, store.Talabat.BranchID...)
	}

	itemBranchAvailabilityMap := make(map[string][]talabatModels.BranchAvailability) // key - product/attribute id
	talabatCategoriesMap := make(map[string]talabatModels.Category)                  // key - category id
	talabatProductsMap := make(map[string]talabatModels.Item)                        // key - product id
	talabatAttributeGroupsMap := make(map[string]talabatModels.ChoiceCategory)       // key - attribute group id
	talabatAttributesMap := make(map[string]talabatModels.Choice)                    // key - attribute id

	for _, store := range stores {
		for _, storeMenu := range store.Menus {
			if storeMenu.Delivery != models.TALABAT.String() || storeMenu.IsDeleted || !storeMenu.IsActive {
				continue
			}

			menu, err := m.menuRepo.Get(ctx, selector.Menu{
				ID: storeMenu.ID,
			})
			if err != nil {
				return talabatModels.Menu{}, errors.Wrapf(err, "get menu error, menu_id: %s", storeMenu.ID)
			}

			attributeGroupMap := make(map[string]models.AttributeGroup, len(menu.AttributesGroups))
			for _, attributeGroup := range menu.AttributesGroups {
				attributeGroupMap[attributeGroup.ExtID] = attributeGroup
			}

			attributeMap := make(map[string]models.Attribute, len(menu.Attributes))
			for _, attribute := range menu.Attributes {
				if !isValidAttribute(attribute) {
					continue
				}
				attributeMap[attribute.ExtID] = attribute
			}

			sort.Sort(menu.Sections)
			for _, section := range menu.Sections {
				if _, ok := talabatCategoriesMap[section.ExtID]; ok {
					continue
				}

				talabatCategoriesMap[section.ExtID] = talabatModels.Category{
					ID:          section.ExtID,
					EnglishName: getLangValueByLangCode(section.Description, "en"),
					ArabicName:  getLangValueByLangCode(section.Description, "ar"),
					SortOrder:   section.SectionOrder,
				}
			}
			for _, product := range menu.Products {
				if !isValidProduct(product) {
					continue
				}
				_, ok := talabatCategoriesMap[product.Section]
				if !ok {
					continue
				}

				if len(stores) > 1 {
					itemBranchAvailabilityMap = setItemBranchesAvailability(itemBranchAvailabilityMap, branches, store.Talabat.BranchID, product.ExtID, product.IsAvailable, product.Price[0].Value)
				}

				for _, ag := range product.AttributesGroups {
					attG, ok := attributeGroupMap[ag]
					if !ok {
						continue
					}

					for _, a := range attG.Attributes {
						attribute, ok := attributeMap[a]
						if !ok {
							continue
						}

						if len(stores) > 1 {
							itemBranchAvailabilityMap = setItemBranchesAvailability(itemBranchAvailabilityMap, branches, store.Talabat.BranchID, attribute.ExtID, attribute.IsAvailable, attribute.Price)
						}

						tempAttributeGroups := make([]string, 0)
						if val, ok := talabatAttributesMap[attribute.ExtID]; ok {
							tempAttributeGroups = val.AttributeGroupIDs
						}
						tempAttributeGroups = append(tempAttributeGroups, ag)
						talabatAttributesMap[attribute.ExtID] = talabatModels.Choice{
							ID:                attribute.ExtID,
							EnglishName:       getLangValueByLangCode(attribute.Description, "en"),
							ArabicName:        getLangValueByLangCode(attribute.Description, "ar"),
							Price:             attribute.Price,
							IsAvailable:       attribute.IsAvailable,
							AttributeGroupIDs: tempAttributeGroups,
						}
					}

					tempProducts := make([]string, 0)
					if val, ok := talabatAttributeGroupsMap[ag]; ok {
						tempProducts = val.ProductIDs
					}
					tempProducts = append(tempProducts, product.ExtID)
					talabatAttributeGroupsMap[ag] = talabatModels.ChoiceCategory{
						ID:                attG.ExtID,
						EnglishName:       getLangValueByLangCode(attG.Description, "en"),
						ArabicName:        getLangValueByLangCode(attG.Description, "ar"),
						MinimumSelections: attG.Min,
						MaximumSelections: attG.Max,
						ProductIDs:        tempProducts,
					}
				}

				tempSections := make([]string, 0)
				if val, ok := talabatAttributeGroupsMap[product.Section]; ok {
					tempSections = val.ProductIDs
				}
				tempSections = append(tempSections, product.Section)

				talabatProductsMap[product.ExtID] = talabatModels.Item{
					ID:                 product.ExtID,
					EnglishName:        getLangValueByLangCode(product.Name, "en"),
					ArabicName:         getLangValueByLangCode(product.Name, "ar"),
					EnglishDescription: getLangValueByLangCode(product.Description, "en"),
					ArabicDescription:  getLangValueByLangCode(product.Description, "ar"),
					Price:              product.Price[0].Value,
					IsAvailable:        product.IsAvailable,
					ImageURL:           product.ImageURLs[0],
					CategoryIDs:        tempSections,
					AvailableDays:      "0,1,2,3,4,5,6",
					AvailableFrom:      "00:00",
					AvailableTo:        "00:00",
				}
			}

		}
	}

	for _, attribute := range talabatAttributesMap {
		for _, attributeGroupID := range unique(attribute.AttributeGroupIDs) {
			if attributeGroup, ok := talabatAttributeGroupsMap[attributeGroupID]; ok {
				attribute.BranchesAvailability = nil
				itemAvailability, exist := itemBranchAvailabilityMap[attribute.ID]
				if exist {
					attribute.BranchesAvailability = itemAvailability
				}

				temp := append(attributeGroup.Choices, attribute)
				talabatAttributeGroupsMap[attributeGroupID] = talabatModels.ChoiceCategory{
					ID:                attributeGroup.ID,
					EnglishName:       attributeGroup.EnglishName,
					ArabicName:        attributeGroup.ArabicName,
					MaximumSelections: attributeGroup.MaximumSelections,
					MinimumSelections: attributeGroup.MinimumSelections,
					SortOrder:         attributeGroup.SortOrder,
					ProductIDs:        attributeGroup.ProductIDs,
					Choices:           temp,
				}
			}
		}
	}

	for _, attributeGroup := range talabatAttributeGroupsMap {
		for _, productID := range unique(attributeGroup.ProductIDs) {
			if product, ok := talabatProductsMap[productID]; ok {
				temp := append(product.ChoiceCategories, attributeGroup)
				talabatProductsMap[productID] = talabatModels.Item{
					ID:                 product.ID,
					EnglishName:        product.EnglishName,
					ArabicName:         product.ArabicName,
					EnglishDescription: product.EnglishDescription,
					ArabicDescription:  product.ArabicDescription,
					Price:              product.Price,
					IsAvailable:        product.IsAvailable,
					ImageURL:           product.ImageURL,
					CategoryIDs:        product.CategoryIDs,
					ChoiceCategories:   temp,
					AvailableTo:        "00:00",
					AvailableFrom:      "00:00",
					AvailableDays:      "0,1,2,3,4,5,6",
				}
			}
		}
	}
	for _, product := range talabatProductsMap {
		for _, categoryID := range unique(product.CategoryIDs) {
			if category, ok := talabatCategoriesMap[categoryID]; ok {
				product.BranchesAvailability = nil
				itemAvailability, exist := itemBranchAvailabilityMap[product.ID]
				if exist {
					product.BranchesAvailability = itemAvailability
				}

				temp := append(category.Items, product)
				talabatCategoriesMap[categoryID] = talabatModels.Category{
					ID:          category.ID,
					EnglishName: category.EnglishName,
					ArabicName:  category.ArabicName,
					SortOrder:   category.SortOrder,
					Items:       temp,
				}
			}
		}
	}

	result := talabatModels.Menu{}

	for _, category := range talabatCategoriesMap {
		result = talabatModels.Menu{
			Categories: append(result.Categories, category),
		}
	}

	scheduledOn, err := time.Now().Add(2 * time.Minute).UTC().MarshalText()
	if err != nil {
		return talabatModels.Menu{}, errors.Wrapf(err, "talabat create menu time error")
	}

	result.ScheduledOn = string(scheduledOn)

	return result, nil
}

func getLangValueByLangCode(data []models.LanguageDescription, code string) string {
	for _, val := range data {
		if val.LanguageCode == code {
			return val.Value
		}
	}
	return ""
}

func isValidProduct(product models.Product) bool {
	if product.ExtID == "" {
		return false
	}
	if product.IsDeleted {
		return false
	}
	if len(product.Price) == 0 {
		return false
	}

	return true
}

func isValidAttribute(attribute models.Attribute) bool {
	if attribute.ExtID == "" {
		return false
	}
	if attribute.IsDeleted {
		return false
	}

	return true
}

func unique(stringSlice []string) []string {
	keys := make(map[string]struct{})
	result := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = struct{}{}
			result = append(result, entry)
		}
	}
	return result
}
func toAvailabilities(storeID string, products models.Products, attributes models.Attributes) []talabatModels.Availability {
	items := make([]talabatModels.ItemStoplist, 0, len(products)+len(attributes))
	for _, product := range products {
		items = append(items, toItem(product.ExtID, product.IsAvailable))
	}
	for _, attribute := range attributes {
		items = append(items, toItem(attribute.ExtID, attribute.IsAvailable))
	}

	return []talabatModels.Availability{
		{
			BranchId: storeID,
			Items:    items,
		},
	}
}

func toItem(itemID string, isAvailable bool) talabatModels.ItemStoplist {
	return talabatModels.ItemStoplist{
		ItemId: itemID,
		Status: toStatus(isAvailable),
	}
}

func toStatus(isAvailable bool) int {
	if isAvailable {
		return 0
	}
	return 1
}

func pointerOfBool(b bool) *bool {
	return &b
}
