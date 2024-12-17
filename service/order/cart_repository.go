package order

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const directCartCollectionName = "qr_menu_carts"
const kwaakaAdminCartCollectionName = "kwaaka_admin_carts"

type CartRepository interface {
	GetQRMenuCartByID(ctx context.Context, cartID string) (models.Cart, error)
	GetKwaakaAdminCartByCartID(ctx context.Context, cartID string) (models.Cart, error)
	GetOldQRMenuCartByID(ctx context.Context, cartID string) (models.OldCart, error)
}

type CartMongoRepository struct {
	qrMenuRepo      *mongo.Collection
	kwaakaAdminRepo *mongo.Collection
}

func NewCartRepository(db *mongo.Database) *CartMongoRepository {
	return &CartMongoRepository{
		qrMenuRepo:      db.Collection(directCartCollectionName),
		kwaakaAdminRepo: db.Collection(kwaakaAdminCartCollectionName),
	}
}

func (r *CartMongoRepository) GetQRMenuCartByID(ctx context.Context, cartID string) (models.Cart, error) {
	oid, err := primitive.ObjectIDFromHex(cartID)
	if err != nil {
		return models.Cart{}, err
	}

	filter := bson.M{"_id": oid}

	var cart models.Cart
	if err := r.qrMenuRepo.FindOne(ctx, filter).Decode(&cart); err != nil {
		return models.Cart{}, errors.ErrorSwitch(err)
	}

	return cart, nil
}

func (r *CartMongoRepository) GetKwaakaAdminCartByCartID(ctx context.Context, cartID string) (models.Cart, error) {
	oid, err := primitive.ObjectIDFromHex(cartID)
	if err != nil {
		return models.Cart{}, err
	}
	filter := bson.M{"_id": oid}

	var cart models.Cart
	if err := r.kwaakaAdminRepo.FindOne(ctx, filter).Decode(&cart); err != nil {
		return models.Cart{}, errors.ErrorSwitch(err)
	}

	return cart, nil
}

func (r *CartMongoRepository) GetOldQRMenuCartByID(ctx context.Context, cartID string) (models.OldCart, error) {
	oid, err := primitive.ObjectIDFromHex(cartID)
	if err != nil {
		return models.OldCart{}, err
	}
	filter := bson.M{"_id": oid}

	var cart models.OldCart
	if err := r.qrMenuRepo.FindOne(ctx, filter).Decode(&cart); err != nil {
		return models.OldCart{}, errors.ErrorSwitch(err)
	}
	return cart, nil
}
