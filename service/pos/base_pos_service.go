package pos

import (
	"context"
	"fmt"
	errs "github.com/kwaaka-team/orders-core/core/errors"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/utils"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	menuCore "github.com/kwaaka-team/orders-core/pkg/menu"
	menuModels "github.com/kwaaka-team/orders-core/pkg/menu/dto"
	menuUtils "github.com/kwaaka-team/orders-core/pkg/menu/utils"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

var constructorError = errors.New("base pos service is nil")

type BasePosService struct {
	anotherBillStoreIDs map[string][]string
}

func (bps *BasePosService) GetOrderStatus(ctx context.Context, order models.Order) (string, error) {
	return "", ErrUnsupportedMethod
}

func (bps *BasePosService) GetBalanceLimit(ctx context.Context, store coreStoreModels.Store) int {
	return 0
}

func (bps *BasePosService) IsStopListByBalance(ctx context.Context, store coreStoreModels.Store) bool {
	return false
}

func (bps *BasePosService) IsAliveStatus(ctx context.Context, store coreStoreModels.Store) (bool, error) {
	return false, nil
}

func (bps *BasePosService) AwakeTerminal(ctx context.Context, store coreStoreModels.Store) error {
	return nil
}

func setPosOrderId(order models.Order, id string) models.Order {
	order.PosOrderID = id
	return order
}

func (bps *BasePosService) SetPosRequestBodyToOrder(order models.Order, req interface{}) (models.Order, error) {
	var err error

	order.LogMessages.ToPos, err = utils.GetJsonFormatFromModel(req)
	if err != nil {
		return order, err
	}

	return order, nil
}

func activeMenuPositions(aggregatorMenu coreMenuModels.Menu) (map[string]string, map[string]string, map[string]float64) {
	productsMap := make(map[string]string, len(aggregatorMenu.Products))
	discountPrices := make(map[string]float64)

	for _, product := range aggregatorMenu.Products {
		if product.PosID == "" {
			productsMap[product.ExtID] = product.ExtID
			continue
		}

		if product.DiscountPrice.IsActive {
			discountPrices[product.ExtID] = product.DiscountPrice.Value
		}

		productsMap[product.ExtID] = product.PosID
	}

	attributesMap := make(map[string]string, len(aggregatorMenu.Attributes))

	for _, attribute := range aggregatorMenu.Attributes {
		if attribute.PosID == "" {
			attributesMap[attribute.ExtID] = attribute.ExtID
			continue
		}
		if attribute.PosID != "" && attribute.ExtID == attribute.PosID {
			attributesMap[attribute.ExtID] = attribute.ExtID
			continue
		}

		if attribute.PosID != "" && attribute.ExtID != attribute.PosID {
			attributesMap[attribute.ExtID] = attribute.ExtID
			continue
		}

		attributesMap[attribute.ExtID] = attribute.PosID
	}

	return productsMap, attributesMap, discountPrices
}
func activeAggregatorsMenuAttributesPosition(aggregatorMenu coreMenuModels.Menu) map[string]coreMenuModels.Attribute {
	attributesMap := make(map[string]coreMenuModels.Attribute, len(aggregatorMenu.Attributes))

	for _, attribute := range aggregatorMenu.Attributes {
		attributesMap[attribute.ExtID] = attribute
	}

	return attributesMap
}

func getPromosMap(ctx context.Context, req models.Order, menuCli menuCore.Client) (models.Order, map[string]menuModels.PromoDiscount, map[string]string, map[string]menuModels.ProductPercentage, error) {
	productsIDs := make([]string, 0, len(req.Products))
	for _, product := range req.Products {
		productsIDs = append(productsIDs, product.ID)
	}

	promosMap := make(map[string]menuModels.PromoDiscount)
	promos, err := menuCli.GetStorePromos(ctx, menuModels.GetPromosSelector{
		StoreID:         req.RestaurantID,
		DeliveryService: req.DeliveryService,
		IsActive:        true,
	})

	if err != nil {
		log.Error().Err(err).Msgf("failed to get promos")
		return req, nil, nil, nil, err
	}

	giftMap := make(map[string]string)
	promoWithPercentMap := make(map[string]menuModels.ProductPercentage)

	for _, promo := range promos {
		if promo.HasAggregatorAndPartnerDiscount {
			req.HasAggregatorAndPartnerDiscount = true
		}

		for _, productID := range promo.ProductIds {
			promosMap[productID] = promo
		}

		for _, gift := range promo.ProductGifts {
			giftMap[gift.PromoId] = gift.ProductId
		}

		if promo.PercentageForEachProduct {
			for _, product := range promo.ProductsPercentage {
				promoWithPercentMap[product.ProductID] = product
			}
		}
	}

	return req, promosMap, giftMap, promoWithPercentMap, nil
}

func applyOrderDiscount(ctx context.Context, req models.Order, promosMap map[string]menuModels.PromoDiscount, giftMap map[string]string, promoWithPercentMap map[string]menuModels.ProductPercentage) models.Order {
	var discounts float64

	if req.PartnerDiscountsProducts.Value == 0 && !req.HasAggregatorAndPartnerDiscount && req.DeliveryService == models.GLOVO.String() {
		return req
	}

	for i := 0; i < len(req.Products); i++ {

		if productID, ok := giftMap[req.Products[i].ID]; ok {
			req.Products[i].ID = productID
			req.Products[i].Price.Value = 0
			continue
		}

		if promo, ok := promosMap[req.Products[i].ID]; ok {

			if promo.Type == models.Discount.String() {

				switch {
				// if the promo has a percentage for each product, the percentages for each product are counted
				case promo.PercentageForEachProduct:
					req.Products[i].Price.Value, discounts = countDiscountForEachPosition(
						req.Products[i].ID,
						req.Products[i].Quantity,
						promo.Percent,
						req.Products[i].Price.Value,
						discounts,
						promoWithPercentMap)

				default:
					discountAmount := float64(promo.Percent) * req.Products[i].Price.Value / 100.0
					req.Products[i].Price.Value -= discountAmount
					discounts += discountAmount * float64(req.Products[i].Quantity)
				}
			}
		}

		for j := 0; j < len(req.Products[i].Attributes); j++ {
			if promo, ok := promosMap[req.Products[i].Attributes[j].ID]; ok {

				if promo.Type == models.Discount.String() {

					switch {
					case promo.PercentageForEachProduct:
						req.Products[i].Attributes[j].Price.Value, discounts = countDiscountForEachPosition(
							req.Products[i].Attributes[j].ID,
							req.Products[i].Attributes[j].Quantity,
							promo.Percent,
							req.Products[i].Attributes[j].Price.Value,
							discounts,
							promoWithPercentMap)

					default:
						discountAmount := float64(promo.Percent) * req.Products[i].Attributes[j].Price.Value / 100.0
						req.Products[i].Attributes[j].Price.Value -= discountAmount
						discounts += discountAmount * float64(req.Products[i].Attributes[j].Quantity)
					}
				}

			}
		}
	}

	if req.PartnerDiscountsProducts.Value == 0 && req.HasAggregatorAndPartnerDiscount {
		req.PartnerDiscountsProducts.Value = discounts
	}

	return req
}

func countDiscountForEachPosition(positionID string, positionQuantity, totalPromoPercent int, positionPrice, discounts float64, promoWithPercentMap map[string]menuModels.ProductPercentage) (float64, float64) {
	if _, ok := promoWithPercentMap[positionID]; ok {
		discountAmount := float64(promoWithPercentMap[positionID].Percent) * positionPrice / 100
		positionPrice -= discountAmount
		discounts += discountAmount * float64(positionQuantity)
		return positionPrice, discounts
	}
	discountAmount := float64(totalPromoPercent) * positionPrice / 100
	positionPrice -= discountAmount
	discounts += discountAmount * float64(positionQuantity)
	return positionPrice, discounts
}

func fullFillProducts(
	req models.Order, store coreStoreModels.Store,
	productsMap map[string]coreMenuModels.Product, attributesMap map[string]coreMenuModels.Attribute,
	attributeGroupsMap map[string]coreMenuModels.AttributeGroup, comboMap map[string]coreMenuModels.Combo, aggregatorMenu coreMenuModels.Menu,
	promosMap map[string]menuModels.PromoDiscount, promoWithPercentMap map[string]menuModels.ProductPercentage,
) (models.Order, float64, error) {

	var (
		serviceFee    float64
		orderProducts = make([]models.OrderProduct, 0, len(req.Products))
	)

	aggregatorProducts, aggregatorAttributes, discountPrices := activeMenuPositions(aggregatorMenu)
	aggregatorAttributesObjects := activeAggregatorsMenuAttributesPosition(aggregatorMenu)

	var cookingTime int32
	for _, product := range req.Products {

		if len(product.Promos) != 0 {
			if product.Promos[0].Type == "GIFT" {
				orderProduct := models.OrderProduct{
					ID:   product.ID,
					Name: product.Name,
					Price: models.Price{
						Value:        0,
						CurrencyCode: store.Settings.Currency,
					},
					Quantity: product.Quantity,
					Promos:   product.Promos,
				}

				orderProducts = append(orderProducts, orderProduct)
				continue
			}
		}

		comboProduct := comboMap[product.ID]

		product.ProgramID = comboProduct.ProgramID
		product.SourceActionID = comboProduct.SourceActionID

		var (
			orderProductAttributes = make([]models.ProductAttribute, 0, len(product.Attributes))
			modifiersPrice         float64
		)

		productPosID, exist := aggregatorProducts[product.ID]
		switch exist {
		case true:
			product.ID = productPosID
		default:
			log.Error().Msgf("product with ID %s %s not matched", product.ID, product.Name)
		}

		menuProduct, ok := productsMap[product.ID]
		if !ok {
			log.Info().Msgf("PRODUCT NOT FOUND IN POS MENU, ID %s, NAME %s", product.ID, product.Name)
			req.FailReason.Code = PRODUCT_MISSED_CODE
			req.FailReason.Message = PRODUCT_MISSED + product.ID
			return req, 0, errors.Wrap(errs.ErrProductNotFound, fmt.Sprintf("PRODUCT NOT FOUND IN POS MENU, ID %s, NAME %s", product.ID, product.Name))
		}

		var countCommonAttributes int

		if menuProduct.MenuDefaultAttributes != nil {
			countCommonAttributes = len(menuProduct.DefaultAttributes)
		}

		modifiersIds := make(map[string]struct{})
		for index, attribute := range product.Attributes {

			if attribute.ID == models.ServiceFee {
				serviceFee += attribute.Price.Value * float64(product.Quantity) * float64(attribute.Quantity)

				if index < len(product.Attributes)-countCommonAttributes {
					modifiersPrice += attribute.Price.Value * float64(product.Quantity) * float64(attribute.Quantity)
				}
				continue
			}

			attributePosID, attributeExistInAggregatorMenu := aggregatorAttributes[attribute.ID]
			attrbuteIdBeforeChange := attribute.ID
			switch attributeExistInAggregatorMenu {
			case true:
				attribute.ID = attributePosID
			default:
				log.Info().Msgf("attribute with ID %s %s not matched", attribute.ID, attribute.Name)
			}

			menuAttribute, attributeExistInPosMenu := attributesMap[attribute.ID]
			if !attributeExistInPosMenu {
				menuAttribute, attributeExistInPosMenu = attributesMap[aggregatorAttributesObjects[attrbuteIdBeforeChange].ExtID]
				if attributeExistInPosMenu {
					if menuAttribute.PosID != "" {
						attribute.ID = menuAttribute.PosID

					} else {
						attribute.ID = menuAttribute.ExtID
					}
				} else {
					menuAttribute, attributeExistInPosMenu = attributesMap[aggregatorAttributesObjects[attrbuteIdBeforeChange].PosID]
					if attributeExistInPosMenu {
						if menuAttribute.PosID != "" {
							attribute.ID = menuAttribute.PosID

						} else {
							attribute.ID = menuAttribute.ExtID
						}
					}
				}

			}

			if attributeExistInPosMenu {
				modifiersIds[attribute.ID] = struct{}{}
				attribute.GroupID = menuAttribute.ParentAttributeGroup
				attribute.IsComboAttribute = menuAttribute.IsComboAttribute

				if store.Settings.PriceSource == models.POSPriceSource {
					attribute.Price.Value = menuAttribute.Price
				}

				if index < len(product.Attributes)-countCommonAttributes {
					modifiersPrice += attribute.Price.Value * float64(attribute.Quantity)
				}
				if menuAttribute.PosID != attribute.ID && menuAttribute.PosID != "" {
					attribute.ID = menuAttribute.PosID
				}
				orderProductAttributes = append(orderProductAttributes, attribute)
				continue
			}

			menuAttributeProduct, productExistInPosMenu := productsMap[attribute.ID]
			if !productExistInPosMenu {
				req.FailReason.Code = ATTRIBUTE_MISSED_CODE
				req.FailReason.Message = ATTRIBUTE_MISSED + attribute.ID
				return req, 0, fmt.Errorf("ATTRIBUTE NOT FOUND IN POS MENU, ID %s, NAME %s", attribute.ID, attribute.Name)
			}

			orderProduct := models.OrderProduct{
				ID:       menuAttributeProduct.ProductID,
				Name:     attribute.Name,
				Quantity: attribute.Quantity * product.Quantity,
				IsCombo:  menuAttributeProduct.IsCombo,
				Price: models.Price{
					Value:        attribute.Price.Value,
					CurrencyCode: store.Settings.Currency,
				},
				ProgramID:      comboProduct.ProgramID,
				SourceActionID: comboProduct.SourceActionID,
			}

			if menuAttributeProduct.SizeID != "" {
				orderProduct.SizeId = menuAttributeProduct.SizeID
			}

			if store.Settings.PriceSource == models.POSPriceSource {
				orderProduct.Price = models.Price{
					Value:        menuAttributeProduct.Price[0].Value,
					CurrencyCode: store.Settings.Currency,
				}
			}

			if index < len(product.Attributes)-countCommonAttributes {
				modifiersPrice += orderProduct.Price.Value * float64(attribute.Quantity)
			}

			orderProducts = append(orderProducts, orderProduct)

			if store.PosType == models.Poster.String() {
				if menuProduct.CookingTime > cookingTime {
					cookingTime = menuProduct.CookingTime
				}
			}
		}

		switch store.Settings.PriceSource {
		case models.POSPriceSource:
			if _, hasPromo := promosMap[menuProduct.ProductID]; hasPromo {

			} else if _, hasPercentPromo := promoWithPercentMap[menuProduct.ProductID]; hasPercentPromo {

			} else if discountPrice, hasDiscountPrice := discountPrices[menuProduct.ProductID]; hasDiscountPrice {
				if discountPrice != product.Price.Value {
					product.Price.Value = menuProduct.Price[0].Value
				}
			} else {
				product.Price.Value = menuProduct.Price[0].Value
			}
		default:
			if req.DeliveryService == models.YANDEX.String() {
				product.Price.Value = product.Price.Value - modifiersPrice
			}
		}

		if menuProduct.ProductID != "" {
			product.ID = menuProduct.ProductID
		}

		if menuProduct.SizeID != "" {
			product.SizeId = menuProduct.SizeID
		}

		if menuProduct.MenuDefaultAttributes != nil && len(menuProduct.MenuDefaultAttributes) > 0 {
			for _, defaultAttribute := range menuProduct.MenuDefaultAttributes {
				if _, has := modifiersIds[defaultAttribute.ExtID]; has {
					continue
				}

				if defaultAttribute.DefaultAmount == 0 {
					defaultAttribute.DefaultAmount = 1
				}

				posAttribute, ok := attributesMap[defaultAttribute.ExtID]
				var attributeGroupID string

				if ok {
					attributeGroupID = menuUtils.FindAttributeGroupID(product.ID, defaultAttribute.ExtID, productsMap, attributeGroupsMap, attributesMap, store.IikoCloud.IsExternalMenu, posAttribute.ParentAttributeGroup)
				}

				orderProductAttributes = append(orderProductAttributes, models.ProductAttribute{
					ID:       defaultAttribute.ExtID,
					Name:     defaultAttribute.Name,
					Quantity: defaultAttribute.DefaultAmount,
					GroupID:  attributeGroupID,
					Price: models.Price{
						Value:        float64(defaultAttribute.Price),
						CurrencyCode: store.Settings.Currency,
					},
					IsComboAttribute: posAttribute.IsComboAttribute,
				})
			}
		}

		product.IsCombo = menuProduct.IsCombo
		product.Attributes = orderProductAttributes
		orderProducts = append(orderProducts, product)

		if store.PosType == models.Poster.String() {
			if menuProduct.CookingTime > cookingTime {
				cookingTime = menuProduct.CookingTime
			}
		}
	}

	switch req.DeliveryService {
	case models.QRMENU.String(), models.KWAAKA_ADMIN.String():
		if store.Kwaaka3PL.DeliveryPosProductId != "" && req.SendCourier {
			extraProduct, exist := productsMap[store.Kwaaka3PL.DeliveryPosProductId]
			if !exist {
				log.Info().Msgf("for order %s, delivery pos product id for 3pl order is nil", req.OrderID)
			} else {
				orderProductAttributes := make([]models.ProductAttribute, 0)
				modifiersPrice := 0.0

				if extraProduct.MenuDefaultAttributes != nil && len(extraProduct.MenuDefaultAttributes) > 0 {
					for _, defaultAttribute := range extraProduct.MenuDefaultAttributes {
						attributeID := defaultAttribute.ExtID
						menuAttribute, attributeExistInPosMenu := attributesMap[attributeID]
						if !attributeExistInPosMenu {
							log.Info().Msgf("attribute with ID %s not found in pos menu", attributeID)
							continue
						}

						orderProductAttributes = append(orderProductAttributes, models.ProductAttribute{
							ID:       attributeID,
							Name:     defaultAttribute.Name,
							Quantity: defaultAttribute.DefaultAmount,
							GroupID:  menuAttribute.ParentAttributeGroup,
							Price: models.Price{
								Value:        menuAttribute.Price,
								CurrencyCode: store.Settings.Currency,
							},
							IsComboAttribute: menuAttribute.IsComboAttribute,
						})

						modifiersPrice += float64(defaultAttribute.Price) * float64(defaultAttribute.DefaultAmount)
					}
				}

				var extraOrderProductPrice float64
				if len(extraProduct.Price) > 0 {
					extraOrderProductPrice = extraProduct.Price[0].Value + modifiersPrice
				} else {
					extraOrderProductPrice = modifiersPrice
				}
				extraOrderProduct := models.OrderProduct{
					ID:       extraProduct.ExtID,
					Name:     extraProduct.Name[0].Value,
					Quantity: 1,
					Price: models.Price{
						Value:        extraOrderProductPrice,
						CurrencyCode: store.Settings.Currency,
					},
					Attributes: orderProductAttributes,
					IsCombo:    extraProduct.IsCombo,
					ImageURLs:  extraProduct.ImageURLs,
				}

				orderProducts = append(orderProducts, extraOrderProduct)
			}
		}
	}

	req.Products = orderProducts
	if store.PosType == models.Poster.String() {
		if cookingTime == 0 {
			cookingTime = store.Poster.CookingTime
		}
	}
	req.CookingCompleteTime = time.Now().UTC().Add(time.Minute * time.Duration(cookingTime))

	return req, serviceFee, nil
}

func (s *BasePosService) isAnotherBill(storeID, deliveryService string) bool {
	for _, service := range s.anotherBillStoreIDs[storeID] {
		if service == deliveryService {
			return true
		}
	}

	return false
}
