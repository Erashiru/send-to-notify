package mongodb

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/integration_api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const externalPosAuthCollectionName = "external_pos_integration_auth"

type repo struct {
	collection *mongo.Collection
}

func NewExternalPosIntegrationAuthRepository(mongoClient *mongo.Client, dbName string) *repo {
	return &repo{
		collection: mongoClient.Database(dbName).Collection(externalPosAuthCollectionName),
	}
}

func (repo *repo) GetAuthInfo(ctx context.Context, token string) (models.AuthInfo, error) {
	filter := bson.D{
		{
			Key:   "token",
			Value: token,
		},
	}

	result := repo.collection.FindOne(ctx, filter)

	var info models.AuthInfo

	if err := result.Decode(&info); err != nil {
		return models.AuthInfo{}, err
	}

	return info, nil
}
