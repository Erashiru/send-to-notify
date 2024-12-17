package mongo

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StoreRepository struct {
	restColl *mongo.Collection
}

func NewStoreRepository(restColl *mongo.Collection) *StoreRepository {
	return &StoreRepository{
		restColl: restColl,
	}
}

func (repo *StoreRepository) ListStoresByTalabatRestautantID(ctx context.Context, restaurantID string) ([]storeModels.Store, error) {
	cur, err := repo.restColl.Find(ctx, bson.D{
		{Key: "talabat.restaurant_id", Value: restaurantID},
	})
	if err != nil {
		return nil, errorSwitch(err)
	}

	stores := make([]storeModels.Store, 0, cur.RemainingBatchLength())
	if err := cur.All(ctx, &stores); err != nil {
		return nil, errorSwitch(err)
	}

	return stores, nil
}

func (repo *StoreRepository) List(ctx context.Context, query selector.Store) ([]storeModels.Store, int64, error) {

	filter, err := repo.filterFrom(query)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find()
	if query.HasPagination() {
		opts.SetSkip(query.Skip()).SetLimit(query.Limit)
	}
	opts.SetSort(repo.sortFrom(query.Sorting))

	count, err := repo.restColl.CountDocuments(ctx, filter)
	if err != nil || count == 0 {
		return []storeModels.Store{}, 0, err
	}

	cur, err := repo.restColl.Find(ctx, filter)
	if err != nil {
		return nil, 0, errorSwitch(err)
	}

	stores := make([]storeModels.Store, 0, cur.RemainingBatchLength())
	if err := cur.All(ctx, &stores); err != nil {
		return nil, 0, errorSwitch(err)
	}

	return stores, count, nil
}

func (repo *StoreRepository) Get(ctx context.Context, query selector.Store) (storeModels.Store, error) {
	filter, err := repo.filterFrom(query)
	if err != nil {
		return storeModels.Store{}, err
	}

	var store storeModels.Store
	if err = repo.restColl.FindOne(ctx, filter).Decode(&store); err != nil {
		return storeModels.Store{}, errorSwitch(err)
	}

	return store, nil
}

func (repo *StoreRepository) Update(ctx context.Context, store storeModels.Store) error {

	oid, err := primitive.ObjectIDFromHex(store.ID)
	if err != nil {
		return drivers.ErrInvalid
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	store.UpdatedAt = coreModels.TimeNow().Time

	set, err := repo.setFrom(store)
	if err != nil {
		return drivers.ErrInvalid
	}

	update := bson.D{
		{Key: "$set", Value: set},
	}

	res, err := repo.restColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return errorSwitch(err)
	}

	if res.MatchedCount == 0 {
		return drivers.ErrNotFound
	}

	return nil
}

func (repo *StoreRepository) filterFrom(query selector.Store) (bson.D, error) {
	filter := make(bson.D, 0, 5)

	if query.HasID() {
		oid, err := primitive.ObjectIDFromHex(query.ID)
		if err != nil {
			return nil, drivers.ErrInvalid
		}
		filter = append(filter, bson.E{
			Key: "_id", Value: oid,
		})
	}

	if query.HasToken() {
		filter = append(filter, bson.E{
			Key: "token", Value: query.Token,
		})
	}

	if query.HasExternalStoreID() && query.HasDeliveryService() {
		filter = append(filter, bson.E{
			Key: fmt.Sprintf("%s.store_id", query.DeliveryService), Value: query.ExternalStoreID,
		})
	}

	if query.HasDeliveryService() {
		filter = append(filter, bson.E{
			Key: "menus.delivery", Value: query.DeliveryService,
		})
	}

	if query.HasAggregatorMenuID() {
		filter = append(filter, bson.E{
			Key: "menus.menu_id", Value: query.AggregatorMenuID,
		})
	}

	if query.HasAggregatorMenuIDs() {

		oIDs := make([]primitive.ObjectID, 0, len(query.AggregatorMenuIDs))
		for _, menuID := range query.AggregatorMenuIDs {
			oid, err := primitive.ObjectIDFromHex(menuID)
			if err != nil {
				continue
			}
			oIDs = append(oIDs, oid)
		}

		if query.HasIsActiveMenu() {
			filter = append(filter, bson.E{
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
			filter = append(filter, bson.E{
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

	return filter, nil
}

func (*StoreRepository) sortFrom(query selector.Sorting) bson.D {
	sort := make(bson.D, 0, 2)

	if query.HasSorting() {
		sort = append(sort, bson.E{Key: query.Param, Value: query.Direction})
	}
	sort = append(sort, bson.E{Key: "_id", Value: 1})

	return sort
}

type menu struct {
	MenuID         primitive.ObjectID `bson:"menu_id"`
	Name           string             `bson:"name"`
	IsActive       bool               `bson:"is_active"`
	IsDeleted      bool               `bson:"is_deleted"`
	IsSync         bool               `bson:"is_sync"`
	SyncAttributes bool               `bson:"sync_attributes"`
	Delivery       string             `bson:"delivery"`
	Timestamp      int                `bson:"timestamp"`
	UpdatedAt      time.Time          `bson:"updated_at"`
}

func (*StoreRepository) setFrom(req storeModels.Store) (bson.D, error) {
	set := bson.D{
		{Key: "token", Value: req.Token},
		{Key: "name", Value: req.Name},
		{Key: "iiko_cloud", Value: req.IikoCloud},
		{Key: "tillypad", Value: req.TillyPad},
		{Key: "rkeeper", Value: req.RKeeper},
		{Key: "moysklad", Value: req.MoySklad},
		{Key: "glovo", Value: req.Glovo},
		{Key: "wolt", Value: req.Wolt},
		{Key: "yandex", Value: req.Yandex},
		{Key: "telegram", Value: req.Telegram},
		{Key: "delivery", Value: req.Delivery},
		{Key: "qr_menu", Value: req.QRMenu},
		{Key: "kwaaka_admin", Value: req.KwaakaAdmin},
		{Key: "callcenter", Value: req.CallCenter},
		{Key: "settings", Value: req.Settings},
		{Key: "integration_date", Value: req.IntegrationDate},
		{Key: "updated_at", Value: req.UpdatedAt},
	}

	if req.MenuID != "" {
		oid, err := primitive.ObjectIDFromHex(req.MenuID)
		if err != nil {
			return nil, drivers.ErrInvalid
		}
		set = append(set, bson.E{Key: "menu_id", Value: oid})
	}

	if len(req.Menus) > 0 {
		menus := make([]menu, 0, len(req.Menus))
		for _, v := range req.Menus {

			oid, err := primitive.ObjectIDFromHex(v.ID)
			if err != nil {
				continue
			}

			menus = append(menus, menu{
				MenuID:         oid,
				Name:           v.Name,
				IsActive:       v.IsActive,
				IsDeleted:      v.IsDeleted,
				IsSync:         v.IsSync,
				SyncAttributes: v.SyncAttributes,
				Delivery:       v.Delivery,
				Timestamp:      v.Timestamp,
				UpdatedAt:      v.UpdatedAt.UTC(),
			})

		}
		set = append(set, bson.E{Key: "menus", Value: menus})
	}

	return set, nil

}
