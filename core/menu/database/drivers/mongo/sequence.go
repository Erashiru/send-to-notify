package mongo

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SequencesRepository struct {
	collection *mongo.Collection
}

var _ drivers.SequencesRepository = (*SequencesRepository)(nil)

func NewSequencesRepository(collection *mongo.Collection) *SequencesRepository {
	return &SequencesRepository{collection: collection}
}

func (s *SequencesRepository) NextSequenceValue(ctx context.Context, name string) (int, error) {
	if name == "" {
		return 0, drivers.ErrEmptySequenceID
	}

	filter := bson.D{
		{Key: "_id", Value: name},
	}
	update := bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "value", Value: 1},
		}},
	}
	opts := options.
		FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var seq models.Sequence
	err := s.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&seq)
	if err != nil {
		return 0, err
	}

	return seq.Value, nil
}
