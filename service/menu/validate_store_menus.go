package menu

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strings"
)

type AttributeGroupReport struct {
	Name                      string            `json:"name"`
	Description               string            `json:"description"`
	PosMin                    int               `json:"pos_menu_min"`
	PosMax                    int               `json:"pos_menu_max"`
	AggregatorAttributesCount int               `json:"aggregator_attributes_count"`
	AggregatorAttributes      []AttributeReport `json:"aggregator_attributes"`
	PosAttributes             []AttributeReport `json:"pos_attributes"`
}

type AttributeReport struct {
	ID                string `json:"id"`
	POSID             string `json:"pos_id"`
	Name              string `json:"name"`
	Position          int    `json:"position"`
	AttributeGroupId  string `json:"-"`
	AttributeGroupMax int    `json:"-"`
	IsDefault         bool   `json:"is_default_attribute"`
}

type MenuAnalytics struct {
	StoreName         string `json:"store_name"`
	ValidMenuCount    int    `json:"menu_valid"`
	NotValidMenuCount int    `json:"menu_with_errors"`
}

type MenuDetail struct {
	RestaurantName string          `json:"restaurant_name"`
	ID             string          `json:"menu_id"`
	Name           string          `json:"menu_name"`
	Delivery       string          `json:"delivery"`
	ProductDetails []ProductDetail `json:"products"`
}

type ProductDetail struct {
	ID                      string                 `json:"id"`
	Name                    string                 `json:"name"`
	Position                int                    `json:"position"`
	Description             string                 `json:"description"`
	AttributeGroupReport    []AttributeGroupReport `json:"attribute_group_report"`
	NonExistingAttributeIds []AttributeReport      `json:"non_existing_attribute_ids_in_pos_menu"`
}

type AttributeWithPosition struct {
	Attribute models.Attribute
	Position  int
}

type VirtualStore struct {
	Store       coreStoreModels.Store
	MenuService *Service
}

func newVirtualStore(store coreStoreModels.Store, menuService *Service) VirtualStore {
	return VirtualStore{
		Store:       store,
		MenuService: menuService,
	}
}

func getAttributeMap(attributes []models.Attribute) map[string]AttributeWithPosition {
	result := make(map[string]AttributeWithPosition, len(attributes))

	for index, attribute := range attributes {
		result[attribute.ExtID] = AttributeWithPosition{
			Attribute: attribute,
			Position:  index,
		}
		if attribute.PosID != "" {
			result[attribute.PosID] = AttributeWithPosition{
				Attribute: attribute,
				Position:  index,
			}
		}

	}

	return result
}

func getProductMap(products []models.Product) map[string]models.Product {
	result := make(map[string]models.Product, len(products))

	for _, product := range products {
		result[product.ExtID] = product
	}

	return result
}

func getAttributeGroupMap(attributeGroups []models.AttributeGroup) map[string]models.AttributeGroup {
	result := make(map[string]models.AttributeGroup, len(attributeGroups))

	for _, group := range attributeGroups {
		result[group.ExtID] = group
	}

	return result
}

func getAllBelongingAttributesOfTheProduct(product models.Product, attributeGroupMap map[string]models.AttributeGroup, defaultAttributes []models.MenuDefaultAttributes, aggregatorAttributeMap map[string]AttributeWithPosition) map[string]AttributeReport {
	attributeIDsMap := make(map[string]AttributeReport, len(defaultAttributes))

	for _, attributeGroupId := range product.AttributesGroups {
		if group, ok := attributeGroupMap[attributeGroupId]; !ok {
			// TODO: return err if not exist
			continue
		} else {
			for _, attributeId := range group.Attributes {
				if attributeId == "service_fee" {
					continue
				}

				if attributeWithPosition, ok := aggregatorAttributeMap[attributeId]; !ok {
					continue
				} else {
					id := attributeWithPosition.Attribute.ExtID
					if attributeWithPosition.Attribute.PosID != "" {
						id = attributeWithPosition.Attribute.PosID
					}

					attributeIDsMap[id] = AttributeReport{
						ID:                id,
						POSID:             attributeWithPosition.Attribute.PosID,
						Name:              attributeWithPosition.Attribute.Name,
						Position:          attributeWithPosition.Position,
						AttributeGroupMax: group.Max,
						AttributeGroupId:  group.ExtID,
					}
				}
			}
		}
	}

	for _, defAttribute := range defaultAttributes {
		if defAttribute.ExtID == "service_fee" {
			continue
		}

		if defAttribute.ByAdmin {
			attributeIDsMap[defAttribute.ExtID] = AttributeReport{
				ID:        defAttribute.ExtID,
				Name:      defAttribute.Name,
				IsDefault: true,
			}
		}
	}

	return attributeIDsMap
}

