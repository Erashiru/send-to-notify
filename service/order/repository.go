package order

import (
	"context"
	errorsGo "errors"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/selector"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const orderCollectionName = "orders"

type Repository interface {
	UpdateOrderStatusByID(ctx context.Context, orderID string, newStatus string) error
	UpdateOrder(ctx context.Context, order models.Order) error
	InsertOrder(ctx context.Context, order models.Order) (models.Order, error)
	FindOrderByPosOrderID(ctx context.Context, posOrderID string) (models.Order, error)
	FindOrderByID(ctx context.Context, id string) (models.Order, error)
	FindOrderByDeliveryOrderId(ctx context.Context, id string) (models.Order, error)
	FindActiveOrdersByPosType(ctx context.Context, posType string, orderTimeFrom time.Time) ([]models.Order, error)
	GetOrderNumber(ctx context.Context, query selector.Order) (int64, error)
	GetAllOrders(ctx context.Context, query selector.Order) ([]models.Order, int, error)
	AddStatusToHistory(ctx context.Context, id, status string) error
	UpdateOrderDeliveryIDByOrderID(ctx context.Context, orderID string, deliveryOrderID string) error
	GetOrdersByStatusesAndPosType(ctx context.Context, posType string, statuses []string) ([]models.Order, error)
	FindOrderByOrderID(ctx context.Context, orderID string) (models.Order, error)
	GetAverageBill(ctx context.Context, order selector.Order) (float64, error)
	Get3plOrdersForCron(ctx context.Context, indriveCallTime int64) ([]models.Order, error)
	SetProposals(ctx context.Context, orderID string, proposals []models.Proposal) error
	SetDeliveryDispatcherPrice(ctx context.Context, orderID string, deliveryDispatcherPrice float64) error
	UpdateOrderDeferStatus(ctx context.Context, status bool, id string) error
	GetOrderBy3plDeliveryID(ctx context.Context, deliveryOrderId string) (models.Order, error)
	Save3plDeliveryHistoryAndSetEmptyDispatcherService(ctx context.Context, id string, cancelled3plDeliveryInfo models.History3plDelivery, newDeliveryAddress models.DeliveryAddress, newCustomer models.Customer) error
	GetIIKO3plOrdersForCron(ctx context.Context, callTime int64) ([]models.Order, error)
	SetCancelledDeliveryDispatcherPrice(ctx context.Context, orderID string, deliveryDispatcherPrice float64) error
	FindOrderByDeliveryOrderID(ctx context.Context, deliveryID string) (models.Order, error)
}

var nonActiveOrderStatuses = []string{
	string(models.STATUS_CANCELLED),
	models.CLOSED.String(),
	models.FAILED.String(),
	models.CANCELLED_BY_POS_SYSTEM.String(),
	string(models.STATUS_CANCELLED_BY_DELIVERY_SERVICE),
	string(models.STATUS_SKIPPED),
}

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) (*MongoRepository, error) {
	r := MongoRepository{
		collection: db.Collection(orderCollectionName),
	}
	return &r, nil
}

