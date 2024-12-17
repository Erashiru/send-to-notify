package repository

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/service/payment/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const subscriptionsCollectionName = "subscriptions"

type SubscriptionsRepository interface {
	InsertSubscription(ctx context.Context, subscription models.PaymentSystemSubscription) (models.PaymentSystemSubscription, error)
	UpdateSubscription(ctx context.Context, subscription models.PaymentSystemSubscription) error
	GetSubscriptionByID(ctx context.Context, id string) (models.PaymentSystemSubscription, error)
}

type SubscriptionsMongoRepository struct {
	collection *mongo.Collection
}

func NewSubscriptionsMongoRepository(db *mongo.Database) (*SubscriptionsMongoRepository, error) {
	r := SubscriptionsMongoRepository{
		collection: db.Collection(subscriptionsCollectionName),
	}
	return &r, nil
}

func (r *SubscriptionsMongoRepository) InsertSubscription(ctx context.Context, subscription models.PaymentSystemSubscription) (models.PaymentSystemSubscription, error) {
	subscription.CreatedAt = time.Now().UTC()
	res, err := r.collection.InsertOne(ctx, subscription)
	if err != nil {
		return models.PaymentSystemSubscription{}, errors.ErrorSwitch(err)
	}

	subscription.ID = res.InsertedID.(primitive.ObjectID).Hex()

	return subscription, nil
}

func (r *SubscriptionsMongoRepository) UpdateSubscription(ctx context.Context, subscription models.PaymentSystemSubscription) error {
	oid, err := primitive.ObjectIDFromHex(subscription.ID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	subscription.ID = ""
	subscription.UpdatedAt = time.Now().UTC()

	update := bson.D{
		{Key: "$set", Value: subscription},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.ErrorSwitch(err)
	}

	if res.MatchedCount == 0 {
		return err
	}
	return nil
}

func (r *SubscriptionsMongoRepository) GetSubscriptionByID(ctx context.Context, id string) (models.PaymentSystemSubscription, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.PaymentSystemSubscription{}, err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	var subscription models.PaymentSystemSubscription

	if err := r.collection.FindOne(ctx, filter).Decode(&subscription); err != nil {
		return models.PaymentSystemSubscription{}, err
	}

	return subscription, nil
}