func validateProductRestrictions(posProductsMap map[string]models.Product, aggregatorProducts []models.Product,
	aggregatorAttributeGroupsMap, posAttributeGroupsMap map[string]models.AttributeGroup,
	posAttributeMap, aggregatorAttributeMap map[string]AttributeWithPosition) []ProductDetail {

	productDetails := make([]ProductDetail, 0)

	for index, aggregatorProduct := range aggregatorProducts {
		if aggregatorProduct.IsDeleted || !aggregatorProduct.IsSync {
			continue
		}

		id := aggregatorProduct.ExtID
		if aggregatorProduct.PosID != "" {
			id = aggregatorProduct.PosID
		}

		if posProduct, exist := posProductsMap[id]; !exist {
			continue
		} else {
			attributeGroupReport, attributeReport := validateAttributeGroupsRestrictions(posProduct, aggregatorProduct, aggregatorAttributeGroupsMap, posAttributeGroupsMap, posAttributeMap, aggregatorAttributeMap, posProductsMap)

			var name string
			if len(aggregatorProduct.Name) != 0 {
				name = aggregatorProduct.Name[0].Value
			}

			detail := ProductDetail{
				ID:       aggregatorProduct.ExtID,
				Name:     name,
				Position: index,
			}

			var hasError bool
			if len(attributeGroupReport) != 0 {
				detail.AttributeGroupReport = attributeGroupReport
				hasError = true
			}

			if len(attributeReport) != 0 {
				detail.NonExistingAttributeIds = attributeReport
				hasError = true
			}

			if hasError {
				productDetails = append(productDetails, detail)
			}
		}
	}

	return productDetails
}

