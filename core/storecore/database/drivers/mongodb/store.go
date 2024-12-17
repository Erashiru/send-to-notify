package mongodb

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/models"
	drivers2 "github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	models2 "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/pkg/menu/dto"
	models3 "github.com/kwaaka-team/orders-core/service/kwaaka_3pl/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ drivers2.StoreRepository = &StoreRepo{}

type StoreRepo struct {
	collection *mongo.Collection
}

func NewStoreRepository(collection *mongo.Collection) drivers2.StoreRepository {
	return &StoreRepo{
		collection: collection,
	}
}

func (s *StoreRepo) UpdateStoreByFields(ctx context.Context, store models2.UpdateStore) error {
	if store.ID == nil {
		return errors.Wrap(drivers2.ErrInvalid, "storeID is nil")
	}

	oid, err := primitive.ObjectIDFromHex(*store.ID)
	if err != nil {
		return errors.Wrap(drivers2.ErrInvalid, "storeID error")
	}

	filter := bson.D{{Key: "_id", Value: oid}}

	update, err := s.setFields(store)
	if err != nil {
		return err
	}

	res, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errorSwitch(err)
	}

	if res.MatchedCount == 0 {
		return errors.Wrap(drivers2.ErrNotFound, "not found error")
	}

	return nil
}

//
//func (s *StoreRepo) Update(ctx context.Context, store models2.Store) error {
//	oid, err := primitive.ObjectIDFromHex(store.ID)
//	if err != nil {
//		return errors.Wrap(drivers2.ErrInvalid, "storeID error")
//	}
//
//	filter := bson.D{{Key: "_id", Value: oid}}
//
//	// be careful, it's important thing; without this it won't work
//	store.ID = ""
//
//	update := bson.D{{Key: "$set", Value: store}}
//
//	res, err := s.collection.UpdateOne(ctx, filter, update)
//	if err != nil {
//		return errorSwitch(err)
//	}
//
//	if res.MatchedCount == 0 {
//		return errors.Wrap(drivers2.ErrNotFound, "not found error")
//	}
//
//	return nil
//}

func (s *StoreRepo) Get(ctx context.Context, query selector.Store) (models2.Store, error) {
	filter, err := s.filterFrom(query)
	if err != nil {
		return models2.Store{}, err
	}

	var store models2.Store
	if err = s.collection.FindOne(ctx, filter).Decode(&store); err != nil {
		return models2.Store{}, errorSwitch(err)
	}

	return store, nil
}

func (s *StoreRepo) Create(ctx context.Context, store models2.Store) (string, error) {
	store.IntegrationDate = time.Now().UTC()
	store.CreatedAt = time.Now().UTC()
	store.UpdatedAt = time.Now().UTC()

	res, err := s.collection.InsertOne(ctx, store)
	if err != nil {
		return "", errorSwitch(err)
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *StoreRepo) List(ctx context.Context, query selector.Store) ([]models2.Store, error) {
	filter, err := s.filterFrom(query)
	if err != nil {
		return nil, err
	}

	cur, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, errorSwitch(err)
	}

	var stores []models2.Store
	for cur.Next(ctx) {
		var store models2.Store
		if err := cur.Decode(&store); err != nil {
			return nil, errorSwitch(err)
		}

		stores = append(stores, store)
	}

	return stores, nil
}

func (s *StoreRepo) FindCallCenterStores(ctx context.Context) ([]models2.CallCenterRestaurant, error) {
	var restaurants []models2.CallCenterRestaurant

	filter := bson.M{"kwaaka_admin.is_integrated": true}
	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, errorSwitch(err)
	}
	for cursor.Next(ctx) {
		var restaurant models2.CallCenterRestaurant
		if err := cursor.Decode(&restaurant); err != nil {
			return nil, err
		}
		restaurant.IsIntegrated = true
		restaurants = append(restaurants, restaurant)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return restaurants, nil
}

func (s *StoreRepo) FindDirectStores(ctx context.Context) ([]models2.DirectRestaurant, error) {
	var restaurants []models2.DirectRestaurant

	filter := bson.M{"qr_menu.is_marketplace": true}
	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, errorSwitch(err)
	}
	for cursor.Next(ctx) {
		var restaurant models2.DirectRestaurant
		if err := cursor.Decode(&restaurant); err != nil {
			return nil, err
		}
		restaurant.QRMenuIsMarketplace = true
		restaurants = append(restaurants, restaurant)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return restaurants, nil
}

