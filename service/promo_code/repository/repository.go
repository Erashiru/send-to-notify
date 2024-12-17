package repository

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	"github.com/kwaaka-team/orders-core/service/promo_code/dto"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const promoCodeCollectionName = "promo_codes"

type Repository interface {
	CreatePromo(ctx context.Context, promoCodeRequest models.PromoCode) error
	GetPromoByCodeAndRestaurantID(ctx context.Context, promoCode string, restaurantID string) (models.PromoCode, error)
	UpdatePromo(ctx context.Context, updatePromoCodeRequest models.UpdatePromoCode) error
	GetPromoCodeByID(ctx context.Context, promoCodeID string) (models.PromoCode, error)
	GetPromoCodesByRestaurantId(ctx context.Context, restaurantId string, pagination selector.Pagination) ([]models.PromoCode, error)
	GetAvailablePromoCodeByCode(ctx context.Context, promoCodeValue string) (models.PromoCode, error)
}

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) (*MongoRepository, error) {
	r := MongoRepository{
		collection: db.Collection(promoCodeCollectionName),
	}
	return &r, nil
}

func (r *MongoRepository) CreatePromo(ctx context.Context, promoCodeRequest models.PromoCode) error {

	promoCodeRequest.CreatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, promoCodeRequest)
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoRepository) UpdatePromo(ctx context.Context, updatePromoCodeRequest models.UpdatePromoCode) error {

	objID, err := primitive.ObjectIDFromHex(updatePromoCodeRequest.ID)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objID}}

	update := r.setPromocodeUpdateFields(updatePromoCodeRequest)

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (r *MongoRepository) setPromocodeUpdateFields(promocode models.UpdatePromoCode) bson.D {
	var update primitive.D

	if promocode.Name != nil {
		update = append(update, bson.E{
			Key:   "name",
			Value: *promocode.Name,
		})
	}
	if promocode.Code != nil {
		update = append(update, bson.E{
			Key:   "code",
			Value: *promocode.Code,
		})
	}
	if promocode.Description != nil {
		update = append(update, bson.E{
			Key:   "description",
			Value: *promocode.Description,
		})
	}
	if promocode.Link != nil {
		update = append(update, bson.E{
			Key:   "link",
			Value: *promocode.Link,
		})
	}
	if promocode.RestaurantIDs != nil {
		update = append(update, bson.E{
			Key:   "restaurant_ids",
			Value: *promocode.RestaurantIDs,
		})
	}
	if promocode.UsageTime != nil {
		update = append(update, bson.E{
			Key:   "usage_time",
			Value: *promocode.UsageTime,
		})
	}
	if promocode.DeliveryType != nil {
		update = append(update, bson.E{
			Key:   "delivery_type",
			Value: *promocode.DeliveryType,
		})
	}
	if promocode.ValidFrom != nil {
		update = append(update, bson.E{
			Key:   "valid_from",
			Value: *promocode.ValidFrom,
		})
	}
	if promocode.ValidUntil != nil {
		update = append(update, bson.E{
			Key:   "valid_until",
			Value: *promocode.ValidUntil,
		})
	}
	if promocode.MinimumOrderPrice != nil {
		update = append(update, bson.E{
			Key:   "minimum_order_price",
			Value: *promocode.MinimumOrderPrice,
		})
	}
	if promocode.PromoCodeCategory != nil {
		update = append(update, bson.E{
			Key:   "promo_code_category",
			Value: *promocode.PromoCodeCategory,
		})
	}
	if promocode.SaleType != nil {
		update = append(update, bson.E{
			Key:   "sale_type",
			Value: *promocode.SaleType,
		})
	}
	if promocode.Sale != nil {
		update = append(update, bson.E{
			Key:   "sale",
			Value: *promocode.Sale,
		})
	}
	if promocode.Available != nil {
		update = append(update, bson.E{
			Key:   "available",
			Value: *promocode.Available,
		})
	}
	if promocode.ForAllProduct != nil {
		update = append(update, bson.E{
			Key:   "for_all_product",
			Value: *promocode.ForAllProduct,
		})
	}
	if promocode.Product != nil {
		update = append(update, bson.E{
			Key:   "products",
			Value: *promocode.Product,
		})
	}
	if promocode.IsDeleted != nil {
		update = append(update, bson.E{
			Key:   "is_deleted",
			Value: *promocode.IsDeleted,
		})
	}
	update = append(update, bson.E{
		Key:   "updated_at",
		Value: time.Now().UTC(),
	})

	return bson.D{{Key: "$set", Value: update}}
}

