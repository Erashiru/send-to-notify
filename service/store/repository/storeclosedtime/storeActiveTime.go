package storeclosedtime

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const collectionActiveTime = "store_closed_time"

type Repository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) (Repository, error) {
	r := Repository{
		collection: db.Collection(collectionActiveTime),
	}
	return r, nil
}

func (r Repository) Insert(ctx context.Context, request models.StoreActiveTime) error {
	_, err := r.collection.InsertOne(ctx, request)
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) GetByFilter(ctx context.Context, record models.FilterStoreActiveTime) (models.StoreActiveTime, bool, error) {
	filter := r.filter(record)
	var result models.StoreActiveTime

	doc := r.collection.FindOne(ctx, filter)
	if doc.Err() != nil {
		if errors.Is(doc.Err(), mongo.ErrNoDocuments) {
			return result, false, nil
		}
		return result, false, doc.Err()
	}

	if err := doc.Decode(&result); err != nil {
		return result, false, err
	}

	return result, true, nil
}

func (r Repository) filter(record models.FilterStoreActiveTime) bson.D {
	filter := bson.D{
		primitive.E{Key: "restaurant_id", Value: record.RestaurantID},
		primitive.E{Key: "store_id", Value: record.StoreID},
		primitive.E{Key: "delivery_service", Value: record.DeliveryService},
		primitive.E{Key: "end_time", Value: bson.D{{Key: "$eq", Value: nil}}},
	}
	return filter
}

func (r Repository) UpdateEndTime(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "end_time", Value: time.Now()},
		}},
	}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if result == nil || err != nil {
		return err
	}
	return nil
}
