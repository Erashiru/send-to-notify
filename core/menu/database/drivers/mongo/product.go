package mongo

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (repo *MenuRepository) GetProduct(ctx context.Context, query selector.Menu) (models.Product, error) {

	filter, err := repo.filterFrom(query)
	if err != nil {
		return models.Product{}, err
	}

	var res models.Product
	if err = repo.menuColl.FindOne(ctx, filter).Decode(&res); err != nil {
		return models.Product{}, errorSwitch(err)
	}

	return res, nil
}

func (repo *MenuRepository) UpdateProductByFields(ctx context.Context, menuId string, productID string, req models.ProductUpdateRequest) error {
	oid, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		return errorSwitch(err)
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "products.ext_id", Value: productID},
	}

	var update bson.D

	if req.IsAvailable != nil {
		update = append(update, bson.E{
			Key:   "products.$.available",
			Value: *req.IsAvailable,
		})
	}

	if req.IsDisabled != nil {
		update = append(update, bson.E{
			Key:   "products.$.is_disabled",
			Value: *req.IsDisabled,
		})
	}

	result := bson.D{
		{
			Key:   "$set",
			Value: update,
		},
	}

	res, err := repo.menuColl.UpdateOne(ctx, filter, result)
	if err != nil {
		return errorSwitch(err)
	}

	if res.MatchedCount == 0 {
		return err
	}

	return nil
}

// ListProducts get all products from menu_id
func (repo *MenuRepository) ListProducts(ctx context.Context, query selector.Menu) ([]models.Product, int64, error) {

	oid, err := primitive.ObjectIDFromHex(query.ID)
	if err != nil {
		return nil, 0, drivers.ErrInvalid
	}

	match := bson.D{{Key: "_id", Value: oid}}
	unwind := "$products"
	project := bson.D{{Key: "products", Value: 1}}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$unwind", Value: unwind}},
		{{Key: "$project", Value: project}},
	}

	if query.HasSorting() {
		pipeline = append(pipeline,
			bson.D{{Key: "$sort", Value: repo.sortFrom(query.Sorting)}},
		)

	}

	if query.HasPagination() {
		pipeline = append(pipeline,
			bson.D{{Key: "$skip", Value: query.Skip()}},
			bson.D{{Key: "$limit", Value: query.Pagination.Limit}},
		)
	}

	cur, err := repo.menuColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer closeCur(cur)

	// TODO: testing this case
	products := make([]models.Product, 0, cur.RemainingBatchLength())
	for cur.Next(ctx) {
		var temp struct {
			Product models.Product `bson:"products"`
		}
		if err = cur.Decode(&temp); err != nil {
			return nil, 0, err
		}

		products = append(products, temp.Product)
	}

	return products, 0, nil

}

func (repo *MenuRepository) GetProductsByIDs(ctx context.Context, query selector.Menu, ids []string) ([]models.Product, error) {

	oid, err := primitive.ObjectIDFromHex(query.ID)
	if err != nil {
		return nil, errorSwitch(err)
	}

	match := bson.D{{Key: "_id", Value: oid}}
	unwind := "$products"
	project := bson.D{{Key: "products", Value: 1}}
	in := bson.D{{Key: "$in", Value: ids}}
	matchIDs := bson.D{{Key: "products.ext_id", Value: in}}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$unwind", Value: unwind}},
		{{Key: "$project", Value: project}},
		{{Key: "$match", Value: matchIDs}},
	}

	cur, err := repo.menuColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer closeCur(cur)

	var response []models.Product

	for cur.Next(ctx) {
		var tmp struct {
			Product models.Product `bson:"products"`
		}

		if err = cur.Decode(&tmp); err != nil {
			return nil, err
		}

		response = append(response, tmp.Product)
	}

	return response, nil
}

