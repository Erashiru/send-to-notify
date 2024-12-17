package restaurant_set

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionName = "restaurant_set"

type Repository interface {
	CreateRestaurantSet(ctx context.Context, entity models.RestaurantSet) (string, error)
	GetRestaurantSetById(ctx context.Context, id string) (models.RestaurantSet, error)
	DeleteRestaurantSetById(ctx context.Context, id string) error
	GetRestaurantSetByDomainName(ctx context.Context, domainName string) (models.RestaurantSet, error)
}

type MongoRepository struct {
	coll *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) (*MongoRepository, error) {
	return &MongoRepository{coll: db.Collection(collectionName)}, nil
}

func (m *MongoRepository) CreateRestaurantSet(ctx context.Context, entity models.RestaurantSet) (string, error) {
	res, err := m.coll.InsertOne(ctx, entity)
	if err != nil {
		return "", err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", drivers.ErrInvalid
	}

	return oid.Hex(), nil
}

func (m *MongoRepository) GetRestaurantSetById(ctx context.Context, id string) (models.RestaurantSet, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.RestaurantSet{}, err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	var res models.RestaurantSet
	if err = m.coll.FindOne(ctx, filter).Decode(&res); err != nil {
		return models.RestaurantSet{}, err
	}

	return res, nil
}

func (m *MongoRepository) DeleteRestaurantSetById(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{
		{
			Key: "_id", Value: oid,
		},
	}

	res, err := m.coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("no document found")
	}

	return nil
}

func (m *MongoRepository) GetRestaurantSetByDomainName(ctx context.Context, domainName string) (models.RestaurantSet, error) {
	filter := bson.D{
		{Key: "domain_name", Value: domainName},
	}

	res := models.RestaurantSet{}
	if err := m.coll.FindOne(ctx, filter).Decode(&res); err != nil {
		return models.RestaurantSet{}, err
	}

	return res, nil
}
