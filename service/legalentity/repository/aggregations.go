package repository

import (
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	"github.com/kwaaka-team/orders-core/service/legalentity/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetAllDocumentsByLegalEntityID aggregation pipeline
func (r *Repo) getAllDocumentsByLegalEntityID(legalEntityID primitive.ObjectID, pagination selector.Pagination, filter models.DocumentFilter) []bson.D {
	var pipeline []bson.D

	pipeline = append(pipeline, matchGetAllDocuments(legalEntityID), unwindDocumentsField(), matchAfterUnwindGetAllDocuments())

	if pagination.HasPagination() {
		pipeline = append(pipeline, applyPagination(pagination)...)
	}

	pipeline = append(pipeline, applyFiltersGetAllDocumentsByLegalEntityID(filter), projectionGetAllDocuments())

	return pipeline
}

func applyFiltersGetAllDocumentsByLegalEntityID(filter models.DocumentFilter) bson.D {
	var withFilters bson.D

	docType := bson.E{Key: "$match", Value: bson.D{{Key: "documents.type", Value: filter.DocType}}}

	return append(withFilters, docType)
}

func matchGetAllDocuments(legalEntityID primitive.ObjectID) bson.D {
	return bson.D{{Key: "$match", Value: bson.D{
		{Key: "_id", Value: legalEntityID}}},
	}
}

func matchAfterUnwindGetAllDocuments() bson.D {
	return bson.D{{Key: "$match", Value: bson.D{
		{Key: "documents.status", Value: ACTIVE}}},
	}
}

func unwindDocumentsField() bson.D {
	return bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$documents"}}}}
}

func projectionGetAllDocuments() bson.D {
	return bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 0}, {Key: "documents", Value: 1}}}}
}

// GetListOfStores aggregation pipeline
func (r *Repo) getListOfStores(pagination selector.Pagination, objectID primitive.ObjectID) []bson.D {
	var pipeline []bson.D

	if pagination.HasPagination() {
		pipeline = append(pipeline, applyPagination(pagination)...)
	}

	pipeline = append(pipeline, matchListStores(objectID), addFieldsStoreIDs(), lookupStores(), addFieldsBrandIDs(),
		lookupBrands(), projectionListStores())

	return pipeline
}

func matchListStores(objectID primitive.ObjectID) bson.D {
	return bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: objectID}}}}
}

func projectionListStores() bson.D {
	return bson.D{
		{Key: "$project",
			Value: bson.D{
				{Key: "_id", Value: 1},
				{Key: "stores", Value: 1},
				{Key: "brands", Value: 1},
			},
		},
	}
}

// GetListOfLegalEntities aggregation pipeline
func (r *Repo) getListOfLegalEntities(pagination selector.Pagination, filter models.Filter) []bson.D {
	var pipeline []bson.D

	if pagination.HasPagination() {
		pipeline = append(pipeline, applyPagination(pagination)...)
	}

	pipeline = append(pipeline, addFieldsStoreIDs(), lookupStoresListLegalEntities(), addFieldsBrandIDs(),
		lookupBrands(), projectionList())

	filtered := applyFiltersListOfLegalEntities(filter)
	if filtered != nil {
		pipeline = append(pipeline, filtered)
	}

	return pipeline
}

func applyFiltersListOfLegalEntities(filter models.Filter) bson.D {
	var withFilters bson.D

	if filter.Search != "" {
		or := bson.E{Key: "$or", Value: bson.A{
			bson.D{{Key: "name", Value: bson.D{{Key: "$regex", Value: filter.Search}, {Key: "$options", Value: "i"}}}},
			bson.D{{Key: "brands", Value: bson.D{{Key: "$elemMatch", Value: bson.D{{Key: "name", Value: bson.D{{Key: "$regex", Value: filter.Search}, {Key: "$options", Value: "i"}}}}}}}},
		}}

		withFilters = append(withFilters, or)
	}

	if filter.ContactName != "" {
		contactName := bson.E{Key: "contacts.full_name", Value: bson.D{{Key: "$regex", Value: filter.ContactName}, {Key: "$options", Value: "i"}}}

		withFilters = append(withFilters, contactName)
	}

	if len(filter.PaymentType) > 0 {
		paymentType := bson.E{Key: "payment_type", Value: bson.D{{Key: "$in", Value: filter.PaymentType}}}

		withFilters = append(withFilters, paymentType)
	}

	if len(filter.Status) > 0 {
		status := bson.E{Key: "status", Value: bson.D{{Key: "$in", Value: filter.Status}}}

		withFilters = append(withFilters, status)
	}

	if len(withFilters) != 0 {
		return bson.D{{Key: "$match", Value: withFilters}}
	}

	return nil
}

