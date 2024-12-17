package mongo

import (
	"context"
	coreErrors "github.com/kwaaka-team/orders-core/core/errors"
	"github.com/pkg/errors"

	"github.com/kwaaka-team/orders-core/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BKOfferRepository struct {
	collection *mongo.Collection
}

func NewBKOfferRepository(bkOfferCollection *mongo.Collection) *BKOfferRepository {
	return &BKOfferRepository{
		collection: bkOfferCollection,
	}
}

func NewBKOfferRepository2(db *mongo.Database) (*BKOfferRepository, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	return &BKOfferRepository{
		collection: db.Collection(bkOfferCollectionName),
	}, nil
}

func (bkOfferRepo *BKOfferRepository) GetActiveOffers(ctx context.Context) ([]models.BKOffer, error) {
	filter := bson.D{
		{Key: "is_active", Value: true},
	}

	offersCursor, err := bkOfferRepo.collection.Find(ctx, filter)
	if err != nil {
		return nil, coreErrors.ErrorSwitch(err)
	}

	defer coreErrors.CloseCur(offersCursor)

	var offers []models.BKOffer

	if err = offersCursor.All(ctx, &offers); err != nil {
		return nil, coreErrors.ErrorSwitch(err)
	}

	return offers, nil
}
