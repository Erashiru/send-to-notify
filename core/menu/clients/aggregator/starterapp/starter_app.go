package starterapp

import (
	"context"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kwaaka-team/orders-core/config/menu"
	constErrors "github.com/kwaaka-team/orders-core/core/menu/clients/aggregator/errors"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/starterapp"
	starterappCli "github.com/kwaaka-team/orders-core/pkg/starterapp/clients"
	"github.com/kwaaka-team/orders-core/pkg/starterapp/clients/dto"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
)

type StarterApp interface {
	BulkUpdate(ctx context.Context, restaurantID, storeID string, products models.Products, attributes models.Attributes, store storeModels.Store) (string, error)
	ModifyProduct(ctx context.Context, storeID string, product models.Product) (models.ProductModifyResponse, error)
	UploadMenu(ctx context.Context, menuId, extStoreId string, menu models.Menu, store storeModels.Store, offers []models.BkOffers, sv3 *s3.S3, userRole string) (models.ExtTransaction, error)
	VerifyMenu(ctx context.Context, transaction models.ExtTransaction) (models.Status, error)
	BulkAttribute(ctx context.Context, storeID string, attributes models.Attributes) (string, error)
	GetMenu(ctx context.Context, extStoreId string) (models.Menu, error)
	ValidateMenu(ctx context.Context, request models.MenuValidateRequest) (models.MenuUploadTransaction, error)
}

type starterAppImpl struct {
	cli starterappCli.StarterApp
	//s3Bucket menu.S3_BUCKET
	menuRepo drivers.MenuRepository
	//storeCli store.Client
}

func NewStarterAppMenuManager(ctx context.Context, apiKey string, cfg menu.Configuration, menuRepo drivers.MenuRepository) (StarterApp, error) {
	cli, err := starterapp.NewStarterAppClient(&starterappCli.Config{
		Protocol: "http",
		BaseURL:  cfg.StarterApp.BaseUrl,
		ApiKey:   apiKey,
	})
	if err != nil {
		log.Trace().Err(err).Msg("can't initialize StarterApp menu client ")
		return nil, err
	}

	return &starterAppImpl{
		cli:      cli,
		menuRepo: menuRepo,
	}, nil
}

func (s starterAppImpl) BulkUpdate(ctx context.Context, restaurantID, storeID string, products models.Products, attributes models.Attributes, store storeModels.Store) (string, error) {
	return "", constErrors.ErrNotImplemented
}

func (s starterAppImpl) ModifyProduct(ctx context.Context, storeID string, product models.Product) (models.ProductModifyResponse, error) {
	return models.ProductModifyResponse{}, constErrors.ErrNotImplemented
}

