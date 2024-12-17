package menu

import (
	"context"
	"fmt"
	"time"

	customeErrors "github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
)

const collectionMenuName = "menus"

type Repository interface {
	FindById(ctx context.Context, menuID string) (*models.Menu, error)
	UpdateProductStopListStatus(ctx context.Context, menuId string, productID string, req models.ProductUpdateRequest) error
	UpdateAttributeStopListStatus(ctx context.Context, menuId string, attributeID string, isAvailable *bool, isDisabled *bool) error
	UpdateStopList(ctx context.Context, menuId string, stopListProducts []string) error
	ListProductsByMenuId(ctx context.Context, menuId string) ([]models.Product, int64, error)
	GetCombosByMenuId(ctx context.Context, menuId string) ([]models.Combo, int64, error)
	SearchProduct(ctx context.Context, menuId, productName string) ([]models.Product, error)
	BulkUpdateAttributesAvailability(ctx context.Context, menuId string, attributeIds []string, availability bool) error
	BulkUpdateProductsAvailability(ctx context.Context, menuId string, productIds []string, availability bool) error
	BulkUpdateProductsDisabledStatus(ctx context.Context, menuId string, productIds []string, isDisabled bool) error
	BulkUpdateAttributesDisabledStatus(ctx context.Context, menuId string, attributeIds []string, isDisabled bool) error
	Insert(ctx context.Context, menu models.Menu) (string, error)
	GetProductsByMenuIDAndExtIds(ctx context.Context, menuId string, productsExtIds []string) (models.Products, error)
	GetProductsByMenuIDAndSectionID(ctx context.Context, menuId, sectionId string) (models.Products, error)
	UpdateMenuEntities(ctx context.Context, menuId string, menu models.Menu) error
	BulkUpdateAttributesIsDeleted(ctx context.Context, menuId string, attributeIds []string, isDeleted bool, reason string) error
	BulkUpdateProductsIsDeleted(ctx context.Context, menuId string, productIds []string, isDeleted bool, reason string) error
	UpdateProductsImageAndDescription(ctx context.Context, menuID string, req []models.UpdateProductImageAndDescription) error
	AddNameInProduct(ctx context.Context, req models.AddLanguageDescriptionRequest) error
	AddDescriptionInProduct(ctx context.Context, req models.AddLanguageDescriptionRequest) error
	AddNameInSection(ctx context.Context, req models.AddLanguageDescriptionRequest) error
	AddDescriptionInSection(ctx context.Context, req models.AddLanguageDescriptionRequest) error
	AddNameInAttributeGroup(ctx context.Context, req models.AddLanguageDescriptionRequest) error
	AddNameInAttribute(ctx context.Context, req models.AddLanguageDescriptionRequest) error
	ChangeNameInProduct(ctx context.Context, req models.AddLanguageDescriptionRequest) error
	ChangeDescriptionInProduct(ctx context.Context, req models.AddLanguageDescriptionRequest) error
	ChangeNameInSection(ctx context.Context, req models.AddLanguageDescriptionRequest) error
	ChangeDescriptionInSection(ctx context.Context, req models.AddLanguageDescriptionRequest) error
	ChangeNameInAttributeGroup(ctx context.Context, req models.AddLanguageDescriptionRequest) error
	ChangeNameInAttribute(ctx context.Context, req models.AddLanguageDescriptionRequest) error
	AddRegulatoryInformation(ctx context.Context, req models.RegulatoryInformationRequest) error
	ChangeRegulatoryInformation(ctx context.Context, req models.RegulatoryInformationRequest) error
	DeleteAttributesFromAttributeGroup(ctx context.Context, menuID string, attributeIDs []string) error
	UpdateExcludedFromMenuProduct(ctx context.Context, menuID string, productIDs []string) error
	DeleteAttrGroupFromProduct(ctx context.Context, menuID, productID, attrGroupID string) error
	UpdateProductsDisabledByValidation(ctx context.Context, menuID string, productIDs []string, disabledByValidation bool) error
	UpdateAttributesPrice(ctx context.Context, menuID string, req []models.UpdateAttributePrice) error
}

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMenuMongoRepository(db *mongo.Database) (*MongoRepository, error) {
	r := MongoRepository{
		collection: db.Collection(collectionMenuName),
	}
	return &r, nil
}

func (r *MongoRepository) updateTo(ctx context.Context, menu models.Menu) bson.D {
	updateTo := bson.D{
		{
			Key:   "updated_at",
			Value: time.Now().UTC(),
		},
	}

	if len(menu.Products) != 0 {
		updateTo = append(updateTo, bson.E{
			Key:   "products",
			Value: menu.Products,
		})
	}

	if len(menu.Attributes) != 0 {
		updateTo = append(updateTo, bson.E{
			Key:   "attributes",
			Value: menu.Attributes,
		})
	}

	if len(menu.AttributesGroups) != 0 {
		updateTo = append(updateTo, bson.E{
			Key:   "attributes_groups",
			Value: menu.AttributesGroups,
		})
	}

	if len(menu.Sections) != 0 {
		updateTo = append(updateTo, bson.E{
			Key:   "sections",
			Value: menu.Sections,
		})
	}

	if len(menu.Combos) != 0 {
		updateTo = append(updateTo, bson.E{
			Key:   "combos",
			Value: menu.Combos,
		})
	}

	if len(menu.Collections) != 0 {
		updateTo = append(updateTo, bson.E{
			Key:   "collections",
			Value: menu.Collections,
		})
	}

	if len(menu.SuperCollections) != 0 {
		updateTo = append(updateTo, bson.E{
			Key:   "super_collections",
			Value: menu.SuperCollections,
		})
	}

	return updateTo
}

