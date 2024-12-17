package mongo

import (
	"context"
	"errors"
	errors2 "github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RestaurantGroupMenuRepository struct {
	coll *mongo.Collection
}

func NewRestaurantGroupMenuRepository(coll *mongo.Collection) *RestaurantGroupMenuRepository {
	return &RestaurantGroupMenuRepository{coll: coll}
}

func (r *RestaurantGroupMenuRepository) GetMenuByRestGroupId(ctx context.Context, restGroupId string) (models.RestGroupMenu, error) {

	var result models.RestGroupMenu
	if err := r.coll.FindOne(ctx, bson.M{"rest_group_id": restGroupId}).Decode(&result); err != nil {
		switch {
		case errors.Is(errorSwitch(err), drivers.ErrNotFound):
			return models.RestGroupMenu{}, errors2.ErrNotFound
		default:
			return models.RestGroupMenu{}, err
		}
	}

	return result, nil
}

func (r *RestaurantGroupMenuRepository) UpdateOrCreateMenu(ctx context.Context, restGroupId string, newMenu models.RestGroupMenu) error {

	filter := bson.M{"rest_group_id": restGroupId}
	update := bson.M{
		"$set": newMenu,
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}