func (r *MongoRepository) UpdateOrderDeliveryIDByOrderID(ctx context.Context, orderID string, deliveryOrderID string) error {

	oid, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return err
	}

	filter := bson.D{{Key: "_id", Value: oid}}

	updateFields := bson.D{{Key: "delivery_order_id", Value: deliveryOrderID},
		{Key: "updated_at", Value: time.Now().UTC()}}

	update := bson.D{{Key: "$set", Value: updateFields}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *MongoRepository) FindActiveOrdersByPosType(ctx context.Context, posType string, orderTimeFrom time.Time) ([]models.Order, error) {
	filter := bson.D{
		{
			Key: "status",
			Value: bson.D{
				{
					Key:   "$nin",
					Value: nonActiveOrderStatuses,
				},
			},
		},
		{
			Key:   "pos_type",
			Value: posType,
		},
		{
			Key: "order_time.value",
			Value: bson.M{
				"$gte": primitive.NewDateTimeFromTime(orderTimeFrom),
			},
		},
		{
			Key: "pos_order_id",
			Value: bson.M{
				"$ne":     "",
				"$exists": true,
			},
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	orders := make([]models.Order, 0, cursor.RemainingBatchLength())
	if err = cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *MongoRepository) FindOrderByPosOrderID(ctx context.Context, posOrderID string) (models.Order, error) {
	query := selector.EmptyOrderSearch().SetPosOrderID(posOrderID)

	filter, err := r.filterFrom(query)
	if err != nil {
		return models.Order{}, errors.ErrorSwitch(err)
	}

	var order models.Order
	err = r.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		return models.Order{}, errors.ErrorSwitch(err)
	}

	return order, nil

}

func (r *MongoRepository) GetOrderNumber(ctx context.Context, query selector.Order) (int64, error) {
	filter, err := r.filterFrom(query)
	if err != nil {
		return 0, errors.ErrorSwitch(err)
	}

	result, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (r *MongoRepository) GetAllOrders(ctx context.Context, query selector.Order) ([]models.Order, int, error) {
	filter, err := r.filterFrom(query)
	if err != nil {
		return nil, 0, err
	}

	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	option := options.Find()
	if query.HasSorting() {
		option.SetSort(r.sortFrom(query.Sorting))
	}

	if query.HasPagination() {
		option.SetLimit(query.Pagination.Limit).SetSkip(query.Skip())
	}

	cursor, err := r.collection.Find(ctx, filter, option)
	if err != nil {
		return nil, 0, err
	}

	orders := make([]models.Order, 0, cursor.RemainingBatchLength())
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, 0, err
	}

	return orders, int(totalCount), nil
}

func (r *MongoRepository) FindOrderByID(ctx context.Context, id string) (models.Order, error) {
	query := selector.EmptyOrderSearch().SetID(id)

	filter, err := r.filterFrom(query)
	if err != nil {
		return models.Order{}, errors.ErrorSwitch(err)
	}

	var order models.Order
	err = r.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		return models.Order{}, errors.ErrorSwitch(err)
	}

	return order, nil

}

func (r *MongoRepository) FindOrderByDeliveryOrderId(ctx context.Context, id string) (models.Order, error) {
	query := selector.EmptyOrderSearch().SetDeliveryOrderId(id)

	filter, err := r.filterFrom(query)
	if err != nil {
		return models.Order{}, errors.ErrorSwitch(err)
	}

	var order models.Order
	err = r.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		return models.Order{}, errors.ErrorSwitch(err)
	}

	return order, nil
}

func (r *MongoRepository) InsertOrder(ctx context.Context, order models.Order) (models.Order, error) {

	order.IsActive = true
	order.CreatedAt = models.TimeNow()
	order.UpdatedAt = models.TimeNow()
	res, err := r.collection.InsertOne(ctx, order)
	if err != nil {
		return order, errors.ErrorSwitch(err)
	}

	order.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return order, nil
}

