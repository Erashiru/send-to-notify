package menu

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/disintegration/imaging"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	notifyQueue "github.com/kwaaka-team/orders-core/pkg/que"
	woltModels "github.com/kwaaka-team/orders-core/pkg/wolt/clients/dto"
	"github.com/kwaaka-team/orders-core/service/aws_s3"
	storeServicePkg "github.com/kwaaka-team/orders-core/service/store"
	"github.com/pkg/errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Service struct {
	repo         Repository
	storeService storeServicePkg.Service
	s3Service    aws_s3.Service
	sqsClient    notifyQueue.SQSInterface
}

func NewMenuService(repo Repository, storeService storeServicePkg.Service, s3Service aws_s3.Service, sqsClient notifyQueue.SQSInterface) (*Service, error) {
	if repo == nil {
		return nil, errors.New("menu repository is nil")
	}
	return &Service{
		repo:         repo,
		s3Service:    s3Service,
		sqsClient:    sqsClient,
		storeService: storeService,
	}, nil
}

const (
	minWidth  = 1000
	minHeight = 563
	maxWidth  = 10000
	maxHeight = 10000
)

func (s *Service) GenerateNewAggregatorMenuFromPosMenu(ctx context.Context, store storeModels.Store, delivery string) (string, error) {
	posMenu, err := s.repo.FindById(ctx, store.MenuID)
	if err != nil {
		return "", err
	}

	aggregatorMenu := posMenu

	aggregatorMenu.ID = ""
	aggregatorMenu.Delivery = delivery

	id, err := s.repo.Insert(ctx, *aggregatorMenu)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Service) SetMarkupToAggregatorMenu(ctx context.Context, aggregatorMenuId, posMenuId, currencyCode string, markupPercent int) error {
	posMenu, err := s.repo.FindById(ctx, posMenuId)
	if err != nil {
		return errors.Wrap(err, "pos menu not found")
	}

	aggregatorMenu, err := s.repo.FindById(ctx, aggregatorMenuId)
	if err != nil {
		return errors.Wrap(err, "aggregator menu not found")
	}

	posMenuProductsMap := make(map[string]float64, len(posMenu.Products))

	for _, product := range posMenu.Products {
		if len(product.Price) == 0 {
			continue
		}

		posMenuProductsMap[product.ExtID] = product.Price[0].Value
	}

	for i, product := range aggregatorMenu.Products {
		id := product.ExtID
		if product.PosID != "" {
			id = product.PosID
		}

		posPrice, ok := posMenuProductsMap[id]
		if !ok {
			continue
		}

		aggregatorMenu.Products[i].Price = []coreMenuModels.Price{
			{
				Value:        float64(int(posPrice) * (100 + markupPercent) / 100),
				CurrencyCode: currencyCode,
			},
		}
	}

	if err = s.repo.UpdateMenuEntities(ctx, aggregatorMenuId, *aggregatorMenu); err != nil {
		return err
	}

	return nil
}

func (s *Service) UploadImagesInWoltFormat(ctx context.Context, menuId string) error {
	systemMenu, err := s.repo.FindById(ctx, menuId)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, pr := range systemMenu.Products {
		if len(pr.ImageURLs) == 0 {
			continue
		}

		wg.Add(1)

		go func(product coreMenuModels.Product) {
			defer wg.Done()

			resp, err := http.Get(product.ImageURLs[0])
			if err != nil {
				return
			}
			defer resp.Body.Close()

			// Read the image content
			imageData, err := io.ReadAll(resp.Body)
			if err != nil {
				return
			}

			// resizing
			img, err := imaging.Decode(bytes.NewReader(imageData))
			if err != nil {
				return
			}

			width := img.Bounds().Dx()
			height := img.Bounds().Dy()

			if width < minWidth || height < minHeight || width > maxWidth || height > maxHeight {
				newWidth, newHeight := calculateNewDimensions(width, height, minWidth, minHeight, maxWidth, maxHeight)
				img = imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)
			}

			var buf bytes.Buffer
			if err := imaging.Encode(&buf, img, imaging.PNG); err != nil {
				return
			}

			// File name in S3
			link := strings.TrimSpace(fmt.Sprintf("s3://%v/small/wolt/%s.png", os.Getenv(coreMenuModels.S3_BUCKET), product.ExtID))

			if err = s.s3Service.PutObjectFromBytes(link, buf.Bytes(), os.Getenv(coreMenuModels.S3_BUCKET), "image/png"); err != nil {
				return
			}
		}(pr)

		time.Sleep(3 * time.Millisecond)
	}

	wg.Wait()

	return nil
}