func validateAttributeGroupsRestrictions(posProduct, aggregatorProduct models.Product,
	attributeGroupMap, posAttributeGroupsMap map[string]models.AttributeGroup,
	posAttributeMap, aggregatorAttributeMap map[string]AttributeWithPosition, posProductsMap map[string]models.Product) ([]AttributeGroupReport, []AttributeReport) {

	groupReports := make([]AttributeGroupReport, 0)

	aggregatorProductAttributeIDsMap := getAllBelongingAttributesOfTheProduct(aggregatorProduct, attributeGroupMap, posProduct.MenuDefaultAttributes, aggregatorAttributeMap)

	for _, groupId := range posProduct.AttributesGroups {
		count := 0

		group, exist := posAttributeGroupsMap[groupId]
		if !exist {
			continue
		}

		var (
			posAttributesName        = make([]AttributeReport, 0, len(group.Attributes))
			aggregatorAttributesName = make([]AttributeReport, 0, len(aggregatorProductAttributeIDsMap))
		)

		unique := make(map[string]struct{})

		for _, attributeId := range group.Attributes {
			if report, ok := aggregatorProductAttributeIDsMap[attributeId]; ok {
				if report.IsDefault {
					count++
				} else if _, has := unique[report.AttributeGroupId]; !has {
					unique[report.AttributeGroupId] = struct{}{}
					count += report.AttributeGroupMax
				}

				if attributeWithPosition, existId := aggregatorAttributeMap[attributeId]; existId {
					aggregatorAttributesName = append(aggregatorAttributesName, AttributeReport{
						ID:        attributeWithPosition.Attribute.ExtID,
						Name:      attributeWithPosition.Attribute.Name,
						Position:  attributeWithPosition.Position,
						IsDefault: report.IsDefault,
					})
				}
				if attributeWithPosition, existId := posAttributeMap[attributeId]; existId {
					posAttributesName = append(posAttributesName, AttributeReport{
						ID:       attributeWithPosition.Attribute.ExtID,
						Name:     attributeWithPosition.Attribute.Name,
						Position: attributeWithPosition.Position,
					})
				}
				delete(aggregatorProductAttributeIDsMap, attributeId)
			} else {
				if attributeWithPosition, existId := posAttributeMap[attributeId]; existId {
					posAttributesName = append(posAttributesName, AttributeReport{
						ID:       attributeWithPosition.Attribute.ExtID,
						Name:     attributeWithPosition.Attribute.Name,
						Position: attributeWithPosition.Position,
					})
				}
			}
		}

		if count < group.Min {
			groupReport := AttributeGroupReport{
				Name:                      group.Name,
				Description:               "Общее количество привязанных атрибутов к продукту в агрегатор меню меньше min атрибут группы в POS меню",
				PosMax:                    group.Max,
				AggregatorAttributesCount: count,
				AggregatorAttributes:      aggregatorAttributesName,
				PosAttributes:             posAttributesName,
			}

			groupReports = append(groupReports, groupReport)
			// TODO: add ProductDetail количество атрибутов меньше положенного
		}

		if count > group.Max {
			groupReport := AttributeGroupReport{
				Name:                      group.Name,
				Description:               "Общее количество привязанных атрибутов к продукту в агрегатор меню больше max атрибут группы в POS меню",
				PosMax:                    group.Max,
				AggregatorAttributesCount: count,
				AggregatorAttributes:      aggregatorAttributesName,
				PosAttributes:             posAttributesName,
			}

			groupReports = append(groupReports, groupReport)
		}
	}

	for _, attributeId := range posProduct.Attributes {
		delete(aggregatorProductAttributeIDsMap, attributeId)
	}

	attributeReports := make([]AttributeReport, 0, len(aggregatorProductAttributeIDsMap))

	if len(aggregatorProductAttributeIDsMap) != 0 {
		for _, attributeReport := range aggregatorProductAttributeIDsMap {
			if _, ok := posProductsMap[attributeReport.ID]; ok {
				continue
			}
			if _, ok := posProductsMap[attributeReport.POSID]; ok {
				continue
			}
			attributeReports = append(attributeReports, attributeReport)
		}
	}

	return groupReports, attributeReports
}

func validateProductMatching(posProductsMap map[string]models.Product, aggregatorProducts []models.Product, posType string) []ProductDetail {
	ProductDetails := make([]ProductDetail, 0)

	for index, product := range aggregatorProducts {
		if product.IsDeleted || !product.IsSync {
			continue
		}

		id := product.ExtID
		if product.PosID != "" {
			id = product.PosID
		}

		var name string
		if len(product.Name) != 0 {
			name = product.Name[0].Value
		}

		if posProduct, exist := posProductsMap[id]; !exist {

			ProductDetails = append(ProductDetails, ProductDetail{
				ID:          product.ExtID,
				Name:        name,
				Position:    index,
				Description: "Продукт не существует в POS меню",
			})
		} else {
			if posProduct.IsDeleted {
				ProductDetails = append(ProductDetails, ProductDetail{
					ID:          product.ExtID,
					Name:        name,
					Position:    index,
					Description: "Продукт удален в POS меню",
				})
			}

			if !posProduct.IsIncludedInMenu && posType == models.IIKO.String() {
				ProductDetails = append(ProductDetails, ProductDetail{
					ID:          product.ExtID,
					Name:        name,
					Position:    index,
					Description: "Продукт не включен в POS меню",
				})
			}

		}

	}

	return ProductDetails
}

