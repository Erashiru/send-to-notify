package repository

import (
	"context"
	customeErrors "github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/service/pos/models/poster"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const posterStoreAuthsCollectionName = "poster_store_auths"

type PosterStoreAuthsRepository interface {
	InsertCustomer(ctx context.Context, posterStore poster.PosterStoreAuth) (poster.PosterStoreAuth, error)
}

type PosterStoreAuthsMongoRepository struct {
	collection *mongo.Collection
}

func NewPosterStoreAuthsMongoRepository(db *mongo.Database) (*PosterStoreAuthsMongoRepository, error) {
	r := PosterStoreAuthsMongoRepository{
		collection: db.Collection(posterStoreAuthsCollectionName),
	}
	return &r, nil
}

func (r *PosterStoreAuthsMongoRepository) InsertCustomer(ctx context.Context, posterStore poster.PosterStoreAuth) (poster.PosterStoreAuth, error) {
	posterStore.CreatedAt = time.Now().UTC()
	res, err := r.collection.InsertOne(ctx, posterStore)
	if err != nil {
		return poster.PosterStoreAuth{}, customeErrors.ErrorSwitch(err)
	}

	posterStore.ID = res.InsertedID.(primitive.ObjectID).Hex()

	return posterStore, nil
}
