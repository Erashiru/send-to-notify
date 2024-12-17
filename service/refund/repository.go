package refund

import (
	"context"
	refundModels "github.com/kwaaka-team/orders-core/service/refund/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const collectionName = "refunds"

type Repository interface {
	InsertRefundInfo(ctx context.Context, refund refundModels.Refund) error
	GetRefund(ctx context.Context, orderID string) (refundModels.Refund, error)
}

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) (*MongoRepository, error) {
	return &MongoRepository{collection: db.Collection(collectionName)}, nil
}

func (m *MongoRepository) InsertRefundInfo(ctx context.Context, request refundModels.Refund) error {
	request.CreatedAt = time.Now().UTC()

	oid, err := primitive.ObjectIDFromHex(request.ID)
	if err != nil {
		return err
	}
	create := bson.D{
		{Key: "_id", Value: oid},
		{Key: "amount", Value: request.Amount},
		{Key: "reason", Value: request.Reason},
		{Key: "order_id", Value: request.OrderID},
		{Key: "payment_id", Value: request.PaymentID},
		{Key: "payment_system", Value: request.PaymentSystem},
		{Key: "created_at", Value: request.CreatedAt},
	}
	_, err = m.collection.InsertOne(ctx, create)
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoRepository) GetRefund(ctx context.Context, orderID string) (refundModels.Refund, error) {
	filter := bson.D{
		{Key: "order_id", Value: orderID},
	}
	var result refundModels.Refund

	err := m.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return refundModels.Refund{}, err
	}
	return result, nil
}
