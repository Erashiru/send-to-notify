package mongo

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PromoRepository struct {
	collection *mongo.Collection
}

var _ drivers.PromoRepository = (*PromoRepository)(nil)

func NewPromoRepository(collection *mongo.Collection) *PromoRepository {
	return &PromoRepository{collection: collection}
}

func (repo *PromoRepository) GetPromos(ctx context.Context, query selector.Promo) (models.Promo, error) {
	filter, err := repo.filterFrom(query)
	if err != nil {
		return models.Promo{}, err
	}

	var promo models.Promo
	if err := repo.collection.FindOne(ctx, filter).Decode(&promo); err != nil {
		return models.Promo{}, errorSwitch(err)
	}

	return promo, nil
}

func (repo *PromoRepository) FindPromos(ctx context.Context, query selector.Promo) ([]models.Promo, error) {
	// Find many promos in DB by filter

	filter, err := repo.filterFrom(query)
	if err != nil {
		return nil, err
	}

	cur, err := repo.collection.Find(ctx, filter)
	if err != nil {
		return nil, errorSwitch(err)
	}

	promos := make([]models.Promo, 0, cur.RemainingBatchLength())
	if err := cur.All(ctx, &promos); err != nil {
		return nil, errorSwitch(err)
	}

	return promos, nil
}

func (repo *PromoRepository) filterFrom(query selector.Promo) (bson.D, error) {
	result := make(bson.D, 0, 6)

	if query.HasID() {
		oid, err := primitive.ObjectIDFromHex(query.ID)
		if err != nil {
			return nil, err
		}

		result = append(result, bson.E{
			Key:   "_id",
			Value: oid,
		})
	}

	if query.HasIsActive() {
		result = append(result, bson.E{
			Key:   "is_active",
			Value: query.IsActive,
		})
	}

	if query.HasStoreID() {
		result = append(result, bson.E{
			Key:   "restaurant_ids",
			Value: query.StoreID,
		})
	}

	if query.HasMenuID() {
		oid, err := primitive.ObjectIDFromHex(query.MenuID)
		if err != nil {
			return nil, err
		}

		result = append(result, bson.E{
			Key:   "menu_id",
			Value: oid,
		})
	}

	if query.HasDeliveryService() {
		result = append(result, bson.E{
			Key:   "delivery_service",
			Value: query.DeliveryService,
		})
	}

	if query.HasPOS() {
		result = append(result, bson.E{
			Key:   "pos_type",
			Value: query.POS,
		})
	}

	if query.HasProductIDs() {
		result = append(result, bson.E{
			Key: "product_ids", Value: bson.D{
				{Key: "$in", Value: query.ProductIDs},
			},
		})
	}

	return result, nil
}