func (s *StoreRepo) DeleteStore(ctx context.Context, storeId string) error {
	oid, err := primitive.ObjectIDFromHex(storeId)
	if err != nil {
		return errors.Wrap(drivers2.ErrNotFound, "not found error")
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	_, err = s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (s *StoreRepo) UpdateYandexConfig(ctx context.Context, storeID string, yandexConfig models2.UpdateStoreYandexConfig) error {
	oid, err := primitive.ObjectIDFromHex(storeID)
	if err != nil {
		return errorSwitch(err)
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "external.type", Value: models.YANDEX.String()},
	}

	update := s.setYandexConfig(yandexConfig)

	res, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errorSwitch(err)
	}

	if res.MatchedCount == 0 {
		return errors.Wrap(drivers2.ErrNotFound, "not found error")
	}

	return nil
}

func (s *StoreRepo) CreateYandexConfig(ctx context.Context, storeID string, yandexConfig models2.YandexConfig) error {
	oid, err := primitive.ObjectIDFromHex(storeID)
	if err != nil {
		return errorSwitch(err)
	}

	filterExternalArray := bson.M{"_id": oid, "external": bson.M{"$type": "array"}}
	updateExternalArray := bson.M{"$push": bson.M{"external": yandexConfig}}

	externalArrayRes, err := s.collection.UpdateOne(ctx, filterExternalArray, updateExternalArray)

	if err != nil {
		return err
	}

	if externalArrayRes.ModifiedCount == 0 {
		filterExternalNotArray := bson.M{"_id": oid}
		updateExternalNotArray := bson.M{"$set": bson.M{"external": []models2.YandexConfig{yandexConfig}}}

		externalNotArrayRes, err := s.collection.UpdateOne(ctx, filterExternalNotArray, updateExternalNotArray)
		if err != nil {
			return err
		}

		if externalNotArrayRes.ModifiedCount == 0 {
			return errors.Wrap(drivers2.ErrNotFound, "not found error")
		}
	}

	return nil
}

