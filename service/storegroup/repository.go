package storegroup

import (
	"context"
	errors2 "github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	drivers2 "github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	selector2 "github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/service/legalentity/models"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collectionStoreGroupName = "restaurant_groups"

type Repository interface {
	GetStoreGroupByID(ctx context.Context, storeGroupID string) (storeModels.StoreGroup, error)
	GetStoreGroupByStoreID(ctx context.Context, storeID string) (storeModels.StoreGroup, error)
	GetStoreGroupsWithFilter(ctx context.Context, query selector2.StoreGroup) ([]storeModels.StoreGroup, error)
	CreateStoreGroup(ctx context.Context, group storeModels.StoreGroup) (string, error)
	UpdateStoreGroup(ctx context.Context, group storeModels.UpdateStoreGroup) error
	GetStoreGroupLegalEntities(ctx context.Context, id string) ([]models.LegalEntityView, error)
	GetAllStoreGroupsIdsAndNames(ctx context.Context) ([]storeModels.StoreGroupIdAndName, error)
	AddBrandInfo(ctx context.Context, brandInfo storeModels.BrandInfo, restGroupID string) error

	CreateDirectPromoBanners(ctx context.Context, storeGroupID string, banner storeModels.DirectPromoBanner) error
	UpdateDirectPromoBanners(ctx context.Context, storeGroupID string, banner storeModels.UpdateDirectPromoBanner) error
	GetAllDirectPromoBannersByStoreGroup(ctx context.Context, storeGroupID string) ([]storeModels.DirectPromoBanner, error)
	DeleteDirectPromoBannerByID(ctx context.Context, storeGroupID, directPromoID string) error
}

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) (*MongoRepository, error) {
	r := MongoRepository{
		collection: db.Collection(collectionStoreGroupName),
	}
	return &r, nil
}

func (r *MongoRepository) GetStoreGroupByStoreID(ctx context.Context, storeID string) (storeModels.StoreGroup, error) {
	filter := bson.D{
		{Key: "restaurant_ids", Value: storeID},
	}

	res := r.collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return storeModels.StoreGroup{}, errNoDocs
		}
		return storeModels.StoreGroup{}, err
	}

	var result storeModels.StoreGroup

	if err := res.Decode(&result); err != nil {
		return storeModels.StoreGroup{}, err
	}

	return result, nil
}

func (r *MongoRepository) GetStoreGroupByID(ctx context.Context, storeGroupID string) (storeModels.StoreGroup, error) {
	objID, err := primitive.ObjectIDFromHex(storeGroupID)
	if err != nil {
		return storeModels.StoreGroup{}, err
	}
	filter := bson.D{
		{Key: "_id", Value: objID},
	}

	res := r.collection.FindOne(ctx, filter)
	if err = res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return storeModels.StoreGroup{}, errNoDocs
		}
		return storeModels.StoreGroup{}, err
	}

	var result storeModels.StoreGroup

	if err = res.Decode(&result); err != nil {
		return storeModels.StoreGroup{}, err
	}
	return result, nil
}

func (r *MongoRepository) GetStoreGroupsWithFilter(ctx context.Context, query selector2.StoreGroup) ([]storeModels.StoreGroup, error) {
	pipeline := mongo.Pipeline{}

	filter, err := r.filterFrom(query)
	if err != nil {
		return []storeModels.StoreGroup{}, err
	}

	if len(filter) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: filter}})
	}

	sort := make(bson.D, 0, 2)
	if query.HasSorting() {
		if query.Param == "restaurant_ids" {
			pipeline = append(pipeline, bson.D{
				{Key: "$addFields", Value: bson.D{
					{Key: "restaurant_ids_len", Value: bson.D{
						{Key: "$size", Value: bson.D{
							{Key: "$ifNull", Value: bson.A{
								"$restaurant_ids",
								bson.A{},
							}},
						}},
					}},
				}},
			})

			sort = append(sort, bson.E{Key: "restaurant_ids_len", Value: query.Direction})
		} else {
			sort = append(sort, bson.E{Key: query.Param, Value: query.Direction})
		}
	}

	if len(sort) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$sort", Value: sort}})
	}

	if query.HasPagination() {
		pipeline = append(pipeline, bson.D{{Key: "$skip", Value: query.Pagination.Skip()}})
		pipeline = append(pipeline, bson.D{{Key: "$limit", Value: query.Pagination.Limit}})
	}

	cur, err := r.collection.Aggregate(ctx, pipeline)

	if err != nil {
		return []storeModels.StoreGroup{}, err
	}

	var storeGroups []storeModels.StoreGroup
	if err = cur.All(ctx, &storeGroups); err != nil {
		return []storeModels.StoreGroup{}, err
	}

	return storeGroups, nil
}

