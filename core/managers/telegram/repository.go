package telegram

import (
	"context"
	errorsGo "errors"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	telegramCustomersCollection = "telegram_customers"
	ordersCollection            = "orders"

	NEW         = "new"
	WAIT_REVIEW = "wait_review"
	REVIEWED    = "reviewed"
)

type User struct {
	FirstName             string  `bson:"first_name"`
	ChatId                string  `bson:"chat_id"`
	Status                string  `bson:"status"`
	ReviewingOrderID      *string `bson:"reviewing_order_id,omitempty"`
	ReviewingRestaurantID *string `bson:"reviewing_restaurant_id,omitempty"`
}

type Repository struct {
	tgCustomersCollection *mongo.Collection
	ordersCollection      *mongo.Collection
}

func NewTelegramRepo(db *mongo.Database) *Repository {
	return &Repository{
		tgCustomersCollection: db.Collection(telegramCustomersCollection),
		ordersCollection:      db.Collection(ordersCollection),
	}
}

func (r *Repository) GetTelegramReviewRatingFromOrder(ctx context.Context, orderID string) (float32, error) {
	var order models.Order

	oid, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return 0, errors.ErrorSwitch(err)
	}

	filter := bson.D{{Key: "_id", Value: oid}}

	err = r.ordersCollection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		log.Err(err).Msgf("error decoding order: %s", orderID)
		return 0, errors.ErrorSwitch(err)
	}
	if order.Review.Rating == 0 {
		return 0, errorsGo.New("no review rating")
	}

	return order.Review.Rating, nil
}

func (r *Repository) SaveTelegramReviewRating(ctx context.Context, orderID string, rating float32) error {
	oid, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return errors.ErrorSwitch(err)
	}

	filter := bson.D{{Key: "_id", Value: oid}}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "review.rating", Value: rating}}}}

	opts := options.Update().SetUpsert(true)

	res, err := r.ordersCollection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return errors.ErrorSwitch(err)
	}
	if res.MatchedCount == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (r *Repository) GetReviewingOrderID(ctx context.Context, chatID string) (string, error) {
	var user User

	filter := bson.D{{Key: "chat_id", Value: chatID}}

	err := r.tgCustomersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return "", err
	}
	if user.ReviewingOrderID == nil {
		return "", errorsGo.New("error getting reviewing order ID (it's empty)")
	}

	return *user.ReviewingOrderID, nil
}

func (r *Repository) GetReviewingRestaurantID(ctx context.Context, chatID string) (string, error) {
	var user User

	filter := bson.D{{Key: "chat_id", Value: chatID}}

	err := r.tgCustomersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return "", err
	}
	if user.ReviewingRestaurantID == nil {
		return "", errorsGo.New("error getting reviewing restaurant ID (it's empty)")
	}

	return *user.ReviewingRestaurantID, nil
}

func (r *Repository) UpdateReviewingOrderInfo(ctx context.Context, chatID, orderID, restID string) error {
	filter := bson.D{{Key: "chat_id", Value: chatID}}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "reviewing_order_id", Value: orderID}, {Key: "reviewing_restaurant_id", Value: restID}}}}

	res, err := r.tgCustomersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *Repository) GetUserStatus(ctx context.Context, chatID string) (string, error) {
	var user User

	filter := bson.D{{Key: "chat_id", Value: chatID}}

	err := r.tgCustomersCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return "", err
	}

	return user.Status, nil
}

func (r *Repository) UpdateUserStatus(ctx context.Context, chatID string, status string) error {
	filter := bson.M{"chat_id": chatID}
	update := bson.M{"$set": bson.M{"status": status}}

	res, err := r.tgCustomersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *Repository) InsertUser(ctx context.Context, firstName string, chatID string) error {
	user := User{
		FirstName: firstName,
		ChatId:    chatID,
		Status:    NEW,
	}

	_, err := r.tgCustomersCollection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetUserChatID(ctx context.Context) (string, error) {
	// TODO: hardcoded now, but should implement logic of getting user's chat id (probably would be stored in 'order' struct)
	return "", nil
}