func lookupStoresListLegalEntities() bson.D {
	return bson.D{
		{Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: storesCollection},
				{Key: "let", Value: bson.D{{Key: "st_ids", Value: "$stIds"}}},
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
													"$$st_ids",
												},
											},
										},
									},
								},
							},
						},
						bson.D{
							{Key: "$project",
								Value: bson.D{
									{Key: "_id", Value: 0},
									{Key: "name", Value: 1},
									{Key: "restaurant_group_id", Value: 1},
									{Key: "fare", Value: 1},
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

func projectionList() bson.D {
	return bson.D{
		{Key: "$project",
			Value: bson.D{
				{Key: "_id", Value: 1},
				{Key: "name", Value: 1},
				{Key: "brands", Value: 1},
				{Key: "payment_type", Value: 1},
				{Key: "contacts", Value: 1},
				{Key: "status", Value: 1},
				{Key: "stores", Value: 1},
			},
		},
	}
}

// GetLegalEntityViewPipeline aggregation pipeline
func (r *Repo) getLegalEntityView(objectId primitive.ObjectID) []bson.D {
	var pipeline []bson.D

	pipeline = append(pipeline, matchGetLegalEntityView(objectId), addFieldsStoreIDs(), lookupStores(), addFieldsBrandIDs(), lookupBrands(),
		lookupManagers(), lookupSales(), projectionViewLegalEntity())

	return pipeline
}

func matchGetLegalEntityView(objectId primitive.ObjectID) bson.D {
	return bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: objectId}}}}
}

func addFieldsStoreIDs() bson.D {
	return bson.D{
		{Key: "$addFields",
			Value: bson.D{
				{Key: "stIds",
					Value: bson.D{
						{Key: "$map",
							Value: bson.D{
								{Key: "input", Value: "$store_ids"},
								{Key: "in", Value: bson.D{{Key: "$toObjectId", Value: "$$this"}}},
							},
						},
					},
				},
			},
		},
	}
}