func calculateNewDimensions(width, height, minWidth, minHeight, maxWidth, maxHeight int) (newWidth, newHeight int) {
	aspectRatio := float64(width) / float64(height)

	if width < minWidth || height < minHeight {
		if width < minWidth {
			newWidth = minWidth
			newHeight = int(float64(newWidth) / aspectRatio)
		}
		if newHeight < minHeight {
			newHeight = minHeight
			newWidth = int(float64(newHeight) * aspectRatio)
		}
	} else {
		newWidth = width
		newHeight = height
	}

	if newWidth > maxWidth || newHeight > maxHeight {
		if newWidth > maxWidth {
			newWidth = maxWidth
			newHeight = int(float64(newWidth) / aspectRatio)
		}
		if newHeight > maxHeight {
			newHeight = maxHeight
			newWidth = int(float64(newHeight) * aspectRatio)
		}
	}

	return newWidth, newHeight
}

func (s *Service) ConvertMenuToWoltCsv(ctx context.Context, menuId string, needUploadImage bool, file multipart.File) (string, error) {
	systemMenu, err := s.repo.FindById(ctx, menuId)
	if err != nil {
		return "", err
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	woltCsvProducts := make(woltModels.WoltCsvProducts, 0, len(systemMenu.Products))

	for _, product := range systemMenu.Products {
		if len(product.Name) == 0 || len(product.Price) == 0 {
			continue
		}

		var description string

		if len(product.Description) > 0 {
			description = product.Description[0].Value
		}

		enabled := "NO"
		if product.IsAvailable {
			enabled = "YES"
		}

		merchantSku, _ := strconv.Atoi(product.ExtID)

		woltCsvProduct := woltModels.WoltCsvProduct{
			//Gtin:                   int64(gtin),
			MerchantSku:            merchantSku,
			Name:                   product.Name[0].Value,
			Price:                  product.Price[0].Value,
			DiscountedPrice:        product.DiscountPrice.Value,
			AlcoholPercentage:      product.AlcoholPercentage,
			Description:            description,
			PosId:                  product.ExtID,
			Enabled:                enabled,
			UseInventoryManagement: "NO",
			UseLimitedQuantity:     "YES",
			MinQuantityPerPurchase: 1,
			MaxQuantityPerPurchase: 5,
			IsAgeRestrictedItem:    "YES",
			AgeLimit:               6,
			CaffeineContentUnits:   "mg_per_100_ml",
			CategoryId:             s.convertIdToCategoryId(product, *systemMenu, records),
		}

		woltCsvProducts = append(woltCsvProducts, woltCsvProduct)
	}

	bytes, err := woltCsvProducts.ToCSVBytes()
	if err != nil {
		return "", err
	}

	currentDate := fmt.Sprintf("%v", time.Now().Format("2006-01-02"))

	currentTime := fmt.Sprintf("%v", time.Now().Format("15:04:05"))

	link := strings.TrimSpace(fmt.Sprintf("s3://%v/small_wolt_csv/%s/%s/%s", os.Getenv(coreMenuModels.S3_BUCKET), currentDate, currentTime, menuId))

	if err = s.s3Service.PutObjectFromBytes(link, bytes, os.Getenv(coreMenuModels.S3_BUCKET), "text/csv"); err != nil {
		return "", err
	}

	if needUploadImage {
		if err = s.sqsClient.SendMessageToUploadWoltImagesToS3("small_wolt_images", menuId); err != nil {
			return "", err
		}
	}

	return "", nil
}

// TODO: reimplement in future
func (s *Service) convertIdToCategoryId(product coreMenuModels.Product, menu coreMenuModels.Menu, records [][]string) string {
	categoryMap := make(map[string]string)
	for i, r := range records {
		if i == 0 {
			continue
		}
		if len(r) < 2 {
			continue
		}
		categoryID := r[0]
		categoryName := r[1]
		categoryMap[categoryName] = categoryID
	}

	var savedSec coreMenuModels.Section
	for _, section := range menu.Sections {
		if section.ExtID == product.Section {
			savedSec = section
		}
	}

	var savedColl coreMenuModels.MenuCollection
	for _, collection := range menu.Collections {
		if collection.ExtID == savedSec.Collection {
			savedColl = collection
		}
	}

	var savedSuperColl coreMenuModels.MenuSuperCollection
	for _, coll := range menu.SuperCollections {
		if coll.ExtID == savedColl.SuperCollection {
			savedSuperColl = coll
		}
	}

	catId, ok := categoryMap[savedSuperColl.Name]
	if !ok {
		return ""
	}

	return catId
}

func (s *Service) createMenuInDb(ctx context.Context, posMenu coreMenuModels.Menu) (string, error) {
	return s.repo.Insert(ctx, posMenu)
}

func (s *Service) updateStoreMenuId(ctx context.Context, storeId, menuId string) error {
	return s.storeService.UpdateMenuId(ctx, storeId, menuId)
}

func (s *Service) ifPosMenuIdExist(store storeModels.Store) bool {
	return store.MenuID != ""
}

func (s *Service) setSystemProductInfo(ctx context.Context, systemMenu coreMenuModels.Menu, posMenu coreMenuModels.Menu) (coreMenuModels.Menu, error) {
	systemProducts := systemMenu.Products

	existProductMap := make(map[string]coreMenuModels.Product)

	for _, systemProduct := range systemProducts {
		existProductMap[systemProduct.ExtID] = systemProduct
	}

	for index, posProduct := range posMenu.Products {
		systemProduct, ok := existProductMap[posProduct.ExtID]
		if !ok {
			continue
		}

		// TODO: в некоторых POS ext_id не равен product_id они генерятся рандомно, нужно написать скрипт который ext_id сделает равным product_id (где отсутствует size id) - (paloma, rkeeper)

		// add default attributes
		if len(systemProduct.MenuDefaultAttributes) != 0 {
			unique := make(map[string]bool)

			// если дефолт атрибут уже добавлен в UpsertMenu то не нужно добавлять его
			for _, defaultAttribute := range posMenu.Products[index].MenuDefaultAttributes {
				if defaultAttribute.ByAdmin {
					unique[defaultAttribute.ExtID] = true
				}
			}

			for _, defaultAttribute := range systemProduct.MenuDefaultAttributes {
				if defaultAttribute.ByAdmin && !unique[defaultAttribute.ExtID] {
					posMenu.Products[index].MenuDefaultAttributes = append(posMenu.Products[index].MenuDefaultAttributes, defaultAttribute)
				}
			}
		}
	}

	return posMenu, nil
}

func (s *Service) UpdateMenuEntities(ctx context.Context, menuID string, menu coreMenuModels.Menu) error {
	return s.repo.UpdateMenuEntities(ctx, menuID, menu)
}

func (s *Service) UpsertMenu(ctx context.Context, store storeModels.Store, systemMenu coreMenuModels.Menu, externalPosMenu coreMenuModels.Menu) (string, error) {
	if s.ifPosMenuIdExist(store) {
		posMenu, err := s.setSystemProductInfo(ctx, systemMenu, externalPosMenu)
		if err != nil {
			return "", err
		}

		if err = s.repo.UpdateMenuEntities(ctx, systemMenu.ID, posMenu); err != nil {
			return "", err
		}

		return store.MenuID, nil
	}

	menuId, err := s.createMenuInDb(ctx, externalPosMenu)
	if err != nil {
		return "", err
	}

	if err = s.storeService.UpdateMenuId(ctx, store.ID, menuId); err != nil {
		return "", err
	}

	return menuId, nil
}

func (s *Service) IsMenuExists(store storeModels.Store, deliveryService string) bool {
	_, isExist := s.getActiveMenuID(store, deliveryService)
	return isExist
}

func (s *Service) FindById(ctx context.Context, menuID string) (*coreMenuModels.Menu, error) {
	return s.repo.FindById(ctx, menuID)
}

func (s *Service) GetAggregatorMenuIfExists(ctx context.Context, store storeModels.Store, deliveryService string) (coreMenuModels.Menu, error) {
	menuID, isExist := s.getActiveMenuID(store, deliveryService)
	if !isExist {
		return coreMenuModels.Menu{}, nil
	}

	m, err := s.repo.FindById(ctx, menuID)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	return *m, nil
}

func (s *Service) getActiveMenuID(store storeModels.Store, deliveryService string) (string, bool) {
	for _, m := range store.Menus {
		if m.IsActive && m.Delivery == deliveryService {
			return m.ID, true
		}
	}

	return "", false
}

func (s *Service) UpdateProductsDeletedStatus(ctx context.Context, menuId string, productIds []string, isDeleted bool, reason string) error {
	return s.repo.BulkUpdateProductsIsDeleted(ctx, menuId, productIds, isDeleted, reason)
}

func (s *Service) UpdateAttributesDeletedStatus(ctx context.Context, menuId string, attributeIds []string, isDeleted bool, reason string) error {
	return s.repo.BulkUpdateAttributesIsDeleted(ctx, menuId, attributeIds, isDeleted, reason)
}

func (s *Service) UpdateProductsAvailabilityStatus(ctx context.Context, menuId string, productIds []string, availability bool) error {
	return s.repo.BulkUpdateProductsAvailability(ctx, menuId, productIds, availability)
}

func (s *Service) UpdateAttributesAvailabilityStatus(ctx context.Context, menuId string, attributeIds []string, availability bool) error {
	return s.repo.BulkUpdateAttributesAvailability(ctx, menuId, attributeIds, availability)
}

func (s *Service) UpdateProductsDisabledStatus(ctx context.Context, menuId string, productIds []string, isDisabled bool) error {
	return s.repo.BulkUpdateProductsDisabledStatus(ctx, menuId, productIds, isDisabled)
}

func (s *Service) UpdateAttributesDisabledStatus(ctx context.Context, menuId string, attributeIds []string, isDisabled bool) error {
	return s.repo.BulkUpdateAttributesDisabledStatus(ctx, menuId, attributeIds, isDisabled)
}

func (s *Service) UpdateProductStopListStatus(ctx context.Context, menuId string, productID string, isAvailable *bool, isDisabled *bool) error {
	upd := coreMenuModels.ProductUpdateRequest{
		IsAvailable: isAvailable,
		IsDisabled:  isDisabled,
	}
	return s.repo.UpdateProductStopListStatus(ctx, menuId, productID, upd)
}

func (s *Service) UpdateAttributeStopListStatus(ctx context.Context, menuId string, attributeID string, isAvailable *bool, isDisabled *bool) error {
	return s.repo.UpdateAttributeStopListStatus(ctx, menuId, attributeID, isAvailable, isDisabled)
}

func (s *Service) UpdateStopList(ctx context.Context, menuId string, stopListProducts []string) error {
	return s.repo.UpdateStopList(ctx, menuId, stopListProducts)
}

func (s *Service) ListProductsByMenuId(ctx context.Context, menuId string) ([]coreMenuModels.Product, int64, error) {
	return s.repo.ListProductsByMenuId(ctx, menuId)
}

func (s *Service) GetCombosByMenuId(ctx context.Context, menuId string) ([]coreMenuModels.Combo, int64, error) {
	return s.repo.GetCombosByMenuId(ctx, menuId)
}

func (s *Service) SearchProduct(ctx context.Context, menuId, productName string) ([]coreMenuModels.Product, error) {
	return s.repo.SearchProduct(ctx, menuId, productName)
}

func (s *Service) getNewAggregatorMenu(posMenu coreMenuModels.Menu, aggregatorMenu coreMenuModels.Menu) coreMenuModels.Menu {
	newAggregatorMenu := coreMenuModels.Menu{
		Attributes:  posMenu.Attributes,
		Sections:    aggregatorMenu.Sections,
		Collections: aggregatorMenu.Collections,
	}

	return newAggregatorMenu
}

func (s *Service) getUsedAttributeGroupsMapInProducts(aggregatorMenu coreMenuModels.Menu, posProductsMap map[string]coreMenuModels.Product) map[string]struct{} {
	unique := make(map[string]struct{})

	for _, product := range aggregatorMenu.Products {
		id := product.ExtID
		if product.PosID != "" {
			id = product.PosID
		}

		if posProduct, ok := posProductsMap[id]; !ok {
			continue
		} else {
			for _, groupId := range posProduct.AttributesGroups {
				unique[groupId] = struct{}{}
			}

		}
	}

	return unique
}

func (s *Service) matchAggregatorAndPosProducts(aggregatorMenu coreMenuModels.Menu, posProductsMap map[string]coreMenuModels.Product, nonDeliveryAttributeGroupIdsMap map[string]struct{}) []coreMenuModels.Product {
	products := make([]coreMenuModels.Product, 0)

	for _, product := range aggregatorMenu.Products {
		id := product.ExtID
		if product.PosID != "" {
			id = product.PosID
		}

		if posProduct, ok := posProductsMap[id]; !ok {
			continue
		} else {
			if len(posProduct.Price) == 0 {
				continue
			}

			if posProduct.Price[0].Value == 0 {
				continue
			}

			attributeGroupIds := make([]string, 0, 4)

			for _, attributeGroupId := range posProduct.AttributesGroups {
				if _, exist := nonDeliveryAttributeGroupIdsMap[attributeGroupId]; exist {
					continue
				}

				attributeGroupIds = append(attributeGroupIds, attributeGroupId)
			}

			product.AttributesGroups = attributeGroupIds
			product.MenuDefaultAttributes = nil
			product.PosID = posProduct.ExtID
			product.IsAvailable = posProduct.IsAvailable
			product.Price = posProduct.Price

			products = append(products, product)
		}
	}

	return products
}

func (s *Service) removeUnnecessaryAttributeGroups(attributeGroups []coreMenuModels.AttributeGroup, used map[string]struct{}, posAttributesMap map[string]coreMenuModels.Attribute) []coreMenuModels.AttributeGroup {
	newAttributeGroups := make([]coreMenuModels.AttributeGroup, 0)

	for _, group := range attributeGroups {
		if _, ok := used[group.ExtID]; !ok {
			continue
		}

		newAttributes := make([]string, 0)

		for _, attributeId := range group.Attributes {
			if _, exist := posAttributesMap[attributeId]; !exist {
				continue
			} else {
				newAttributes = append(newAttributes, attributeId)
			}
		}

		group.Attributes = newAttributes

		newAttributeGroups = append(newAttributeGroups, group)
	}

	return newAttributeGroups
}

func (s *Service) getPosAttributesMap(posMenu coreMenuModels.Menu) map[string]coreMenuModels.Attribute {
	posAttributesMap := make(map[string]coreMenuModels.Attribute)

	for _, attribute := range posMenu.Attributes {
		posAttributesMap[attribute.ExtID] = attribute
	}

	return posAttributesMap
}

func (s *Service) getPosProductsMap(posMenu coreMenuModels.Menu) map[string]coreMenuModels.Product {
	posProductsMap := make(map[string]coreMenuModels.Product)

	for _, product := range posMenu.Products {
		posProductsMap[product.ExtID] = product
	}

	return posProductsMap
}

func (s *Service) AutoUpdateAggregatorMenu(ctx context.Context, store storeModels.Store, aggregatorMenuId string) error {
	posMenu, err := s.repo.FindById(ctx, store.MenuID)
	if err != nil {
		return err
	}

	aggregatorMenu, err := s.repo.FindById(ctx, aggregatorMenuId)
	if err != nil {
		return err
	}

	posProductsMap := s.getPosProductsMap(*posMenu)
	posAttributesMap := s.getPosAttributesMap(*posMenu)

	// coffee boom
	nonDeliveryAttributeGroupIdsMap := map[string]struct{}{
		"071cf1a2-ec6d-4750-9873-aba4ef2d1103": {},
		"a14a78e1-5b34-4d68-b46f-3dff03dfa4d3": {},
		"d0b3ddad-ddba-4e46-88c4-f152b13f2338": {},
		"c901dbb1-7bbc-49ba-b562-3570fe13c8c7": {},
	}

	products := s.matchAggregatorAndPosProducts(*aggregatorMenu, posProductsMap, nonDeliveryAttributeGroupIdsMap)

	used := s.getUsedAttributeGroupsMapInProducts(*aggregatorMenu, posProductsMap)

	attributeGroups := s.removeUnnecessaryAttributeGroups(posMenu.AttributesGroups, used, posAttributesMap)

	newAggregatorMenu := s.getNewAggregatorMenu(*posMenu, *aggregatorMenu)

	newAggregatorMenu.Products = products
	newAggregatorMenu.AttributesGroups = attributeGroups

	return s.repo.UpdateMenuEntities(ctx, aggregatorMenuId, newAggregatorMenu)
}

func (s *Service) GenerateAggregatorMenuFromPosMenu(ctx context.Context, store storeModels.Store, aggregatorMenuId, deliveryService string) (string, error) {
	posMenu, err := s.repo.FindById(ctx, store.MenuID)
	if err != nil {
		return "", err
	}

	aggregatorMenu, err := s.repo.FindById(ctx, aggregatorMenuId)
	if err != nil {
		return "", err
	}

	posProductsMap := s.getPosProductsMap(*posMenu)
	posAttributesMap := s.getPosAttributesMap(*posMenu)

	// coffee boom
	nonDeliveryAttributeGroupIdsMap := map[string]struct{}{
		"071cf1a2-ec6d-4750-9873-aba4ef2d1103": {},
		"a14a78e1-5b34-4d68-b46f-3dff03dfa4d3": {},
		"d0b3ddad-ddba-4e46-88c4-f152b13f2338": {},
		"c901dbb1-7bbc-49ba-b562-3570fe13c8c7": {},
	}

	products := s.matchAggregatorAndPosProducts(*aggregatorMenu, posProductsMap, nonDeliveryAttributeGroupIdsMap)

	used := s.getUsedAttributeGroupsMapInProducts(*aggregatorMenu, posProductsMap)

	attributeGroups := s.removeUnnecessaryAttributeGroups(posMenu.AttributesGroups, used, posAttributesMap)

	newAggregatorMenu := s.getNewAggregatorMenu(*posMenu, *aggregatorMenu)

	newAggregatorMenu.Products = products
	newAggregatorMenu.AttributesGroups = attributeGroups

	return s.repo.Insert(ctx, newAggregatorMenu)
}

func (s *Service) GetLongestCookingTimeByProductIds(ctx context.Context, menuId string, productIds []string) (int32, error) {
	products, err := s.GetProductsByMenuIDAndExtIds(ctx, menuId, productIds)
	if err != nil {
		return 0, err
	}
	var res int32
	for _, item := range products {
		if item.CookingTime > res {
			res = item.CookingTime
		}
	}

	return res, nil
}

func (s *Service) GetProductsByMenuIDAndExtIds(ctx context.Context, menuId string, productExtIds []string) (coreMenuModels.Products, error) {
	return s.repo.GetProductsByMenuIDAndExtIds(ctx, menuId, productExtIds)
}

func (s *Service) GetProductsBySectionID(ctx context.Context, menuID, sectionID string) (coreMenuModels.Products, error) {
	return s.repo.GetProductsByMenuIDAndSectionID(ctx, menuID, sectionID)
}

func (s *Service) UpdateProductsImageAndDescription(ctx context.Context, menuID string, req []coreMenuModels.UpdateProductImageAndDescription) error {
	return s.repo.UpdateProductsImageAndDescription(ctx, menuID, req)
}

func (s *Service) UpdateAttributesPrice(ctx context.Context, menuID string, req []coreMenuModels.UpdateAttributePrice) error {
	return s.repo.UpdateAttributesPrice(ctx, menuID, req)
}

func (s *Service) GetMenuById(ctx context.Context, menuId string) (coreMenuModels.Menu, error) {
	menu, err := s.repo.FindById(ctx, menuId)
	if err != nil {
		return coreMenuModels.Menu{}, err
	}

	return *menu, nil
}

func (s *Service) AddNameInProduct(ctx context.Context, req coreMenuModels.AddLanguageDescriptionRequest) error {
	if err := s.repo.AddNameInProduct(ctx, req); err != nil {
		return err
	}
	return nil
}

func (s *Service) AddDescriptionInProduct(ctx context.Context, req coreMenuModels.AddLanguageDescriptionRequest) error {
	if err := s.repo.AddDescriptionInProduct(ctx, req); err != nil {
		return err
	}
	return nil
}

func (s *Service) AddNameInSection(ctx context.Context, req coreMenuModels.AddLanguageDescriptionRequest) error {
	if err := s.repo.AddNameInSection(ctx, req); err != nil {
		return err
	}
	return nil
}

func (s *Service) AddDescriptionInSection(ctx context.Context, req coreMenuModels.AddLanguageDescriptionRequest) error {
	if err := s.repo.AddDescriptionInSection(ctx, req); err != nil {
		return err
	}
	return nil
}

func (s *Service) AddNameInAttributeGroup(ctx context.Context, req coreMenuModels.AddLanguageDescriptionRequest) error {
	if err := s.repo.AddNameInAttributeGroup(ctx, req); err != nil {
		return err
	}
	return nil
}

func (s *Service) AddNameInAttribute(ctx context.Context, req coreMenuModels.AddLanguageDescriptionRequest) error {
	if err := s.repo.AddNameInAttribute(ctx, req); err != nil {
		return err
	}
	return nil
}

func (s *Service) ChangeNameInProduct(ctx context.Context, req coreMenuModels.AddLanguageDescriptionRequest) error {
	if err := s.repo.ChangeNameInProduct(ctx, req); err != nil {
		return err
	}
	return nil
}

func (s *Service) ChangeDescriptionInProduct(ctx context.Context, req coreMenuModels.AddLanguageDescriptionRequest) error {
	if err := s.repo.ChangeDescriptionInProduct(ctx, req); err != nil {
		return err
	}
	return nil
}

func (s *Service) ChangeNameInSection(ctx context.Context, req coreMenuModels.AddLanguageDescriptionRequest) error {
	if err := s.repo.ChangeNameInSection(ctx, req); err != nil {
		return err
	}
	return nil
}

func (s *Service) ChangeDescriptionInSection(ctx context.Context, req coreMenuModels.AddLanguageDescriptionRequest) error {
	if err := s.repo.ChangeDescriptionInSection(ctx, req); err != nil {
		return err
	}
	return nil
}

func (s *Service) ChangeNameInAttributeGroup(ctx context.Context, req coreMenuModels.AddLanguageDescriptionRequest) error {
	if err := s.repo.ChangeNameInAttributeGroup(ctx, req); err != nil {
		return err
	}
	return nil
}

func (s *Service) ChangeNameInAttribute(ctx context.Context, req coreMenuModels.AddLanguageDescriptionRequest) error {
	if err := s.repo.ChangeNameInAttribute(ctx, req); err != nil {
		return err
	}
	return nil
}

func (s *Service) AddRegulatoryInformation(ctx context.Context, req coreMenuModels.RegulatoryInformationRequest) error {
	if err := s.repo.AddRegulatoryInformation(ctx, req); err != nil {
		return err
	}
	return nil
}

func (s *Service) ChangeRegulatoryInformation(ctx context.Context, req coreMenuModels.RegulatoryInformationRequest) error {
	if err := s.repo.ChangeRegulatoryInformation(ctx, req); err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateExcludedFromMenuProduct(ctx context.Context, menuID string, productIDs []string) error {
	if err := s.repo.UpdateExcludedFromMenuProduct(ctx, menuID, productIDs); err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateProductsDisabledByValidation(ctx context.Context, menuID string, productIDs []string, disabledByValidation bool) error {
	if err := s.repo.UpdateProductsDisabledByValidation(ctx, menuID, productIDs, disabledByValidation); err != nil {
		return err
	}
	return nil
}
