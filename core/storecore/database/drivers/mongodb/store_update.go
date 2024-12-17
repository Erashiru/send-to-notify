package mongodb

import (
	"fmt"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *StoreRepo) setRKeeper7XMLConfig(update bson.D, rkeeper7XMLCfg *models.UpdateStoreRKeeper7XMLConfig) bson.D {
	if rkeeper7XMLCfg.Domain != nil {
		update = append(update, bson.E{
			Key:   "rkeeper7_xml.domain",
			Value: *rkeeper7XMLCfg.Domain,
		})
	}

	if rkeeper7XMLCfg.SeqNumber != nil {
		update = append(update, bson.E{
			Key:   "rkeeper7_xml.seq_number",
			Value: *rkeeper7XMLCfg.SeqNumber,
		})
	}

	return update
}

func (s *StoreRepo) setGlovoConfig(update bson.D, glovoCfg *models.UpdateStoreGlovoConfig) bson.D {
	if glovoCfg.StoreID != nil {
		update = append(update, bson.E{
			Key:   "glovo.store_id",
			Value: glovoCfg.StoreID,
		})
	}

	if glovoCfg.MenuUrl != nil {
		update = append(update, bson.E{
			Key:   "glovo.menu_url",
			Value: *glovoCfg.MenuUrl,
		})
	}

	if glovoCfg.SendToPos != nil {
		update = append(update, bson.E{
			Key:   "glovo.send_to_pos",
			Value: *glovoCfg.SendToPos,
		})
	}

	if glovoCfg.IsMarketplace != nil {
		update = append(update, bson.E{
			Key:   "glovo.is_marketplace",
			Value: *glovoCfg.IsMarketplace,
		})
	}

	if glovoCfg.IsOpen != nil {
		update = append(update, bson.E{
			Key:   "glovo.is_open",
			Value: *glovoCfg.IsOpen,
		})
	}

	if glovoCfg.AdditionalPreparationTimeInMinutes != nil {
		update = append(update, bson.E{
			Key:   "glovo.additional_preparation_time_in_minutes",
			Value: *glovoCfg.AdditionalPreparationTimeInMinutes,
		})
	}

	if glovoCfg.PaymentTypes != nil {
		update = s.setPaymentTypesFields(update, "glovo", glovoCfg.PaymentTypes)
	}

	if glovoCfg.PurchaseTypes != nil {
		update = s.setPurchaseTypes(update, "glovo", glovoCfg.PurchaseTypes)
	}

	return update
}

func (s *StoreRepo) setWoltConfig(update bson.D, woltCfg *models.UpdateStoreWoltConfig) bson.D {
	if woltCfg.StoreID != nil {
		update = append(update, bson.E{
			Key:   "wolt.store_id",
			Value: woltCfg.StoreID,
		})
	}

	if woltCfg.MenuUsername != nil {
		update = append(update, bson.E{
			Key:   "wolt.menu_username",
			Value: *woltCfg.MenuUsername,
		})
	}

	if woltCfg.MenuPassword != nil {
		update = append(update, bson.E{
			Key:   "wolt.menu_password",
			Value: *woltCfg.MenuPassword,
		})
	}

	if woltCfg.ApiKey != nil {
		update = append(update, bson.E{
			Key:   "wolt.api_key",
			Value: *woltCfg.ApiKey,
		})
	}

	if woltCfg.AdjustedPickupMinutes != nil {
		update = append(update, bson.E{
			Key:   "wolt.adjusted_pickup_minutes",
			Value: *woltCfg.AdjustedPickupMinutes,
		})
	}

	if woltCfg.MenuUrl != nil {
		update = append(update, bson.E{
			Key:   "wolt.menu_url",
			Value: *woltCfg.MenuUrl,
		})
	}

	if woltCfg.SendToPos != nil {
		update = append(update, bson.E{
			Key:   "wolt.send_to_pos",
			Value: *woltCfg.SendToPos,
		})
	}

	if woltCfg.IsMarketplace != nil {
		update = append(update, bson.E{
			Key:   "wolt.is_marketplace",
			Value: *woltCfg.IsMarketplace,
		})
	}

	if woltCfg.IsOpen != nil {
		update = append(update, bson.E{
			Key:   "wolt.is_open",
			Value: *woltCfg.IsOpen,
		})
	}
	if woltCfg.IgnoreStatusUpdate != nil {
		update = append(update, bson.E{
			Key:   "wolt.ignore_status_update",
			Value: *woltCfg.IgnoreStatusUpdate,
		})
	}

	if woltCfg.AutoAcceptOn != nil {
		update = append(update, bson.E{
			Key:   "wolt.auto_accept_on",
			Value: *woltCfg.AutoAcceptOn,
		})
	}

	if woltCfg.PaymentTypes != nil {
		update = s.setPaymentTypesFields(update, "wolt", woltCfg.PaymentTypes)
	}

	if woltCfg.PurchaseTypes != nil {
		update = s.setPurchaseTypes(update, "wolt", woltCfg.PurchaseTypes)
	}

	return update
}

