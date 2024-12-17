package repository

import (
	"context"
	customeErrors "github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/service/payment/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	paymentsCollectionName = "payments"
	maxNotificationCount   = 1
)

type PaymentsRepository interface {
	InsertPaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error)
	UpdatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) error
	UpdatePaymentOrderStatusHistory(ctx context.Context, paymentOrder models.PaymentOrder) error
	GetPaymentOrderByPaymentOrderID(ctx context.Context, paymentOrderID string) (models.PaymentOrder, error)
	GetPaymentOrderByOrderID(ctx context.Context, cartID string) (models.PaymentOrder, error)
	GetUnpaidPayments(ctx context.Context, minutes int) ([]models.PaymentOrder, error)
	SetNotificationCount(ctx context.Context, id string, count int) error
	GetUnpaidPaymentsByPaymentSystem(ctx context.Context, minutes int, paymentSystem string) ([]models.PaymentOrder, error)
}

type PaymentsMongoRepository struct {
	collection *mongo.Collection
}

func NewPaymentsMongoRepository(db *mongo.Database) (*PaymentsMongoRepository, error) {
	r := PaymentsMongoRepository{
		collection: db.Collection(paymentsCollectionName),
	}
	return &r, nil
}

func (r *PaymentsMongoRepository) InsertPaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) (models.PaymentOrder, error) {
	paymentOrder.CreatedAt = time.Now().UTC()

	res, err := r.collection.InsertOne(ctx, paymentOrder)
	if err != nil {
		return models.PaymentOrder{}, err
	}

	paymentOrder.ExternalID = res.InsertedID.(primitive.ObjectID).Hex()

	return paymentOrder, nil
}

func (r *PaymentsMongoRepository) GetPaymentOrderByPaymentOrderID(ctx context.Context, paymentOrderID string) (models.PaymentOrder, error) {
	filter := bson.D{
		{Key: "payment_order_id", Value: paymentOrderID},
	}

	var res models.PaymentOrder

	if err := r.collection.FindOne(ctx, filter).Decode(&res); err != nil {
		return models.PaymentOrder{}, customeErrors.ErrorSwitch(err)
	}

	return res, nil
}

func (r *PaymentsMongoRepository) UpdatePaymentOrder(ctx context.Context, paymentOrder models.PaymentOrder) error {
	oid, err := primitive.ObjectIDFromHex(paymentOrder.ExternalID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	paymentOrder.ExternalID = ""
	paymentOrder.UpdatedAt = time.Now().UTC()

	update := bson.D{
		{Key: "$set", Value: paymentOrder},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return customeErrors.ErrorSwitch(err)
	}

	if res.MatchedCount == 0 {
		return errors.New("matched count is 0")
	}

	return nil
}

func (r *PaymentsMongoRepository) UpdatePaymentOrderStatusHistory(ctx context.Context, paymentOrder models.PaymentOrder) error {
	oid, err := primitive.ObjectIDFromHex(paymentOrder.ExternalID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	paymentOrder.ExternalID = ""
	paymentOrder.UpdatedAt = time.Now().UTC()

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "status", Value: paymentOrder.PaymentOrderStatus},
			{Key: "updated_at", Value: paymentOrder.UpdatedAt},
		}},
		{Key: "$push", Value: bson.D{
			{Key: "payment_order_status_history", Value: models.StatusHistory{
				Status: paymentOrder.PaymentOrderStatus,
				Time:   time.Now().UTC(),
			}},
		}},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return customeErrors.ErrorSwitch(err)
	}

	if res.MatchedCount == 0 {
		return errors.New("matched count is 0")
	}

	return nil
}

func (r *PaymentsMongoRepository) GetPaymentOrderByOrderID(ctx context.Context, cartID string) (models.PaymentOrder, error) {
	filter := bson.D{
		{Key: "order_id", Value: cartID},
	}

	var res models.PaymentOrder

	if err := r.collection.FindOne(ctx, filter).Decode(&res); err != nil {
		return models.PaymentOrder{}, customeErrors.ErrorSwitch(err)
	}

	return res, nil
}

func (r *PaymentsMongoRepository) GetUnpaidPayments(ctx context.Context, minutes int) ([]models.PaymentOrder, error) {
	var payments []models.PaymentOrder

	timeBefore := time.Now().Add(time.Duration(-minutes) * time.Minute).UTC()
	filter := bson.D{
		{Key: "status", Value: models.UNPAID},
		{Key: "created_at", Value: bson.D{
			{Key: "$lte", Value: timeBefore},
			{Key: "$gte", Value: timeBefore.Add(time.Duration(-1) * time.Hour)},
		}},
		{Key: "notification_count", Value: bson.D{
			{Key: "$lt", Value: maxNotificationCount},
		}},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, customeErrors.ErrorSwitch(err)
	}
	defer customeErrors.CloseCur(cursor)

	for cursor.Next(ctx) {
		var payment models.PaymentOrder
		if err = cursor.Decode(&payment); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	if err = cursor.Err(); err != nil {
		return nil, customeErrors.ErrorSwitch(err)
	}

	return payments, nil
}

func (r *PaymentsMongoRepository) GetUnpaidPaymentsByPaymentSystem(ctx context.Context, minutes int, paymentSystem string) ([]models.PaymentOrder, error) {
	var payments []models.PaymentOrder

	timeBefore := time.Now().Add(time.Duration(-minutes) * time.Minute).UTC()

	filter := bson.D{
		{Key: "payment_system", Value: paymentSystem},
		{Key: "status", Value: models.UNPAID},
		{Key: "created_at", Value: bson.D{
			{Key: "$gte", Value: timeBefore},
		}},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, customeErrors.ErrorSwitch(err)
	}
	defer customeErrors.CloseCur(cursor)

	for cursor.Next(ctx) {
		var payment models.PaymentOrder
		if err = cursor.Decode(&payment); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	if err = cursor.Err(); err != nil {
		return nil, customeErrors.ErrorSwitch(err)
	}

	return payments, nil
}

func (r *PaymentsMongoRepository) SetNotificationCount(ctx context.Context, id string, count int) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "notification_count", Value: count}}}}
	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return customeErrors.ErrorSwitch(err)
	}
	if res.ModifiedCount == 0 {
		return customeErrors.ErrNotFound
	}

	return nil
}