func (r *MongoRepository) UpdateOrder(ctx context.Context, order models.Order) error {
	oid, err := primitive.ObjectIDFromHex(order.ID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	// TODO: models.UpdateOrder, if you want update order and set, it will be error if order.ID was string while setting in primitive.ObjectID
	order.ID = ""
	order.UpdatedAt = models.TimeNow()

	update := bson.D{
		{Key: "$set", Value: order},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.ErrorSwitch(err)
	}

	if res.MatchedCount == 0 {
		return err
	}
	return nil
}

func (r *MongoRepository) UpdateOrderDeferStatus(ctx context.Context, status bool, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.M{"$set": bson.M{"is_defer_submission": status}}

	if _, err := r.collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (r *MongoRepository) UpdateOrderStatusByID(ctx context.Context, orderID string, newStatus string) error {
	query := selector.EmptyOrderSearch().SetID(orderID)
	filter, err := r.filterFrom(query)
	if err != nil {
		return errors.ErrorSwitch(err)
	}

	updateFields := bson.D{{Key: "status", Value: newStatus}}

	statusHistory := models.OrderStatusUpdate{
		Name: newStatus,
		Time: time.Now(),
	}

	update := bson.D{{Key: "$set", Value: updateFields}, {Key: "$push", Value: bson.M{"statuses_history": statusHistory}}}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (r *MongoRepository) filterFrom(query selector.Order) (bson.D, error) {
	result := make(bson.D, 0, 7)

	if query.HasID() {
		oid, err := primitive.ObjectIDFromHex(query.ID)
		if err != nil {
			return nil, errors.ErrorSwitch(err)
		}

		result = append(result, bson.E{
			Key:   "_id",
			Value: oid,
		})
	}

	if query.HasIsParentOrder() {
		result = append(result, bson.E{
			Key:   "is_parent_order",
			Value: query.HasIsParentOrder(),
		})
	}

	if query.HasIsDeferSubmission() {
		result = append(result, bson.E{
			Key:   "is_defer_submission",
			Value: query.HasIsDeferSubmission(),
		})
	}

	if query.HasPickupTimeFrom() {
		result = append(result, bson.E{
			Key: "estimated_pickup_time.value",
			Value: bson.M{
				"$gte": primitive.NewDateTimeFromTime(query.PickupTimeFrom),
			},
		})
	}

	if query.HasOrderCode() {
		result = append(result, bson.E{
			Key:   "order_code",
			Value: query.OrderCode,
		})
	}

	if query.HasPickupTimeTo() {
		result = append(result, bson.E{
			Key: "estimated_pickup_time.value",
			Value: bson.M{
				"$lte": primitive.NewDateTimeFromTime(query.PickupTimeTo),
			},
		})
	}

	if query.HasType() {
		result = append(result, bson.E{
			Key:   "type",
			Value: query.Type,
		})
	}

	if query.IsSentToPos() {
		result = append(result, bson.E{
			Key:   "sent_to_pos",
			Value: query.SentToPos,
		})
	}

	if query.HasOrderTimeFrom() {
		result = append(result, bson.E{
			Key: "order_time.value",
			Value: bson.M{
				"$gte": primitive.NewDateTimeFromTime(query.OrderTimeFrom),
			},
		})
	}

	if query.HasStatus() {
		result = append(result, bson.E{
			Key:   "status",
			Value: query.Status,
		})
	}

	if query.HasOnlyActive() {
		result = append(result, bson.E{
			Key: "status",
			Value: bson.D{
				{Key: "$nin", Value: []string{
					string(models.STATUS_CANCELLED),
					string(models.STATUS_DELIVERED),
					string(models.STATUS_CLOSED),
					string(models.STATUS_FAILED),
					string(models.STATUS_CANCELLED_BY_POS_SYSTEM),
					string(models.STATUS_CANCELLED_BY_DELIVERY_SERVICE),
					string(models.STATUS_SKIPPED)}},
			},
		})
	}

	if query.HasCustomerNumber() {
		result = append(result, bson.E{
			Key:   "customer.phone_number",
			Value: query.CustomerNumber,
		})
	}

	if query.HasOrderTimeTo() {
		result = append(result, bson.E{
			Key: "order_time.value",
			Value: bson.M{
				"$lte": primitive.NewDateTimeFromTime(query.OrderTimeTo),
			},
		})
	}

	if query.HasIgnoreStatus() {
		result = append(result, bson.E{
			Key: "status",
			Value: bson.M{
				"$ne": query.IgnoreStatus,
			},
		})
	}

	if query.HasPosType() {
		result = append(result, bson.E{
			Key:   "pos_type",
			Value: query.PosType,
		})
	}

	if query.HasStoreID() {
		result = append(result, bson.E{
			Key:   "restaurant_id",
			Value: query.StoreID,
		})
	}

	if query.HasOrderID() {
		result = append(result, bson.E{
			Key:   "order_id",
			Value: query.OrderID,
		})
	}

	if query.HasPosOrderID() {
		result = append(result, bson.E{
			Key:   "pos_order_id",
			Value: query.PosOrderID,
		})
	}

	if query.HasDeliveryService() && query.HasOrderID() {
		result = append(result, bson.E{
			Key:   "delivery_service",
			Value: query.DeliveryService,
		})
		result = append(result, bson.E{
			Key:   "order_id",
			Value: query.OrderID,
		})
	}

	if query.HasDeliveryService() {
		result = append(result, bson.E{
			Key:   "delivery_service",
			Value: query.DeliveryService,
		})
	}

	if query.HasRestaurants() {
		result = append(result, bson.E{
			Key:   "restaurant_id",
			Value: bson.M{"$in": query.Restaurants},
		})
	}

	if query.HasFailedReasonTimeoutCodes() {
		result = append(result, bson.E{
			Key:   "fail_reason.code",
			Value: bson.M{"$in": query.FailedReasonTimeoutCodes},
		})
	}

	if query.HasDeliveryOrderId() {
		result = append(result, bson.E{
			Key:   "delivery_order_id",
			Value: query.DeliveryOrderID,
		})
	}

	if query.HasDeliveryDispatcher() {
		result = append(result, bson.E{
			Key:   "delivery_dispatcher",
			Value: query.DeliveryDispatcher,
		})
	}

	if query.HasSearchForReport() {
		result = append(result, bson.E{
			Key: "$or",
			Value: []bson.M{
				{"order_id": query.SearchForReport},
				{"restaurant_name": bson.M{"$regex": query.SearchForReport, "$options": "i"}},
				{"allergy_info": bson.M{"$regex": query.SearchForReport, "$options": "i"}},
			},
		})
	}

	if query.HasDeliveryServices() {
		result = append(result, bson.E{
			Key: "delivery_service",
			Value: bson.D{
				{Key: "$in", Value: query.DeliveryServices},
			},
		})
	}

	if query.HasStatusFailedOrFailReasonNotEmpty() {
		result = append(result, bson.E{
			Key: "$or",
			Value: bson.A{
				bson.D{
					{Key: "$and",
						Value: bson.A{
							bson.D{{Key: "fail_reason", Value: bson.D{{Key: "$exists", Value: true}}}},
							bson.D{
								{Key: "$or",
									Value: bson.A{
										bson.D{{Key: "fail_reason.code", Value: bson.D{{Key: "$ne", Value: ""}}}},
										bson.D{{Key: "fail_reason.message", Value: bson.D{{Key: "$ne", Value: ""}}}},
									},
								},
							},
							bson.D{
								{Key: "status",
									Value: bson.D{
										{Key: "$nin",
											Value: bson.A{
												"FAILED",
												"SKIPPED",
											},
										},
									},
								},
							},
						},
					},
				},
				bson.D{{Key: "status", Value: "FAILED"}},
			},
		})
	}

	return result, nil
}

func (r *MongoRepository) sortFrom(query selector.Sorting) bson.D {
	sort := bson.D{
		{Key: query.Param, Value: query.Direction},
	}
	return sort
}

func (r *MongoRepository) GetOrdersByStatusesAndPosType(ctx context.Context, posType string, statuses []string) ([]models.Order, error) {
	filter := bson.M{
		"status":   bson.M{"$in": statuses},
		"pos_type": posType,
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	orders := make([]models.Order, 0, cursor.RemainingBatchLength())
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *MongoRepository) AddStatusToHistory(ctx context.Context, id, status string) error {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: objID},
	}

	update := bson.M{
		"$push": bson.M{
			"statuses_history": bson.M{
				"name": status,
				"time": time.Now(),
			},
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoRepository) FindOrderByOrderID(ctx context.Context, orderID string) (models.Order, error) {
	query := selector.EmptyOrderSearch().SetOrderID(orderID)

	filter, err := r.filterFrom(query)
	if err != nil {
		return models.Order{}, errors.ErrorSwitch(err)
	}

	var order models.Order
	err = r.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		return models.Order{}, errors.ErrorSwitch(err)
	}

	return order, nil
}

func (r *MongoRepository) FindOrderByDeliveryOrderID(ctx context.Context, deliveryID string) (models.Order, error) {

	query := selector.EmptyOrderSearch().SetDeliveryOrderID(deliveryID)

	filter, err := r.filterFrom(query)
	if err != nil {
		return models.Order{}, errors.ErrorSwitch(err)
	}

	var order models.Order
	err = r.collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		return models.Order{}, errors.ErrorSwitch(err)
	}

	return order, nil
}

func (r *MongoRepository) GetAverageBill(ctx context.Context, query selector.Order) (float64, error) {
	var result struct {
		AverageBill float64 `bson:"average_bill"`
	}

	pipeline := r.filterFromAverageBill(query)

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
	}

	return result.AverageBill, nil
}

func (r *MongoRepository) filterFromAverageBill(query selector.Order) []bson.M {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"customer.phone_number": query.CustomerNumber,
				"restaurant_id": bson.M{
					"$in": query.Restaurants,
				},
			},
		},

		{
			"$group": bson.M{
				"_id":                  "",
				"total_estimate_price": bson.M{"$sum": "$estimated_total_price.value"},
				"total_orders":         bson.M{"$sum": 1},
			},
		},

		{
			"$project": bson.M{
				"average_bill": bson.M{"$divide": bson.A{"$total_estimate_price", "$total_orders"}},
			},
		},
	}

	return pipeline
}

