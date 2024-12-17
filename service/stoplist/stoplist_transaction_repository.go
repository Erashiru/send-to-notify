package stoplist

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const collectionStoplistTransactionName = "stoplist_transaction"

type Repository interface {
	InsertStopListTransaction(ctx context.Context, req models.StopListTransaction) error
	GetLastTransactionByStoreIDForProducts(ctx context.Context, storeID, deliveryService string) (models.StopListTransaction, error)
	GetLastTransactionByStoreIDForAttributes(ctx context.Context, storeID, deliveryService string) (models.StopListTransaction, error)
	UpdateByTransactionID(ctx context.Context, transactionId, status string, products models.StopListProducts, attributes models.StopListAttributes) error
}

type MongoRepository struct {
	stoplistTransactionCollection *mongo.Collection
}

func NewStoplistTransactionMongoRepository(db *mongo.Database) (*MongoRepository, error) {
	r := MongoRepository{
		stoplistTransactionCollection: db.Collection(collectionStoplistTransactionName),
	}
	return &r, nil
}

func (r *MongoRepository) InsertStopListTransaction(ctx context.Context, req models.StopListTransaction) error {
	req.CreatedAt = time.Now().UTC()
	req.UpdatedAt = time.Now().UTC()

	_, err := r.stoplistTransactionCollection.InsertOne(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoRepository) GetLastTransactionByStoreIDForProducts(ctx context.Context, storeID, deliveryService string) (models.StopListTransaction, error) {

	filter := bson.M{
		"restaurant_id":         storeID,
		"transactions.delivery": deliveryService,
		"transactions.products": bson.M{"$ne": nil},
	}

	var transaction models.StopListTransaction
	if err := r.stoplistTransactionCollection.FindOne(ctx, filter,
		&options.FindOneOptions{Sort: bson.D{{Key: "created_at", Value: -1}}}).Decode(&transaction); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Info().Msgf("GetLastTransactionByStoreIDForProducts not found transaction for store: %s", storeID)
		}
		return models.StopListTransaction{}, err
	}

	return transaction, nil
}

func (r *MongoRepository) GetLastTransactionByStoreIDForAttributes(ctx context.Context, storeID, deliveryService string) (models.StopListTransaction, error) {
	filter := bson.M{
		"restaurant_id":         storeID,
		"attributes":            bson.M{"$ne": nil},
		"transactions.delivery": deliveryService,
	}

	var transaction models.StopListTransaction
	if err := r.stoplistTransactionCollection.FindOne(ctx, filter,
		&options.FindOneOptions{Sort: bson.D{{Key: "created_at", Value: -1}}}).Decode(&transaction); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Info().Msgf("GetLastTransactionByStoreIDForProducts not found transaction for store: %s", storeID)
		}
		return models.StopListTransaction{}, err
	}
	return transaction, nil
}

func (r *MongoRepository) UpdateByTransactionID(ctx context.Context, transactionId, status string, products models.StopListProducts, attributes models.StopListAttributes) error {
	if transactionId == "" {
		return nil
	}

	objID, err := primitive.ObjectIDFromHex(transactionId)
	if err != nil {
		return err
	}

	filter := bson.M{
		"_id": objID,
	}
	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "transactions.$[elem].status", Value: status},
				{Key: "transactions.$[elem].products", Value: products},
				{Key: "transactions.$[elem].attributes", Value: attributes},
				{Key: "updated_at", Value: time.Now().UTC()},
			},
		},
	}
	arrayFilters := options.UpdateOptions{
		ArrayFilters: &options.ArrayFilters{
			Filters: []interface{}{
				bson.M{"elem.delivery": "yandex"},
			},
		},
	}

	res, err := r.stoplistTransactionCollection.UpdateOne(ctx, filter, update, &arrayFilters)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		log.Info().Msgf("transaction with id: %s not found to update", transactionId)
	}

	return nil
}