func (r *MongoRepository) CreateStoreGroup(ctx context.Context, group storeModels.StoreGroup) (string, error) {
	res, err := r.collection.InsertOne(ctx, group)
	if err != nil {
		return "", err
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)

	if !ok {
		return "", drivers.ErrInvalid
	}

	return oid.Hex(), nil
}

func (r *MongoRepository) UpdateStoreGroup(ctx context.Context, group storeModels.UpdateStoreGroup) error {
	oid, err := primitive.ObjectIDFromHex(*group.ID)

	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	*group.ID = ""
	// group.UpdatedAt = models.TimeNow()

	update := r.updateFrom(group)

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors2.ErrorSwitch(err)
	}

	if res.MatchedCount == 0 {
		return errors2.ErrNotFound
	}

	return nil
}

func (r *MongoRepository) GetStoreGroupLegalEntities(ctx context.Context, id string) ([]models.LegalEntityView, error) {
	oid, idErr := primitive.ObjectIDFromHex(id)
	if idErr != nil {
		return []models.LegalEntityView{}, idErr
	}

	legalEntitiesFieldName := "legal_entities"

	pipeline := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: oid}}}},
		bson.D{
			{Key: "$lookup",
				Value: bson.D{
					{Key: "from", Value: "legal_entity"},
					{Key: "localField", Value: "restaurant_ids"},
					{Key: "foreignField", Value: "store_ids"},
					{Key: "as", Value: legalEntitiesFieldName},
				},
			},
		},
	}

	cur, aggErr := r.collection.Aggregate(ctx, pipeline)
	if aggErr != nil {
		return []models.LegalEntityView{}, aggErr
	}

	var res []bson.M
	if curErr := cur.All(ctx, &res); curErr != nil {
		return []models.LegalEntityView{}, curErr
	}

	var legalEntities []models.LegalEntityView
	for _, el := range res {
		legalEntitiesField, ok := el[legalEntitiesFieldName].(bson.A)

		if !ok {
			return []models.LegalEntityView{}, errors.Errorf("%s field not found", legalEntitiesFieldName)
		}

		for _, legalEntityField := range legalEntitiesField {

			legalEntityFieldBytes, marshalErr := bson.Marshal(legalEntityField)

			if marshalErr != nil {
				return []models.LegalEntityView{}, marshalErr
			}

			var legalEntity models.LegalEntityView
			if unmarshalErr := bson.Unmarshal(legalEntityFieldBytes, &legalEntity); unmarshalErr != nil {
				return []models.LegalEntityView{}, unmarshalErr
			}

			legalEntities = append(legalEntities, legalEntity)
		}
	}

	return legalEntities, nil
}

