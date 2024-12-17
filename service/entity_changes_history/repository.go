package entity_changes_history

import (
	"context"
	"github.com/kwaaka-team/orders-core/service/entity_changes_history/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Repository interface {
	InsertHistory(ctx context.Context, req models.EntityChangesHistory) (string, error)
	DeleteHistory(ctx context.Context, interval time.Duration) error
}

type repoImpl struct {
	collection *mongo.Collection
}

const collectionName = "entity_changes_history"

func NewEntityChangesHistoryMongoRepository(db *mongo.Database) (*repoImpl, error) {
	r := repoImpl{
		collection: db.Collection(collectionName),
	}

	return &r, nil
}

func (repository *repoImpl) DeleteHistory(ctx context.Context, interval time.Duration) error {
	filter := bson.D{
		{
			Key: "modified_at",
			Value: bson.M{
				"$lte": primitive.NewDateTimeFromTime(time.Now().UTC().Add(-interval)),
			},
		},
	}

	_, err := repository.collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (repository *repoImpl) InsertHistory(ctx context.Context, req models.EntityChangesHistory) (string, error) {
	req.ModifiedAt = time.Now().UTC()

	res, err := repository.collection.InsertOne(ctx, req)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}