func (r *MongoRepository) Get3plOrdersForCron(ctx context.Context, callTime int64) ([]models.Order, error) {

	log.Info().Msgf("start to get orders for cron [Get3plOrdersForCron] BulkCreate3plOrder")

	callTimeInMiliSeconds := callTime * 60 * 1000

	pipeline := bson.A{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "$or", Value: bson.A{
						bson.D{{Key: "order_time.value", Value: bson.M{"$gte": primitive.NewDateTimeFromTime(time.Now().UTC().Add(-1 * time.Hour))}}},
						bson.D{{Key: "preorder.time.value", Value: bson.M{"$gte": primitive.NewDateTimeFromTime(time.Now().UTC().Add(-1 * time.Hour))}}},
					}},
				},
			},
		},
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "delivery_service", Value: bson.M{"$in": bson.A{"qr_menu", "kwaaka_admin"}}},
					{Key: "delivery_dispatcher", Value: bson.M{"$ne": ""}},
					{Key: "delivery_order_id", Value: ""},
					{Key: "restaurant_id",
						Value: bson.M{"$in": bson.A{
							"6683fc9339a3222785df695f", //Hani Иманбаева
							"6683fd3f0077b538b9497c24", //Hani Туркестан 28
							"6691122b5aafae72c5da35cb", //Hani Туркестан 20
							"669112e0076d30366ac63add", //Hani Тауелсиздик
							"669113bdcb3c0bfc666a8335", //Hani Мангилик ел
							"64d9fbb3cbc9cbed1106e4aa", //Своя кухня / Аксай
							"6544bef48a35272533e8cb63", //Своя кухня / Аль-Фараби (Нурлы-тау)
							"65d32cf9e7edea4cbdee782f", //Своя кухня / Шевченко
							"65d32df5e7edea4cbdee7830", //Своя кухня / Сагадат Нурмагамбетова
							"65f026cfc4b7b61f83f5d4d6", //Своя кухня / Гагарина
							"66b4a791f0012cd5fc9b4e8f", //Eki
							"6642084b19bb854b24399c35", //STAGE: Chiko /  Алматы Розыбакиева
						},
						},
					},
				},
			},
		},
		bson.D{
			{Key: "$addFields",
				Value: bson.D{
					{Key: "timeDifference",
						Value: bson.D{
							{Key: "$subtract",
								Value: bson.A{
									"$estimated_pickup_time.value",
									bson.D{{Key: "$toDate", Value: "$$NOW"}},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "timeDifference", Value: bson.D{{Key: "$lte", Value: callTimeInMiliSeconds}}},
					{Key: "timeDifference", Value: bson.D{{Key: "$gte", Value: 0}}},
				},
			},
		},
	}

	cur, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var res []models.Order

	if err = cur.All(ctx, &res); err != nil {
		return nil, err
	}

	log.Info().Msgf("list of orders for cron [Get3plOrdersForCron]  BulkCreate3plOrder: %v", res)

	return res, nil
}

func (r *MongoRepository) GetIIKO3plOrdersForCron(ctx context.Context, callTime int64) ([]models.Order, error) {

	log.Info().Msgf("start to get orders [GetIIKO3plOrdersForCron] for cron BulkCreate3plOrder")

	callTimeInMiliSeconds := callTime * 60 * 1000

	pipeline := bson.A{
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "$or", Value: bson.A{
						bson.D{{Key: "order_time.value", Value: bson.M{"$gte": primitive.NewDateTimeFromTime(time.Now().UTC().Add(-1 * time.Hour))}}},
						bson.D{{Key: "preorder.time.value", Value: bson.M{"$gte": primitive.NewDateTimeFromTime(time.Now().UTC().Add(-1 * time.Hour))}}},
					}},
				},
			},
		},
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "delivery_service", Value: bson.M{"$in": bson.A{"qr_menu", "kwaaka_admin"}}},
					{Key: "delivery_dispatcher", Value: bson.M{"$ne": ""}},
					{Key: "delivery_order_id", Value: ""},
					{Key: "pos_type", Value: "iiko"},
					{Key: "status", Value: "COOKING_STARTED"},
				},
			},
		},
		bson.D{
			{Key: "$addFields", Value: bson.D{
				{Key: "cooking_start_time", Value: bson.D{
					{Key: "$arrayElemAt", Value: bson.A{
						bson.D{
							{Key: "$filter", Value: bson.D{
								{Key: "input", Value: "$statuses_history"},
								{Key: "as", Value: "status"},
								{Key: "cond", Value: bson.D{
									{Key: "$eq", Value: bson.A{"$$status.name", "COOKING_STARTED"}},
								}},
							}},
						},
						0,
					}},
				}},
			}},
		},
		bson.D{
			{Key: "$addFields", Value: bson.D{
				{Key: "cooking_start_time", Value: "$cooking_start_time.time"},
			}},
		},
		bson.D{
			{Key: "$addFields", Value: bson.D{
				{Key: "cooking_end_time", Value: bson.D{
					{Key: "$dateAdd", Value: bson.D{
						{Key: "startDate", Value: "$cooking_start_time"},
						{Key: "unit", Value: "minute"},
						{Key: "amount", Value: "$cooking_time"},
					}},
				}},
			}},
		},
		bson.D{
			{Key: "$addFields",
				Value: bson.D{
					{Key: "timeDifference",
						Value: bson.D{
							{Key: "$subtract",
								Value: bson.A{
									"$cooking_end_time",
									bson.D{{Key: "$toDate", Value: "$$NOW"}},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{Key: "$match",
				Value: bson.D{
					{Key: "timeDifference", Value: bson.D{{Key: "$lte", Value: callTimeInMiliSeconds}}},
					{Key: "timeDifference", Value: bson.D{{Key: "$gte", Value: 0}}},
				},
			},
		},
	}

	cur, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var res []models.Order

	if err = cur.All(ctx, &res); err != nil {
		return nil, err
	}

	log.Info().Msgf("list of orders [GetIIKO3plOrdersForCron] for cron BulkCreate3plOrder: %v", res)

	return res, nil
}

func (r *MongoRepository) SetProposals(ctx context.Context, orderID string, proposals []models.Proposal) error {
	filter := bson.D{
		{Key: "order_id", Value: orderID},
	}

	if proposals == nil {
		return errorsGo.New("no proposals to set")
	}

	update := bson.D{
		{Key: "$set", Value: bson.M{"proposals": proposals}},
	}

	if _, err := r.collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (r *MongoRepository) SetDeliveryDispatcherPrice(ctx context.Context, orderID string, deliveryDispatcherPrice float64) error {
	oid, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "delivery_dispatcher_price", Value: deliveryDispatcherPrice}}},
	}

	if _, err := r.collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (r *MongoRepository) SetCancelledDeliveryDispatcherPrice(ctx context.Context, orderID string, deliveryDispatcherPrice float64) error {
	oid, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "history_3pl_delivery_info.delivery_dispatcher_price", Value: deliveryDispatcherPrice}}},
	}

	if _, err := r.collection.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}
