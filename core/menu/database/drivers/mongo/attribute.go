package mongo

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAttributes get all attributes from menu_id
func (repo *MenuRepository) GetAttributes(ctx context.Context, query selector.Menu) (models.Attributes, int, error) {
	oid, err := primitive.ObjectIDFromHex(query.ID)
	if err != nil {
		return nil, 0, err
	}

	match := bson.D{{Key: "_id", Value: oid}}
	unwind := "$attributes"
	project := bson.D{{Key: "attributes", Value: 1}}
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$unwind", Value: unwind}},
		{{Key: "$project", Value: project}},
	}

	countPipeline := pipeline
	countPipeline = append(countPipeline, bson.D{{Key: "$count", Value: "total_count"}})

	count, err := repo.menuColl.Aggregate(ctx, countPipeline)
	if err != nil {
		return nil, 0, err
	}

	var totalCount struct {
		TotalCount int `bson:"total_count"`
	}

	for count.Next(ctx) {
		if err = count.Decode(&totalCount); err != nil {
			return nil, 0, err
		}
	}

	if query.HasPagination() {
		pipeline = append(pipeline,
			bson.D{{Key: "$skip", Value: query.Skip()}},
			bson.D{{Key: "$limit", Value: query.Pagination.Limit}})
	}

	cur, err := repo.menuColl.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, err
	}
	defer closeCur(cur)

	result := make(models.Attributes, 0, cur.RemainingBatchLength())

	for cur.Next(ctx) {
		var tmp struct {
			Attribute models.Attribute `bson:"attributes"`
		}

		if err = cur.Decode(&tmp); err != nil {
			return nil, 0, err
		}

		result = append(result, tmp.Attribute)
	}

	return result, totalCount.TotalCount, nil
}

// GetAttributeGroups get all attribute_groups from menu_id
func (repo *MenuRepository) GetAttributeGroups(ctx context.Context, query selector.Menu) (models.AttributeGroups, error) {
	oid, err := primitive.ObjectIDFromHex(query.ID)
	if err != nil {
		return nil, err
	}

	match := bson.D{{Key: "_id", Value: oid}}
	unwind := "$attributes_groups"
	project := bson.D{{Key: "attributes_groups", Value: 1}}
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

	result := make(models.AttributeGroups, 0, cur.RemainingBatchLength())

	for cur.Next(ctx) {
		var tmp struct {
			AttributeGroup models.AttributeGroup `bson:"attributes_groups"`
		}

		if err = cur.Decode(&tmp); err != nil {
			return nil, err
		}

		result = append(result, tmp.AttributeGroup)
	}

	return result, nil
}

func (repo *MenuRepository) DeleteAttributeGroupFromDB(ctx context.Context, menuId string, attrGroupExtId string) error {
	oid, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: oid}}

	update := bson.D{
		{Key: "$pull", Value: bson.D{{Key: "attributes_groups", Value: bson.D{{Key: "ext_id", Value: attrGroupExtId}}}}},
	}

	_, err = repo.menuColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (repo *MenuRepository) ValidateAttributeGroupName(ctx context.Context, menuId, name string) (bool, error) {

	oid, err := primitive.ObjectIDFromHex(menuId)
	if err != nil {
		return false, err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "attributes_groups.name", Value: name},
	}

	var result models.Menu
	err = repo.menuColl.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (repo *MenuRepository) CreateAttributeGroup(ctx context.Context, menuID string, attribute models.Attribute) (string, error) {

	oid, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return "", err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.M{
		"$push": bson.M{
			"attributes_groups": attribute,
		},
	}

	_, err = repo.menuColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return "", err
	}

	return attribute.ExtID, nil
}