func (repo *MenuRepository) GetPromoProducts(ctx context.Context, query selector.Menu) ([]models.Product, error) {
	oid, err := primitive.ObjectIDFromHex(query.ID)
	if err != nil {
		return nil, errorSwitch(err)
	}

	match := bson.D{{Key: "_id", Value: oid}}
	unwind := "$products"
	project := bson.D{{Key: "products", Value: 1}}
	isPromo := bson.D{{Key: "products.is_promo", Value: true}}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$unwind", Value: unwind}},
		{{Key: "$project", Value: project}},
		{{Key: "$match", Value: isPromo}},
	}

	cur, err := repo.menuColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	defer closeCur(cur)

	var response []models.Product

	for cur.Next(ctx) {
		var tmp struct {
			Product models.Product `bson:"products"`
		}

		if err = cur.Decode(&tmp); err != nil {
			return nil, err
		}

		response = append(response, tmp.Product)
	}

	return response, nil
}

func (repo *MenuRepository) DeleteProducts(ctx context.Context, menuId string, productsIds []string) error {
	oid, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		return err
	}

	for _, productId := range productsIds {
		filter := bson.D{
			{Key: "_id", Value: oid},
			{Key: "products.ext_id", Value: productId},
		}

		update := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "products.$.is_deleted", Value: true},
				{Key: "products.$.available", Value: false},
			}},
		}

		_, err := repo.menuColl.UpdateOne(ctx, filter, update)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *MenuRepository) DeleteProductsFromDB(ctx context.Context, menuId string, productsIds []string) error {
	oid, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: oid}}

	update := bson.D{
		{Key: "$pull", Value: bson.D{{Key: "products", Value: bson.D{{Key: "ext_id", Value: bson.D{{Key: "$in", Value: productsIds}}}}}}},
	}

	_, err = repo.menuColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (repo *MenuRepository) UpdateProductForMatching(ctx context.Context, req models.MatchingProducts) error {

	oid, err := primitive.ObjectIDFromHex(req.MenuID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "products.ext_id", Value: req.ProductToChange},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "products.$.ext_id", Value: req.AggregatorProductID},
			{Key: "products.$.pos_id", Value: req.PosProductID},
			{Key: "products.$.sync", Value: req.IsSync},
			{Key: "products.$.available", Value: req.IsSync},
		}}}

	res, err := repo.menuColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return errorSwitch(err)
	}

	if res.MatchedCount == 0 {
		return err
	}
	return nil
}

func (repo *MenuRepository) GetEmptyProducts(ctx context.Context, menuID string, pagination selector.Pagination) ([]models.Product, int, error) {

	oid, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return nil, 0, err
	}

	match := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "_id", Value: oid},
		}},
	}
	unwind := bson.D{
		{Key: "$unwind", Value: "$products"},
	}
	project := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "products", Value: 1},
		}},
	}
	filter := bson.D{
		{Key: "$match", Value: bson.D{
			{Key: "$or", Value: bson.A{
				bson.D{{Key: "products.description", Value: nil}},
				bson.D{{Key: "products.image_urls", Value: nil}},
			}},
		}},
	}
	skip := bson.D{
		{Key: "$skip", Value: pagination.Skip()},
	}
	limit := bson.D{
		{Key: "$limit", Value: pagination.Limit},
	}

	pipeline := mongo.Pipeline{
		match,
		unwind,
		project,
		filter,
		skip,
		limit,
	}

	cur, err := repo.menuColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer closeCur(cur)

	var products []models.Product

	for cur.Next(ctx) {
		var tmp struct {
			Product models.Product `bson:"products"`
		}
		if err = cur.Decode(&tmp); err != nil {
			return nil, 0, err
		}
		products = append(products, tmp.Product)
	}

	return products, len(products), nil
}

func (repo *MenuRepository) UpdateProductAvailableStatus(ctx context.Context, menuID, productID string, status bool) error {

	oid, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "products.ext_id", Value: productID},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "products.$.available", Value: status},
		}},
	}

	res, err := repo.menuColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("could not update product")
	}

	return nil
}
