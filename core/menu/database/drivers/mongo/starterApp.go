package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (repo *MenuRepository) UpdateAttributeStarterAppIDByExtID(ctx context.Context, menuID, extID, starterAppID string) error {
	objID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objID}}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "attributes.$[attribute].starter_app_id", Value: starterAppID},
		}},
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{Key: "attribute.ext_id", Value: extID}},
		},
	})

	_, err = repo.menuColl.UpdateOne(ctx, filter, update, arrayFilters)

	return err
}

func (repo *MenuRepository) UpdateAttributeStarterAppOfferIDByExtID(ctx context.Context, menuID, extID, starterAppOfferID string) error {
	objID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objID}}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "attributes.$[attribute].starter_app_offer_id", Value: starterAppOfferID},
		}},
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{Key: "attribute.ext_id", Value: extID}},
		},
	})

	_, err = repo.menuColl.UpdateOne(ctx, filter, update, arrayFilters)

	return err
}

func (repo *MenuRepository) UpdateProductStarterAppIDByExtID(ctx context.Context, menuID, extID, starterAppID string) error {
	objID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objID}}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "products.$[product].starter_app_id", Value: starterAppID},
		}},
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{Key: "product.ext_id", Value: extID}},
		},
	})

	_, err = repo.menuColl.UpdateOne(ctx, filter, update, arrayFilters)

	return err
}

func (repo *MenuRepository) UpdateProductStarterAppOfferIDByExtID(ctx context.Context, menuID, extID, starterAppOfferID string) error {
	objID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objID}}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "products.$[product].starter_app_offer_id", Value: starterAppOfferID},
		}},
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{Key: "product.ext_id", Value: extID}},
		},
	})

	_, err = repo.menuColl.UpdateOne(ctx, filter, update, arrayFilters)

	return err
}

func (repo *MenuRepository) UpdateAttributeGroupStarterAppIDByExtID(ctx context.Context, menuID, extID, starterAppID string) error {
	objID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objID}}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "attributes_groups.$[ag].starter_app_id", Value: starterAppID},
		}},
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{Key: "ag.ext_id", Value: extID}},
		},
	})

	_, err = repo.menuColl.UpdateOne(ctx, filter, update, arrayFilters)

	return err
}

func (repo *MenuRepository) UpdateSectionStarterAppIDByExtID(ctx context.Context, menuID, extID, starterAppID string) error {
	objID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objID}}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "sections.$[section].starter_app_id", Value: starterAppID},
		}},
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{Key: "section.ext_id", Value: extID}},
		},
	})

	_, err = repo.menuColl.UpdateOne(ctx, filter, update, arrayFilters)
	return err
}

func (repo *MenuRepository) UpdateCollectionStarterAppIDByExtID(ctx context.Context, menuID, extID, starterAppID string) error {
	objID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objID}}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "collections.$[collection].starter_app_id", Value: starterAppID},
		}},
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{Key: "collection.ext_id", Value: extID}},
		},
	})

	_, err = repo.menuColl.UpdateOne(ctx, filter, update, arrayFilters)

	return err
}

func (repo *MenuRepository) UpdateSuperCollectionStarterAppIDByExtID(ctx context.Context, menuID, extID, starterAppID string) error {
	objID, err := primitive.ObjectIDFromHex(menuID)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: objID}}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "super_collections.$[sc].starter_app_id", Value: starterAppID},
		}},
	}

	arrayFilters := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{Key: "sc.ext_id", Value: extID}},
		},
	})

	_, err = repo.menuColl.UpdateOne(ctx, filter, update, arrayFilters)

	return err
}