func (r *MongoRepository) UpdateMenuEntities(ctx context.Context, menuId string, menu models.Menu) error {
	filter, err := r.filterFrom(selector.EmptyMenuSearch().SetMenuID(menuId))
	if err != nil {
		return err
	}

	update := bson.D{
		{
			Key:   "$set",
			Value: r.updateTo(ctx, menu),
		},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.New("matched count 0")
	}

	return nil
}

func (r *MongoRepository) Insert(ctx context.Context, menu models.Menu) (string, error) {
	res, err := r.collection.InsertOne(ctx, menu)
	if err != nil {
		return "", errorSwitch(err)
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *MongoRepository) BulkUpdateAttributesAvailability(ctx context.Context, menuId string, attributeIds []string, availability bool) error {
	oid, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		return err
	}

	filter := bson.D{
		{
			Key:   "_id",
			Value: oid,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"attributes.$[element].available": availability,
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{bson.M{"element.ext_id": bson.M{"$in": attributeIds}}},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return errorSwitch(err)
	}

	if result.MatchedCount == 0 {
		return errors.New("matched count is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) BulkUpdateAttributesDisabledStatus(ctx context.Context, menuId string, attributeIds []string, isDisabled bool) error {
	oid, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		return err
	}

	filter := bson.D{
		{
			Key:   "_id",
			Value: oid,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"attributes.$[element].is_disabled": isDisabled,
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{bson.M{"element.ext_id": bson.M{"$in": attributeIds}}},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return errorSwitch(err)
	}

	if result.MatchedCount == 0 {
		return errors.New("matched count is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) BulkUpdateAttributesIsDeleted(ctx context.Context, menuId string, attributeIds []string, isDeleted bool, reason string) error {
	oid, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		return err
	}

	filter := bson.D{
		{
			Key:   "_id",
			Value: oid,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"attributes.$[element].is_deleted":        isDeleted,
			"attributes.$[element].is_deleted_reason": reason,
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{bson.M{"element.ext_id": bson.M{"$in": attributeIds}}},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return errorSwitch(err)
	}

	if result.MatchedCount == 0 {
		return errors.New("matched count is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) BulkUpdateProductsIsDeleted(ctx context.Context, menuId string, productIds []string, isDeleted bool, reason string) error {
	oid, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		return err
	}

	filter := bson.D{
		{
			Key:   "_id",
			Value: oid,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"products.$[element].is_deleted":        isDeleted,
			"products.$[element].is_deleted_reason": reason,
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{bson.M{"element.ext_id": bson.M{"$in": productIds}}},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return errorSwitch(err)
	}

	if result.MatchedCount == 0 {
		return errors.New("matched count is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) BulkUpdateProductsAvailability(ctx context.Context, menuId string, productIds []string, availability bool) error {
	oid, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		return err
	}

	filter := bson.D{
		{
			Key:   "_id",
			Value: oid,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"products.$[element].available": availability,
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{bson.M{"element.ext_id": bson.M{"$in": productIds}}},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return errorSwitch(err)
	}

	if result.MatchedCount == 0 {
		return errors.New("matched count is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) BulkUpdateProductsDisabledStatus(ctx context.Context, menuId string, productIds []string, isDisabled bool) error {
	oid, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		return err
	}

	filter := bson.D{
		{
			Key:   "_id",
			Value: oid,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"products.$[element].is_disabled": isDisabled,
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{bson.M{"element.ext_id": bson.M{"$in": productIds}}},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return errorSwitch(err)
	}

	if result.MatchedCount == 0 {
		return errors.New("matched count is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) GetCombosByMenuId(ctx context.Context, menuId string) ([]models.Combo, int64, error) {
	query := selector.EmptyMenuSearch().SetMenuID(menuId)

	oid, err := primitive.ObjectIDFromHex(query.ID)
	if err != nil {
		return nil, 0, drivers.ErrInvalid
	}

	match := bson.D{{Key: "_id", Value: oid}}
	unwind := "$combos"
	project := bson.D{{Key: "combos", Value: 1}}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$unwind", Value: unwind}},
		{{Key: "$project", Value: project}},
		{{Key: "$sort", Value: r.sortFrom(query.Sorting)}},
	}

	if query.HasPagination() {
		pipeline = append(pipeline,
			bson.D{{Key: "$skip", Value: query.Skip()}},
			bson.D{{Key: "$limit", Value: query.Pagination.Limit}},
		)
	}

	cur, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer closeCur(cur)

	combos := make([]models.Combo, 0, cur.RemainingBatchLength())

	for cur.Next(ctx) {
		var temp struct {
			Combo models.Combo `bson:"combos"`
		}

		if err = cur.Decode(&temp); err != nil {
			return nil, 0, err
		}

		combos = append(combos, temp.Combo)
	}

	return combos, 0, nil
}

func (r *MongoRepository) ListProductsByMenuId(ctx context.Context, menuId string) ([]models.Product, int64, error) {
	query := selector.EmptyMenuSearch().SetMenuID(menuId)

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
		{{Key: "$sort", Value: r.sortFrom(query.Sorting)}},
	}

	if query.HasPagination() {
		pipeline = append(pipeline,
			bson.D{{Key: "$skip", Value: query.Skip()}},
			bson.D{{Key: "$limit", Value: query.Pagination.Limit}},
		)
	}

	cur, err := r.collection.Aggregate(ctx, pipeline)
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

func closeCur(cur *mongo.Cursor) {
	if err := cur.Close(context.Background()); err != nil {
		log.Err(err).Msg("closing cursor:")
	}
}

func (r *MongoRepository) sortFrom(query selector.Sorting) bson.D {
	sort := make(bson.D, 0, 2)

	if query.HasSorting() {
		sort = append(sort, bson.E{Key: query.Param, Value: query.Direction})
	}
	sort = append(sort, bson.E{Key: "_id", Value: 1})

	return sort
}

func (r *MongoRepository) FindById(ctx context.Context, menuID string) (*models.Menu, error) {
	query := selector.EmptyMenuSearch().
		SetMenuID(menuID)
	filter, err := r.filterFrom(query)
	if err != nil {
		return nil, drivers.ErrInvalid
	}
	var menu models.Menu
	if err = r.collection.FindOne(ctx, filter).Decode(&menu); err != nil {
		return nil, errorSwitch(err)
	}

	return &menu, nil

}

func (r *MongoRepository) filterFrom(query selector.Menu) (bson.D, error) {

	filter := make(bson.D, 0, 4)

	if query.HasMenuID() {
		oid, err := primitive.ObjectIDFromHex(query.ID)
		if err != nil {
			return nil, drivers.ErrInvalid
		}
		filter = append(filter, bson.E{
			Key: "_id", Value: oid,
		})
	}

	if query.HasMenuName() {
		filter = append(filter, bson.E{
			Key: "name", Value: query.Name,
		})
	}

	if query.HasSectionID() {
		filter = append(filter, bson.E{
			Key: "section_id", Value: query.Name,
		})
	}

	if query.HasProductExtID() && !query.HasProductIsAvailable() {
		filter = append(filter, bson.E{
			Key: "products.ext_id", Value: query.ProductExtID,
		})
	}

	if query.HasProductExtID() && query.HasProductIsAvailable() {
		filter = append(filter, bson.E{
			Key: "products", Value: bson.D{
				{Key: "$elemMatch", Value: bson.D{
					{Key: "ext_id", Value: query.ProductExtID},
					{Key: "available", Value: query.ProductAvailable()},
				}},
			},
		})
	}

	return filter, nil
}

func errorSwitch(err error) error {
	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		return drivers.ErrNotFound
	case mongo.IsDuplicateKeyError(err):
		return drivers.ErrAlreadyExist
	default:
		return err
	}
}

func (r *MongoRepository) UpdateProductStopListStatus(ctx context.Context, menuId string, productID string, req models.ProductUpdateRequest) error {
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

	res, err := r.collection.UpdateOne(ctx, filter, result)
	if err != nil {
		return errorSwitch(err)
	}

	if res.MatchedCount == 0 {
		return errors.New("matched count 0")
	}

	return nil
}

func (r *MongoRepository) UpdateAttributeStopListStatus(ctx context.Context, menuId string, attributeID string, isAvailable *bool, isDisabled *bool) error {
	oid, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "attributes.ext_id", Value: attributeID},
	}

	var update bson.D

	if isAvailable != nil {
		update = append(update, bson.E{
			Key:   "attributes.$.available",
			Value: *isAvailable,
		})
	}

	if isDisabled != nil {
		update = append(update, bson.E{
			Key:   "attributes.$.is_disabled",
			Value: *isDisabled,
		})
	}

	result := bson.D{
		{
			Key:   "$set",
			Value: update,
		},
	}

	res, err := r.collection.UpdateOne(ctx, filter, result)
	if err != nil {
		return errorSwitch(err)
	}
	if res.MatchedCount == 0 {
		return errors.New("matched count 0")
	}

	return nil
}

func (r *MongoRepository) UpdateStopList(ctx context.Context, menuId string, stopListProducts []string) error {
	oid, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.D{
		{Key: "stoplist", Value: stopListProducts},
	}

	result := bson.D{
		{Key: "$set", Value: update},
	}

	res, err := r.collection.UpdateOne(ctx, filter, result)
	if err != nil {
		return errorSwitch(err)
	}
	if res.MatchedCount == 0 {
		return errors.New("matched count 0")
	}

	return nil
}

func (r *MongoRepository) SearchProduct(ctx context.Context, menuId, productName string) ([]models.Product, error) {
	oid, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		return nil, fmt.Errorf("menu_id is not valid: %s ", err)
	}
	matchMenu := bson.D{
		{Key: "_id", Value: oid},
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: matchMenu}},
		{{Key: "$unwind", Value: "$products"}},
		{{Key: "$match", Value: bson.M{"products.name.value": bson.M{"$regex": productName, "$options": "i"}}}},
	}

	res, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var result []models.Product

	type Output struct {
		Product models.Product `bson:"products"`
	}
	for res.Next(ctx) {
		var output Output

		if err = res.Decode(&output); err != nil {
			return nil, err
		}

		result = append(result, output.Product)
	}

	return result, nil
}

func (r *MongoRepository) GetProductsByMenuIDAndExtIds(ctx context.Context, menuId string, productsExtIds []string) (models.Products, error) {
	oid, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		return nil, errorSwitch(err)
	}

	match := bson.D{{Key: "_id", Value: oid}}
	unwind := "$products"
	project := bson.D{{Key: "products", Value: 1}}
	in := bson.D{{Key: "$in", Value: productsExtIds}}
	matchIDs := bson.D{{Key: "products.ext_id", Value: in}}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$unwind", Value: unwind}},
		{{Key: "$project", Value: project}},
		{{Key: "$match", Value: matchIDs}},
	}

	cur, err := r.collection.Aggregate(ctx, pipeline)
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

func (r *MongoRepository) GetProductsByMenuIDAndSectionID(ctx context.Context, menuId, sectionId string) (models.Products, error) {
	oid, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		return models.Products{}, errorSwitch(err)
	}

	matchMenuID := bson.D{{Key: "_id", Value: oid}}
	unwind := "$products"
	project := bson.D{{Key: "products", Value: 1}}
	matchSectionID := bson.D{{Key: "products.section", Value: sectionId}}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: matchMenuID}},
		{{Key: "$unwind", Value: unwind}},
		{{Key: "$project", Value: project}},
		{{Key: "$match", Value: matchSectionID}},
	}

	cur, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return models.Products{}, err
	}

	defer closeCur(cur)

	var response models.Products
	for cur.Next(ctx) {
		var tmp struct {
			Product models.Product `bson:"products"`
		}

		if err = cur.Decode(&tmp); err != nil {
			return models.Products{}, err
		}

		response = append(response, tmp.Product)
	}

	return response, err
}

