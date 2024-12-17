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
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TapRestaurantRepository struct {
	collection *mongo.Collection
}

func NewTapRestaurantRepository(collection *mongo.Collection) drivers2.TapRestaurantRepository {
	return &TapRestaurantRepository{
		collection: collection,
	}
}

func (tr *TapRestaurantRepository) Create(ctx context.Context, req models2.TapRestaurant) (string, error) {
	res, err := tr.collection.InsertOne(ctx, req)
	if err != nil {
		return "", errorSwitch(err)
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (tr *TapRestaurantRepository) GetList(ctx context.Context, query selector.TapRestaurant) ([]models2.TapRestaurant, int, error) {
	filter, err := tr.filterFrom(query)
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find()
	opts.SetSort(tr.sortFrom(query.Sorting))
	if query.HasPagination() {
		opts.SetSkip(query.Skip()).SetLimit(query.Pagination.Limit)
	}

	count, err := tr.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, errorSwitch(err)
	}

	cur, err := tr.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, errorSwitch(err)
	}

	defer closeCur(cur)

	tapRestaurants := make([]models2.TapRestaurant, cur.RemainingBatchLength())
	if err = cur.All(ctx, &tapRestaurants); err != nil {
		return nil, 0, errorSwitch(err)
	}

	return tapRestaurants, int(count), nil
}

func (tr *TapRestaurantRepository) GetByID(ctx context.Context, id string) (models2.TapRestaurant, error) {
	var tapRest models2.TapRestaurant
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return tapRest, errors.Wrap(drivers2.ErrInvalid, "ID error")
	}
	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	err = tr.collection.FindOne(ctx, filter).Decode(&tapRest)
	if err != nil {
		return tapRest, errorSwitch(err)
	}
	return tapRest, nil
}

func (tr *TapRestaurantRepository) GetByQuery(ctx context.Context, query selector.TapRestaurant) (models2.TapRestaurant, error) {
	var tapRest models2.TapRestaurant
	filter, err := tr.filterFrom(query)
	if err != nil {
		return tapRest, err
	}

	err = tr.collection.FindOne(ctx, filter).Decode(&tapRest)
	if err != nil {
		return tapRest, err
	}

	return tapRest, nil
}

func (tr *TapRestaurantRepository) Update(ctx context.Context, req models2.UpdateTapRestaurant) error {
	if req.ID == nil {
		return errors.Wrap(drivers2.ErrNotFound, "id is nil")
	}

	oid, err := primitive.ObjectIDFromHex(*req.ID)
	if err != nil {
		return errors.Wrap(drivers2.ErrInvalid, "id error")
	}

	filter := bson.D{{Key: "_id", Value: oid}}

	update, err := tr.setFields(req)
	if err != nil {
		return err
	}

	res, err := tr.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.Wrap(drivers2.ErrNotFound, "not found error")
	}

	return nil
}

func (tr *TapRestaurantRepository) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.Wrap(drivers2.ErrNotFound, "not found error")
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	_, err = tr.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (tr *TapRestaurantRepository) setFields(req models2.UpdateTapRestaurant) (bson.D, error) {
	update := make(bson.D, 0, 3)

	if req.Name != nil {
		if *req.Name != "" {
			update = append(update, bson.E{
				Key:   "name",
				Value: *req.Name,
			})
		}
	}

	if req.Description != nil {
		if *req.Description != "" {
			update = append(update, bson.E{
				Key:   "description",
				Value: *req.Description,
			})
		}
	}

	if req.Image != nil {
		if *req.Image != "" {
			update = append(update, bson.E{
				Key:   "img",
				Value: *req.Image,
			})
		}
	}

	if req.QRMenuLink != nil {
		if *req.QRMenuLink != "" {
			update = append(update, bson.E{
				Key:   "qr_menu_link",
				Value: *req.QRMenuLink,
			})
		}
	}

	if req.Tel != nil {
		if *req.Tel != "" {
			update = append(update, bson.E{
				Key:   "tel",
				Value: *req.Tel,
			})
		}
	}

	if req.Instagram != nil {
		if *req.Instagram != "" {
			update = append(update, bson.E{
				Key:   "instagram",
				Value: *req.Instagram,
			})
		}
	}

	if req.Website != nil {
		if *req.Website != "" {
			update = append(update, bson.E{
				Key:   "website",
				Value: *req.Website,
			})
		}
	}

	result := bson.D{
		{Key: "$set", Value: update},
	}

	return result, nil
}

func (tr *TapRestaurantRepository) sortFrom(query selector.Sorting) bson.D {
	sort := make(bson.D, 0, 2)

	if query.HasSorting() {
		sort = append(sort, bson.E{Key: query.Param, Value: query.Direction})
	}
	sort = append(sort, bson.E{Key: "_id", Value: 1})

	return sort
}

func (tr *TapRestaurantRepository) filterFrom(query selector.TapRestaurant) (bson.D, error) {
	filter := make(bson.D, 0, 3)

	if query.HasName() {
		filter = append(filter, bson.E{
			Key:   "name",
			Value: query.Name,
		})
	}
	return filter, nil
}
