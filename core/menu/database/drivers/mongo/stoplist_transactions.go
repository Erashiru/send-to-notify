package mongo

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StopListTransaction struct {
	stColl *mongo.Collection
}

func NewStopListTransaction(stColl *mongo.Collection) *StopListTransaction {
	return &StopListTransaction{
		stColl: stColl,
	}
}

func (repo StopListTransaction) Insert(ctx context.Context, req models.StopListTransaction) (string, error) {

	res, err := repo.stColl.InsertOne(ctx, req)
	if err != nil {
		return "", errorSwitch(err)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", drivers.ErrInvalid
	}

	return oid.Hex(), nil
}
