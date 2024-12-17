package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type AnalyticsRepository struct {
	collection *mongo.Collection
}

func NewAnalyticsRepository(orderCollection *mongo.Collection) *AnalyticsRepository {
	return &AnalyticsRepository{
		collection: orderCollection,
	}
}