func (s *StoreRepo) setFields(store models2.UpdateStore) (bson.D, error) {
	update := make(bson.D, 0, 3)

	if store.Token != nil {
		update = append(update, bson.E{
			Key:   "token",
			Value: *store.Token,
		})
	}

	if store.Name != nil {
		update = append(update, bson.E{
			Key:   "name",
			Value: *store.Name,
		})
	}

	if store.MenuID != nil {
		oid, err := primitive.ObjectIDFromHex(*store.MenuID)
		if err != nil {
			return nil, errors.Wrap(drivers2.ErrInvalid, "store.MenuID error")
		}

		update = append(update, bson.E{
			Key:   "menu_id",
			Value: oid,
		})
	}

	if store.PosType != nil {
		update = append(update, bson.E{
			Key:   "pos_type",
			Value: *store.PosType,
		})
	}

	if store.IntegrationDate != nil {
		update = append(update, bson.E{
			Key:   "integration_date",
			Value: *store.IntegrationDate,
		})
	}

	if store.UpdatedAt != nil {
		update = append(update, bson.E{
			Key:   "updated_at",
			Value: *store.UpdatedAt,
		})
	}

	if store.CreatedAt != nil {
		update = append(update, bson.E{
			Key:   "created_at",
			Value: *store.CreatedAt,
		})
	}

	if store.RestaurantGroupID != nil {
		update = append(update, bson.E{
			Key:   "restaurant_group_id",
			Value: *store.RestaurantGroupID,
		})
	}

	if store.LegalEntityId != nil {
		update = append(update, bson.E{
			Key:   "legal_entity_id",
			Value: *store.LegalEntityId,
		})
	}

	if store.SalesManagerId != nil {
		update = append(update, bson.E{
			Key:   "sales_manager_id",
			Value: *store.SalesManagerId,
		})
	}

	if store.AccountManagerId != nil {
		update = append(update, bson.E{
			Key:   "account_manager_id",
			Value: *store.AccountManagerId,
		})
	}

	if store.Address != nil {
		if store.Address.Street != nil {
			update = append(update, bson.E{
				Key:   "address.street",
				Value: *store.Address.Street,
			})
		}

		if store.Address.City != nil {
			update = append(update, bson.E{
				Key:   "address.city",
				Value: *store.Address.City,
			})
		}

		if store.Address.UpdateCoordinates != nil {
			update = append(update, bson.E{
				Key:   "address.coordinates",
				Value: *store.Address.UpdateCoordinates,
			})
		}
		if store.Address.Entrance != nil {
			update = append(update, bson.E{
				Key:   "address.entrance",
				Value: *store.Address.Entrance,
			})
		}
	}

	if store.Glovo != nil {
		update = s.setGlovoConfig(update, store.Glovo)
	}

	if store.Wolt != nil {
		update = s.setWoltConfig(update, store.Wolt)
	}

	if store.External != nil {
		update = append(update, bson.E{
			Key:   "external",
			Value: store.External,
		})
	}

	if store.QRMenu != nil {
		update = s.setQRMenuConfig(update, store.QRMenu)
	}

	if store.KwaakaAdmin != nil {
		update = s.setKwaakaAdminConfig(update, store.KwaakaAdmin)
	}

	if store.IikoCloud != nil {
		update = s.setIIKOConfig(update, store.IikoCloud)
	}

	if store.Paloma != nil {
		update = append(update, bson.E{
			Key:   "paloma",
			Value: store.Paloma,
		})
	}

	if store.RKeeper7XML != nil {
		update = s.setRKeeper7XMLConfig(update, store.RKeeper7XML)
	}

	if store.Menus != nil {
		update = append(update, bson.E{
			Key:   "menus",
			Value: store.Menus,
		})
	}

	if store.Links != nil {
		update = append(update, bson.E{
			Key:   "links",
			Value: store.Links,
		})
	}

	if store.SocialMediaLinks != nil {
		update = append(update, bson.E{
			Key:   "social_media_links",
			Value: store.SocialMediaLinks,
		})
	}

	if store.Contacts != nil {
		update = append(update, bson.E{
			Key:   "contacts",
			Value: store.Contacts,
		})
	}

	if store.BillParameter != nil {
		if store.BillParameter.IsActive != nil {
			update = append(update, bson.E{
				Key:   "bill_parameters.is_active",
				Value: store.BillParameter.IsActive,
			})
		}
		if store.BillParameter.UpdateBillParameters.AddAddress != nil {
			update = append(update, bson.E{
				Key:   "bill_parameters.parameters.add_address",
				Value: store.BillParameter.UpdateBillParameters.AddAddress,
			})
		}
		if store.BillParameter.UpdateBillParameters.AddComments != nil {
			update = append(update, bson.E{
				Key:   "bill_parameters.parameters.add_comments",
				Value: store.BillParameter.UpdateBillParameters.AddComments,
			})
		}
		if store.BillParameter.UpdateBillParameters.AddDelivery != nil {
			update = append(update, bson.E{
				Key:   "bill_parameters.parameters.add_delivery",
				Value: store.BillParameter.UpdateBillParameters.AddDelivery,
			})
		}
		if store.BillParameter.UpdateBillParameters.AddOrderCode != nil {
			update = append(update, bson.E{
				Key:   "bill_parameters.parameters.add_orded_code",
				Value: store.BillParameter.UpdateBillParameters.AddOrderCode,
			})
		}
		if store.BillParameter.UpdateBillParameters.AddPaymentType != nil {
			update = append(update, bson.E{
				Key:   "bill_parameters.parameters.add_payment_type",
				Value: store.BillParameter.UpdateBillParameters.AddPaymentType,
			})
		}
		if store.BillParameter.UpdateBillParameters.AddQuantityPersons != nil {
			update = append(update, bson.E{
				Key:   "bill_parameters.parameters.add_quantity_persons",
				Value: store.BillParameter.UpdateBillParameters.AddQuantityPersons,
			})
		}

	}

	if store.Settings != nil {
		if store.Settings.HasVirtualStore != nil {
			update = append(update, bson.E{
				Key:   "settings.has_virtual_store",
				Value: store.Settings.HasVirtualStore,
			})
		}

		if store.Settings.StopListClosingActions != nil {
			update = append(update, bson.E{
				Key:   "settings.stoplist_closing_actions",
				Value: store.Settings.StopListClosingActions,
			})
		}
		if store.Settings.StopListClosingActions != nil {
			update = append(update, bson.E{
				Key:   "settings.stoplist_closing_actions",
				Value: store.Settings.StopListClosingActions,
			})
		}
		if store.Settings.Currency != nil {
			update = append(update, bson.E{
				Key:   "settings.currency",
				Value: store.Settings.Currency,
			})
		}
		if store.Settings.LanguageCode != nil {
			update = append(update, bson.E{
				Key:   "settings.language_code",
				Value: store.Settings.LanguageCode,
			})
		}
	}

	if store.StoreSchedule != nil {
		if store.StoreSchedule.GlovoSchedule != nil {
			update = append(update, bson.E{
				Key:   "store_schedule.glovo_schedule",
				Value: store.StoreSchedule.GlovoSchedule,
			})
		}
		if store.StoreSchedule.WoltSchedule != nil {
			update = append(update, bson.E{
				Key:   "store_schedule.wolt_schedule",
				Value: store.StoreSchedule.WoltSchedule,
			})
		}

		if store.StoreSchedule.DirectSchedule != nil {
			update = append(update, bson.E{
				Key:   "store_schedule.direct_schedule",
				Value: store.StoreSchedule.DirectSchedule,
			})
		}

	}

	if store.IsDeleted != nil {
		update = append(update, bson.E{
			Key:   "is_deleted",
			Value: *store.IsDeleted,
		})
	}

	if store.CompensationCount != nil {
		if *store.CompensationCount != 0 {
			update = append(update, bson.E{
				Key:   "compensation_count",
				Value: store.CompensationCount,
			})
		}
	}

	result := bson.D{
		{Key: "$set", Value: update},
	}

	return result, nil
}