func (s *StoreRepo) setYandexConfig(yandexCfg models.UpdateStoreYandexConfig) bson.D {
	var update bson.D

	if yandexCfg.StoreID != nil {
		update = append(update, bson.E{
			Key:   "external.$.store_id",
			Value: yandexCfg.StoreID,
		})
	}

	if yandexCfg.MenuUrl != nil {
		update = append(update, bson.E{
			Key:   "external.$.menu_url",
			Value: yandexCfg.MenuUrl,
		})
	}

	if yandexCfg.SendToPos != nil {
		update = append(update, bson.E{
			Key:   "external.$.send_to_pos",
			Value: yandexCfg.SendToPos,
		})
	}

	if yandexCfg.IsMarketplace != nil {
		update = append(update, bson.E{
			Key:   "external.$.is_marketplace",
			Value: yandexCfg.IsMarketplace,
		})
	}

	if yandexCfg.PaymentTypes != nil {
		update = s.setYandexPaymentTypes(update, yandexCfg.PaymentTypes)
	}

	if yandexCfg.ClientSecret != nil {
		update = append(update, bson.E{
			Key:   "external.$.client_secret",
			Value: yandexCfg.ClientSecret,
		})
	}

	result := bson.D{
		{Key: "$set", Value: update},
	}

	return result
}

func (s *StoreRepo) setQRMenuConfig(update bson.D, qrmenuCfg *models.UpdateStoreQRMenuConfig) bson.D {
	if qrmenuCfg.StoreID != nil {
		update = append(update, bson.E{
			Key:   "qr_menu.store_id",
			Value: qrmenuCfg.StoreID,
		})
	}

	if qrmenuCfg.URL != nil {
		update = append(update, bson.E{
			Key:   "qr_menu.url",
			Value: *qrmenuCfg.URL,
		})
	}

	if qrmenuCfg.IsIntegrated != nil {
		update = append(update, bson.E{
			Key:   "qr_menu.is_integrated",
			Value: *qrmenuCfg.IsIntegrated,
		})
	}

	if qrmenuCfg.PaymentTypes != nil {
		update = s.setPaymentTypesFields(update, "qr_menu", qrmenuCfg.PaymentTypes)
	}

	if qrmenuCfg.Hash != nil {
		update = append(update, bson.E{
			Key:   "qr_menu.hash",
			Value: *qrmenuCfg.Hash,
		})
	}

	if qrmenuCfg.CookingTime != nil {
		update = append(update, bson.E{
			Key:   "qr_menu.cooking_time",
			Value: *qrmenuCfg.CookingTime,
		})
	}

	if qrmenuCfg.DeliveryTime != nil {
		update = append(update, bson.E{
			Key:   "qr_menu.delivery_time",
			Value: *qrmenuCfg.DeliveryTime,
		})
	}

	if qrmenuCfg.NoTable != nil {
		update = append(update, bson.E{
			Key:   "qr_menu.no_table",
			Value: *qrmenuCfg.NoTable,
		})
	}

	if qrmenuCfg.Theme != nil {
		update = append(update, bson.E{
			Key:   "qr_menu.theme",
			Value: *qrmenuCfg.Theme,
		})
	}

	if qrmenuCfg.IsMarketplace != nil {
		update = append(update, bson.E{
			Key:   "qr_menu.is_marketplace",
			Value: *qrmenuCfg.IsMarketplace,
		})
	}

	if qrmenuCfg.SendToPos != nil {
		update = append(update, bson.E{
			Key:   "qr_menu.send_to_pos",
			Value: *qrmenuCfg.SendToPos,
		})
	}

	if qrmenuCfg.IgnoreStatusUpdate != nil {
		update = append(update, bson.E{
			Key:   "qr_menu.ignore_status_update",
			Value: *qrmenuCfg.IgnoreStatusUpdate,
		})
	}

	if qrmenuCfg.BusyMode != nil {
		update = append(update, bson.E{
			Key:   "qr_menu.busy_mode",
			Value: *qrmenuCfg.BusyMode,
		})
	}

	if qrmenuCfg.AdjustedPickupMinutes != nil {
		update = append(update, bson.E{
			Key:   "qr_menu.adjusted_pickup_minutes",
			Value: *qrmenuCfg.AdjustedPickupMinutes,
		})
	}

	return update
}

