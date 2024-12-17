package pos

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const anotherBillCollectionName = "another_bill_params"

type AnotherBillRepository interface {
	GetStoreIDs() ([]models.AnotherBillParams, error)
}

type MongoAnotherBillRepository struct {
	collection *mongo.Collection
}

func NewMongoAnotherBillRepository(collection *mongo.Database) (*MongoAnotherBillRepository, error) {
	r := MongoAnotherBillRepository{
		collection: collection.Collection(anotherBillCollectionName),
	}
	return &r, nil
}

func (r *MongoAnotherBillRepository) GetStoreIDs() ([]models.AnotherBillParams, error) {
	var result []models.AnotherBillParams
	cur, err := r.collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	if err := cur.All(context.Background(), &result); err != nil {
		return []models.AnotherBillParams{}, err
	}

	return result, nil
}