func (s *StoreRepo) filterFrom(query selector.Store) (bson.D, error) {
	var result bson.D

	if query.HasName() {
		result = append(result, bson.E{
			Key:   "name",
			Value: query.Name,
		})
	}

	if query.HasVirtualStore() {
		result = append(result, bson.E{
			Key:   "settings.has_virtual_store",
			Value: *query.IsVirtualStore,
		})
	}

	if query.HasTimezone() {
		result = append(result, bson.E{
			Key:   "settings.timezone.tz",
			Value: query.Timezone,
		})
	}

	if query.HasUtcOffset() {
		result = append(result, bson.E{
			Key:   "settings.timezone.utc_offset",
			Value: query.UtcOffset,
		})
	}

	if query.HasCurrency() {
		result = append(result, bson.E{
			Key:   "settings.currency",
			Value: query.Currency,
		})
	}

	if query.HasStoreGroupId() {
		result = append(result, bson.E{
			Key:   "restaurant_group_id",
			Value: query.StoreGroupId,
		})
	}

	if query.HasLanguageCode() {
		result = append(result, bson.E{
			Key:   "settings.language_code",
			Value: query.LanguageCode,
		})
	}

	if query.HasStreet() {
		result = append(result, bson.E{
			Key:   "address.street",
			Value: query.Street,
		})
	}

	if query.HasCity() {
		result = append(result, bson.E{
			Key:   "address.city",
			Value: query.City,
		})
	}

	if query.HasID() {
		oid, err := primitive.ObjectIDFromHex(query.ID)
		if err != nil {
			return nil, errors.Wrap(drivers2.ErrInvalid, "query.ID error")
		}

		result = append(result, bson.E{
			Key:   "_id",
			Value: oid,
		})
	}

	if query.HasToken() {
		result = append(result, bson.E{
			Key:   "token",
			Value: query.Token,
		})
	}

	if query.HasAggregatorMenuID() {
		oIDs := make([]primitive.ObjectID, 0, len(query.AggregatorMenuIDs))
		for _, menuID := range query.AggregatorMenuIDs {
			oid, err := primitive.ObjectIDFromHex(menuID)
			if err != nil {
				continue
			}
			oIDs = append(oIDs, oid)
		}

		if query.HasIsActiveMenu() {
			result = append(result, bson.E{
				Key: "menus", Value: bson.D{
					{Key: "$elemMatch", Value: bson.D{
						{Key: "menu_id", Value: bson.D{
							{Key: "$in", Value: oIDs},
						}},
						{Key: "is_active", Value: query.ActiveMenu()},
					}},
				},
			})
		} else {
			result = append(result, bson.E{
				Key: "menus", Value: bson.D{
					{Key: "$elemMatch", Value: bson.D{
						{Key: "menu_id", Value: bson.D{
							{Key: "$in", Value: oIDs},
						}},
					}},
				},
			})
		}
	}

	if query.HasClientSecret() {
		result = append(result, bson.E{
			Key:   "external.client_secret",
			Value: query.ClientSecret,
		})
	}

	if query.HasPosType() {
		result = append(result, bson.E{
			Key:   "pos_type",
			Value: query.PosType,
		})
	}

	if query.HasHash() {
		result = append(result, bson.E{
			Key:   "hash",
			Value: query.Hash,
		})
	}

	if query.HasDeliveryService() && query.HasExternalStoreID() {
		result = append(result, bson.E{
			Key:   fmt.Sprintf("%s.store_id", query.DeliveryService),
			Value: query.ExternalStoreID,
		})
	}

	if query.HasExternalDeliveryService() && query.HasExternalStoreID() {
		result = append(result, bson.E{
			Key:   "external.type",
			Value: query.ExternalDeliveryService,
		})

		result = append(result, bson.E{
			Key:   "external.store_id",
			Value: query.ExternalStoreID,
		})
	}

	if query.HasPosOrganizationID() {
		result = append(result, bson.E{
			Key:   "iiko_cloud.organization_id",
			Value: query.PosOrganizationID,
		})
	}

	if query.HasStoreIDs() {
		objIDs, err := toObjectIDs(query.IDs)
		if err != nil {
			return nil, err
		}

		result = append(result, bson.E{
			Key: "_id", Value: bson.D{
				{Key: "$in", Value: objIDs},
			}})
	}

	if query.HasExpress24StoreId() {
		result = append(result, bson.E{
			Key:   "express24.store_id",
			Value: query.Express24StoreId,
		})
	}

	if query.HasAccountNumber() {
		result = append(result, bson.E{
			Key:   "poster.account_number_string",
			Value: query.PosterAccountNumber,
		})
	}

	if query.HasTalabatRemoreBranchId() {
		result = append(result, bson.E{
			Key:   "talabat.remote_branch_id",
			Value: query.TalabatRemoteBranchId,
		})
	}
	if query.HasDeliveryService() {
		if query.DeliveryService == models2.GLOVO.String() {
			result = append(result, bson.E{
				Key:   "glovo.store_id",
				Value: bson.D{{Key: "$ne", Value: nil}},
			})
		}
	}
	if query.HasScheduledStatusChange() {
		result = append(result, bson.E{
			Key:   "settings.scheduled_status_change.is_active",
			Value: query.ScheduledStatusChange,
		})
	}
	if query.HasYarosStoreId() {
		result = append(result, bson.E{
			Key:   "yaros.store_id",
			Value: query.YarosStoreId,
		})
	}
	if query.HasIsChildStore() {
		result = append(result, bson.E{
			Key: "settings.is_child_store",
			Value: bson.D{
				{Key: "$exists", Value: true},
				{Key: "$eq", Value: *query.IsChildStore},
			},
		})
	}
	if query.HasOrderAutoClose() {
		result = append(result, bson.E{
			Key:   "order_auto_close_settings.order_auto_close",
			Value: query.OrderAutoClose,
		})
	}

	return result, nil
}

