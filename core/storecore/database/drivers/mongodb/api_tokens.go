package mongodb

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/storecore/database/drivers"
	"github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApiTokensRepo struct {
	collection *mongo.Collection
}

func NewApiTokensRepository(collection *mongo.Collection) drivers.ApiTokensRepository {
	return &ApiTokensRepo{
		collection: collection,
	}
}

func (a *ApiTokensRepo) GetStores(ctx context.Context, query selector.ApiToken) ([]models.Store, error) {
	pipe, err := a.filterFrom(query)
	if err != nil {
		return nil, err
	}

	cur, err := a.collection.Aggregate(ctx, pipe)
	if err != nil {
		return nil, errorSwitch(err)
	}
	defer closeCur(cur)

	type Output struct {
		Stores []models.Store `bson:"stores"`
	}
	var output Output
	var res []models.Store
	for cur.Next(ctx) {
		if err := cur.Decode(&output); err != nil {
			return nil, errorSwitch(err)
		}
		res = append(res, output.Stores...)
	}

	return res, nil
}

func (s *ApiTokensRepo) filterFrom(query selector.ApiToken) ([]bson.D, error) {
	var res []bson.D

	if query.HasToken() {
		res = append(res, bson.D{{Key: "$match", Value: bson.D{{Key: "token", Value: query.Token}}}})
	}

	return append(res, restaurantGroupsLookup(), unionOrgIds(), restaurantsLookup(),
		unionDocumentsGroup(), unionRestaurants()), nil
}

func restaurantGroupsLookup() bson.D {
	return bson.D{
		{Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "restaurant_groups"},
				{Key: "let", Value: bson.D{{Key: "vid", Value: "$organization_group_ids"}}},
				{Key: "pipeline", Value: bson.A{
					bson.D{
						{Key: "$match",
							Value: bson.D{
								{Key: "$expr",
									Value: bson.D{
										{Key: "$in",
											Value: bson.A{
												"$_id",
												bson.D{
													{Key: "$map",
														Value: bson.D{
															{Key: "input", Value: "$$vid"},
															{Key: "in", Value: bson.D{{Key: "$toObjectId", Value: "$$this"}}},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				},
				{Key: "as", Value: "restaurant_groups"},
			},
		},
	}
}

func unionOrgIds() bson.D {
	return bson.D{
		{Key: "$project",
			Value: bson.D{
				{Key: "token", Value: "$token"},
				{Key: "name", Value: "$name"},
				{Key: "result",
					Value: bson.D{
						{Key: "$setUnion",
							Value: bson.A{
								"$organization_ids",
								bson.D{
									{Key: "$reduce",
										Value: bson.D{
											{Key: "input",
												Value: bson.D{
													{Key: "$concatArrays",
														Value: bson.A{
															"$restaurant_groups.restaurant_ids",
														},
													},
												},
											},
											{Key: "initialValue", Value: bson.A{}},
											{Key: "in",
												Value: bson.D{
													{Key: "$concatArrays",
														Value: bson.A{
															"$$value",
															"$$this",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func restaurantsLookup() bson.D {
	return bson.D{
		{Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: "restaurants"},
				{Key: "let", Value: bson.D{{Key: "vid", Value: "$result"}}},
				{Key: "pipeline",
					Value: bson.A{
						bson.D{
							{Key: "$match",
								Value: bson.D{
									{Key: "$expr",
										Value: bson.D{
											{Key: "$in",
												Value: bson.A{
													"$_id",
													bson.D{
														{Key: "$map",
															Value: bson.D{
																{Key: "input", Value: "$$vid"},
																{Key: "in", Value: bson.D{{Key: "$toObjectId", Value: "$$this"}}},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				{Key: "as", Value: "stores"},
			},
		},
	}
}

func unionDocumentsGroup() bson.D {
	return bson.D{
		{Key: "$group",
			Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "stores", Value: bson.D{{Key: "$push", Value: "$stores"}}},
			},
		},
	}
}

func unionRestaurants() bson.D {
	return bson.D{
		{Key: "$project",
			Value: bson.D{
				{Key: "stores",
					Value: bson.D{
						{Key: "$reduce",
							Value: bson.D{
								{Key: "input", Value: "$stores"},
								{Key: "initialValue", Value: bson.A{}},
								{Key: "in",
									Value: bson.D{
										{Key: "$setUnion",
											Value: bson.A{
												"$$value",
												"$$this",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
