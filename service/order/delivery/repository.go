package delivery

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/selector"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const delivery3plCollectionName = "3pl_deliveries"

type Repository interface {
	GetDeliveries(ctx context.Context, deliveryIDs []string) ([]models.Delivery3plOrder, error)
	GetDeliveryByDeliveryID(ctx context.Context, deliveryID string) (models.Delivery3plOrder, error)
	GetAllDeliveries(ctx context.Context, query selector.Delivery3plOrder) ([]models.Delivery3plOrder, error)
	InsertChangeHistory(ctx context.Context, deliveryID, username, action string) error
}

type DeliveryMongoRepository struct {
	collection *mongo.Collection
}

func NewDeliveryMongoRepository(db *mongo.Database) (*DeliveryMongoRepository, error) {
	return &DeliveryMongoRepository{
		collection: db.Collection(delivery3plCollectionName),
	}, nil
}

func (r *DeliveryMongoRepository) GetDeliveries(ctx context.Context, deliveryIDs []string) ([]models.Delivery3plOrder, error) {

	var objIDs []primitive.ObjectID
	for _, deliveryID := range deliveryIDs {
		oid, err := primitive.ObjectIDFromHex(deliveryID)
		if err != nil {
			log.Info().Msgf("error: convert delivery id: %s to ObjectID", deliveryID)
			continue
		}
		objIDs = append(objIDs, oid)
	}

	filter := bson.D{
		{Key: "_id", Value: bson.D{
			{Key: "$in", Value: objIDs}},
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.Info().Msg("error: find document from DB")
		return nil, err
	}

	deliveries := make([]models.Delivery3plOrder, 0, cursor.RemainingBatchLength())
	if err = cursor.All(ctx, &deliveries); err != nil {
		log.Info().Msg("error: convert data to models.Delivery3plOrder")
		return nil, err
	}

	return deliveries, nil
}

func (r *DeliveryMongoRepository) GetDeliveryByDeliveryID(ctx context.Context, deliveryID string) (models.Delivery3plOrder, error) {

	oid, err := primitive.ObjectIDFromHex(deliveryID)
	if err != nil {
		log.Err(err).Msgf("error: convert string delivery id: %s to ObjectID()", deliveryID)
		return models.Delivery3plOrder{}, err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	var delivery models.Delivery3plOrder
	if err := r.collection.FindOne(ctx, filter).Decode(&delivery); err != nil {
		log.Err(err).Msg("error: find delivery from DB")
		return models.Delivery3plOrder{}, err
	}

	return delivery, nil
}

func (r *DeliveryMongoRepository) GetAllDeliveries(ctx context.Context, query selector.Delivery3plOrder) ([]models.Delivery3plOrder, error) {
	filter := r.filterFrom(query)

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	deliveries := make([]models.Delivery3plOrder, 0, cursor.RemainingBatchLength())
	if err := cursor.All(ctx, &deliveries); err != nil {
		return nil, err
	}

	return deliveries, nil
}

func (r *DeliveryMongoRepository) filterFrom(query selector.Delivery3plOrder) bson.D {
	result := make(bson.D, 0, 9)

	if query.HasStatus() {
		result = append(result, bson.E{
			Key:   "status",
			Value: query.Status,
		})
	}

	if query.HasUpdatedTimeTo() {
		result = append(result, bson.E{
			Key: "updated_at",
			Value: bson.M{
				"$lte": primitive.NewDateTimeFromTime(query.UpdatedTimeTo),
			},
		})
	}

	if query.HasCreatedTimeFrom() {
		result = append(result, bson.E{
			Key: "created_at",
			Value: bson.M{
				"$gte": primitive.NewDateTimeFromTime(query.CreatedTimeFrom),
			},
		})
	}

	if query.HasCreatedTimeTo() {
		result = append(result, bson.E{
			Key: "created_at",
			Value: bson.M{
				"$lte": primitive.NewDateTimeFromTime(query.CreatedTimeTo),
			},
		})
	}

	return result
}

func (r *DeliveryMongoRepository) InsertChangeHistory(ctx context.Context, deliveryID, username, action string) error {

	oid, err := primitive.ObjectIDFromHex(deliveryID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "changes_history", Value: models.ChangesHistory{
				Username:  username,
				Action:    action,
				UpdatedAt: time.Now(),
			}},
		}},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