func (r *MongoRepository) GetOrderBy3plDeliveryID(ctx context.Context, deliveryOrderId string) (models.Order, error) {

	filter := bson.D{
		{Key: "delivery_order_id", Value: deliveryOrderId},
	}

	res := r.collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		return models.Order{}, fmt.Errorf("repository - fn GetOrderBy3plDeliveryID - fn FindOne: %w", err)
	}

	var order models.Order

	if err := res.Decode(&order); err != nil {
		return models.Order{}, fmt.Errorf("repository - fn GetOrderBy3plDeliveryID - fn Decode: %w", err)
	}

	return order, nil
}

func (r *MongoRepository) Save3plDeliveryHistoryAndSetEmptyDispatcherService(ctx context.Context, id string, cancelled3plDeliveryInfo models.History3plDelivery, newDeliveryAddress models.DeliveryAddress, newCustomer models.Customer) error {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	updateSet := bson.D{
		{Key: "delivery_dispatcher", Value: ""},
		{Key: "delivery_order_id", Value: ""},
		{Key: "full_delivery_price", Value: 0},
		{Key: "restaurant_pay_delivery_price", Value: 0},
		{Key: "kwaaka_charged_delivery_price", Value: 0},
	}

	if newCustomer.PhoneNumber != "" && newDeliveryAddress.Street != "" {
		updateSet = append(updateSet, bson.E{Key: "customer", Value: newCustomer})
		updateSet = append(updateSet, bson.E{Key: "delivery_address", Value: newDeliveryAddress})
	}

	update := bson.D{
		{
			Key:   "$set",
			Value: updateSet,
		},
		{
			Key: "$push",
			Value: bson.D{
				{Key: "history_3pl_delivery_info", Value: cancelled3plDeliveryInfo},
			},
		},
	}

	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("repository - fn Save3plDeliveryHistoryAndSetEmptyDispatcherService - fn UpdateOne: %w", err)
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("order with id %s is not exist in database", id)
	}

	return nil
}
