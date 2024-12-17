package mongodb

import (
	"context"
	drivers2 "github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	models2 "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StoreGroupRepo struct {
	collection *mongo.Collection
}

func NewStoreGroupRepository(collection *mongo.Collection) drivers2.StoreGroupRepository {
	return &StoreGroupRepo{
		collection: collection,
	}
}

func (s *StoreGroupRepo) Get(ctx context.Context, query selector.StoreGroup) (models2.StoreGroup, error) {
	var storeGroup models2.StoreGroup

	filter, err := filterFrom(query)
	if err != nil {
		return storeGroup, err
	}

	err = s.collection.FindOne(ctx, filter).Decode(&storeGroup)
	if err != nil {
		return storeGroup, errorSwitch(err)
	}

	return storeGroup, nil
}

func filterFrom(query selector.StoreGroup) (bson.D, error) {
	filter := make(bson.D, 0, 2)

	if query.HasID() {
		oid, err := primitive.ObjectIDFromHex(query.ID)
		if err != nil {
			return nil, errors.Wrap(drivers2.ErrInvalid, "query.ID error")
		}
		filter = append(filter, bson.E{Key: "_id", Value: oid})
	}

	if query.HasName() {
		filter = append(filter, bson.E{Key: "name", Value: query.Name})
	}

	if query.HasStoreIDs() {
		filter = append(filter, bson.E{Key: "restaurant_ids", Value: bson.M{"$in": query.StoreIDs}})
	}

	if query.HasCountry() {
		filter = append(filter, bson.E{Key: "country", Value: query.Country})
	}

	if query.HasCategory() {
		filter = append(filter, bson.E{Key: "category", Value: query.Category})
	}

	if query.HasStatus() {
		filter = append(filter, bson.E{Key: "status", Value: query.Status})
	}

	if query.HasDomainName() {
		filter = append(filter, bson.E{Key: "domain_name", Value: query.DomainName})
	}

	return filter, nil
}
func (s *StoreGroupRepo) Create(ctx context.Context, storeGroup models2.StoreGroup) (string, error) {
	res, err := s.collection.InsertOne(ctx, storeGroup)
	if err != nil {
		return "", errorSwitch(err)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", errors.Wrap(drivers2.ErrInvalid, "oid error")
	}

	return oid.Hex(), nil
}

func (s *StoreGroupRepo) UpdateByFields(ctx context.Context, storeGroup models2.UpdateStoreGroup) (int64, error) {

	if storeGroup.ID == nil {
		return 0, errors.Wrap(drivers2.ErrInvalid, "storeGroup.ID is nil")
	}
	oid, err := primitive.ObjectIDFromHex(*storeGroup.ID)
	if err != nil {
		return 0, errors.Wrap(drivers2.ErrNotFound, "not found error")
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update, err := s.setFields(storeGroup)
	if err != nil {
		return 0, errorSwitch(err)
	}

	res, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, errorSwitch(err)
	}

	if res.MatchedCount == 0 {
		return 0, errors.Wrap(drivers2.ErrNotFound, "not found error")
	}
	return res.ModifiedCount, nil
}

func (s *StoreGroupRepo) setFields(storeGroup models2.UpdateStoreGroup) (bson.D, error) {
	update := make(bson.D, 0, 2)

	if storeGroup.Name != nil {
		update = append(update, bson.E{Key: "name", Value: storeGroup.Name})
	}

	if storeGroup.StoreIds != nil {
		update = append(update, bson.E{
			Key:   "restaurant_ids",
			Value: storeGroup.StoreIds,
		})
	}

	if storeGroup.RetryCount != nil {
		update = append(update, bson.E{
			Key:   "retry_count",
			Value: storeGroup.RetryCount,
		})
	}

	if storeGroup.ColumnView != nil {
		update = append(update, bson.E{
			Key:   "column_view",
			Value: storeGroup.ColumnView,
		})
	}

	result := bson.D{{Key: "$set", Value: update}}

	return result, nil
}

func (s *StoreGroupRepo) List(ctx context.Context, query selector.StoreGroup) ([]models2.StoreGroup, error) {
	return []models2.StoreGroup{}, nil
}

func (s *StoreGroupRepo) All(ctx context.Context) ([]models2.StoreGroup, error) {
	var results []models2.StoreGroup

	cur, err := s.collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		return nil, errorSwitch(err)
	}

	for cur.Next(context.TODO()) {
		//Create a value into which the single document can be decoded
		var elem models2.StoreGroup
		err := cur.Decode(&elem)
		if err != nil {
			return nil, errorSwitch(err)
		}

		results = append(results, elem)
	}
	return results, nil
}
