package mongo

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/pkg/menu/dto"
	"github.com/kwaaka-team/orders-core/service/entity_changes_history"
	entityChangesHistoryModels "github.com/kwaaka-team/orders-core/service/entity_changes_history/models"
	"github.com/rs/zerolog/log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MenuRepository struct {
	menuColl                 *mongo.Collection
	entityChangesHistoryRepo entity_changes_history.Repository
}

func NewMenuRepository(menuColl *mongo.Collection, entityChangesHistoryRepo entity_changes_history.Repository) *MenuRepository {
	return &MenuRepository{
		menuColl:                 menuColl,
		entityChangesHistoryRepo: entityChangesHistoryRepo,
	}
}

func (repo *MenuRepository) Get(ctx context.Context, query selector.Menu) (models.Menu, error) {
	filter, err := repo.filterFrom(query)
	if err != nil {
		return models.Menu{}, drivers.ErrInvalid
	}

	var menu models.Menu
	if err = repo.menuColl.FindOne(ctx, filter).Decode(&menu); err != nil {
		return models.Menu{}, errorSwitch(err)
	}

	return menu, nil
}

func (repo *MenuRepository) List(ctx context.Context, query selector.Menu) ([]models.Menu, error) {
	filter, err := repo.filterFrom(query)
	if err != nil {
		return nil, drivers.ErrInvalid
	}

	cur, err := repo.menuColl.Find(ctx, filter)
	if err != nil {
		return nil, errorSwitch(err)
	}

	menus := make([]models.Menu, 0, cur.RemainingBatchLength())
	for cur.Next(ctx) {
		var menu models.Menu

		if err = cur.Decode(&menu); err != nil {
			return nil, errorSwitch(err)
		}

		menus = append(menus, menu)
	}

	return menus, nil

}

func (repo *MenuRepository) GetMenuIDs(ctx context.Context, query selector.Menu) ([]string, error) {

	filter, err := repo.filterFrom(query)
	if err != nil {
		return nil, drivers.ErrInvalid
	}

	curs, err := repo.menuColl.Distinct(ctx, "_id", filter)
	if err != nil {
		return nil, errorSwitch(err)
	}

	if len(curs) == 0 {
		return nil, drivers.ErrNotFound
	}

	ids := make([]string, 0, len(curs))
	for _, cur := range curs {
		id, ok := cur.(primitive.ObjectID)
		if !ok {
			continue
		}
		ids = append(ids, id.Hex())
	}

	return ids, nil

}

func (repo *MenuRepository) GetIDByName(ctx context.Context, name string) (string, error) {

	filter := bson.D{
		{Key: "name", Value: name},
	}

	opts := options.FindOne().
		SetProjection(bson.D{
			{Key: "_id", Value: 1},
		})

	var id string
	if err := repo.menuColl.FindOne(ctx, filter, opts).Decode(&id); err != nil {
		return "", errorSwitch(err)
	}

	return id, nil
}

func (repo *MenuRepository) Insert(ctx context.Context, menu models.Menu) (string, error) {

	menu.CreatedAt = coreModels.TimeNow()
	menu.UpdatedAt = coreModels.TimeNow()

	log.Info().Msgf("(MenuRepository) Insert menu: %#v", menu)

	res, err := repo.menuColl.InsertOne(ctx, menu)
	if err != nil {
		log.Err(err).Msgf("(MenuRepository) Insert error")
		return "", errorSwitch(err)
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", drivers.ErrInvalid
	}

	return oid.Hex(), nil
}

func (repo *MenuRepository) setMenuHistory(ctx context.Context, menuId string, history entityChangesHistoryModels.EntityChangesHistory, operationType, repositoryMethod string) {
	oldMenu, err := repo.Get(ctx, selector.EmptyMenuSearch().SetMenuID(menuId))
	if err != nil {
		log.Err(err).Msgf("(MenuRepository) Get menu for entity changes history error")
	} else {
		history.OldBody = oldMenu
		history.ModifiedAt = time.Now().UTC()
		history.OperationType = operationType
		history.RepositoryMethod = repositoryMethod
		history.CollectionName = "menus"

		_, err = repo.entityChangesHistoryRepo.InsertHistory(ctx, history)
		if err != nil {
			log.Err(err).Msgf("insert entity changes history error")
		}
	}
}

func (repo *MenuRepository) Update(ctx context.Context, menu models.Menu, history entityChangesHistoryModels.EntityChangesHistory) error {
	repo.setMenuHistory(ctx, menu.ID, history, "update", "Update")

	oid, err := primitive.ObjectIDFromHex(menu.ID)
	if err != nil {
		return drivers.ErrInvalid
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	menu.UpdatedAt = coreModels.TimeNow()
	update := bson.D{
		{Key: "$set", Value: menu.ToUpdate()},
	}

	log.Info().Msgf("(MenuRepository) Update menu: %#v", menu)

	res, err := repo.menuColl.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Err(err).Msgf("(MenuRepository) Update error")
		return errorSwitch(err)
	}

	if res.MatchedCount == 0 {
		return drivers.ErrNotFound
	}
	return nil
}