func (s *StoreRepo) Get3plRestaurantStatus(ctx context.Context, storeID string) (bool, error) {

	oid, err := primitive.ObjectIDFromHex(storeID)
	if err != nil {
		return false, nil
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	var store models2.Store
	if err := s.collection.FindOne(ctx, filter).Decode(&store); err != nil {
		return false, err
	}

	return store.Kwaaka3PL.Is3pl, nil
}

func (s *StoreRepo) Update3plRestaurantStatus(ctx context.Context, query models2.Update3plRestaurantStatus, indriveStoreID string) error {

	oid, err := primitive.ObjectIDFromHex(*query.RestaurantId)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	result, err := s.collection.UpdateOne(ctx, filter, s.set3plConfig(query.Is3pl, indriveStoreID))
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.Wrap(drivers2.ErrNotFound, "not found error while update 3pl restaurants status")
	}

	return nil
}

func (s *StoreRepo) UpdateRestaurantPolygons(ctx context.Context, restaurantID string, polygons []models2.Polygon) error {

	oid, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	result, err := s.collection.UpdateOne(ctx, filter, s.set3PlPolygons(polygons))
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.Wrap(drivers2.ErrNotFound, "not found error while update restaurants 3pl polygons")
	}

	return nil
}

func (s *StoreRepo) UpdateDynamicPolygon(ctx context.Context, restaurantID string, isDynamic bool, cpo float64) error {
	oid, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	result, err := s.collection.UpdateOne(ctx, filter, s.set3PlDynamicPolygon(isDynamic, cpo))
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.Wrap(drivers2.ErrNotFound, "not found error while update restaurants 3pl dynamic polygons")
	}

	return nil
}