func (r *MongoRepository) GetAllStoreGroupsIdsAndNames(ctx context.Context) ([]storeModels.StoreGroupIdAndName, error) {
	project := options.Find().SetProjection(map[string]int{"_id": 1, "name": 1})

	curs, err := r.collection.Find(ctx, bson.D{}, project)
	if err != nil {
		return nil, err
	}

	resp := make([]storeModels.StoreGroupIdAndName, 0, curs.RemainingBatchLength())
	if err = curs.All(ctx, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *MongoRepository) filterFrom(query selector2.StoreGroup) (bson.D, error) {
	filter := make(bson.D, 0, 2)

	if query.HasID() {
		oid, err := primitive.ObjectIDFromHex(query.ID)
		if err != nil {
			return nil, drivers.ErrInvalid
		}

		filter = append(filter, bson.E{Key: "_id", Value: oid})
	}

	if query.HasCategory() {
		filter = append(filter, bson.E{Key: "category", Value: query.Category})
	}

	if query.HasCountry() {
		filter = append(filter, bson.E{Key: "country", Value: query.Country})
	}

	if query.HasName() {
		filter = append(filter, bson.E{Key: "name", Value: query.Name})
	}

	if query.HasStatus() {
		filter = append(filter, bson.E{Key: "status", Value: query.Status})
	}

	if query.HasStoreIDs() {
		filter = append(filter, bson.E{Key: "restaurant_ids", Value: bson.D{{Key: "$in", Value: query.StoreIDs}}})
	}

	if query.HasCountries() {
		filter = append(filter, bson.E{Key: "country", Value: bson.D{{Key: "$in", Value: query.Countries}}})
	}

	if query.HasCategories() {
		filter = append(filter, bson.E{Key: "category", Value: bson.D{{Key: "$in", Value: query.Categories}}})
	}

	if query.HasStatuses() {
		filter = append(filter, bson.E{Key: "status", Value: bson.D{{Key: "$in", Value: query.Statuses}}})
	}

	return filter, nil
}

func (r *MongoRepository) updateFrom(group storeModels.UpdateStoreGroup) bson.D {
	update := make(bson.D, 0, 10)

	if group.Name != nil && *group.Name != "" {
		update = append(update, bson.E{Key: "name", Value: *group.Name})
	}

	if group.StoreIds != nil {
		update = append(update, bson.E{Key: "restaurant_ids", Value: group.StoreIds})
	}

	if group.Locations != nil {
		update = append(update, bson.E{Key: "locations", Value: group.Locations})
	}

	if group.IsTopPartner != nil {
		update = append(update, bson.E{Key: "is_top_partner", Value: *group.IsTopPartner})
	}

	if group.RetryCount != nil {
		update = append(update, bson.E{Key: "retry_count", Value: *group.RetryCount})
	}

	if group.ColumnView != nil {
		update = append(update, bson.E{Key: "column_view", Value: *group.ColumnView})
	}

	if group.Logo != nil && *group.Logo != "" {
		update = append(update, bson.E{Key: "logo", Value: *group.Logo})
	}

	if group.Country != nil && *group.Country != "" {
		update = append(update, bson.E{Key: "country", Value: *group.Country})
	}

	if group.Category != nil && *group.Category != "" {
		update = append(update, bson.E{Key: "category", Value: *group.Category})
	}

	if group.BrandType != nil && *group.BrandType != "" {
		update = append(update, bson.E{Key: "brand_type", Value: *group.BrandType})
	}

	if group.Category != nil && *group.Category != "" {
		update = append(update, bson.E{Key: "category", Value: *group.Category})
	}

	if group.Chats != nil {
		update = append(update, bson.E{Key: "chats", Value: group.Chats})
	}

	if group.Contacts != nil {
		update = append(update, bson.E{Key: "contacts", Value: group.Contacts})
	}

	if group.SalesComments != nil && *group.SalesComments != "" {
		update = append(update, bson.E{Key: "sales_comments", Value: *group.SalesComments})
	}

	if group.Status != nil && *group.Status != "" {
		update = append(update, bson.E{Key: "status", Value: *group.Status})
	}

	if group.Description != nil && *group.Description != "" {
		update = append(update, bson.E{Key: "description", Value: *group.Description})
	}

	if group.HeaderImage != nil && *group.HeaderImage != "" {
		update = append(update, bson.E{Key: "header_image", Value: *group.HeaderImage})
	}

	if group.WorkSchedule != nil && len(group.WorkSchedule) > 0 {
		update = append(update, bson.E{Key: "work_schedule", Value: group.WorkSchedule})
	}

	if group.SocialMediaLinks != nil && len(group.SocialMediaLinks) > 0 {
		update = append(update, bson.E{Key: "social_media_links", Value: group.SocialMediaLinks})
	}

	if group.DomainName != nil && *group.DomainName != "" {
		update = append(update, bson.E{Key: "domain_name", Value: *group.DomainName})
	}

	if group.DefaultRestaurantId != nil && *group.DefaultRestaurantId != "" {
		update = append(update, bson.E{Key: "default_restaurant_id", Value: *group.DefaultRestaurantId})
	}

	res := bson.D{{
		Key: "$set", Value: update,
	}}

	return res
}

func (s *MongoRepository) AddBrandInfo(ctx context.Context, brandInfo storeModels.BrandInfo, restGroupID string) error {
	oid, err := primitive.ObjectIDFromHex(restGroupID)
	if err != nil {
		return err
	}
	filter := bson.D{
		{Key: "_id", Value: oid},
	}
	update := bson.D{
		{
			Key: "$set", Value: bson.D{
				{Key: "brand_info", Value: brandInfo},
			},
		},
	}
	res, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return drivers2.ErrNotFound
	}
	return nil
}

