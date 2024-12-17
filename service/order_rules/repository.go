package order_rules

import (
	"context"
	"github.com/kwaaka-team/orders-core/service/order_rules/models"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionOrderRulesName = "order_rules"

type Repository interface {
	FindOrderRulesByRestaurantId(ctx context.Context, restaurantId string) ([]models.OrderRule, error)
}

type repoImpl struct {
	collection *mongo.Collection
}

func NewOrderRulesMongoRepository(db *mongo.Database) (*repoImpl, error) {
	r := repoImpl{
		collection: db.Collection(collectionOrderRulesName),
	}
	return &r, nil
}

func (m *repoImpl) FindOrderRulesByRestaurantId(ctx context.Context, restaurantId string) ([]models.OrderRule, error) {
	filter := bson.D{
		{
			Key:   "restaurant_ids",
			Value: restaurantId,
		},
	}

	res, err := m.collection.Find(ctx, filter)
	if err != nil {
		return nil, errorSwitch(err)
	}
	defer res.Close(ctx)

	orderRules := make([]models.OrderRule, 0, res.RemainingBatchLength())

	for res.Next(ctx) {
		var orderRule models.OrderRule

		if err = res.Decode(&orderRule); err != nil {
			log.Err(err).Msgf("error decoding into models.OrderRule model")
			continue
		}

		orderRules = append(orderRules, orderRule)
	}

	return orderRules, nil
}