func (repo *MenuRepository) AddRowToAttributeGroup(ctx context.Context, menuId string, attributeMinMax []models.AttributeIdMinMax, attributeGroupID string) error {

	oid, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "attributes_groups.ext_id", Value: attributeGroupID},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "attributes_groups.$.attribute_min_max", Value: attributeMinMax},
		}},
	}

	res, err := repo.menuColl.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Err(err).Msgf("(MenuRepository) AddRowToAttributeGroup error")
		return err
	}

	if res.MatchedCount == 0 {
		return drivers.ErrNotFound
	}

	return nil
}

// TODO: need to test
func (repo *MenuRepository) Upsert(ctx context.Context, req models.Menu) (models.Menu, error) {

	filter := bson.M{
		"name": req.Name,
	}

	update := bson.M{
		"$set": req,
	}

	opts := options.
		FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var menu models.Menu
	if err := repo.menuColl.FindOneAndUpdate(ctx, filter, update, opts).Decode(&menu); err != nil {
		log.Err(err).Msgf("(MenuRepository) Upsert error")
		return models.Menu{}, errorSwitch(err)
	}

	return menu, nil
}

func (repo *MenuRepository) Delete(ctx context.Context, menuID string) error {
	oid, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return drivers.ErrInvalid
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}
	_, err = repo.menuColl.DeleteOne(ctx, filter)
	if err != nil {
		return errorSwitch(err)
	}

	return nil
}

func (repo *MenuRepository) GetGroups(ctx context.Context, query selector.Menu) (models.Groups, error) {

	oid, err := primitive.ObjectIDFromHex(query.ID)
	if err != nil {
		return nil, drivers.ErrInvalid
	}

	match := bson.D{{Key: "_id", Value: oid}}
	unwind := "$groups"
	project := bson.D{{Key: "groups", Value: 1}}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$unwind", Value: unwind}},
		{{Key: "$project", Value: project}},
	}

	cur, err := repo.menuColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer closeCur(cur)

	groups := make(models.Groups, 0, cur.RemainingBatchLength())

	for cur.Next(ctx) {
		var temp struct {
			Group models.Group `bson:"groups"`
		}
		if err = cur.Decode(&temp); err != nil {
			return nil, err
		}

		groups = append(groups, temp.Group)
	}

	return groups, nil
}

func (repo *MenuRepository) filterFrom(query selector.Menu) (bson.D, error) {

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

func (repo *MenuRepository) sortFrom(query selector.Sorting) bson.D {
	sort := make(bson.D, 0, 2)

	if query.HasSorting() {
		sort = append(sort, bson.E{Key: query.Param, Value: query.Direction})
	}
	sort = append(sort, bson.E{Key: "_id", Value: 1})

	return sort
}

func (repo *MenuRepository) UpdateMenuName(ctx context.Context, query models.UpdateMenuName) error {
	oid, err := primitive.ObjectIDFromHex(query.MenuID)
	if err != nil {
		return drivers.ErrInvalid
	}

	filter := bson.D{{Key: "_id", Value: oid}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: query.MenuName}}}}

	res, err := repo.menuColl.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Err(err).Msgf("(MenuRepository) UpdateMenuName error")
		return errorSwitch(err)
	}

	if res.MatchedCount == 0 {
		return drivers.ErrNotFound
	}
	return nil
}

func (repo *MenuRepository) CreateGlovoSuperCollection(ctx context.Context, menuId string, superCollections dto.MenuSuperCollections) error {
	oid, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		return drivers.ErrInvalid
	}

	for i := range superCollections {
		superCollections[i].ExtId = primitive.NewObjectID().Hex()
	}

	filter := bson.D{{Key: "_id", Value: oid}}

	fields := repo.constructSuperCollectionUpdateBody(superCollections)
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "super_collections", Value: fields}}}}

	res, err := repo.menuColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return errorSwitch(err)
	}

	if res.MatchedCount == 0 {
		return drivers.ErrNotFound
	}

	return nil
}

func (repo *MenuRepository) constructSuperCollectionUpdateBody(superCollections dto.MenuSuperCollections) bson.A {
	var res bson.A

	for _, sc := range superCollections {
		res = append(res, bson.D{
			bson.E{Key: "ext_id", Value: sc.ExtId},
			bson.E{Key: "name", Value: sc.Name},
			bson.E{Key: "img_url", Value: sc.ImgUrl},
			bson.E{Key: "order", Value: sc.SuperCollectionOrder},
		})
	}

	return res
}