func (s *StoreRepo) UpdateDispatchDeliveryStatus(ctx context.Context, query models2.UpdateDispatchDeliveryAvailable) error {

	oid, err := primitive.ObjectIDFromHex(query.RestaurantID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	result, err := s.collection.UpdateOne(ctx, filter, s.setDispatchStatus(query.DeliveryService, query.Available, query.Is3pl))
	if err != nil {
		return nil
	}

	if result.MatchedCount == 0 {
		return errors.Wrap(drivers2.ErrNotFound, "not found error while update 3pl delivery status")
	}

	return nil
}

func (s *StoreRepo) setDispatchStatus(deliveryService string, status bool, threePLStatus bool) bson.D {

	var service string
	switch deliveryService {
	case models3.IndriveDelivery:
		service = "kwaaka_3pl.indrive_available"
	case dto.Yandex.String():
		service = "kwaaka_3pl.yandex_available"
	case models3.WoltDelivery:
		service = "kwaaka_3pl.wolt_drive_available"
	}

	update := bson.D{
		{Key: service, Value: status},
		//{Key: "kwaaka_3pl.is_3pl", Value: threePLStatus},
	}

	result := bson.D{
		{Key: "$set", Value: update},
	}
	return result
}

func (s *StoreRepo) AppendMenuToStoreMenus(ctx context.Context, storeId string, menu models2.StoreDSMenu) error {
	oid, err := primitive.ObjectIDFromHex(storeId)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": oid}
	update := bson.M{"$push": bson.M{"menus": menu}}

	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *StoreRepo) set3PlDynamicPolygon(isDynamic bool, cpo float64) bson.D {
	update := bson.D{
		{Key: "kwaaka_3pl.is_dynamic", Value: isDynamic},
		{Key: "kwaaka_3pl.cpo", Value: cpo},
	}

	result := bson.D{
		{Key: "$set", Value: update},
	}

	return result
}

func (s *StoreRepo) set3PlPolygons(polygons []models2.Polygon) bson.D {
	update := bson.D{
		{Key: "kwaaka_3pl.polygons", Value: polygons},
	}

	result := bson.D{
		{Key: "$set", Value: update},
	}

	return result
}

func (s *StoreRepo) set3plConfig(status *bool, indriveStoreID string) bson.D {

	update := bson.D{
		{Key: "kwaaka_3pl.is_3pl", Value: status},
		//{Key: "kwaaka_3pl.indrive_available", Value: status},
		//{Key: "kwaaka_3pl.wolt_drive_available", Value: status},
		//{Key: "kwaaka_3pl.yandex_available", Value: status},
		//{Key: "kwaaka_3pl.kwaaka_charge_absolute", Value: 20},
		//{Key: "kwaaka_3pl.kwaaka_charge_percentage", Value: 5},
	}

	//if indriveStoreID != "" {
	//	update = append(update, bson.E{
	//		Key: "kwaaka_3pl.indrive_store_id", Value: indriveStoreID,
	//	})
	//}

	result := bson.D{
		{Key: "$set", Value: update},
	}
	return result
}

func (s *StoreRepo) UpdateWoltBusyMode(ctx context.Context, storeID string, busyMode bool, busyModeTime int) error {

	filter := bson.D{
		{Key: "wolt.store_id", Value: storeID},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "wolt.busy_mode", Value: busyMode},
			{Key: "wolt.adjusted_pickup_minutes", Value: busyModeTime},
		}},
	}

	res, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return drivers2.ErrNotFound
	}
	return nil
}

func (s *StoreRepo) UpdateDirectBusyMode(ctx context.Context, storeID string, busyMode bool, busyModeTime int) error {
	filter := bson.D{
		{Key: "qr_menu.store_id", Value: storeID},
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "qr_menu.busy_mode", Value: busyMode},
			{Key: "qr_menu.adjusted_pickup_minutes", Value: busyModeTime},
		}},
	}

	res, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return drivers2.ErrNotFound
	}
	return nil
}

