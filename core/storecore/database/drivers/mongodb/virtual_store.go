package mongodb

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	"github.com/kwaaka-team/orders-core/core/storecore/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ drivers.VirtualRepository = &VirtualStoreRepo{}

type VirtualStoreRepo struct {
	collection *mongo.Collection
}

func NewVirtualRepository(virtualStoreCollection *mongo.Collection) *VirtualStoreRepo {
	return &VirtualStoreRepo{
		collection: virtualStoreCollection,
	}
}

func (repo *VirtualStoreRepo) GetVirtualStore(ctx context.Context, query selector.VirtualStore) (models.VirtualStore, error) {
	filter, err := repo.filterFrom(ctx, query)
	if err != nil {
		return models.VirtualStore{}, errorSwitch(err)
	}

	res := repo.collection.FindOne(ctx, filter)

	var virtualStore models.VirtualStore

	if err = res.Decode(&virtualStore); err != nil {
		return models.VirtualStore{}, errorSwitch(err)
	}

	return virtualStore, nil
}

func (repo *VirtualStoreRepo) filterFrom(ctx context.Context, query selector.VirtualStore) (bson.D, error) {
	var result bson.D

	if query.HasRestaurantID() {
		result = append(result, bson.E{
			Key:   "virtual_store_restaurant_id",
			Value: query.RestaurantID,
		})
	}

	if query.HasName() {
		result = append(result, bson.E{
			Key:   "name",
			Value: query.Name,
		})
	}

	if query.HasChildRestaurantID() {
		result = append(result, bson.E{
			Key:   "restaurant_ids",
			Value: query.ChildRestaurantID,
		})
	}

	if query.HasExternalStoreID() {
		result = append(result, bson.E{
			Key:   "store_ids",
			Value: query.ExternalStoreID,
		})
	}

	if query.HasVirtualStoreType() {
		result = append(result, bson.E{
			Key:   "store_type",
			Value: query.VirtualStoreType,
		})
	}

	if query.HasDeliveryService() {
		result = append(result, bson.E{
			Key:   "delivery_service",
			Value: query.DeliveryService,
		})
	}

	if query.HasClientSecret() {
		result = append(result, bson.E{
			Key:   "client_secret",
			Value: query.ClientSecret,
		})
	}

	return result, nil
}