func (s *StoreRepo) setKwaakaAdminConfig(update bson.D, kwaakaAdminCfg *models.UpdateStoreKwaakaAdminConfig) bson.D {

	if kwaakaAdminCfg.StoreID != nil {
		update = append(update, bson.E{
			Key:   "kwaaka_admin.store_id",
			Value: kwaakaAdminCfg.StoreID,
		})
	}

	if kwaakaAdminCfg.IsIntegrated != nil {
		update = append(update, bson.E{
			Key:   "kwaaka_admin.is_integrated",
			Value: *kwaakaAdminCfg.IsIntegrated,
		})
	}

	if kwaakaAdminCfg.SendToPos != nil {
		update = append(update, bson.E{
			Key:   "kwaaka_admin.send_to_pos",
			Value: *kwaakaAdminCfg.SendToPos,
		})
	}

	if kwaakaAdminCfg.CookingTime != nil {
		update = append(update, bson.E{
			Key:   "kwaaka_admin.cooking_time",
			Value: *kwaakaAdminCfg.CookingTime,
		})
	}

	if kwaakaAdminCfg.IsActive != nil {
		update = append(update, bson.E{
			Key:   "kwaaka_admin.is_active",
			Value: *kwaakaAdminCfg.IsActive,
		})
	}

	return update
}

func (s *StoreRepo) setPurchaseTypes(update bson.D, delivery string, purchaseTypes *models.UpdatePurchaseTypes) bson.D {
	if purchaseTypes.Instant != nil {
		update = append(update, bson.E{
			Key:   fmt.Sprintf("%s.purchase_types.instant", delivery),
			Value: purchaseTypes.Instant,
		})
	}

	if purchaseTypes.Preorder != nil {
		update = append(update, bson.E{
			Key:   fmt.Sprintf("%s.purchase_types.preorder", delivery),
			Value: purchaseTypes.Preorder,
		})
	}

	if purchaseTypes.TakeAway != nil {
		update = append(update, bson.E{
			Key:   fmt.Sprintf("%s.purchase_types.takeaway", delivery),
			Value: purchaseTypes.TakeAway,
		})
	}

	return update
}

// setPaymentTypesFields - function to set payment_types fields, receive update bson.D, delivery (glovo, wolt, etc...). It will be nameof main field in DB.
func (s *StoreRepo) setPaymentTypesFields(update bson.D, delivery string, paymentTypes *models.UpdateDeliveryServicePaymentType) bson.D {

	if paymentTypes.CASH != nil {
		if paymentTypes.CASH.IikoPaymentTypeID != nil {
			update = append(update, bson.E{
				Key:   fmt.Sprintf("%s.payment_types.CASH.iiko_payment_type_id", delivery),
				Value: *paymentTypes.CASH.IikoPaymentTypeID,
			})
		}

		if paymentTypes.CASH.IikoPaymentTypeKind != nil {
			update = append(update, bson.E{
				Key:   fmt.Sprintf("%s.payment_types.CASH.iiko_payment_type_kind", delivery),
				Value: *paymentTypes.CASH.IikoPaymentTypeKind,
			})
		}

		if paymentTypes.CASH.OrderType != nil {
			update = append(update, bson.E{
				Key:   fmt.Sprintf("%s.payment_types.CASH.order_type", delivery),
				Value: *paymentTypes.CASH.OrderType,
			})
		}
	}

	if paymentTypes.DELAYED.IikoPaymentTypeID != nil {
		update = append(update, bson.E{
			Key:   fmt.Sprintf("%s.payment_types.DELAYED.iiko_payment_type_id", delivery),
			Value: *paymentTypes.DELAYED.IikoPaymentTypeID,
		})
	}

	if paymentTypes.DELAYED.IikoPaymentTypeKind != nil {
		update = append(update, bson.E{
			Key:   fmt.Sprintf("%s.payment_types.DELAYED.iiko_payment_type_kind", delivery),
			Value: *paymentTypes.DELAYED.IikoPaymentTypeKind,
		})
	}

	if paymentTypes.DELAYED.OrderType != nil {
		update = append(update, bson.E{
			Key:   fmt.Sprintf("%s.payment_types.DELAYED.order_type", delivery),
			Value: *paymentTypes.DELAYED.OrderType,
		})
	}

	return update
}