func getActiveDeliveries(store coreStoreModels.Store) map[string]struct{} {
	deliveries := make(map[string]struct{})

	if store.Glovo.SendToPos {
		deliveries["glovo"] = struct{}{}
	}

	if store.Wolt.SendToPos {
		deliveries["wolt"] = struct{}{}
	}

	if store.Express24.SendToPos {
		deliveries["express24"] = struct{}{}
	}

	for _, external := range store.ExternalConfig {
		if external.SendToPos {
			deliveries[external.Type] = struct{}{}
		}
	}

	return deliveries
}

func (s *Service) updateDeletedStatus(ctx context.Context, productMatching []ProductDetail, menuId string, store coreStoreModels.Store) {
	productIds := make([]string, 0, len(productMatching))

	for _, productId := range productMatching {
		productIds = append(productIds, productId.ID)
	}

	if err := s.UpdateProductsDeletedStatus(ctx, menuId, productIds, true, "validation_menu"); err != nil {
		log.Err(err).Msgf("update products deleted status error for restaurant name %s, id %s", store.Name, store.ID)
	}
}

func (s *Service) updateExcludedStatus(ctx context.Context, productMatching []ProductDetail, menuID string, store coreStoreModels.Store) {
	var productIds []string

	for _, productId := range productMatching {
		if productId.Description == "Продукт не включен в POS меню" {
			productIds = append(productIds, productId.ID)
		}
	}

	if err := s.UpdateExcludedFromMenuProduct(ctx, menuID, productIds); err != nil {
		log.Err(err).Msgf("update products excluded from menu status error for store: %s, %s", store.ID, store.Name)
	}
}

func (s *Service) ValidateStoreMenus(ctx context.Context, store coreStoreModels.Store) ([]MenuDetail, MenuAnalytics, error) {
	posMenu, err := s.repo.FindById(ctx, store.MenuID)
	if err != nil {
		return nil, MenuAnalytics{}, err
	}

	posProductsMap := getProductMap(posMenu.Products)
	posAttributeGroupsMap := getAttributeGroupMap(posMenu.AttributesGroups)
	posAttributeMapWithPosition := getAttributeMap(posMenu.Attributes)

	var (
		allErrs       error
		menuDetails   = make([]MenuDetail, 0)
		menuAnalytics = MenuAnalytics{
			StoreName: store.Name,
		}
	)

	deliveries := getActiveDeliveries(store)

	for _, menuDS := range store.Menus {
		if _, ok := deliveries[menuDS.Delivery]; !ok || !menuDS.IsActive {
			continue
		}

		productDetails := make([]ProductDetail, 0)

		aggregatorMenu, err := s.repo.FindById(ctx, menuDS.ID)
		if err != nil {
			allErrs = errors.Wrap(allErrs, err.Error())
			continue
		}

		if s.isOneShotVirtual(store.ID) {
			virtualStore := newVirtualStore(store, s)
			aggregatorMenu, posProductsMap, posAttributeGroupsMap, posAttributeMapWithPosition, err = virtualStore.getVirtualStoreItems(ctx, menuDS.ID)
			if err != nil {
				log.Err(err).Msgf("get virtual items error, storeID: %s, aggregatorMenuID: %s", store.ID, menuDS.ID)
				continue
			}
		}

		productMatching := validateProductMatching(posProductsMap, aggregatorMenu.Products, store.PosType)
		if len(productMatching) != 0 {
			s.updateDeletedStatus(ctx, productMatching, menuDS.ID, store)
			s.updateExcludedStatus(ctx, productMatching, menuDS.ID, store)

			productDetails = append(productDetails, productMatching...)
		}

		aggregatorAttributeGroupsMap := getAttributeGroupMap(aggregatorMenu.AttributesGroups)
		aggregatorAttributeMapWithPosition := getAttributeMap(aggregatorMenu.Attributes)

		productRestrictions := validateProductRestrictions(posProductsMap, aggregatorMenu.Products, aggregatorAttributeGroupsMap, posAttributeGroupsMap, posAttributeMapWithPosition, aggregatorAttributeMapWithPosition)
		if len(productRestrictions) != 0 {
			productDetails = append(productDetails, productRestrictions...)
		}

		for _, productRestriction := range productRestrictions {
			if productRestriction.NonExistingAttributeIds != nil && len(productRestriction.NonExistingAttributeIds) > 0 {
				if err := s.delAttrAndAttrGroupFromAggrMenu(ctx, *aggregatorMenu, productRestriction, aggregatorAttributeGroupsMap); err != nil {
					log.Err(err).Msgf("del attribute, attribute group from aggregator menu error, menu id: %s", aggregatorMenu.ID)
				}
			}
		}

		if len(productDetails) != 0 {
			menuDetails = append(menuDetails, MenuDetail{
				RestaurantName: store.Name,
				ID:             menuDS.ID,
				Name:           menuDS.Name,
				Delivery:       menuDS.Delivery,
				ProductDetails: productDetails,
			})
			menuAnalytics.NotValidMenuCount++
		} else {
			menuAnalytics.ValidMenuCount++
		}
	}

	return menuDetails, menuAnalytics, allErrs
}
func (s *Service) isOneShotVirtual(storeID string) bool {
	if storeID == "66a72f2dce5e86021dd7d739" {
		return true
	}
	return false
}

