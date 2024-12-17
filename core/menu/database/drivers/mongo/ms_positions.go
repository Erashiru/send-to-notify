package mongo

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MSPositionsRepository struct {
	collection *mongo.Collection
}

func (ms MSPositionsRepository) GetPositions(ctx context.Context, query selector.MoySklad) ([]models.MoySkladPosition, error) {
	filter := bson.M{"restaurant_id": query.RestaurantID, "is_deleted": query.IsDeleted}

	cur, err := ms.collection.Find(ctx, filter)
	if err != nil {
		return nil, errorSwitch(err)
	}
	defer closeCur(cur)

	result := make([]models.MoySkladPosition, 0, cur.RemainingBatchLength())
	if err = cur.All(ctx, &result); err != nil {
		return nil, errorSwitch(err)
	}

	return result, nil
}

func (ms MSPositionsRepository) RemovePosition(ctx context.Context, query selector.MoySklad) error {
	filter := bson.M{"restaurant_id": query.RestaurantID, "order_id": query.OrderID, "position_id": query.ID}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "is_deleted", Value: query.IsDeleted},
			{Key: "updated_at", Value: query.UpdatedAt},
		},
		}}
	log.Info().Msgf("repo remove position data %+v", query)
	res, err := ms.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return ErrAttributeGroupExtIDNotFound
	}
	return nil
}

func (ms MSPositionsRepository) CreatePosition(ctx context.Context, position models.MoySkladPosition) error {
	_, err := ms.collection.InsertOne(ctx, position)
	if err != nil {
		return err
	}
	return nil
}

var _ drivers.MSPositionsRepository = (*MSPositionsRepository)(nil)

func NewMSPositionsRepository(collection *mongo.Collection) *MSPositionsRepository {
	return &MSPositionsRepository{collection: collection}
}