func (s *StoreRepo) UpdateCookingTimeWolt(ctx context.Context, restaurantID string, cookingTime int) error {

	oid, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "wolt.cooking_time", Value: cookingTime},
		}},
	}

	res, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return drivers2.ErrNotFound
	}

	return nil
}

func (s *StoreRepo) GetStoresIDsAndNamesByGroupId(ctx context.Context, groupID string) ([]models2.StoreIdAndName, error) {
	filter := bson.D{
		{
			Key:   "restaurant_group_id",
			Value: groupID,
		},
	}

	project := options.Find().SetProjection(map[string]int{"_id": 1, "name": 1})

	curs, err := s.collection.Find(ctx, filter, project)
	if err != nil {
		return nil, err
	}

	res := make([]models2.StoreIdAndName, 0, curs.RemainingBatchLength())

	if err = curs.All(ctx, &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (s *StoreRepo) SetTwoGisLink(ctx context.Context, twoGisLink, restID string) error {
	objID, err := primitive.ObjectIDFromHex(restID)
	if err != nil {
		return err
	}

	var result struct {
		ExternalLinks []struct {
			Name string `bson:"name"`
			URL  string `bson:"url"`
		} `bson:"external_links"`
	}
	err = s.collection.FindOne(ctx, bson.D{{Key: "_id", Value: objID}}).Decode(&result)
	if err != nil {
		return err
	}

	if result.ExternalLinks == nil {
		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "external_links", Value: bson.A{}},
			}},
		}
		_, err := s.collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: objID}}, update)
		if err != nil {
			return err
		}
	}

	for _, link := range result.ExternalLinks {
		if link.Name == "2gis" {
			if link.URL == twoGisLink {
				return nil
			} else {
				filter := bson.D{
					{Key: "_id", Value: objID},
					{Key: "external_links.name", Value: "2gis"},
				}
				update := bson.D{
					{Key: "$set", Value: bson.D{
						{Key: "external_links.$.url", Value: twoGisLink},
					}},
				}
				_, err := s.collection.UpdateOne(ctx, filter, update)
				if err != nil {
					return err
				}
				return nil
			}
		}
	}
	pushUpdate := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "external_links", Value: bson.D{
				{Key: "name", Value: "2gis"},
				{Key: "url", Value: twoGisLink},
			}},
		}},
	}
	_, err = s.collection.UpdateOne(ctx, bson.D{{Key: "_id", Value: objID}}, pushUpdate)
	if err != nil {
		return err
	}
	return nil
}

func (s *StoreRepo) UpdateRestaurantCharge(ctx context.Context, req models2.UpdateRestaurantCharge, restID string) error {
	oid, err := primitive.ObjectIDFromHex(restID)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	update := s.updateFromRestaurantCharge(req)

	updateResult, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if updateResult.ModifiedCount == 0 {
		return drivers2.ErrNotFound
	}

	return nil
}

func (s *StoreRepo) updateFromRestaurantCharge(req models2.UpdateRestaurantCharge) bson.D {
	update := make(bson.D, 0, 3)

	if req.IsRestaurantChargeOn != nil {
		update = append(update, bson.E{Key: "restaurant_charge.is_restaurant_charge_on", Value: *req.IsRestaurantChargeOn})
	}

	if req.MinRestaurantCharge != nil {
		update = append(update, bson.E{Key: "restaurant_charge.min_restaurant_charge", Value: *req.MinRestaurantCharge})
	}

	if req.MaxRestaurantCharge != nil {
		update = append(update, bson.E{Key: "restaurant_charge.max_restaurant_charge", Value: *req.MaxRestaurantCharge})
	}

	if req.RestaurantChargePercent != nil {
		update = append(update, bson.E{Key: "restaurant_charge.restaurant_charge_percent", Value: *req.RestaurantChargePercent})
	}

	update = bson.D{{
		Key:   "$set",
		Value: update,
	}}

	return update
}

func (s *StoreRepo) AddAddressCoordinates(ctx context.Context, storeID string, long, lat float64) error {
	objID, err := primitive.ObjectIDFromHex(storeID)
	if err != nil {
		return err
	}
	filter := bson.D{
		{Key: "_id", Value: objID},
	}
	update := bson.D{
		{
			Key: "$set", Value: bson.D{
				{Key: "address.coordinates.longitude", Value: long},
				{Key: "address.coordinates.latitude", Value: lat},
			},
		},
	}
	res, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return drivers2.ErrNotFound
	}
	return nil
}
