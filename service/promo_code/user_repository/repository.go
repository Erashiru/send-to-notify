package user_repository

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const userAndPromoCodeCollectionName = "user_promo_codes"

type Repository interface {
	CreateUserUsePromoCodeTime(ctx context.Context, userId, promoCode string, restIds []string) error
	GetUsageCountForUser(ctx context.Context, userId, promoCode, restaurantId string) (int, error)
	UpdateUsageTimeForUser(ctx context.Context, userId, promoCode, restaurantId string, usageTime int) error
}

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) (*MongoRepository, error) {
	r := MongoRepository{
		collection: db.Collection(userAndPromoCodeCollectionName),
	}
	return &r, nil
}

func (r *MongoRepository) GetUsageCountForUser(ctx context.Context, userId, promoCode, restaurantId string) (int, error) {

	filter := bson.D{
		{Key: "user_id", Value: userId},
		{Key: "promo_code", Value: promoCode},
		{Key: "restaurant_ids", Value: restaurantId},
	}

	res := r.collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, nil
		}
		return 0, err
	}

	var result models.UserAndPromoCode
	if err := res.Decode(&result); err != nil {
		return 0, err
	}

	return result.UsageTime, nil
}

func (r *MongoRepository) CreateUserUsePromoCodeTime(ctx context.Context, userId, promoCode string, restIds []string) error {

	userPromoCode := models.UserAndPromoCode{
		UserId:        userId,
		UsageTime:     1,
		RestaurantIds: restIds,
		PromoCode:     promoCode,
	}

	_, err := r.collection.InsertOne(ctx, userPromoCode)
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoRepository) UpdateUsageTimeForUser(ctx context.Context, userId, promoCode, restaurantId string, usageTime int) error {

	filter := bson.D{
		{Key: "user_id", Value: userId},
		{Key: "promo_code", Value: promoCode},
		{Key: "restaurant_ids", Value: restaurantId},
	}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "usage_time", Value: usageTime + 1}}}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return err
	}
	return nil
}
