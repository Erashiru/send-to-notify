package store

import (
	"context"
	"fmt"
	customErrors "github.com/kwaaka-team/orders-core/core/errors"
	kwaakaAdminModels "github.com/kwaaka-team/orders-core/core/kwaaka_admin/models"
	"github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionStoreName = "restaurants"

type Repository interface {
	GetById(ctx context.Context, storeID string) (*models.Store, error)
	GetByExternalIdAndAggregator(ctx context.Context, externalStoreID string, deliveryService string) (*models.Store, error)
	GetStoresByDeliveryService(ctx context.Context, deliveryService string) ([]models.Store, error)
	GetStoresByToken(ctx context.Context, token string) ([]models.Store, error)
	FindStoresByPosType(ctx context.Context, posType string) ([]models.Store, error)
	FindAllStores(ctx context.Context) ([]models.Store, error)
	FindStoresByStoreGroupID(ctx context.Context, storeGroupID string) ([]models.Store, error)
	FindStoresByTimeZone(ctx context.Context, timeZone string) ([]models.Store, error)
	UpdateTZ(ctx context.Context, storeID string, targetTZ string) error
	UpdateOffset(ctx context.Context, storeID string, targetOffset int) error
	GetByYarosRestaurantID(ctx context.Context, restaurantID string) (models.Store, error)
	GetRestaurantsByGroupId(ctx context.Context, pagination selector.Pagination, restGroupId string) ([]models.Store, error)
	UpdateMenus(ctx context.Context, storeId string, menus []models.StoreDSMenu) error
	FindStoreInRestGroupByName(ctx context.Context, name string, restaurantGroupId string) ([]models.Store, error)
	UpdateStoreSchedule(ctx context.Context, storeId string, schedule models.AggregatorSchedule, queryPrefix string) error
	UpdateMenuId(ctx context.Context, storeId, menuId string) error
	CreatePolygon(ctx context.Context, restaurantID string, request models.Geometry) error
	UpdatePolygon(ctx context.Context, restaurantID string, request models.Geometry) error
	GetPolygonByRestaurantID(ctx context.Context, restaurantID string) (models.Geometry, error)
	CreateStorePhoneEmail(ctx context.Context, restaurantID string, request kwaakaAdminModels.StorePhoneEmail) error
	UpdateKwaakaAdminBusyMode(ctx context.Context, storeID string, busyMode bool, busyModeTime int) error
	GetStoresByIIKOOrganizationId(ctx context.Context, organizationId string) ([]models.Store, error)
	GetStoresByWppPhoneNumber(ctx context.Context, phoneNum string) ([]models.Store, error)
	UpdateStoreByFields(ctx context.Context, store storeModels.UpdateStore) error
	GetStoresBySelectorFilter(ctx context.Context, query selector.Store) ([]models.Store, error)
}

type MongoRepository struct {
	collection                    *mongo.Collection
	aggregatorsWithSeparateFields map[string]interface{}
}

func NewStoreMongoRepository(db *mongo.Database) (*MongoRepository, error) {
	r := MongoRepository{
		collection: db.Collection(collectionStoreName),
		aggregatorsWithSeparateFields: map[string]interface{}{
			"glovo":        struct{}{},
			"wolt":         struct{}{},
			"chocofood":    struct{}{},
			"express24":    struct{}{},
			"moysklad":     struct{}{},
			"qr_menu":      struct{}{},
			"deliveroo":    struct{}{},
			"talabat":      struct{}{},
			"kwaaka_admin": struct{}{},
			"starter_app":  struct{}{},
		},
	}
	return &r, nil
}

func (r *MongoRepository) GetStoresByIIKOOrganizationId(ctx context.Context, organizationId string) ([]models.Store, error) {
	filter := bson.D{
		{
			Key:   "iiko_cloud.organization_id",
			Value: organizationId,
		},
	}

	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var stores []models.Store
	for cur.Next(ctx) {
		var store models.Store
		if err = cur.Decode(&store); err != nil {
			return nil, err
		}

		stores = append(stores, store)
	}

	return stores, nil
}

func (r *MongoRepository) GetStoresByWppPhoneNumber(ctx context.Context, phoneNum string) ([]models.Store, error) {
	filter := bson.D{
		{
			Key:   "whatsapp_config.phone_number",
			Value: phoneNum,
		},
	}

	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var res []models.Store
	for cur.Next(ctx) {
		var store models.Store
		if err := cur.Decode(&store); err != nil {
			return nil, err
		}
		res = append(res, store)
	}

	return res, nil
}

func (r *MongoRepository) UpdateMenuId(ctx context.Context, storeId, menuId string) error {
	filter, err := r.filterFrom(selector.NewEmptyStoreSearch().SetID(storeId))
	if err != nil {
		return errorSwitch(err)
	}

	oid, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		return errorSwitch(err)
	}

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "menu_id", Value: oid},
			},
		},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("matched count is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) UpdateStoreSchedule(ctx context.Context, storeId string, schedule models.AggregatorSchedule, queryPrefix string) error {
	filter, err := r.filterFrom(selector.NewEmptyStoreSearch().SetID(storeId))
	if err != nil {
		return err
	}

	var update bson.D

	update = append(update, bson.E{
		Key: "$set",
		Value: bson.D{
			{Key: queryPrefix, Value: schedule},
		},
	})

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("matched count is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) UpdateMenus(ctx context.Context, storeId string, menus []models.StoreDSMenu) error {
	oid, err := primitive.ObjectIDFromHex(storeId)
	if err != nil {
		return err
	}

	filter := bson.D{
		{
			Key:   "_id",
			Value: oid,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"menus": menus,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errorSwitch(err)
	}

	if result.MatchedCount == 0 {
		return errors.New("matched count is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) UpdateOffset(ctx context.Context, storeID string, targetOffset int) error {
	oid, err := primitive.ObjectIDFromHex(storeID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{
			Key:   "_id",
			Value: oid,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"settings.timezone.utc_offset": targetOffset,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errorSwitch(err)
	}

	if result.MatchedCount == 0 {
		return errors.New("matched count is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) UpdateTZ(ctx context.Context, storeID string, targetTZ string) error {
	oid, err := primitive.ObjectIDFromHex(storeID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{
			Key:   "_id",
			Value: oid,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"settings.timezone.tz": targetTZ,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errorSwitch(err)
	}

	if result.MatchedCount == 0 {
		return errors.New("matched count is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) GetById(ctx context.Context, storeID string) (*models.Store, error) {
	query := selector.NewEmptyStoreSearch().
		SetID(storeID)

	filter, err := r.filterFrom(query)
	if err != nil {
		return nil, err
	}

	var store models.Store
	if err = r.collection.FindOne(ctx, filter).Decode(&store); err != nil {
		return nil, errorSwitch(err)
	}

	return &store, nil
}

func (r *MongoRepository) GetByExternalIdAndAggregator(ctx context.Context, externalStoreID string, deliveryService string) (*models.Store, error) {
	query := selector.NewEmptyStoreSearch().
		SetExternalStoreID(externalStoreID)

	if r.isHasSeparateField(deliveryService) {
		query = query.SetDeliveryService(deliveryService)
	} else {
		query = query.SetExternalDeliveryService(deliveryService)
	}

	filter, err := r.filterFrom(query)
	if err != nil {
		return nil, err
	}

	var store models.Store
	if err = r.collection.FindOne(ctx, filter).Decode(&store); err != nil {
		return nil, errorSwitch(err)
	}

	if r.isHasSeparateField(deliveryService) {
		return &store, nil
	}

	for _, external := range store.ExternalConfig {
		if external.Type != deliveryService {
			continue
		}
		for _, storeID := range external.StoreID {
			if externalStoreID != storeID {
				continue
			}
			return &store, nil
		}
	}

	return nil, errors.New("restaurant is not exists")
}

func (r *MongoRepository) filterFrom(query selector.Store) (bson.D, error) {
	result := bson.D{}

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
			return nil, errors.Wrap(drivers.ErrInvalid, "query.ID error")
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

	if query.HasDeferSubmission() {
		result = append(result, bson.E{
			Key:   "defer_submission.is_defer_submission",
			Value: query.DeferSubmission,
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

func (r *MongoRepository) isHasSeparateField(aggregator string) bool {
	if _, ok := r.aggregatorsWithSeparateFields[aggregator]; ok {
		return true
	}
	return false
}

func toObjectIDs(ids []string) ([]primitive.ObjectID, error) {
	var objIDs []primitive.ObjectID
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		objIDs = append(objIDs, oid)
	}
	return objIDs, nil
}

func errorSwitch(err error) error {
	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		return errors.Wrap(customErrors.ErrStoreNotFound, err.Error())
	case mongo.IsDuplicateKeyError(err):
		return errors.Wrap(drivers.ErrAlreadyExist, "duplicate key error")
	default:
		return err
	}
}

func (r *MongoRepository) GetStoresByDeliveryService(ctx context.Context, deliveryService string) ([]models.Store, error) {
	filter := r.filterByDeliveryService(deliveryService)

	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var stores []models.Store
	for cur.Next(ctx) {
		var store models.Store
		if err = cur.Decode(&store); err != nil {
			return nil, err
		}

		stores = append(stores, store)
	}

	return stores, nil
}

func (r *MongoRepository) filterByDeliveryService(deliveryService string) bson.D {
	var filter bson.D
	filter = append(filter, bson.E{
		Key: fmt.Sprintf("%s.store_id", deliveryService),
		Value: bson.D{
			{Key: "$exists", Value: true},
			{Key: "$ne", Value: bson.A{}},
		},
	})
	filter = append(filter, bson.E{
		Key:   fmt.Sprintf("%s.send_to_pos", deliveryService),
		Value: true,
	})
	return filter
}

func (r *MongoRepository) GetStoresByToken(ctx context.Context, token string) ([]models.Store, error) {
	query := selector.NewEmptyStoreSearch().
		SetToken(token)

	filter, err := r.filterFrom(query)
	if err != nil {
		return nil, err
	}

	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var stores []models.Store
	for cur.Next(ctx) {
		var store models.Store
		if err = cur.Decode(&store); err != nil {
			return nil, err
		}

		stores = append(stores, store)
	}

	return stores, nil
}

func (r *MongoRepository) FindStoresByTimeZone(ctx context.Context, timeZone string) ([]models.Store, error) {
	var filter bson.D

	filter = append(filter, bson.E{
		Key:   "settings.timezone.tz",
		Value: timeZone,
	})
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var stores []models.Store
	for cur.Next(ctx) {
		var store models.Store
		if err = cur.Decode(&store); err != nil {
			return nil, err
		}

		stores = append(stores, store)
	}

	return stores, nil
}

func (r *MongoRepository) FindStoresByPosType(ctx context.Context, posType string) ([]models.Store, error) {
	var filter bson.D

	filter = append(filter, bson.E{
		Key:   "pos_type",
		Value: posType,
	})
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var stores []models.Store
	for cur.Next(ctx) {
		var store models.Store
		if err = cur.Decode(&store); err != nil {
			return nil, err
		}

		stores = append(stores, store)
	}

	return stores, nil
}

func (r *MongoRepository) FindStoresByStoreGroupID(ctx context.Context, storeGroupID string) ([]models.Store, error) {
	var filter bson.D

	filter = append(filter, bson.E{
		Key:   "restaurant_group_id",
		Value: storeGroupID,
	})
	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var stores []models.Store
	for cur.Next(ctx) {
		var store models.Store
		if err = cur.Decode(&store); err != nil {
			return nil, err
		}

		stores = append(stores, store)
	}

	return stores, nil
}

func (r *MongoRepository) GetByYarosRestaurantID(ctx context.Context, restaurantID string) (models.Store, error) {
	var filter bson.D

	filter = append(filter, bson.E{
		Key:   "yaros.store_id",
		Value: restaurantID,
	})
	var store models.Store
	if err := r.collection.FindOne(ctx, filter).Decode(&store); err != nil {
		return models.Store{}, errorSwitch(err)
	}
	return store, nil
}

func (r *MongoRepository) FindAllStores(ctx context.Context) ([]models.Store, error) {

	cur, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var stores []models.Store
	for cur.Next(ctx) {
		var store models.Store
		if err = cur.Decode(&store); err != nil {
			return nil, err
		}

		stores = append(stores, store)
	}
	return stores, nil
}

func (r *MongoRepository) GetRestaurantsByGroupId(ctx context.Context, pagination selector.Pagination, restGroupId string) ([]models.Store, error) {
	var stores []models.Store
	var cursor *mongo.Cursor
	var err error

	filter := bson.M{}

	if restGroupId != "" {
		filter["restaurant_group_id"] = restGroupId
	}

	if pagination.HasPagination() {
		cursor, err = r.collection.Find(ctx, filter, options.Find().SetLimit(pagination.Limit).SetSkip(pagination.Limit*(pagination.Page-1)))
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)
	} else {
		cursor, err = r.collection.Find(ctx, filter, options.Find())
		if err != nil {
			return nil, err
		}
		defer cursor.Close(ctx)
	}

	for cursor.Next(ctx) {
		var store models.Store
		if err := cursor.Decode(&store); err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}

	if err := cursor.Err(); err != nil {
		return []models.Store{}, err
	}

	return stores, nil
}

func (r *MongoRepository) FindStoreInRestGroupByName(ctx context.Context, name string, restaurantGroupId string) ([]models.Store, error) {
	nameRegex := `(?i)` + name
	filter := bson.M{
		"restaurant_group_id": restaurantGroupId,
		"name":                bson.M{"$regex": nameRegex},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var stores []models.Store
	for cursor.Next(ctx) {
		var store models.Store
		if err := cursor.Decode(&store); err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return stores, nil
}

func (r *MongoRepository) CreatePolygon(ctx context.Context, restaurantID string, request models.Geometry) error {

	oid, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "address.geometry", Value: request},
		}},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil

}

func (r *MongoRepository) UpdatePolygon(ctx context.Context, restaurantID string, request models.Geometry) error {

	oid, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := r.setFields(request)

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *MongoRepository) setFields(request models.Geometry) bson.D {

	update := make(bson.D, 0, 2)

	if request.Coordinates != nil {
		update = append(update, bson.E{Key: "address.geometry.coordinates", Value: request.Coordinates})
	}

	update = append(update,
		bson.E{Key: "address.geometry.percentage_modifier", Value: request.PercentageModifier},
	)

	result := bson.D{
		{Key: "$set", Value: update},
	}

	return result
}

func (r *MongoRepository) GetPolygonByRestaurantID(ctx context.Context, restaurantID string) (models.Geometry, error) {

	oid, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return models.Geometry{}, err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	var store models.Store
	if err := r.collection.FindOne(ctx, filter).Decode(&store); err != nil {
		return models.Geometry{}, err
	}

	return store.Address.Geometry, nil
}

func (r *MongoRepository) CreateStorePhoneEmail(ctx context.Context, restaurantID string, request kwaakaAdminModels.StorePhoneEmail) error {

	oid, err := primitive.ObjectIDFromHex(restaurantID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "settings.email", Value: request.Email},
			{Key: "store_phone_number", Value: request.PhoneNumber},
		}},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *MongoRepository) UpdateKwaakaAdminBusyMode(ctx context.Context, storeID string, busyMode bool, busyModeTime int) error {
	filter := bson.D{
		{Key: "kwaaka_admin.store_id", Value: storeID},
	}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "kwaaka_admin.busy_mode", Value: busyMode},
			{Key: "kwaaka_admin.adjusted_pickup_minutes", Value: busyModeTime},
		}},
	}
	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *MongoRepository) UpdateStoreByFields(ctx context.Context, store storeModels.UpdateStore) error {
	if store.ID == nil {
		return fmt.Errorf("storeID is nil")
	}

	oid, err := primitive.ObjectIDFromHex(*store.ID)
	if err != nil {
		return errors.Wrap(err, "storeID error")
	}

	filter := bson.D{{Key: "_id", Value: oid}}

	update, err := r.setAllFields(store)
	if err != nil {
		return err
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("not found")
	}

	return nil
}

func (r *MongoRepository) GetStoresBySelectorFilter(ctx context.Context, query selector.Store) ([]models.Store, error) {
	filter, err := r.filterFrom(query)
	if err != nil {
		return nil, err
	}

	cur, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var stores []models.Store
	for cur.Next(ctx) {
		var store models.Store
		if err := cur.Decode(&store); err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}

	return stores, nil
}

func (r *MongoRepository) setAllFields(store storeModels.UpdateStore) (bson.D, error) {
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
			return nil, errors.Wrap(err, "store.MenuID error")
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
		update = r.setGlovoConfig(update, store.Glovo)
	}

	if store.Wolt != nil {
		update = r.setWoltConfig(update, store.Wolt)
	}

	if store.External != nil {
		update = append(update, bson.E{
			Key:   "external",
			Value: store.External,
		})
	}

	if store.QRMenu != nil {
		update = r.setQRMenuConfig(update, store.QRMenu)
	}

	if store.KwaakaAdmin != nil {
		update = r.setKwaakaAdminConfig(update, store.KwaakaAdmin)
	}

	if store.IikoCloud != nil {
		update = r.setIIKOConfig(update, store.IikoCloud)
	}

	if store.Paloma != nil {
		update = append(update, bson.E{
			Key:   "paloma",
			Value: store.Paloma,
		})
	}

	if store.RKeeper7XML != nil {
		update = r.setRKeeper7XMLConfig(update, store.RKeeper7XML)
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

	result := bson.D{
		{Key: "$set", Value: update},
	}

	return result, nil
}

func (r *MongoRepository) setGlovoConfig(update bson.D, glovoCfg *models.UpdateStoreGlovoConfig) bson.D {
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
		update = r.setPaymentTypesFields(update, "glovo", glovoCfg.PaymentTypes)
	}

	if glovoCfg.PurchaseTypes != nil {
		update = r.setPurchaseTypes(update, "glovo", glovoCfg.PurchaseTypes)
	}

	return update
}

func (r *MongoRepository) setWoltConfig(update bson.D, woltCfg *models.UpdateStoreWoltConfig) bson.D {
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
		update = r.setPaymentTypesFields(update, "wolt", woltCfg.PaymentTypes)
	}

	if woltCfg.PurchaseTypes != nil {
		update = r.setPurchaseTypes(update, "wolt", woltCfg.PurchaseTypes)
	}

	return update
}

func (r *MongoRepository) setQRMenuConfig(update bson.D, qrmenuCfg *models.UpdateStoreQRMenuConfig) bson.D {
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
		update = r.setPaymentTypesFields(update, "qr_menu", qrmenuCfg.PaymentTypes)
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

func (r *MongoRepository) setKwaakaAdminConfig(update bson.D, kwaakaAdminCfg *models.UpdateStoreKwaakaAdminConfig) bson.D {

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

func (r *MongoRepository) setPurchaseTypes(update bson.D, delivery string, purchaseTypes *models.UpdatePurchaseTypes) bson.D {
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
func (r *MongoRepository) setPaymentTypesFields(update bson.D, delivery string, paymentTypes *models.UpdateDeliveryServicePaymentType) bson.D {

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

func (r *MongoRepository) setYandexPaymentTypes(update bson.D, yandexPaymentTypes *models.UpdateDeliveryServicePaymentType) bson.D {
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

func (r *MongoRepository) setIIKOConfig(update bson.D, iikoConfigs *models.UpdateStoreIikoConfig) bson.D {
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

func (r *MongoRepository) setRKeeper7XMLConfig(update bson.D, rkeeper7XMLCfg *models.UpdateStoreRKeeper7XMLConfig) bson.D {
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
