package mongo

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthClientRepository struct {
	collection *mongo.Collection
}

func NewAuthClientRepository(authClientCollection *mongo.Collection) *AuthClientRepository {
	return &AuthClientRepository{
		collection: authClientCollection,
	}
}

func (authClientRepo *AuthClientRepository) FindByIDAndSecret(ctx context.Context, clientID string, clientSecret string) (models.AuthClient, error) {
	filter := bson.D{
		{Key: "client_id", Value: clientID},
		{Key: "client_secret", Value: clientSecret},
	}

	var authClient models.AuthClient
	err := authClientRepo.collection.FindOne(ctx, filter).Decode(&authClient)
	if err != nil {
		return models.AuthClient{}, errors.ErrorSwitch(err)
	}

	return authClient, nil
}

func (authClientRepo *AuthClientRepository) FindByID(ctx context.Context, clientID string) (models.AuthClient, error) {
	filter := bson.D{
		{Key: "client_id", Value: clientID},
	}

	var authClient models.AuthClient
	err := authClientRepo.collection.FindOne(ctx, filter).Decode(&authClient)

	if err != nil {
		return models.AuthClient{}, errors.ErrorSwitch(err)
	}

	return authClient, nil
}

func (authClientRepo *AuthClientRepository) CreateAuthClient(ctx context.Context, req models.AuthClient) (string, error) {
	res, err := authClientRepo.collection.InsertOne(ctx, req)
	if err != nil {
		return "", errors.ErrorSwitch(err)
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (authClientRepo *AuthClientRepository) AuthClientExist(ctx context.Context, clientID string, clientSecret string) error {

	filter := bson.D{
		{Key: "client_id", Value: clientID},
		{Key: "client_secret", Value: clientSecret},
	}
	count, err := authClientRepo.collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("document by this client_id %s and client_secret already exist", clientID)
	}
	return nil
}

func (authClientRepo *AuthClientRepository) GetListID(ctx context.Context) ([]models.AuthClient, error) {
	filter := bson.M{}

	cursor, err := authClientRepo.collection.Find(ctx, filter)
	if err != nil {
		return []models.AuthClient{}, err
	}
	defer cursor.Close(ctx)

	var res []models.AuthClient

	for cursor.Next(ctx) {
		var temp models.AuthClient
		if err := cursor.Decode(&temp); err != nil {
			return []models.AuthClient{}, err
		}
		res = append(res, temp)

	}
	return res, nil
}

func (authClientRepo *AuthClientRepository) GetAuthClientByID(ctx context.Context, authID string) (models.AuthClient, error) {
	id, err := primitive.ObjectIDFromHex(authID)
	if err != nil {
		return models.AuthClient{}, err
	}

	filter := bson.M{"_id": id}

	var result models.AuthClient
	err = authClientRepo.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return models.AuthClient{}, err
	}
	return result, nil
}
