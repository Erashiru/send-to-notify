package yaros

import (
	"github.com/google/uuid"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	models2 "github.com/kwaaka-team/orders-core/core/models"
	yarosModels "github.com/kwaaka-team/orders-core/pkg/yaros/models"
	"strconv"
)

func menuFromClient(items yarosModels.GetItemsResponse, categories yarosModels.GetCategoriesResponse) (models.Menu, error) {
	otherCollection := models.MenuCollection{
		ExtID: uuid.New().String(),
		Name:  "Other",
	}

	otherSection := models.Section{
		ExtID:      uuid.New().String(),
		Collection: otherCollection.ExtID,
		Name:       "Other",
	}

	sections, collections := toSectionCollection(categories.Categories)

	products, otherSectionUsed, err := toProduct(items.Items, sections, otherSection)
	if err != nil {
		return models.Menu{}, err
	}
	if otherSectionUsed {
		sections = append(sections, otherSection)
		collections = append(collections, otherCollection)
	}

	return models.Menu{
		Collections: collections,
		Name:        models.YAROS.String(),
		CreatedAt:   models2.TimeNow(),
		UpdatedAt:   models2.TimeNow(),
		Sections:    sections,
		Products:    products,
	}, nil
}

func toProduct(items []yarosModels.Item, sections models.Sections, otherSection models.Section) (models.Products, bool, error) {
	products := make(models.Products, 0, len(items))

	var otherSectionUsed bool

	sectionsMap := make(map[string]models.Section)
	for _, section := range sections {
		sectionsMap[section.ExtID] = section
	}

	for _, item := range items {
		if item.ImageUrl == "" {
			item.ImageUrl = "https://kwaaka-menu-files.s3.eu-west-1.amazonaws.com/images/default_image_for_product/7b05c838-2f84-423c-bc61-af2ab55d50c3.jpg"
		}
		product := models.Product{
			ExtID:            item.Id,
			PosID:            item.Id,
			Section:          item.CategoryId,
			IsIncludedInMenu: true,
			Name: []models.LanguageDescription{
				{
					Value:        item.Title,
					LanguageCode: "ru",
				},
			},
			ImageURLs: []string{item.ImageUrl},
			Description: []models.LanguageDescription{
				{
					Value:        item.Description,
					LanguageCode: "ru",
				},
			},
			MeasureUnit: item.Measure,
		}
		if item.Price == "" {
			item.Price = "0"
		}
		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return models.Products{}, false, err
		}

		product.Price = []models.Price{
			{
				Value:        price,
				CurrencyCode: "", // what currency?
			},
		}
		if _, ok := sectionsMap[item.CategoryId]; !ok {
			product.Section = otherSection.ExtID
			otherSectionUsed = true
		}
		products = append(products, product)
	}
	return products, otherSectionUsed, nil
}

func toSectionCollection(categories []yarosModels.Category) (models.Sections, models.MenuCollections) {
	collections := make([]models.MenuCollection, 0, 4)
	sections := make([]models.Section, 0, 4)

	for _, category := range categories {
		if category.ParentId == "" {
			collections = append(collections, models.MenuCollection{
				ExtID:           category.Id,
				Name:            category.Title,
				CollectionOrder: category.SortPriority,
			})
			continue
		}
		sections = append(sections, models.Section{
			Collection:   category.ParentId,
			Name:         category.Title,
			ExtID:        category.Id,
			SectionOrder: category.SortPriority,
		})
	}
	return sections, collections
}