func (r *MongoRepository) UpdateProductsImageAndDescription(ctx context.Context, menuID string, req []models.UpdateProductImageAndDescription) error {
	oid, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return err
	}

	var updates []mongo.WriteModel

	for _, product := range req {
		filter := bson.M{"_id": oid, "products.pos_id": product.PosID}
		update := bson.M{}
		if product.ImageURLs != nil {
			update["products.$.image_urls"] = product.ImageURLs
		}

		if product.Description != nil && len(product.Description) != 0 {
			if product.Description[0].Value != "" {
				update["products.$.description.0.Value"] = product.Description[0].Value
			}
		}

		update["products.$.weight"] = product.Weight
		update["products.$.measure_unit"] = product.MeasureUnit
		update["products.$.price"] = product.Price
		if len(update) > 0 {
			updateModel := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(bson.M{"$set": update})
			updates = append(updates, updateModel)
		}
	}

	tx, cb, err := r.StartSession(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = cb(err)
	}()

	opts := options.BulkWrite()

	if _, err = r.collection.BulkWrite(tx, updates, opts); err != nil {
		return err
	}

	return nil
}

func (r *MongoRepository) StartSession(ctx context.Context) (context.Context, drivers.TxCallback, error) {
	wc := writeconcern.Majority()
	rc := readconcern.Snapshot()
	txOpts := options.Transaction().
		SetWriteConcern(wc).
		SetReadConcern(rc).
		SetReadPreference(readpref.Primary())

	session, err := r.collection.Database().Client().StartSession()
	if err != nil {
		return nil, nil, err
	}

	if err = session.StartTransaction(txOpts); err != nil {
		return nil, nil, err
	}

	return mongo.NewSessionContext(ctx, session), callback(session), nil
}