func lookupStores() bson.D {
	return bson.D{
		{Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: storesCollection},
				{Key: "let", Value: bson.D{{Key: "st_ids", Value: "$stIds"}}},
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
													"$$st_ids",
												},
											},
										},
									},
								},
							},
						},
						bson.D{
							{Key: "$project",
								Value: bson.D{
									{Key: "_id", Value: 0},
									{Key: "name", Value: 1},
									{Key: "restaurant_group_id", Value: 1},
									{Key: "fare", Value: 1},
									{Key: "wolt_exists",
										Value: bson.D{
											{Key: "$cond",
												Value: bson.A{
													bson.D{
														{Key: "$ifNull",
															Value: bson.A{
																bson.D{
																	{Key: "$size",
																		Value: bson.D{
																			{Key: "$ifNull",
																				Value: bson.A{
																					"$wolt.store_id",
																					bson.A{},
																				},
																			},
																		},
																	},
																},
																false,
															},
														},
													},
													true,
													false,
												},
											},
										},
									},
									{Key: "glovo_exists",
										Value: bson.D{
											{Key: "$cond",
												Value: bson.A{
													bson.D{
														{Key: "$ifNull",
															Value: bson.A{
																bson.D{
																	{Key: "$size",
																		Value: bson.D{
																			{Key: "$ifNull",
																				Value: bson.A{
																					"$glovo.store_id",
																					bson.A{},
																				},
																			},
																		},
																	},
																},
																false,
															},
														},
													},
													true,
													false,
												},
											},
										},
									},
									{Key: "yandex_exists",
										Value: bson.D{
											{Key: "$cond",
												Value: bson.A{
													bson.D{
														{Key: "$gt",
															Value: bson.A{
																bson.D{
																	{Key: "$size",
																		Value: bson.D{
																			{Key: "$filter",
																				Value: bson.D{
																					{Key: "input",
																						Value: bson.D{
																							{Key: "$ifNull",
																								Value: bson.A{
																									"$external",
																									bson.A{},
																								},
																							},
																						},
																					},
																					{Key: "as", Value: "ext"},
																					{Key: "cond",
																						Value: bson.D{
																							{Key: "$eq",
																								Value: bson.A{
																									"$$ext.type",
																									"yandex",
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																},
																0,
															},
														},
													},
													true,
													false,
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

func addFieldsBrandIDs() bson.D {
	return bson.D{
		{Key: "$addFields",
			Value: bson.D{
				{Key: "brandIds",
					Value: bson.D{
						{Key: "$map",
							Value: bson.D{
								{Key: "input", Value: "$stores"},
								{Key: "as", Value: "store"},
								{Key: "in", Value: bson.D{{Key: "$toObjectId", Value: "$$store.restaurant_group_id"}}},
							},
						},
					},
				},
			},
		},
	}
}

func lookupBrands() bson.D {
	return bson.D{
		{Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: brandsCollection},
				{Key: "let", Value: bson.D{{Key: "b_ids", Value: "$brandIds"}}},
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
													"$$b_ids",
												},
											},
										},
									},
								},
							},
						},
						bson.D{
							{Key: "$project",
								Value: bson.D{
									{Key: "name", Value: 1},
									{Key: "_id", Value: 0},
								},
							},
						},
					},
				},
				{Key: "as", Value: "brands"},
			},
		},
	}
}

func lookupManagers() bson.D {
	return bson.D{
		{Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: managersCollection},
				{Key: "let", Value: bson.D{{Key: "m_id", Value: "$linked_acc_manager"}}},
				{Key: "pipeline",
					Value: bson.A{
						bson.D{{Key: "$addFields", Value: bson.D{{Key: "mid", Value: bson.D{{Key: "$toObjectId", Value: "$$m_id"}}}}}},
						bson.D{
							{Key: "$match",
								Value: bson.D{
									{Key: "$expr",
										Value: bson.D{
											{Key: "$eq",
												Value: bson.A{
													"$mid",
													"$_id",
												},
											},
										},
									},
								},
							},
						},
						bson.D{{Key: "$project", Value: bson.D{{Key: "mid", Value: 0}}}},
					},
				},
				{Key: "as", Value: "manager"},
			},
		},
	}
}

func lookupSales() bson.D {
	return bson.D{
		{Key: "$lookup",
			Value: bson.D{
				{Key: "from", Value: salesCollection},
				{Key: "let", Value: bson.D{{Key: "s_id", Value: "$sales_id"}}},
				{Key: "pipeline",
					Value: bson.A{
						bson.D{{Key: "$addFields", Value: bson.D{{Key: "sid", Value: bson.D{{Key: "$toObjectId", Value: "$$s_id"}}}}}},
						bson.D{
							{Key: "$match",
								Value: bson.D{
									{Key: "$expr",
										Value: bson.D{
											{Key: "$eq",
												Value: bson.A{
													"$sid",
													"$_id",
												},
											},
										},
									},
								},
							},
						},
						bson.D{{Key: "$project", Value: bson.D{{Key: "sid", Value: 0}}}},
					},
				},
				{Key: "as", Value: "sales"},
			},
		},
	}
}

func projectionViewLegalEntity() bson.D {
	return bson.D{
		{Key: "$project",
			Value: bson.D{
				{Key: "_id", Value: 1},
				{Key: "name", Value: 1},
				{Key: "bin", Value: 1},
				{Key: "knp", Value: 1},
				{Key: "payment_type", Value: 1},
				{Key: "contacts", Value: 1},
				{Key: "sales_comment", Value: 1},
				{Key: "status", Value: 1},
				{Key: "payment_cycle", Value: 1},
				{Key: "first_payment_at", Value: 1},
				{Key: "manager", Value: 1},
				{Key: "sales", Value: 1},
				{Key: "stores", Value: 1},
				{Key: "brands", Value: 1},
			},
		},
	}
}
