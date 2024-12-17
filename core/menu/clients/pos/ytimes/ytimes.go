package ytimes

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/clients/pos/base"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/ytimes"
	"github.com/kwaaka-team/orders-core/pkg/ytimes/clients"
	ytimesModels "github.com/kwaaka-team/orders-core/pkg/ytimes/clients/models"
	"github.com/pkg/errors"
)

type manager struct {
	ytimesCli clients.Client
	store     coreStoreModels.Store
}

func NewYtimesManager(baseURL string, store coreStoreModels.Store) (base.Manager, error) {
	if baseURL == "" {
		return nil, errors.New("ytimes base url is empty")
	}
	if store.YTimes.AuthToken == "" {
		return nil, errors.New("ytimes store token is empty")
	}

	ytimesClient := ytimes.New(clients.Config{
		BaseUrl: baseURL,
		Token:   store.YTimes.AuthToken,
	})

	return &manager{
		ytimesCli: ytimesClient,
		store:     store,
	}, nil
}

func (m manager) GetMenu(ctx context.Context, store coreStoreModels.Store) (models.Menu, error) {
	responseMenu, err := m.ytimesCli.GetMenu(ctx, store.YTimes.PointId)
	if err != nil {
		return models.Menu{}, err
	}

	responseSupplementList, err := m.ytimesCli.GetSupplementList(ctx, store.YTimes.PointId)
	if err != nil {
		return models.Menu{}, err
	}

	systemMenu := m.toSystemMenu(responseMenu, responseSupplementList)

	return systemMenu, nil
}

func (ks manager) toSystemMenu(posMenu ytimesModels.Menu, posModifierGroupList ytimesModels.SupplementList) models.Menu {
	systemMenu := models.Menu{
		Sections: ks.getSections(posMenu.Rows),
	}

	systemProducts := ks.toSystemProducts(posMenu.Rows)

	systemAttributeGroups, systemAttributes := ks.toSystemAttributesAndAttributeGroups(posModifierGroupList.Rows)

	systemMenu.Products = systemProducts
	systemMenu.Attributes = systemAttributes
	systemMenu.AttributesGroups = systemAttributeGroups

	return systemMenu
}

func (ks manager) getSections(rows []ytimesModels.MenuRow) []models.Section {
	sections := make([]models.Section, 0, len(rows))

	for _, row := range rows {
		sections = append(sections, models.Section{
			ExtID:        row.Guid,
			Name:         row.Name,
			SectionOrder: row.Priority,
			ImageUrl:     row.ImageLink,
		})
	}

	return sections
}

func (ks manager) toSystemProducts(rows []ytimesModels.MenuRow) []models.Product {
	var (
		systemProducts = make([]models.Product, 0, 4)
	)

	for _, row := range rows {
		for _, categoryList := range row.CategoryList {
			for _, item := range categoryList.ItemList {
				systemProducts = append(systemProducts, ks.itemToSystemProducts(item, categoryList.Guid)...)
			}

			for _, good := range categoryList.GoodsList {
				systemProduct := ks.goodsItemToSystemProduct(good, categoryList.Guid)
				systemProducts = append(systemProducts, systemProduct)
			}
		}

		for _, item := range row.ItemList {
			systemProducts = append(systemProducts, ks.itemToSystemProducts(item, row.Guid)...)
		}

		for _, goodsItem := range row.GoodsList {
			systemProduct := ks.goodsItemToSystemProduct(goodsItem, row.Guid)
			systemProducts = append(systemProducts, systemProduct)
		}
	}

	return systemProducts
}

func (ks manager) itemToSystemProducts(item ytimesModels.ItemList, sectionId string) []models.Product {
	systemProducts := make([]models.Product, 0, len(item.TypeList))

	productName := item.Name

	systemProduct := models.Product{
		ExtID:     item.Guid,
		ProductID: item.Guid,
		ImageURLs: []string{item.ImageLink},
		Description: []models.LanguageDescription{
			{
				Value: item.Description,
			},
		},
		Section:          sectionId,
		IsAvailable:      true,
		IsIncludedInMenu: true,
		IsSync:           true,
	}

	for key := range item.SupplementCategoryToFreeCount {
		systemProduct.AttributesGroups = append(systemProduct.AttributesGroups, key)
	}

	if len(item.TypeList) > 1 {
		for i := 0; i < len(item.TypeList); i++ {
			systemProduct.ExtID = item.Guid + item.TypeList[i].Guid
			systemProduct.SizeID = item.TypeList[i].Guid
			systemProduct.Name = []models.LanguageDescription{
				{
					Value: item.TypeList[i].Name + " " + productName,
				},
			}
			systemProduct.Price = []models.Price{
				{
					Value: item.TypeList[i].Price,
				},
			}
			systemProducts = append(systemProducts, systemProduct)
		}
	} else if len(item.TypeList) == 1 {
		systemProduct.Name = []models.LanguageDescription{
			{
				Value: item.TypeList[0].Name + " " + productName,
			},
		}
		systemProduct.SizeID = item.TypeList[0].Guid
		systemProduct.Price = []models.Price{
			{
				Value: item.TypeList[0].Price,
			},
		}
		systemProducts = append(systemProducts, systemProduct)
	}

	// TODO: default supplements docs is empty

	return systemProducts
}

func (ks manager) goodsItemToSystemProduct(goodsItem ytimesModels.GoodsList, sectionId string) models.Product {
	return models.Product{
		ExtID:     goodsItem.Guid,
		ProductID: goodsItem.Guid,
		Name: []models.LanguageDescription{
			{
				Value: goodsItem.Name,
			},
		},
		ImageURLs: []string{goodsItem.ImageLink},
		Description: []models.LanguageDescription{
			{
				Value: goodsItem.Description,
			},
		},
		Price: []models.Price{
			{
				Value: float64(goodsItem.Price),
			},
		},
		Section:          sectionId,
		IsAvailable:      true,
		IsIncludedInMenu: true,
		IsSync:           true,
	}
}

func (ks manager) toSystemAttributesAndAttributeGroups(rows []ytimesModels.SupplementRow) ([]models.AttributeGroup, []models.Attribute) {
	var (
		attributeGroups    = make([]models.AttributeGroup, 0, len(rows))
		attributes         = make([]models.Attribute, 0, 4)
		existingAttributes = make(map[string]bool)
	)

	for _, row := range rows {
		attributeIds := make([]string, 0, len(row.ItemList))

		for _, attribute := range row.ItemList {
			if !existingAttributes[attribute.Guid] {
				attributes = append(attributes, models.Attribute{
					ExtID: attribute.Guid,
					Name:  attribute.Name,
					Price: float64(attribute.DefaultPrice),
				})

				existingAttributes[attribute.Guid] = true
			}

			attributeIds = append(attributeIds, attribute.Guid)
		}

		attributeGroups = append(attributeGroups, models.AttributeGroup{
			ExtID:          row.Guid,
			Name:           row.Name,
			Min:            0,
			Max:            row.MaxSelectedCount,
			MultiSelection: row.AllowSeveralItem,
			Attributes:     attributeIds,
		})

	}

	return attributeGroups, attributes
}

func (m manager) GetAggMenu(ctx context.Context, store coreStoreModels.Store) ([]models.Menu, error) {
	return nil, errors.New("method not implemented")
}