func (virt *VirtualStore) getVirtualStoreItems(ctx context.Context, aggrMenuID string) (*models.Menu, map[string]models.Product, map[string]models.AttributeGroup, map[string]AttributeWithPosition, error) {
	aggrMenu, virtualStoreMap, err := virt.getAggregatorMenuAndStoreMap(ctx, aggrMenuID)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	resPosProductMap := make(map[string]models.Product)
	resPosAttributeGroupsMap := make(map[string]models.AttributeGroup)
	resPosAttributeMapWithPosition := make(map[string]AttributeWithPosition)

	for _, storeID := range virtualStoreMap {
		store, err := virt.MenuService.storeService.GetByID(ctx, storeID)
		if err != nil {
			log.Err(err).Msgf("get store by id error, storeID: %s", storeID)
			return nil, nil, nil, nil, err
		}

		posMenu, err := virt.MenuService.repo.FindById(ctx, store.MenuID)
		if err != nil {
			log.Err(err).Msgf("get pos menu error, storeID: %s", storeID)
			return nil, nil, nil, nil, err
		}

		posProductsMap := getProductMap(posMenu.Products)
		posAttributeGroupsMap := getAttributeGroupMap(posMenu.AttributesGroups)
		posAttributeMapWithPosition := getAttributeMap(posMenu.Attributes)

		for k, v := range posProductsMap {
			if _, ok := resPosProductMap[k]; ok {
				continue
			}
			resPosProductMap[k] = v
		}

		for k, v := range posAttributeGroupsMap {
			if _, ok := resPosAttributeGroupsMap[k]; ok {
				continue
			}
			resPosAttributeGroupsMap[k] = v
		}

		for k, v := range posAttributeMapWithPosition {
			if _, ok := resPosAttributeMapWithPosition[k]; ok {
				continue
			}
			resPosAttributeMapWithPosition[k] = v
		}
	}

	return aggrMenu, resPosProductMap, resPosAttributeGroupsMap, resPosAttributeMapWithPosition, nil
}