func (s *StoreRepo) setYandexPaymentTypes(update bson.D, yandexPaymentTypes *models.UpdateDeliveryServicePaymentType) bson.D {
	if yandexPaymentTypes.CASH != nil {
		if yandexPaymentTypes.CASH.IikoPaymentTypeID != nil {
			update = append(update, bson.E{
				Key:   "external.$.payment_types.CASH.iiko_payment_type_id",
				Value: yandexPaymentTypes.CASH.IikoPaymentTypeID,
			})
		}
		if yandexPaymentTypes.CASH.IikoPaymentTypeKind != nil {
			update = append(update, bson.E{
				Key:   "external.$.payment_types.CASH.iiko_payment_type_kind",
				Value: yandexPaymentTypes.CASH.IikoPaymentTypeKind,
			})
		}
		if yandexPaymentTypes.CASH.OrderType != nil {
			update = append(update, bson.E{
				Key:   "external.$.payment_types.CASH.order_type",
				Value: yandexPaymentTypes.CASH.OrderType,
			})
		}
	}

	if yandexPaymentTypes.DELAYED != nil {
		if yandexPaymentTypes.DELAYED.IikoPaymentTypeID != nil {
			update = append(update, bson.E{
				Key:   "external.$.payment_types.DELAYED.iiko_payment_type_id",
				Value: yandexPaymentTypes.DELAYED.IikoPaymentTypeID,
			})
		}
		if yandexPaymentTypes.DELAYED.IikoPaymentTypeKind != nil {
			update = append(update, bson.E{
				Key:   "external.$.payment_types.DELAYED.iiko_payment_type_kind",
				Value: yandexPaymentTypes.DELAYED.IikoPaymentTypeKind,
			})
		}
		if yandexPaymentTypes.DELAYED.OrderType != nil {
			update = append(update, bson.E{
				Key:   "external.$.payment_types.DELAYED.order_type",
				Value: yandexPaymentTypes.DELAYED.OrderType,
			})
		}
	}
	return update
}

func (s *StoreRepo) setIIKOConfig(update bson.D, iikoConfigs *models.UpdateStoreIikoConfig) bson.D {
	if iikoConfigs.Key != nil {
		update = append(update, bson.E{
			Key:   "iiko_cloud.key",
			Value: iikoConfigs.Key,
		})
	}

	if iikoConfigs.TerminalID != nil {
		update = append(update, bson.E{
			Key:   "iiko_cloud.terminal_id",
			Value: iikoConfigs.TerminalID,
		})
	}

	if iikoConfigs.OrganizationID != nil {
		update = append(update, bson.E{
			Key:   "iiko_cloud.organization_id",
			Value: iikoConfigs.OrganizationID,
		})
	}

	if iikoConfigs.IsExternalMenu != nil {
		update = append(update, bson.E{
			Key:   "iiko_cloud.is_external_menu",
			Value: iikoConfigs.IsExternalMenu,
		})
	}

	if iikoConfigs.ExternalMenuID != nil {
		update = append(update, bson.E{
			Key:   "iiko_cloud.external_menu_id",
			Value: iikoConfigs.ExternalMenuID,
		})
	}

	if iikoConfigs.PriceCategory != nil {
		update = append(update, bson.E{
			Key:   "iiko_cloud.price_category",
			Value: iikoConfigs.PriceCategory,
		})
	}

	if iikoConfigs.StopListByBalance != nil {
		update = append(update, bson.E{
			Key:   "iiko_cloud.stoplist_by_balance",
			Value: iikoConfigs.StopListByBalance,
		})
	}

	if iikoConfigs.StopListBalanceLimit != nil {
		update = append(update, bson.E{
			Key:   "iiko_cloud.stoplist_balance_limit",
			Value: iikoConfigs.StopListBalanceLimit,
		})
	}

	return update
}
