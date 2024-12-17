package mongo

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MenuUploadTransaction struct {
	mtColl *mongo.Collection
}

func NewMenuUploadTransaction(mtColl *mongo.Collection) *MenuUploadTransaction {
	return &MenuUploadTransaction{
		mtColl: mtColl,
	}
}

func (repo *MenuUploadTransaction) Get(ctx context.Context, query selector.MenuUploadTransaction) (models.MenuUploadTransaction, error) {

	filter, err := repo.filterFrom(query)
	if err != nil {
		return models.MenuUploadTransaction{}, err
	}

	opts := options.FindOne()
	opts.SetSort(repo.sortFrom(query.Sorting))

	var res models.MenuUploadTransaction
	if err = repo.mtColl.FindOne(ctx, filter, opts).Decode(&res); err != nil {
		return models.MenuUploadTransaction{}, errorSwitch(err)
	}

	return res, nil
}

func (repo *MenuUploadTransaction) List(ctx context.Context, query selector.MenuUploadTransaction) ([]models.MenuUploadTransaction, int64, error) {

	filter, err := repo.filterFrom(query)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find()
	if query.HasPagination() {
		opts.SetSkip(query.Skip()).SetLimit(query.Limit)
	}
	opts.SetSort(repo.sortFrom(query.Sorting))

	count, err := repo.mtColl.CountDocuments(ctx, filter)
	if err != nil || count == 0 {
		return nil, 0, errorSwitch(err)
	}

	cur, err := repo.mtColl.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, errorSwitch(err)
	}
	defer closeCur(cur)

	res := make([]models.MenuUploadTransaction, 0, cur.RemainingBatchLength())
	if err = cur.All(ctx, &res); err != nil {
		return nil, 0, errorSwitch(err)
	}

	return res, count, nil
}

func (repo *MenuUploadTransaction) Insert(ctx context.Context, req models.MenuUploadTransaction) (string, error) {

	req.CreatedAt.Value = coreModels.TimeNow()
	req.UpdatedAt.Value = coreModels.TimeNow()

	res, err := repo.mtColl.InsertOne(ctx, req)
	if err != nil {
		return "", errorSwitch(err)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", drivers.ErrInvalid
	}

	return oid.Hex(), nil
}

func (repo *MenuUploadTransaction) Update(ctx context.Context, req models.MenuUploadTransaction) error {
	oid, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return drivers.ErrInvalid
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	req.UpdatedAt.Value = coreModels.TimeNow()

	update := bson.D{
		{Key: "$set", Value: models.ToUpdateMenuTransactions(req)},
	}

	res, err := repo.mtColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return errorSwitch(err)
	}

	if res.MatchedCount == 0 {
		return drivers.ErrNotFound
	}

	return nil
}

func (repo *MenuUploadTransaction) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return drivers.ErrInvalid
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}
	_, err = repo.mtColl.DeleteOne(ctx, filter)
	if err != nil {
		return errorSwitch(err)
	}

	return nil
}

func (repo *MenuUploadTransaction) sortFrom(query selector.Sorting) bson.D {
	sort := make(bson.D, 0, 2)

	if query.HasSorting() {
		sort = append(sort, bson.E{Key: query.Param, Value: query.Direction})
	}
	sort = append(sort, bson.E{Key: "_id", Value: 1})

	return sort
}

func (repo *MenuUploadTransaction) filterFrom(query selector.MenuUploadTransaction) (bson.D, error) {
	filter := make(bson.D, 0, 8)

	if query.HasID() {
		oid, err := primitive.ObjectIDFromHex(query.ID)
		if err != nil {
			return nil, drivers.ErrInvalid
		}
		filter = append(filter, bson.E{
			Key: "_id", Value: oid,
		})
	}

	if query.HasMenuID() {
		filter = append(filter, bson.E{
			Key: "ext_transactions.menu_id", Value: query.MenuID,
		})
	}

	if query.HasExtTransactionID() {
		filter = append(filter, bson.E{
			Key: "ext_transactions.id", Value: query.ExtTransactionID,
		})
	}

	if query.HasStoreID() {
		//oid, err := primitive.ObjectIDFromHex(query.StoreID)
		//if err != nil {
		//	return nil, drivers.ErrInvalid
		//}
		filter = append(filter, bson.E{
			Key: "restaurant_id", Value: query.StoreID,
		})
	}

	if query.HasService() {
		filter = append(filter, bson.E{
			Key: "service", Value: query.Service,
		})
	}

	if query.HasStatus() {
		filter = append(filter, bson.E{
			Key: "status", Value: query.Status,
		})
	}

	if query.HasCreatedFrom() {
		filter = append(filter, bson.E{
			Key: "created_at.value",
			Value: bson.M{
				"$gte": primitive.NewDateTimeFromTime(query.CreatedFrom),
			},
		})
	}

	if query.HasCreatedTo() {
		filter = append(filter, bson.E{
			Key: "created_at.value",
			Value: bson.M{
				"$lte": primitive.NewDateTimeFromTime(query.CreatedTo),
			},
		})
	}

	return filter, nil
}
