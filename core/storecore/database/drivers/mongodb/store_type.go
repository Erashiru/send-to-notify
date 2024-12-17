package mongodb

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var _ drivers.StoreTypeRepository = &StoreTypeRepo{}

type StoreTypeRepo struct {
	collection *mongo.Collection
}

func NewStoreTypeRepository(collection *mongo.Collection) drivers.StoreTypeRepository {
	return &StoreTypeRepo{
		collection: collection,
	}
}

func (s *StoreTypeRepo) GetList(ctx context.Context) ([]models.StoreType, error) {

	cur, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, errorSwitch(err)
	}

	var stores []models.StoreType
	for cur.Next(ctx) {
		var store models.StoreType
		if err = cur.Decode(&store); err != nil {
			return nil, errorSwitch(err)
		}
		stores = append(stores, store)
	}

	return stores, nil
}
