package mongo

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (repo *MenuRepository) GetMenuSection(ctx context.Context, query selector.Menu) ([]*models.Section, error) {
	filter, err := repo.filterFrom(query)
	if err != nil {
		return nil, err
	}

	opts := options.FindOne().SetProjection(bson.D{
		{Key: "sections", Value: bson.D{
			{Key: "$elemMatch", Value: bson.D{
				{Key: "ext_id", Value: query.SectionID},
			}}}},
	})

	var result struct {
		Sections []*models.Section `bson:"sections"`
	}

	if err = repo.menuColl.FindOne(ctx, filter, opts).Decode(&result); err != nil {
		return nil, errorSwitch(err)
	}
	return result.Sections, nil
}

func (repo *MenuRepository) GetMenuSections(ctx context.Context, query selector.Menu) ([]*models.Section, error) {
	filter, err := repo.filterFrom(query)
	if err != nil {
		return nil, err
	}

	opts := options.FindOne().SetProjection(bson.D{
		{Key: "sections", Value: 1},
	})

	var result struct {
		Sections []*models.Section `bson:"sections"`
	}

	if err = repo.menuColl.FindOne(ctx, filter, opts).Decode(&result); err != nil {
		return nil, errorSwitch(err)
	}
	return result.Sections, nil
}

func (repo *MenuRepository) UpdateSection(ctx context.Context, menuID string, section models.Section) error {
	oid, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return drivers.ErrInvalid
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "sections.ext_id", Value: section.ExtID},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "sections.$", Value: section}},
		}}

	res, err := repo.menuColl.UpdateOne(ctx, filter, update)
	if err != nil {
		return errorSwitch(err)
	}

	if res.MatchedCount == 0 {
		return drivers.ErrNotFound
	}

	return nil
}
