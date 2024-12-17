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

const paymentSystemCustomersCollectionName = "payment_system_customers"

type CustomersRepository interface {
	InsertCustomer(ctx context.Context, customer models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error)
	UpdateCustomer(ctx context.Context, customer models.PaymentSystemCustomer) error
	GetCustomerByEmail(ctx context.Context, email string) (models.PaymentSystemCustomer, error)
	GetCustomerByPaymentSystemCustomerID(ctx context.Context, paymentSystemCustomerID string) (models.PaymentSystemCustomer, error)
}

type CustomersMongoRepository struct {
	collection *mongo.Collection
}

func NewCustomersMongoRepository(db *mongo.Database) (*CustomersMongoRepository, error) {
	r := CustomersMongoRepository{
		collection: db.Collection(paymentSystemCustomersCollectionName),
	}
	return &r, nil
}
func (r *CustomersMongoRepository) GetCustomerByPaymentSystemCustomerID(ctx context.Context, paymentSystemCustomerID string) (models.PaymentSystemCustomer, error) {
	filter := bson.D{
		{Key: "id", Value: paymentSystemCustomerID},
	}

	var customer models.PaymentSystemCustomer

	if err := r.collection.FindOne(ctx, filter).Decode(&customer); err != nil {
		return models.PaymentSystemCustomer{}, err
	}

	return customer, nil
}

func (r *CustomersMongoRepository) GetCustomerByEmail(ctx context.Context, email string) (models.PaymentSystemCustomer, error) {
	filter := bson.D{
		{Key: "email", Value: email},
	}

	var customer models.PaymentSystemCustomer

	if err := r.collection.FindOne(ctx, filter).Decode(&customer); err != nil {
		return models.PaymentSystemCustomer{}, err
	}

	return customer, nil
}

func (r *CustomersMongoRepository) InsertCustomer(ctx context.Context, customer models.PaymentSystemCustomer) (models.PaymentSystemCustomer, error) {
	customer.CreatedAt = time.Now().UTC()
	res, err := r.collection.InsertOne(ctx, customer)
	if err != nil {
		return models.PaymentSystemCustomer{}, customeErrors.ErrorSwitch(err)
	}

	customer.ExternalID = res.InsertedID.(primitive.ObjectID).Hex()

	return customer, nil
}

func (r *CustomersMongoRepository) UpdateCustomer(ctx context.Context, customer models.PaymentSystemCustomer) error {
	oid, err := primitive.ObjectIDFromHex(customer.ExternalID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	customer.ExternalID = ""
	customer.UpdatedAt = time.Now().UTC()

	update := bson.D{
		{Key: "$set", Value: customer},
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