func (s starterAppImpl) UploadMenu(ctx context.Context, menuId, extStoreId string, menu models.Menu, store storeModels.Store, offers []models.BkOffers, sv3 *s3.S3, userRole string) (models.ExtTransaction, error) {
	superCollMap := make(map[string]int64)
	collMap := make(map[string]int64)
	sectionsMap := make(map[string]int64)
	attributesMap := make(map[string]int64)
	attributeGroupsMap := make(map[string]int64)
	productsMap := make(map[string]int64)

	err := s.syncSuperCollections(ctx, superCollMap, menu.SuperCollections, menuId)
	if err != nil {
		return models.ExtTransaction{
			Status:     models.NOT_PROCESSED.String(),
			MenuID:     menuId,
			ExtStoreID: extStoreId,
			Details:    []string{errors.Wrap(err, "can't sync super collections").Error()},
		}, errors.Wrap(err, "can't sync super collections")
	}

	err = s.syncCollections(ctx, superCollMap, collMap, menu.Collections, menuId)
	if err != nil {
		return models.ExtTransaction{
			Status:     models.NOT_PROCESSED.String(),
			MenuID:     menuId,
			ExtStoreID: extStoreId,
			Details:    []string{errors.Wrap(err, "can't sync collections").Error()},
		}, errors.Wrap(err, "can't sync collections")
	}

	err = s.syncSections(ctx, collMap, sectionsMap, menu.Sections, menuId)
	if err != nil {
		return models.ExtTransaction{
			Status:     models.NOT_PROCESSED.String(),
			MenuID:     menuId,
			ExtStoreID: extStoreId,
			Details:    []string{errors.Wrap(err, "can't sync sections").Error()},
		}, errors.Wrap(err, "can't sync sections")
	}

	err = s.syncAttributes(ctx, attributesMap, menu.Attributes, menuId)
	if err != nil {
		return models.ExtTransaction{
			Status:     models.NOT_PROCESSED.String(),
			MenuID:     menuId,
			ExtStoreID: extStoreId,
			Details:    []string{errors.Wrap(err, "can't sync attributes").Error()},
		}, errors.Wrap(err, "can't sync attributes")
	}

	err = s.syncAttributesGroups(ctx, attributeGroupsMap, attributesMap, menu.AttributesGroups, menuId)
	if err != nil {
		return models.ExtTransaction{
			Status:     models.NOT_PROCESSED.String(),
			MenuID:     menuId,
			ExtStoreID: extStoreId,
			Details:    []string{errors.Wrap(err, "can't sync attributes groups").Error()},
		}, errors.Wrap(err, "can't sync attributes groups")
	}

	err = s.syncProducts(ctx, productsMap, attributeGroupsMap, sectionsMap, menu.Products, menuId)
	if err != nil {
		return models.ExtTransaction{
			Status:     models.NOT_PROCESSED.String(),
			MenuID:     menuId,
			ExtStoreID: extStoreId,
			Details:    []string{errors.Wrap(err, "can't sync products").Error()},
		}, errors.Wrap(err, "can't sync products")
	}

	err = s.syncMealsOffers(ctx, extStoreId, menuId, menu.Products, productsMap)
	if err != nil {
		return models.ExtTransaction{
			Status:     models.NOT_PROCESSED.String(),
			MenuID:     menuId,
			ExtStoreID: extStoreId,
			Details:    []string{errors.Wrap(err, "can't sync meal offers").Error()},
		}, errors.Wrap(err, "can't sync meal offers")
	}

	err = s.syncModifierOffers(ctx, extStoreId, menuId, menu.Attributes, attributesMap)
	if err != nil {
		return models.ExtTransaction{
			Status:     models.NOT_PROCESSED.String(),
			MenuID:     menuId,
			ExtStoreID: extStoreId,
			Details:    []string{errors.Wrap(err, "can't sync modifier offers").Error()},
		}, errors.Wrap(err, "can't sync modifier offers")
	}

	return models.ExtTransaction{
		Status:     models.SUCCESS.String(),
		MenuID:     menuId,
		ExtStoreID: extStoreId,
		Details:    []string{},
	}, nil
}

func (s starterAppImpl) VerifyMenu(ctx context.Context, transaction models.ExtTransaction) (models.Status, error) {
	return "", nil
}

func (s starterAppImpl) BulkAttribute(ctx context.Context, storeID string, attributes models.Attributes) (string, error) {
	return "", constErrors.ErrNotImplemented
}

func (s starterAppImpl) GetMenu(ctx context.Context, extStoreId string) (models.Menu, error) {
	return models.Menu{}, constErrors.ErrNotImplemented
}

func (s starterAppImpl) ValidateMenu(ctx context.Context, request models.MenuValidateRequest) (models.MenuUploadTransaction, error) {
	return models.MenuUploadTransaction{}, nil
}

