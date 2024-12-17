package mongodb

import (
	"context"
	drivers2 "github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserStoreRepo struct {
	collection *mongo.Collection
}

func NewUserStoreRepository(collection *mongo.Collection) drivers2.UserStoreRepository {
	return &UserStoreRepo{
		collection: collection,
	}
}

func (u *UserStoreRepo) Insert(ctx context.Context, userStores []models.UserStore) error {
	for _, userStore := range userStores {
		_, err := u.collection.InsertOne(
			context.TODO(),
			userStore,
			options.InsertOne().SetBypassDocumentValidation(false),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *UserStoreRepo) FindUsers(ctx context.Context, user selector.User) ([]models.UserStore, error) {
	filter, err := u.filterFrom(user)
	if err != nil {
		return nil, err
	}

	cur, err := u.collection.Find(ctx, filter)
	if err != nil {
		return nil, errorSwitch(err)
	}

	var users []models.UserStore
	for cur.Next(ctx) {
		var store models.UserStore
		if err := cur.Decode(&store); err != nil {
			return nil, errorSwitch(err)
		}

		users = append(users, store)
	}

	return users, nil
}

func (u *UserStoreRepo) filterFrom(query selector.User) (bson.D, error) {
	var result bson.D

	if query.HasID() {
		oid, err := primitive.ObjectIDFromHex(query.ID)
		if err != nil {
			return nil, errors.Wrap(drivers2.ErrInvalid, "query.ID error")
		}
		result = append(result, bson.E{
			Key:   "_id",
			Value: oid,
		})
	}
	if query.HasStoreGroupId() {
		result = append(result, bson.E{
			Key:   "restaurant_group_id",
			Value: query.StoreGroupId,
		})
	}
	if query.HasUsername() {
		result = append(result, bson.E{
			Key:   "username",
			Value: query.Username,
		})
	}
	if query.HasStoreID() {
		result = append(result, bson.E{
			Key:   "restaurant_id",
			Value: query.StoreId,
		})
	}

	if query.SendNotification {
		result = append(result, bson.E{
			Key:   "send_notification",
			Value: true,
		})
	}
	return result, nil
}

func (u *UserStoreRepo) Delete(ctx context.Context, user selector.User) error {
	filter, err := u.filterFrom(user)
	if err != nil {
		return err
	}

	_, err = u.collection.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserStoreRepo) UpdateUserOrderNotifications(ctx context.Context, username, fcmToken string, stores []string) error {

	update := bson.M{
		"$set": bson.M{
			"send_notification": true,
		},
	}
	if fcmToken != "" {
		update["$addToSet"] = bson.M{
			"fcm_tokens": fcmToken,
		}
	}

	_, err := u.collection.UpdateMany(ctx, bson.M{
		"username": bson.M{"$eq": username},
		//"restaurant_id": bson.M{"$in": stores},
	}, update)
	if err != nil {
		return err
	}

	//_, err = u.collection.UpdateMany(ctx, bson.M{
	//	"username":      bson.M{"$eq": username},
	//	"restaurant_id": bson.M{"$nin": stores},
	//}, bson.M{
	//	"$set": bson.M{"send_notification": false},
	//})

	if err != nil {
		return err
	}
	return nil
}