func (s *MongoRepository) CreateDirectPromoBanners(ctx context.Context, storeGroupID string, banner storeModels.DirectPromoBanner) error {
	oid, err := primitive.ObjectIDFromHex(storeGroupID)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: oid}}

	update := bson.D{
		{Key: "$addToSet", Value: bson.D{
			{Key: "direct_promo_banners", Value: banner},
		}},
	}

	res, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return drivers2.ErrNotFound
	}

	return nil
}

func (s *MongoRepository) UpdateDirectPromoBanners(ctx context.Context, storeGroupID string, banner storeModels.UpdateDirectPromoBanner) error {
	oid, err := primitive.ObjectIDFromHex(storeGroupID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update, arrayFilters := s.filterFromForBanner(banner)

	updateOpts := options.Update().SetArrayFilters(options.ArrayFilters{Filters: arrayFilters})

	if _, err = s.collection.UpdateOne(ctx, filter, update, updateOpts); err != nil {
		return err
	}

	return nil
}

func (s *MongoRepository) filterFromForBanner(banner storeModels.UpdateDirectPromoBanner) (bson.D, bson.A) {
	update := bson.D{}

	if banner.Image != nil && *banner.Image != "" {
		update = append(update, bson.E{Key: "direct_promo_banners.$[elem].image", Value: *banner.Image})
	}

	if banner.RestaurantIDs != nil && len(*banner.RestaurantIDs) != 0 {
		update = append(update, bson.E{Key: "direct_promo_banners.$[elem].restaurant_ids", Value: *banner.RestaurantIDs})
	}

	if banner.IsActive != nil {
		update = append(update, bson.E{Key: "direct_promo_banners.$[elem].is_active", Value: *banner.IsActive})
	}

	arrayFilters := bson.A{bson.M{"elem.promo_id": banner.ID}}

	return bson.D{{Key: "$set", Value: update}}, arrayFilters
}

func (s *MongoRepository) GetAllDirectPromoBannersByStoreGroup(ctx context.Context, storeGroupID string) ([]storeModels.DirectPromoBanner, error) {

	storeGroup := storeModels.StoreGroup{}

	id, err := primitive.ObjectIDFromHex(storeGroupID)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{Key: "_id", Value: id}}

	err = s.collection.FindOne(ctx, filter).Decode(&storeGroup)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, drivers2.ErrNotFound
		}
		return nil, err
	}

	return storeGroup.DirectPromoBanners, nil
}

func (s *MongoRepository) DeleteDirectPromoBannerByID(ctx context.Context, storeGroupID, directPromoID string) error {

	oid, err := primitive.ObjectIDFromHex(storeGroupID)
	if err != nil {
		return err
	}
	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "direct_promo_banners.promo_id", Value: directPromoID},
	}

	update := bson.D{
		{
			Key: "$pull",
			Value: bson.D{
				{Key: "direct_promo_banners", Value: bson.D{{Key: "promo_id", Value: directPromoID}}},
			},
		},
	}
	res, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return drivers2.ErrNotFound
	}
	return nil
}