func (virt *VirtualStore) getAggregatorMenuAndStoreMap(ctx context.Context, aggrMenuID string) (*models.Menu, map[string]string, error) {
	virtualStoreMap := make(map[string]string)

	aggrMenu, err := virt.MenuService.repo.FindById(ctx, aggrMenuID)
	if err != nil {
		return nil, nil, err
	}

	for i := range aggrMenu.Products {
		extID := aggrMenu.Products[i].ExtID
		storeID, productExtID, isOk := virt.isVirtualStoreExtID(extID)
		if isOk {
			aggrMenu.Products[i].ExtID = productExtID

			if _, ok := virtualStoreMap[storeID]; !ok {
				virtualStoreMap[storeID] = storeID
			}
		}
	}

	for i := range aggrMenu.Attributes {
		extID := aggrMenu.Attributes[i].ExtID

		storeID, attributeExtID, isOk := virt.isVirtualStoreExtID(extID)
		if isOk {
			aggrMenu.Attributes[i].ExtID = attributeExtID

			if _, ok := virtualStoreMap[storeID]; !ok {
				virtualStoreMap[storeID] = storeID
			}
		}
	}

	return aggrMenu, virtualStoreMap, nil
}

func (virt *VirtualStore) isVirtualStoreExtID(extID string) (string, string, bool) {
	if len(extID) > 36 && strings.Contains(extID, "_") {
		storeID, newExtID := virt.getStoreIDExtIDByVirtualExtID(extID)
		return storeID, newExtID, true
	}
	return "", "", false
}

func (virt *VirtualStore) getStoreIDExtIDByVirtualExtID(virtualExtID string) (string, string) {
	part := strings.Split(virtualExtID, "_")
	storeID := part[0]
	extID := part[1]

	return storeID, extID
}

func (s *Service) delAttrAndAttrGroupFromAggrMenu(ctx context.Context, aggrMenu models.Menu, productRestriction ProductDetail, aggregatorAttributeGroupsMap map[string]models.AttributeGroup) error {
	productMap := make(map[string]models.Product)
	for _, product := range aggrMenu.Products {
		productMap[product.ExtID] = product
	}

	nonExistAttrs := s.getAttrsForDel(productRestriction.NonExistingAttributeIds)

	if _, ok := productMap[productRestriction.ID]; !ok {
		return nil
	}

	var delAttrGroupIDs []string
	for _, attrGroupID := range productMap[productRestriction.ID].AttributesGroups {
		attrGroup, ok := aggregatorAttributeGroupsMap[attrGroupID]
		if !ok {
			continue
		}

		if isSame := compareAttributes(nonExistAttrs, attrGroup.Attributes); !isSame {
			continue
		}

		delAttrGroupIDs = append(delAttrGroupIDs, attrGroup.ExtID)
	}

	switch {
	case len(delAttrGroupIDs) != 0:
		for _, attrGroupID := range delAttrGroupIDs {
			if err := s.repo.DeleteAttrGroupFromProduct(ctx, aggrMenu.ID, productRestriction.ID, attrGroupID); err != nil {
				log.Err(err).Msgf("delete attribute group from product error, menu id: %s", aggrMenu.ID)
			}
		}
	default:
		if err := s.repo.DeleteAttributesFromAttributeGroup(ctx, aggrMenu.ID, nonExistAttrs); err != nil {
			log.Err(err).Msgf("delete attributes from attribute group error, menu id: %s", aggrMenu.ID)
		}
	}

	return nil
}

func compareAttributes(arr1, arr2 []string) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	m := make(map[string]struct{})
	for i := range arr1 {
		m[arr1[i]] = struct{}{}
	}

	for i := range arr2 {
		if _, ok := m[arr2[i]]; ok {
			delete(m, arr2[i])
			continue
		}
		return false
	}

	return len(m) == 0
}

func (s *Service) getAttrsForDel(attrReport []AttributeReport) []string {
	var delAttrs []string
	for i := range attrReport {
		delAttrs = append(delAttrs, attrReport[i].ID)
	}

	return removeDuplicates(delAttrs)
}

func removeDuplicates(arr []string) []string {
	m := make(map[string]struct{}, len(arr))
	res := []string{}

	for _, v := range arr {
		if _, ok := m[v]; !ok {
			m[v] = struct{}{}
			res = append(res, v)
		}
	}

	return res
}