func (s starterAppImpl) syncSuperCollections(ctx context.Context, superCatMap map[string]int64, superCollections models.MenuSuperCollections, menuID string) error {
	if len(superCollections) < 2 {
		return nil
	}
	reqCreate := make([]dto.CategoryRequest, 0, 1)
	reqUpdate := make([]dto.CategoryRequest, 0, 1)
	for i := range superCollections {
		if len(superCollections[i].Collections) == 0 {
			continue
		}
		if superCollections[i].StarterAppID == "" {
			reqCreate = append(reqCreate, dto.CategoryRequest{
				PosId:             superCollections[i].ExtID,
				Name:              superCollections[i].Name,
				Images:            []string{superCollections[i].ImageUrl},
				SortIndex:         superCollections[i].SuperCollectionOrder,
				IsActive:          len(superCollections[i].Collections) > 0,
				ParentCategoryIds: []int64{},
			})
		} else {
			id, err := strconv.Atoi(superCollections[i].StarterAppID)
			if err != nil {
				log.Err(err).Msg("can't convert super collection id to int")
				continue
			}
			superCatMap[superCollections[i].ExtID] = int64(id)

			reqUpdate = append(reqUpdate, dto.CategoryRequest{
				Id:                int64(id),
				PosId:             superCollections[i].ExtID,
				Name:              superCollections[i].Name,
				Images:            []string{superCollections[i].ImageUrl},
				SortIndex:         superCollections[i].SuperCollectionOrder,
				IsActive:          len(superCollections[i].Collections) > 0,
				ParentCategoryIds: []int64{},
			})
		}
	}

	if len(reqCreate) > 0 {
		createResp, err := s.cli.CreateCategories(ctx, reqCreate)
		if err != nil {
			return err
		}
		for i := range createResp.Data {
			superCatMap[createResp.Data[i].PosId] = createResp.Data[i].Id
			err := s.menuRepo.UpdateSuperCollectionStarterAppIDByExtID(ctx, menuID, createResp.Data[i].PosId, strconv.Itoa(int(createResp.Data[i].Id)))
			if err != nil {
				log.Err(err).Msgf("can't update super collection, extID: %s, starterAppID: %d", createResp.Data[i].PosId, createResp.Data[i].Id)
			}
		}
	}

	if len(reqUpdate) > 0 {
		err := s.cli.UpdateCategories(ctx, reqUpdate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s starterAppImpl) syncCollections(ctx context.Context, superCollMap, collMap map[string]int64, collections models.MenuCollections, menuID string) error {
	if len(collections) < 2 {
		return nil
	}
	reqCreate := make([]dto.CategoryRequest, 0, 1)
	reqUpdate := make([]dto.CategoryRequest, 0, 1)
	for i := range collections {
		parents := make([]int64, 0)
		if val, ok := superCollMap[collections[i].SuperCollection]; ok {
			parents = append(parents, val)
		}
		if collections[i].StarterAppID == "" {
			reqCreate = append(reqCreate, dto.CategoryRequest{
				PosId:             collections[i].ExtID,
				Name:              collections[i].Name,
				Images:            []string{collections[i].ImageURL},
				SortIndex:         collections[i].CollectionOrder,
				IsActive:          len(collections[i].Sections) > 0,
				ParentCategoryIds: parents,
			})
		} else {
			id, err := strconv.Atoi(collections[i].StarterAppID)
			if err != nil {
				log.Err(err).Msg("can't convert super collection id to int")
				continue
			}
			collMap[collections[i].ExtID] = int64(id)

			reqUpdate = append(reqUpdate, dto.CategoryRequest{
				Id:                int64(id),
				PosId:             collections[i].ExtID,
				Name:              collections[i].Name,
				Images:            []string{collections[i].ImageURL},
				SortIndex:         collections[i].CollectionOrder,
				IsActive:          len(collections[i].Sections) > 0,
				ParentCategoryIds: parents,
			})
		}
	}

	if len(reqCreate) > 0 {
		createResp, err := s.cli.CreateCategories(ctx, reqCreate)
		if err != nil {
			return err
		}
		for i := range createResp.Data {
			collMap[createResp.Data[i].PosId] = createResp.Data[i].Id
			err := s.menuRepo.UpdateCollectionStarterAppIDByExtID(ctx, menuID, createResp.Data[i].PosId, strconv.Itoa(int(createResp.Data[i].Id)))
			if err != nil {
				log.Err(err).Msgf("can't update collection, extID: %s, starterAppID: %d", createResp.Data[i].PosId, createResp.Data[i].Id)
			}
		}
	}

	if len(reqUpdate) > 0 {
		err := s.cli.UpdateCategories(ctx, reqUpdate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s starterAppImpl) syncSections(ctx context.Context, collMap, sectionsMap map[string]int64, sections models.Sections, menuID string) error {
	if len(sections) == 0 {
		return nil
	}
	reqCreate := make([]dto.CategoryRequest, 0, 1)
	reqUpdate := make([]dto.CategoryRequest, 0, 1)
	for i := range sections {
		if sections[i].IsDeleted {
			continue
		}
		if sections[i].StarterAppID == "" {
			description := ""
			if len(sections[i].Description) > 0 {
				description = sections[i].Description[0].Value
			}
			parents := make([]int64, 0)
			if val, ok := collMap[sections[i].Collection]; ok {
				parents = append(parents, val)
			}

			reqCreate = append(reqCreate, dto.CategoryRequest{
				PosId:             sections[i].ExtID,
				Name:              sections[i].Name,
				Images:            []string{sections[i].ImageUrl},
				SortIndex:         sections[i].SectionOrder,
				IsActive:          !sections[i].IsDeleted,
				Description:       description,
				ParentCategoryIds: parents,
			})
		} else {
			id, err := strconv.Atoi(sections[i].StarterAppID)
			if err != nil {
				log.Err(err).Msg("can't convert super collection id to int")
				continue
			}
			sectionsMap[sections[i].ExtID] = int64(id)

			parents := make([]int64, 0)
			if val, ok := collMap[sections[i].Collection]; ok {
				parents = append(parents, val)
			}
			reqUpdate = append(reqUpdate, dto.CategoryRequest{
				Id:                int64(id),
				PosId:             sections[i].ExtID,
				Name:              sections[i].Name,
				Images:            []string{sections[i].ImageUrl},
				SortIndex:         sections[i].SectionOrder,
				IsActive:          !sections[i].IsDeleted,
				ParentCategoryIds: parents,
			})
		}
	}

	if len(reqCreate) > 0 {
		createResp, err := s.cli.CreateCategories(ctx, reqCreate)
		if err != nil {
			return err
		}
		for i := range createResp.Data {
			sectionsMap[createResp.Data[i].PosId] = createResp.Data[i].Id
			err := s.menuRepo.UpdateSectionStarterAppIDByExtID(ctx, menuID, createResp.Data[i].PosId, strconv.Itoa(int(createResp.Data[i].Id)))
			if err != nil {
				log.Err(err).Msgf("can't update section, extID: %s, starterAppID: %d", createResp.Data[i].PosId, createResp.Data[i].Id)
			}
		}
	}

	if len(reqUpdate) > 0 {
		err := s.cli.UpdateCategories(ctx, reqUpdate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s starterAppImpl) syncAttributes(ctx context.Context, attributesMap map[string]int64, attributes models.Attributes, menuID string) error {
	if len(attributes) == 0 {
		return nil
	}
	reqCreate := make([]dto.ModifiersRequest, 0, 1)
	reqUpdate := make([]dto.ModifiersRequest, 0, 1)
	for i := range attributes {
		if attributes[i].StarterAppID == "" {

			reqCreate = append(reqCreate, dto.ModifiersRequest{
				PosId:     attributes[i].ExtID,
				Name:      attributes[i].Name,
				Price:     attributes[i].Price,
				MaxAmount: attributes[i].Max,
				MinAmount: attributes[i].Min,
				Images:    []string{},
			})
		} else {
			id, err := strconv.Atoi(attributes[i].StarterAppID)
			if err != nil {
				log.Err(err).Msg("can't convert super collection id to int")
				continue
			}
			attributesMap[attributes[i].ExtID] = int64(id)

			reqUpdate = append(reqUpdate, dto.ModifiersRequest{
				Id:        int64(id),
				PosId:     attributes[i].ExtID,
				Name:      attributes[i].Name,
				Price:     attributes[i].Price,
				MaxAmount: attributes[i].Max,
				MinAmount: attributes[i].Min,
				Images:    []string{},
			})
		}
	}

	if len(reqCreate) > 0 {
		createResp, err := s.cli.CreateModifiers(ctx, reqCreate)
		if err != nil {
			return err
		}
		for i := range createResp.Data {
			attributesMap[createResp.Data[i].PosId] = createResp.Data[i].Id
			err := s.menuRepo.UpdateAttributeStarterAppIDByExtID(ctx, menuID, createResp.Data[i].PosId, strconv.Itoa(int(createResp.Data[i].Id)))
			if err != nil {
				log.Err(err).Msgf("can't update attribute, extID: %s, starterAppID: %d", createResp.Data[i].PosId, createResp.Data[i].Id)
			}
		}
	}

	if len(reqUpdate) > 0 {
		err := s.cli.UpdateModifiers(ctx, reqUpdate)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s starterAppImpl) syncAttributesGroups(ctx context.Context, attributeGroupsMap, attributesMap map[string]int64, attributeGroups models.AttributeGroups, menuID string) error {
	if len(attributeGroups) == 0 {
		return nil
	}

	reqCreate := make([]dto.ModifierGroupRequest, 0, 1)
	reqUpdate := make([]dto.ModifierGroupRequest, 0, 1)
	for i := range attributeGroups {
		if attributeGroups[i].StarterAppID == "" {

			reqCreate = append(reqCreate, dto.ModifierGroupRequest{
				PosId:     attributeGroups[i].ExtID,
				Name:      attributeGroups[i].Name,
				MaxAmount: attributeGroups[i].Max,
				MinAmount: attributeGroups[i].Min,
				Modifiers: s.fillModifierGroupsAttributes(attributeGroups[i].Attributes, attributesMap, attributeGroups[i].Min, attributeGroups[i].Max),
			})
		} else {
			id, err := strconv.Atoi(attributeGroups[i].StarterAppID)
			if err != nil {
				log.Err(err).Msg("can't convert super collection id to int")
				continue
			}
			attributeGroupsMap[attributeGroups[i].ExtID] = int64(id)

			reqUpdate = append(reqUpdate, dto.ModifierGroupRequest{
				Id:        int64(id),
				Name:      attributeGroups[i].Name,
				MaxAmount: attributeGroups[i].Max,
				MinAmount: attributeGroups[i].Min,
				Modifiers: s.fillModifierGroupsAttributes(attributeGroups[i].Attributes, attributesMap, attributeGroups[i].Min, attributeGroups[i].Max),
			})
		}
	}

	if len(reqCreate) > 0 {
		createResp, err := s.cli.CreateModifierGroups(ctx, reqCreate)
		if err != nil {
			return err
		}
		for i := range createResp.Data {
			attributeGroupsMap[createResp.Data[i].PosId] = createResp.Data[i].Id
			err := s.menuRepo.UpdateAttributeGroupStarterAppIDByExtID(ctx, menuID, createResp.Data[i].PosId, strconv.Itoa(int(createResp.Data[i].Id)))
			if err != nil {
				log.Err(err).Msgf("can't update attribute group, extID: %s, starterAppID: %d", createResp.Data[i].PosId, createResp.Data[i].Id)
			}
		}
	}

	if len(reqUpdate) > 0 {
		err := s.cli.UpdateModifierGroups(ctx, reqUpdate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s starterAppImpl) fillModifierGroupsAttributes(attributeIDs []string, attributesMap map[string]int64, agMin, agMax int) []dto.ModifierGroupModifier {
	res := make([]dto.ModifierGroupModifier, 0, len(attributeIDs))
	for i := range attributeIDs {
		attributeID, ok := attributesMap[attributeIDs[i]]
		if !ok {
			continue
		}
		res = append(res, dto.ModifierGroupModifier{
			Id:        attributeID,
			MaxAmount: agMax,
			MinAmount: 0,
		})
	}

	return res
}

func (s starterAppImpl) syncProducts(ctx context.Context, productsMap, attributeGroupsMap, sectionMap map[string]int64, products models.Products, menuID string) error {
	if len(products) == 0 {
		return nil
	}

	reqCreate := make([]dto.MealRequest, 0, 1)
	reqUpdate := make([]dto.MealRequest, 0, 1)
	for i := range products {
		imageUrls := products[i].ImageURLs
		if len(imageUrls) == 0 {
			imageUrls = []string{}
		}
		if products[i].StarterAppID == "" {
			name := ""
			if len(products[i].Name) > 0 {
				name = products[i].Name[0].Value
			}
			description := ""
			if len(products[i].Description) > 0 {
				description = products[i].Description[0].Value
			}
			reqCreate = append(reqCreate, dto.MealRequest{
				PosId:                products[i].ExtID,
				Name:                 name,
				Description:          description,
				Images:               imageUrls,
				IsActive:             !products[i].IsDeleted,
				ModifierGroups:       s.fillProductsModifierGroups(products[i].AttributesGroups, attributeGroupsMap),
				CategoryIds:          s.fillProductsSection(products[i].Section, sectionMap),
				DeliveryRestrictions: []string{},
			})
		} else {
			id, err := strconv.Atoi(products[i].StarterAppID)
			if err != nil {
				log.Err(err).Msg("can't convert super collection id to int")
				continue
			}
			productsMap[products[i].ExtID] = int64(id)
			name := ""
			if len(products[i].Name) > 0 {
				name = products[i].Name[0].Value
			}
			description := ""
			if len(products[i].Description) > 0 {
				description = products[i].Description[0].Value
			}

			reqUpdate = append(reqUpdate, dto.MealRequest{
				Id:                   int64(id),
				PosId:                products[i].ExtID,
				Name:                 name,
				Description:          description,
				Images:               imageUrls,
				IsActive:             !products[i].IsDeleted,
				ModifierGroups:       s.fillProductsModifierGroups(products[i].AttributesGroups, attributeGroupsMap),
				CategoryIds:          s.fillProductsSection(products[i].Section, sectionMap),
				DeliveryRestrictions: []string{},
			})
		}
	}

	if len(reqCreate) > 0 {
		createResp, err := s.cli.CreateMeals(ctx, reqCreate)
		if err != nil {
			return err
		}
		for i := range createResp.Data {
			productsMap[createResp.Data[i].PosId] = createResp.Data[i].Id
			err := s.menuRepo.UpdateProductStarterAppIDByExtID(ctx, menuID, createResp.Data[i].PosId, strconv.Itoa(int(createResp.Data[i].Id)))
			if err != nil {
				log.Err(err).Msgf("can't update product, extID: %s, starterAppID: %d", createResp.Data[i].PosId, createResp.Data[i].Id)
			}
		}
	}

	if len(reqUpdate) > 0 {
		err := s.cli.UpdateMeals(ctx, reqUpdate)
		//_, err := s.cli.CreateMeals(ctx, reqUpdate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s starterAppImpl) fillProductsModifierGroups(agIDs []string, attributeGroupsMap map[string]int64) []int {
	res := make([]int, 0, len(agIDs))
	for i := range agIDs {
		agID, ok := attributeGroupsMap[agIDs[i]]
		if !ok {
			continue
		}
		res = append(res, int(agID))
	}

	return res
}

func (s starterAppImpl) fillProductsSection(sectionID string, sectionMap map[string]int64) []int {
	res := make([]int, 0, 1)
	if sectionID != "" {
		sectionID, ok := sectionMap[sectionID]
		if ok {
			res = append(res, int(sectionID))
		}
	}

	return res
}

func (s starterAppImpl) syncMealsOffers(ctx context.Context, shopID, menuID string, products models.Products, productsMap map[string]int64) error {
	shopIDint, err := strconv.Atoi(shopID)
	if err != nil {
		return err
	}
	reqCreate := make([]dto.MealOfferRequest, 0, 1)
	reqUpdate := make([]dto.MealOfferRequest, 0, 1)

	for i := range products {
		if products[i].IsDeleted {
			continue
		}
		mealId, ok := productsMap[products[i].ExtID]
		if !ok {
			continue
		}

		quatity := 0
		if products[i].IsAvailable {
			quatity = 1000
		}

		if len(products[i].Price) == 0 {
			continue
		}

		if products[i].StarterAppOfferID != "" {
			id, err := strconv.Atoi(products[i].StarterAppOfferID)
			if err != nil {
				log.Err(err).Msg("can't convert starterapp product offer id to int")
				continue
			}
			reqUpdate = append(reqUpdate, dto.MealOfferRequest{
				PosId:    products[i].ExtID,
				Quantity: quatity,
				Price:    products[i].Price[0].Value,
				InMenu:   products[i].IsAvailable,
				MealId:   mealId,
				ID:       int64(id),
			})
		} else {
			reqCreate = append(reqCreate, dto.MealOfferRequest{
				PosId:    products[i].ExtID,
				Quantity: quatity,
				Price:    products[i].Price[0].Value,
				InMenu:   products[i].IsAvailable,
				MealId:   mealId,
			})
		}
	}
	if len(reqCreate) > 0 {
		respCreate, err := s.cli.CreateMealOffers(ctx, reqCreate, shopIDint)
		if err != nil {
			return err
		}

		for i := range respCreate.Data {
			err := s.menuRepo.UpdateProductStarterAppOfferIDByExtID(ctx, menuID, respCreate.Data[i].PosId, strconv.Itoa(int(respCreate.Data[i].Id)))
			if err != nil {
				log.Err(err).Msgf("can't update products starterapp offer id, extID: %s, starterAppOfferID: %d", respCreate.Data[i].PosId, respCreate.Data[i].Id)
			}
		}
	}

	if len(reqUpdate) > 0 {
		err := s.cli.UpdateMealOffers(ctx, reqUpdate, shopIDint)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s starterAppImpl) syncModifierOffers(ctx context.Context, shopID, menuID string, attributes models.Attributes, attributesMap map[string]int64) error {
	shopIDint, err := strconv.Atoi(shopID)
	if err != nil {
		return err
	}
	reqCreate := make([]dto.ModifierOfferRequest, 0, 1)
	reqUpdate := make([]dto.ModifierOfferRequest, 0, 1)

	for i := range attributes {
		if attributes[i].IsDeleted {
			continue
		}
		modifierId, ok := attributesMap[attributes[i].ExtID]
		if !ok {
			continue
		}

		quatity := 0
		if attributes[i].IsAvailable {
			quatity = 1000
		}

		if attributes[i].StarterAppOfferID != "" {
			id, err := strconv.Atoi(attributes[i].StarterAppOfferID)
			if err != nil {
				log.Err(err).Msg("can't convert starterapp product offer id to int")
				continue
			}
			reqUpdate = append(reqUpdate, dto.ModifierOfferRequest{
				PosId:      attributes[i].ExtID,
				Quantity:   quatity,
				Price:      attributes[i].Price,
				ModifierId: modifierId,
				ID:         int64(id),
				ShopId:     shopIDint,
			})
		} else {
			reqCreate = append(reqCreate, dto.ModifierOfferRequest{
				PosId:      attributes[i].ExtID,
				Quantity:   quatity,
				Price:      attributes[i].Price,
				ModifierId: modifierId,
				ShopId:     shopIDint,
			})
		}
	}

	if len(reqCreate) > 0 {
		respCreate, err := s.cli.CreateModifierOffers(ctx, reqCreate)
		if err != nil {
			return err
		}

		for i := range respCreate.Data {
			err := s.menuRepo.UpdateAttributeStarterAppOfferIDByExtID(ctx, menuID, respCreate.Data[i].PosId, strconv.Itoa(int(respCreate.Data[i].Id)))
			if err != nil {
				log.Err(err).Msgf("can't update attributes starterapp offer id, extID: %s, starterAppOfferID: %d", respCreate.Data[i].PosId, respCreate.Data[i].Id)
			}
		}
	}

	if len(reqUpdate) > 0 {
		err = s.cli.UpdateModifierOffers(ctx, reqUpdate)
		if err != nil {
			return err
		}
	}

	return nil
}
