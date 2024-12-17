package mongo

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BkOffers struct {
	bkOffersColl *mongo.Collection
}

func NewBkOffers(bkOffersColl *mongo.Collection) *BkOffers {
	return &BkOffers{
		bkOffersColl: bkOffersColl,
	}
}

func (repo *BkOffers) List(ctx context.Context, query selector.BkOffers) ([]models.BkOffers, error) {

	filter, err := repo.filterFrom(query)
	if err != nil {
		return nil, err
	}

	opts := options.Find()

	cur, err := repo.bkOffersColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, errorSwitch(err)
	}
	defer closeCur(cur)

	res := make([]models.BkOffers, 0, cur.RemainingBatchLength())
	if err = cur.All(ctx, &res); err != nil {
		return nil, errorSwitch(err)
	}

	return res, nil
}

func (repo *BkOffers) filterFrom(query selector.BkOffers) (bson.D, error) {
	filter := make(bson.D, 0, 2)

	if query.HasID() {
		oid, err := primitive.ObjectIDFromHex(query.ID)
		if err != nil {
			return nil, drivers.ErrInvalid
		}
		filter = append(filter, bson.E{
			Key: "_id", Value: oid,
		})
	}

	if query.HasIsActive() {
		filter = append(filter, bson.E{
			Key: "is_active", Value: query.IsActive,
		})
	}

	return filter, nil
}
