package mongodb

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/auth/database/datastore/drivers"
	"github.com/kwaaka-team/orders-core/core/auth/models"
	"github.com/kwaaka-team/orders-core/core/auth/models/selector"
	"github.com/kwaaka-team/orders-core/core/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	collection *mongo.Collection
}

func NewUserRepository(userCollection *mongo.Collection) drivers.UserRepository {
	return &UserRepo{
		collection: userCollection,
	}
}

func (u *UserRepo) CreateUser(ctx context.Context, user models.User) error {
	if _, err := u.collection.InsertOne(ctx, user); err != nil {
		return errors.ErrorSwitch(err)
	}

	return nil
}

func (u *UserRepo) GetUser(ctx context.Context, query selector.User) (models.User, error) {
	filter := u.filterFrom(query)

	res := u.collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return models.User{}, errors.ErrorSwitch(err)
	}

	var user models.User
	if err := res.Decode(&user); err != nil {
		return models.User{}, errors.ErrorSwitch(err)
	}

	return user, nil
}

func (u *UserRepo) UpdateUserInfo(ctx context.Context, user models.User) error {
	filter := u.filterFrom(selector.NewEmptyUser().SetUID(user.UID))

	update := bson.M{
		"$set": bson.M{
			"fcm_token": user.FCMToken,
		},
	}

	res, err := u.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.ErrorSwitch(err)
	}

	if res.MatchedCount == 0 {
		return errors.ErrNotFound
	}
	return nil
}

func (u *UserRepo) filterFrom(query selector.User) bson.D {
	var result bson.D

	if query.HasUID() {
		result = append(result, bson.E{
			Key:   "uid",
			Value: query.UID,
		})
	}

	if query.HasName() {
		result = append(result, bson.E{
			Key:   "name",
			Value: query.Name,
		})
	}

	if query.HasPhoneNumber() {
		result = append(result, bson.E{
			Key:   "phone_number",
			Value: query.PhoneNumber,
		})
	}

	return result
}