func (r *MongoRepository) GetPromoByCodeAndRestaurantID(ctx context.Context, promoCode string, restaurantID string) (models.PromoCode, error) {

	filter := bson.D{
		{Key: "code", Value: promoCode},
		{Key: "restaurant_ids", Value: restaurantID},
	}

	res := r.collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.PromoCode{}, dto.ErrPromoCodeNotFound
		}
		return models.PromoCode{}, err
	}

	var result models.PromoCode
	if err := res.Decode(&result); err != nil {
		return models.PromoCode{}, err
	}

	return result, nil
}

func (r *MongoRepository) GetPromoCodeByID(ctx context.Context, promoCodeID string) (models.PromoCode, error) {

	var promoCode models.PromoCode

	objID, err := primitive.ObjectIDFromHex(promoCodeID)
	if err != nil {
		return models.PromoCode{}, dto.ErrInvalidPromoCodeID
	}

	filter := bson.D{{Key: "_id", Value: objID}}

	if err := r.collection.FindOne(ctx, filter).Decode(&promoCode); err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return models.PromoCode{}, dto.ErrPromoCodeNotFound
		default:
			return models.PromoCode{}, err
		}
	}

	return promoCode, nil
}

func (r *MongoRepository) GetAvailablePromoCodeByCode(ctx context.Context, promoCodeValue string) (models.PromoCode, error) {

	var promoCode models.PromoCode

	filter := bson.D{
		{Key: "code", Value: promoCodeValue},
		{Key: "available", Value: true},
		{Key: "is_deleted", Value: false},
	}

	if err := r.collection.FindOne(ctx, filter).Decode(&promoCode); err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return models.PromoCode{}, dto.ErrPromoCodeNotFound
		default:
			return models.PromoCode{}, err
		}
	}

	return promoCode, nil
}

func (r *MongoRepository) GetPromoCodesByRestaurantId(ctx context.Context, restaurantId string, pagination selector.Pagination) ([]models.PromoCode, error) {

	var promoCodes []models.PromoCode

	pipeline := r.getPromoCodesAggregation(restaurantId, pagination)
	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var promoCode models.PromoCode
		if err := cursor.Decode(&promoCode); err != nil {
			return nil, err
		}
		promoCodes = append(promoCodes, promoCode)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return promoCodes, nil
}

func (r *MongoRepository) getPromoCodesAggregation(restaurantId string, pagination selector.Pagination) []bson.D {
	var pipeline []bson.D

	pipeline = append(pipeline, r.matchGetAllPromocodes(restaurantId))

	if pagination.HasPagination() {
		pipeline = append(pipeline, r.applyPagination(pagination)...)
	}

	return pipeline
}

func (r *MongoRepository) matchGetAllPromocodes(restGroupID string) bson.D {
	return bson.D{{Key: "$match", Value: bson.D{
		{Key: "restaurant_ids", Value: restGroupID}}},
	}
}

func (r *MongoRepository) applyPagination(pagination selector.Pagination) []bson.D {
	var withPagination []bson.D

	withPagination = append(withPagination, bson.D{{Key: "$skip", Value: pagination.Skip()}})
	withPagination = append(withPagination, bson.D{{Key: "$limit", Value: pagination.Limit}})

	return withPagination
}