func callback(session mongo.Session) func(err error) error {
	return func(err error) error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		defer session.EndSession(ctx)

		if err == nil {
			err = session.CommitTransaction(ctx)
		}
		if err != nil {
			if abortErr := session.AbortTransaction(ctx); abortErr != nil {
				var errs customeErrors.Error
				errs.Append(err, abortErr)
				err = errs.ErrorOrNil()
			}
		}
		return err
	}
}

func (r *MongoRepository) AddNameInProduct(ctx context.Context, req models.AddLanguageDescriptionRequest) error {
	tx, cb, err := r.StartSession(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = cb(err)
	}()

	oid, err := primitive.ObjectIDFromHex(req.MenuID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.M{
		"$push": bson.M{"products.$[element].name": req.Request},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{bson.M{"element.ext_id": req.ObjectID}},
	}

	result, err := r.collection.UpdateOne(tx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("matched count in aggregator menu is equal 0, not found")
	}

	oidPosMenu, err := primitive.ObjectIDFromHex(req.PosMenuID)
	if err != nil {
		return err
	}

	filterPosMenu := bson.D{
		{Key: "_id", Value: oidPosMenu},
	}

	updatePosMenu := bson.M{
		"$push": bson.M{"products.$[element].name": req.Request},
	}

	arrayFiltersPosMenu := options.ArrayFilters{
		Filters: []interface{}{bson.M{"element.ext_id": req.ObjectID}},
	}

	resultPosMenu, err := r.collection.UpdateOne(tx, filterPosMenu, updatePosMenu, options.Update().SetArrayFilters(arrayFiltersPosMenu))
	if err != nil {
		return err
	}

	if resultPosMenu.MatchedCount == 0 {
		return fmt.Errorf("matched count in pos menu is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) AddDescriptionInProduct(ctx context.Context, req models.AddLanguageDescriptionRequest) error {
	tx, cb, err := r.StartSession(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = cb(err)
	}()

	oid, err := primitive.ObjectIDFromHex(req.MenuID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.M{
		"$push": bson.M{"products.$[element].description": req.Request},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{bson.M{"element.ext_id": req.ObjectID}},
	}

	result, err := r.collection.UpdateOne(tx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("matched count in aggregator menu is equal 0, not found")
	}

	oidPosMenu, err := primitive.ObjectIDFromHex(req.PosMenuID)
	if err != nil {
		return err
	}

	filterPosMenu := bson.D{
		{Key: "_id", Value: oidPosMenu},
	}

	updatePosMenu := bson.M{
		"$push": bson.M{"products.$[element].description": req.Request},
	}

	arrayFiltersPosMenu := options.ArrayFilters{
		Filters: []interface{}{bson.M{"element.ext_id": req.ObjectID}},
	}

	resultPosMenu, err := r.collection.UpdateOne(tx, filterPosMenu, updatePosMenu, options.Update().SetArrayFilters(arrayFiltersPosMenu))
	if err != nil {
		return err
	}

	if resultPosMenu.MatchedCount == 0 {
		return fmt.Errorf("matched count in pos menu is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) AddNameInSection(ctx context.Context, req models.AddLanguageDescriptionRequest) error {
	tx, cb, err := r.StartSession(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = cb(err)
	}()

	oid, err := primitive.ObjectIDFromHex(req.MenuID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	var menu models.Menu
	if err = r.collection.FindOne(tx, filter).Decode(&menu); err != nil {
		return err
	}

	update := bson.M{
		"$push": bson.M{"sections.$[element].names_by_language": req.Request},
	}

	for _, section := range menu.Sections {
		if section.ExtID == req.ObjectID && section.NamesByLanguage == nil {
			update = bson.M{
				"$set": bson.M{"sections.$[element].names_by_language": bson.A{req.Request}},
			}
			break
		}
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{bson.M{"element.ext_id": req.ObjectID}},
	}

	result, err := r.collection.UpdateOne(tx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("matched count in aggregator menu is equal 0, not found")
	}

	oidPosMenu, err := primitive.ObjectIDFromHex(req.PosMenuID)
	if err != nil {
		return err
	}

	filterPosMenu := bson.D{
		{Key: "_id", Value: oidPosMenu},
	}

	var posMenu models.Menu
	if err = r.collection.FindOne(tx, filterPosMenu).Decode(&posMenu); err != nil {
		return err
	}

	updatePosMenu := bson.M{
		"$push": bson.M{"sections.$[element].names_by_language": req.Request},
	}

	for _, section := range posMenu.Sections {
		if section.ExtID == req.ObjectID && section.NamesByLanguage == nil {
			updatePosMenu = bson.M{
				"$set": bson.M{"sections.$[element].names_by_language": bson.A{req.Request}},
			}
			break
		}
	}

	arrayFiltersPosMenu := options.ArrayFilters{
		Filters: []interface{}{bson.M{"element.ext_id": req.ObjectID}},
	}

	resultPosMenu, err := r.collection.UpdateOne(tx, filterPosMenu, updatePosMenu, options.Update().SetArrayFilters(arrayFiltersPosMenu))
	if err != nil {
		return err
	}

	if resultPosMenu.MatchedCount == 0 {
		return fmt.Errorf("matched count in pos menu is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) AddDescriptionInSection(ctx context.Context, req models.AddLanguageDescriptionRequest) error {
	tx, cb, err := r.StartSession(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = cb(err)
	}()

	oid, err := primitive.ObjectIDFromHex(req.MenuID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	var menu models.Menu
	if err = r.collection.FindOne(tx, filter).Decode(&menu); err != nil {
		return err
	}

	update := bson.M{
		"$push": bson.M{"sections.$[element].description": req.Request},
	}

	for _, section := range menu.Sections {
		if section.ExtID == req.ObjectID && section.Description == nil {
			update = bson.M{
				"$set": bson.M{"sections.$[element].description": bson.A{req.Request}},
			}
			break
		}
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{bson.M{"element.ext_id": req.ObjectID}},
	}

	result, err := r.collection.UpdateOne(tx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("matched count in aggregator menu is equal 0, not found")
	}

	oidPosMenu, err := primitive.ObjectIDFromHex(req.PosMenuID)
	if err != nil {
		return err
	}

	filterPosMenu := bson.D{
		{Key: "_id", Value: oidPosMenu},
	}

	var posMenu models.Menu
	if err = r.collection.FindOne(tx, filterPosMenu).Decode(&posMenu); err != nil {
		return err
	}

	updatePosMenu := bson.M{
		"$push": bson.M{"sections.$[element].description": req.Request},
	}

	for _, section := range posMenu.Sections {
		if section.ExtID == req.ObjectID && section.Description == nil {
			updatePosMenu = bson.M{
				"$set": bson.M{"sections.$[element].description": bson.A{req.Request}},
			}
			break
		}
	}

	arrayFiltersPosMenu := options.ArrayFilters{
		Filters: []interface{}{bson.M{"element.ext_id": req.ObjectID}},
	}

	resultPosMenu, err := r.collection.UpdateOne(tx, filterPosMenu, updatePosMenu, options.Update().SetArrayFilters(arrayFiltersPosMenu))
	if err != nil {
		return err
	}

	if resultPosMenu.MatchedCount == 0 {
		return fmt.Errorf("matched count in pos menu is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) AddNameInAttributeGroup(ctx context.Context, req models.AddLanguageDescriptionRequest) error {
	tx, cb, err := r.StartSession(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = cb(err)
	}()

	oid, err := primitive.ObjectIDFromHex(req.MenuID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	var menu models.Menu
	if err = r.collection.FindOne(tx, filter).Decode(&menu); err != nil {
		return err
	}

	update := bson.M{
		"$push": bson.M{"attributes_groups.$[element].names_by_language": req.Request},
	}

	for _, group := range menu.AttributesGroups {
		if group.ExtID == req.ObjectID && group.NamesByLanguage == nil {
			update = bson.M{
				"$set": bson.M{"attributes_groups.$[element].names_by_language": bson.A{req.Request}},
			}
			break
		}
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{bson.M{"element.ext_id": req.ObjectID}},
	}

	result, err := r.collection.UpdateOne(tx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("matched count in aggregator menu is equal 0, not found")
	}

	oidPosMenu, err := primitive.ObjectIDFromHex(req.PosMenuID)
	if err != nil {
		return err
	}

	filterPosMenu := bson.D{
		{Key: "_id", Value: oidPosMenu},
	}

	var posMenu models.Menu
	if err = r.collection.FindOne(tx, filterPosMenu).Decode(&posMenu); err != nil {
		return err
	}

	updatePosMenu := bson.M{
		"$push": bson.M{"attributes_groups.$[element].names_by_language": req.Request},
	}

	for _, group := range posMenu.AttributesGroups {
		if group.ExtID == req.ObjectID && group.NamesByLanguage == nil {
			updatePosMenu = bson.M{
				"$set": bson.M{"attributes_groups.$[element].names_by_language": bson.A{req.Request}},
			}
			break
		}
	}

	arrayFiltersPosMenu := options.ArrayFilters{
		Filters: []interface{}{bson.M{"element.ext_id": req.ObjectID}},
	}

	resultPosMenu, err := r.collection.UpdateOne(tx, filterPosMenu, updatePosMenu, options.Update().SetArrayFilters(arrayFiltersPosMenu))
	if err != nil {
		return err
	}

	if resultPosMenu.MatchedCount == 0 {
		return fmt.Errorf("matched count in pos menu is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) AddNameInAttribute(ctx context.Context, req models.AddLanguageDescriptionRequest) error {
	tx, cb, err := r.StartSession(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = cb(err)
	}()

	oid, err := primitive.ObjectIDFromHex(req.MenuID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	var menu models.Menu
	if err = r.collection.FindOne(tx, filter).Decode(&menu); err != nil {
		return err
	}

	update := bson.M{
		"$push": bson.M{"attributes.$[element].names_by_language": req.Request},
	}

	for _, attr := range menu.Attributes {
		if attr.ExtID == req.ObjectID && attr.NamesByLanguage == nil {
			update = bson.M{
				"$set": bson.M{"attributes.$[element].names_by_language": bson.A{req.Request}},
			}
			break
		}
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{bson.M{"element.ext_id": req.ObjectID}},
	}

	result, err := r.collection.UpdateOne(tx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("matched count in aggregator menu is equal 0, not found")
	}

	oidPosMenu, err := primitive.ObjectIDFromHex(req.PosMenuID)
	if err != nil {
		return err
	}

	filterPosMenu := bson.D{
		{Key: "_id", Value: oidPosMenu},
	}

	var posMenu models.Menu
	if err = r.collection.FindOne(tx, filterPosMenu).Decode(&posMenu); err != nil {
		return err
	}

	updatePosMenu := bson.M{
		"$push": bson.M{"attributes.$[element].names_by_language": req.Request},
	}

	for _, attr := range posMenu.Attributes {
		if attr.ExtID == req.ObjectID && attr.NamesByLanguage == nil {
			updatePosMenu = bson.M{
				"$set": bson.M{"attributes.$[element].names_by_language": bson.A{req.Request}},
			}
			break
		}
	}

	arrayFiltersPosMenu := options.ArrayFilters{
		Filters: []interface{}{bson.M{"element.ext_id": req.ObjectID}},
	}

	resultPosMenu, err := r.collection.UpdateOne(tx, filterPosMenu, updatePosMenu, options.Update().SetArrayFilters(arrayFiltersPosMenu))
	if err != nil {
		return err
	}

	if resultPosMenu.MatchedCount == 0 {
		return fmt.Errorf("matched count in pos menu is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) ChangeNameInProduct(ctx context.Context, req models.AddLanguageDescriptionRequest) error {
	tx, cb, err := r.StartSession(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = cb(err)
	}()

	oid, err := primitive.ObjectIDFromHex(req.MenuID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.M{
		"$set": bson.M{"products.$[product].name.$[name].value": req.Request.Value},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"product.ext_id": req.ObjectID},
			bson.M{"name.language_code": req.Request.LanguageCode},
		},
	}

	result, err := r.collection.UpdateOne(tx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("matched count in aggregator menu is equal 0, not found")
	}

	oidPosMenu, err := primitive.ObjectIDFromHex(req.PosMenuID)
	if err != nil {
		return err
	}

	filterPosMenu := bson.D{
		{Key: "_id", Value: oidPosMenu},
	}

	updatePosMenu := bson.M{
		"$set": bson.M{"products.$[product].name.$[name].value": req.Request.Value},
	}

	arrayFiltersPosMenu := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"product.ext_id": req.ObjectID},
			bson.M{"name.language_code": req.Request.LanguageCode},
		},
	}

	resultPosMenu, err := r.collection.UpdateOne(tx, filterPosMenu, updatePosMenu, options.Update().SetArrayFilters(arrayFiltersPosMenu))
	if err != nil {
		return err
	}

	if resultPosMenu.MatchedCount == 0 {
		return fmt.Errorf("matched count in pos menu is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) ChangeDescriptionInProduct(ctx context.Context, req models.AddLanguageDescriptionRequest) error {
	tx, cb, err := r.StartSession(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = cb(err)
	}()

	oid, err := primitive.ObjectIDFromHex(req.MenuID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.M{
		"$set": bson.M{"products.$[product].description.$[description].value": req.Request.Value},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"product.ext_id": req.ObjectID},
			bson.M{"description.language_code": req.Request.LanguageCode},
		},
	}

	result, err := r.collection.UpdateOne(tx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("matched count in aggregator menu is equal 0, not found")
	}

	oidPosMenu, err := primitive.ObjectIDFromHex(req.PosMenuID)
	if err != nil {
		return err
	}

	filterPosMenu := bson.D{
		{Key: "_id", Value: oidPosMenu},
	}

	updatePosMenu := bson.M{
		"$set": bson.M{"products.$[product].description.$[description].value": req.Request.Value},
	}

	arrayFiltersPosMenu := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"product.ext_id": req.ObjectID},
			bson.M{"description.language_code": req.Request.LanguageCode},
		},
	}

	resultPosMenu, err := r.collection.UpdateOne(tx, filterPosMenu, updatePosMenu, options.Update().SetArrayFilters(arrayFiltersPosMenu))
	if err != nil {
		return err
	}

	if resultPosMenu.MatchedCount == 0 {
		return fmt.Errorf("matched count in pos menu is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) ChangeNameInSection(ctx context.Context, req models.AddLanguageDescriptionRequest) error {
	tx, cb, err := r.StartSession(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = cb(err)
	}()

	oid, err := primitive.ObjectIDFromHex(req.MenuID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.M{
		"$set": bson.M{"sections.$[section].names_by_language.$[name].value": req.Request.Value},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"section.ext_id": req.ObjectID},
			bson.M{"name.language_code": req.Request.LanguageCode},
		},
	}

	result, err := r.collection.UpdateOne(tx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("matched count in aggregator menu is equal 0, not found")
	}

	oidPosMenu, err := primitive.ObjectIDFromHex(req.PosMenuID)
	if err != nil {
		return err
	}

	filterPosMenu := bson.D{
		{Key: "_id", Value: oidPosMenu},
	}

	updatePosMenu := bson.M{
		"$set": bson.M{"sections.$[section].names_by_language.$[name].value": req.Request.Value},
	}

	arrayFiltersPosMenu := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"section.ext_id": req.ObjectID},
			bson.M{"name.language_code": req.Request.LanguageCode},
		},
	}

	resultPosMenu, err := r.collection.UpdateOne(tx, filterPosMenu, updatePosMenu, options.Update().SetArrayFilters(arrayFiltersPosMenu))
	if err != nil {
		return err
	}

	if resultPosMenu.MatchedCount == 0 {
		return fmt.Errorf("matched count in pos menu is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) ChangeDescriptionInSection(ctx context.Context, req models.AddLanguageDescriptionRequest) error {
	tx, cb, err := r.StartSession(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = cb(err)
	}()

	oid, err := primitive.ObjectIDFromHex(req.MenuID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.M{
		"$set": bson.M{"sections.$[section].description.$[description].value": req.Request.Value},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"section.ext_id": req.ObjectID},
			bson.M{"description.language_code": req.Request.LanguageCode},
		},
	}

	result, err := r.collection.UpdateOne(tx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("matched count in aggregator menu is equal 0, not found")
	}

	oidPosMenu, err := primitive.ObjectIDFromHex(req.PosMenuID)
	if err != nil {
		return err
	}

	filterPosMenu := bson.D{
		{Key: "_id", Value: oidPosMenu},
	}

	updatePosMenu := bson.M{
		"$set": bson.M{"sections.$[section].description.$[description].value": req.Request.Value},
	}

	arrayFiltersPosMenu := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"section.ext_id": req.ObjectID},
			bson.M{"description.language_code": req.Request.LanguageCode},
		},
	}

	resultPosMenu, err := r.collection.UpdateOne(tx, filterPosMenu, updatePosMenu, options.Update().SetArrayFilters(arrayFiltersPosMenu))
	if err != nil {
		return err
	}

	if resultPosMenu.MatchedCount == 0 {
		return fmt.Errorf("matched count in pos menu is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) ChangeNameInAttributeGroup(ctx context.Context, req models.AddLanguageDescriptionRequest) error {
	tx, cb, err := r.StartSession(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = cb(err)
	}()

	oid, err := primitive.ObjectIDFromHex(req.MenuID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.M{
		"$set": bson.M{"attributes_groups.$[ag].names_by_language.$[name].value": req.Request.Value},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"ag.ext_id": req.ObjectID},
			bson.M{"name.language_code": req.Request.LanguageCode},
		},
	}

	result, err := r.collection.UpdateOne(tx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("matched count in aggregator menu is equal 0, not found")
	}

	oidPosMenu, err := primitive.ObjectIDFromHex(req.PosMenuID)
	if err != nil {
		return err
	}

	filterPosMenu := bson.D{
		{Key: "_id", Value: oidPosMenu},
	}

	updatePosMenu := bson.M{
		"$set": bson.M{"attributes_groups.$[ag].names_by_language.$[name].value": req.Request.Value},
	}

	arrayFiltersPosMenu := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"ag.ext_id": req.ObjectID},
			bson.M{"name.language_code": req.Request.LanguageCode},
		},
	}

	resultPosMenu, err := r.collection.UpdateOne(tx, filterPosMenu, updatePosMenu, options.Update().SetArrayFilters(arrayFiltersPosMenu))
	if err != nil {
		return err
	}

	if resultPosMenu.MatchedCount == 0 {
		return fmt.Errorf("matched count in pos menu is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) ChangeNameInAttribute(ctx context.Context, req models.AddLanguageDescriptionRequest) error {
	tx, cb, err := r.StartSession(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = cb(err)
	}()

	oid, err := primitive.ObjectIDFromHex(req.MenuID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.M{
		"$set": bson.M{"attributes.$[attribute].names_by_language.$[name].value": req.Request.Value},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"attribute.ext_id": req.ObjectID},
			bson.M{"name.language_code": req.Request.LanguageCode},
		},
	}

	result, err := r.collection.UpdateOne(tx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("matched count in aggregator menu is equal 0, not found")
	}

	oidPosMenu, err := primitive.ObjectIDFromHex(req.PosMenuID)
	if err != nil {
		return err
	}

	filterPosMenu := bson.D{
		{Key: "_id", Value: oidPosMenu},
	}

	updatePosMenu := bson.M{
		"$set": bson.M{"attributes.$[attribute].names_by_language.$[name].value": req.Request.Value},
	}

	arrayFiltersPosMenu := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"attribute.ext_id": req.ObjectID},
			bson.M{"name.language_code": req.Request.LanguageCode},
		},
	}

	resultPosMenu, err := r.collection.UpdateOne(tx, filterPosMenu, updatePosMenu, options.Update().SetArrayFilters(arrayFiltersPosMenu))
	if err != nil {
		return err
	}

	if resultPosMenu.MatchedCount == 0 {
		return fmt.Errorf("matched count in pos menu is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) AddRegulatoryInformation(ctx context.Context, req models.RegulatoryInformationRequest) error {
	tx, cb, err := r.StartSession(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = cb(err)
	}()

	oid, err := primitive.ObjectIDFromHex(req.MenuID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	var menu models.Menu
	if err := r.collection.FindOne(ctx, filter).Decode(&menu); err != nil {
		return err
	}

	update := bson.M{
		"$push": bson.M{"products.$[product].product_information.regulatory_information": req.RegulatoryInformation},
	}

	for _, product := range menu.Products {
		if product.ExtID == req.ProductID && product.ProductInformation.RegulatoryInformation == nil {
			update = bson.M{
				"$set": bson.M{"products.$[product].product_information.regulatory_information": bson.A{req.RegulatoryInformation}},
			}
			break
		}
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{bson.M{"product.ext_id": req.ProductID}},
	}

	result, err := r.collection.UpdateOne(tx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("matched count in aggregator menu is equal 0, not found")
	}

	oidPosMenu, err := primitive.ObjectIDFromHex(req.PosMenuID)
	if err != nil {
		return err
	}

	filterPosMenu := bson.D{
		{Key: "_id", Value: oidPosMenu},
	}

	var posMenu models.Menu
	if err := r.collection.FindOne(ctx, filterPosMenu).Decode(&posMenu); err != nil {
		return err
	}

	updatePosMenu := bson.M{
		"$push": bson.M{"products.$[product].product_information.regulatory_information": req.RegulatoryInformation},
	}

	for _, product := range posMenu.Products {
		if product.ExtID == req.ProductID && product.ProductInformation.RegulatoryInformation == nil {
			updatePosMenu = bson.M{
				"$set": bson.M{"products.$[product].product_information.regulatory_information": bson.A{req.RegulatoryInformation}},
			}
			break
		}
	}

	arrayFiltersPosMenu := options.ArrayFilters{
		Filters: []interface{}{bson.M{"product.ext_id": req.ProductID}},
	}

	resultPosMenu, err := r.collection.UpdateOne(tx, filterPosMenu, updatePosMenu, options.Update().SetArrayFilters(arrayFiltersPosMenu))
	if err != nil {
		return err
	}

	if resultPosMenu.MatchedCount == 0 {
		return fmt.Errorf("matched count in aggregator menu is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) ChangeRegulatoryInformation(ctx context.Context, req models.RegulatoryInformationRequest) error {
	tx, cb, err := r.StartSession(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = cb(err)
	}()

	oid, err := primitive.ObjectIDFromHex(req.MenuID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.M{
		"$set": bson.M{"products.$[product].product_information.regulatory_information.$[info].value": req.RegulatoryInformation.Value},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"product.ext_id": req.ProductID},
			bson.M{"info.name": req.RegulatoryInformation.Name},
		},
	}

	result, err := r.collection.UpdateOne(tx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("matched count in aggregator menu is equal 0, not found")
	}

	oidPosMenu, err := primitive.ObjectIDFromHex(req.PosMenuID)
	if err != nil {
		return err
	}

	filterPosMenu := bson.D{
		{Key: "_id", Value: oidPosMenu},
	}

	updatePosMenu := bson.M{
		"$set": bson.M{"products.$[product].product_information.regulatory_information.$[info].value": req.RegulatoryInformation.Value},
	}

	arrayFiltersPosMenu := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"product.ext_id": req.ProductID},
			bson.M{"info.name": req.RegulatoryInformation.Name},
		},
	}

	resultPosMenu, err := r.collection.UpdateOne(tx, filterPosMenu, updatePosMenu, options.Update().SetArrayFilters(arrayFiltersPosMenu))
	if err != nil {
		return err
	}

	if resultPosMenu.MatchedCount == 0 {
		return fmt.Errorf("matched count in pos menu is equal 0, not found")
	}

	return nil
}
func (r *MongoRepository) DeleteAttributesFromAttributeGroup(ctx context.Context, menuID string, attributeIDs []string) error {
	oid, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.M{
		"$pull": bson.M{"attributes_groups.$[].attributes": bson.M{"$in": attributeIDs}},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("matched count in menu is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) UpdateExcludedFromMenuProduct(ctx context.Context, menuID string, productIDs []string) error {
	oid, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.M{
		"$set": bson.M{
			"products.$[product].included_in_menu": false,
			"products.$[product].available":        false,
			"products.$[product].is_deleted":       true,
		},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"product.ext_id": bson.M{"$in": productIDs}},
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("matched count is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) DeleteAttrGroupFromProduct(ctx context.Context, menuID, productID, attrGroupID string) error {
	oid, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.M{
		"$pull": bson.M{"products.$[product].attributes_groups": attrGroupID},
	}

	arrayFilter := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"product.ext_id": productID},
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetArrayFilters(arrayFilter))
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("matched count is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) UpdateProductsDisabledByValidation(ctx context.Context, menuID string, productIDs []string, disabledByValidation bool) error {
	oid, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.M{
		"$set": bson.M{"products.$[product].disabled_by_validation": disabledByValidation},
	}

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"product.ext_id": bson.M{"$in": productIDs}},
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("matched count is equal 0, not found")
	}

	return nil
}

func (r *MongoRepository) UpdateAttributesPrice(ctx context.Context, menuID string, req []models.UpdateAttributePrice) error {
	oid, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return err
	}

	var updates []mongo.WriteModel

	for _, attribute := range req {
		filter := bson.M{"_id": oid, "attributes.ext_id": attribute.ExtID}
		update := bson.M{}
		update["attributes.$.price_impact"] = attribute.Price
		if len(update) > 0 {
			updateModel := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(bson.M{"$set": update})
			updates = append(updates, updateModel)
		}
	}

	tx, cb, err := r.StartSession(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = cb(err)
	}()

	opts := options.BulkWrite()

	if _, err = r.collection.BulkWrite(tx, updates, opts); err != nil {
		return err
	}

	return nil
}
