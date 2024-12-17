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

func (repo *MenuRepository) GetCombos(ctx context.Context, query selector.Menu) ([]models.Combo, int64, error) {
	oid, err := primitive.ObjectIDFromHex(query.ID)
	if err != nil {
		return nil, 0, drivers.ErrInvalid
	}

	match := bson.D{{Key: "_id", Value: oid}}
	unwind := "$combos"
	project := bson.D{{Key: "combos", Value: 1}}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$unwind", Value: unwind}},
		{{Key: "$project", Value: project}},
		{{Key: "$sort", Value: repo.sortFrom(query.Sorting)}},
	}

	if query.HasPagination() {
		pipeline = append(pipeline,
			bson.D{{Key: "$skip", Value: query.Skip()}},
			bson.D{{Key: "$limit", Value: query.Pagination.Limit}},
		)
	}

	cur, err := repo.menuColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer closeCur(cur)

	// TODO: testing this case
	combos := make([]models.Combo, 0, cur.RemainingBatchLength())

	for cur.Next(ctx) {
		var temp struct {
			Combo models.Combo `bson:"combos"`
		}

		if err = cur.Decode(&temp); err != nil {
			return nil, 0, err
		}

		combos = append(combos, temp.Combo)
	}

	return combos, 0, nil
}
