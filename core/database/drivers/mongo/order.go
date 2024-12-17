package mongo

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/errors"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/core/models/selector"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository struct {
	collection *mongo.Collection
}

func NewOrderRepository(orderCollection *mongo.Collection) *OrderRepository {
	return &OrderRepository{
		collection: orderCollection,
	}
}

func (orderRepo *OrderRepository) GetByID(ctx context.Context, id string) (models.Order, error) {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Order{}, err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	var order models.Order
	err = orderRepo.collection.FindOne(ctx, filter).Decode(&order)

	if err != nil {
		return models.Order{}, errors.ErrorSwitch(err)
	}
	return order, nil
}

func (orderRepo *OrderRepository) InsertOrder(ctx context.Context, order models.Order) (models.Order, error) {

	order.CreatedAt = models.TimeNow()
	order.UpdatedAt = models.TimeNow()
	log.Info().Msgf("repo orderID: %s", order.OrderID)
	res, err := orderRepo.collection.InsertOne(ctx, order)
	if err != nil {
		return order, errors.ErrorSwitch(err)
	}

	order.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return order, nil
}

func (orderRepo *OrderRepository) UpdateOrder(ctx context.Context, order models.Order) error {
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

	res, err := orderRepo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.ErrorSwitch(err)
	}

	if res.MatchedCount == 0 {
		return err
	}
	return nil
}

// add builder
func (orderRepo *OrderRepository) GetOrder(ctx context.Context, query selector.Order) (models.Order, error) {
	filter, err := orderRepo.filterFrom(query)
	if err != nil {
		return models.Order{}, err
	}

	var order models.Order

	if err = orderRepo.collection.FindOne(ctx, filter).Decode(&order); err != nil {
		return models.Order{}, errors.ErrorSwitch(err)
	}

	return order, nil
}

func (orderRepo *OrderRepository) GetOrders(ctx context.Context, query selector.Order) ([]models.Order, int, error) {
	filter, err := orderRepo.filterFrom(query)
	if err != nil {
		return nil, 0, err
	}

	totalCount, err := orderRepo.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	//sortOptions := options.Find().SetSort(bson.D{{Key: "order_time.value", Value: -1}})
	var sortOptions *options.FindOptions
	if query.HasSorting() {
		sortOptions = options.Find().SetSort(bson.D{{Key: query.Sorting.Param, Value: query.Sorting.Direction}})
	}

	option := options.Find()
	if query.HasPagination() {
		option = options.Find().SetLimit(query.Limit).SetSkip(query.Skip())
	}

	cursor, err := orderRepo.collection.Find(ctx, filter, sortOptions, option)
	if err != nil {
		return nil, 0, err
	}

	orders := make([]models.Order, 0, cursor.RemainingBatchLength())
	if err := cursor.All(ctx, &orders); err != nil {
		return nil, 0, err
	}

	return orders, int(totalCount), nil
}

func (orderRepo *OrderRepository) GetAllOrders(ctx context.Context, query selector.Order) ([]models.Order, error) {
	filter, err := orderRepo.filterFrom(query)
	if err != nil {
		return nil, errors.ErrorSwitch(err)
	}

	cur, err := orderRepo.collection.Find(ctx, filter)
	if err != nil {
		return nil, errors.ErrorSwitch(err)
	}

	orders := make([]models.Order, 0, cur.RemainingBatchLength())
	if err := cur.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (orderRepo *OrderRepository) CancelOrder(ctx context.Context, req models.CancelOrder) (models.Order, error) {
	oid, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return models.Order{}, err
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
		{Key: "order_id", Value: req.OrderID},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{{Key: "status", Value: models.STATUS_CANCELLED}}},
		{Key: "$set", Value: bson.D{{Key: "cancel_reason", Value: req.CancelReason}}},
	}

	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	result := orderRepo.collection.FindOneAndUpdate(ctx, filter, update, &opt)
	if err := result.Err(); err != nil {
		return models.Order{}, errors.ErrorSwitch(err)
	}

	var order models.Order

	if err := result.Decode(&order); err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (orderRepo *OrderRepository) GetOrderStatus(ctx context.Context, id string) (models.Order, error) {

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Order{}, errors.ErrorSwitch(err)
	}

	filter := bson.D{
		{Key: "_id", Value: oid},
	}

	var order models.Order

	if err = orderRepo.collection.FindOne(ctx, filter).Decode(&order); err != nil {
		return models.Order{}, errors.ErrorSwitch(err)
	}

	return order, nil
}

func (orderRepo *OrderRepository) SetPaidStatus(ctx context.Context, orderID string) error {
	filter, err := orderRepo.filterFrom(selector.EmptyOrderSearch().SetID(orderID))
	if err != nil {
		return errors.ErrorSwitch(err)
	}
	update := bson.M{"$set": bson.M{"is_paid": true}}

	result, err := orderRepo.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return errors.ErrorSwitch(err)
	}

	if result.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (orderRepo *OrderRepository) UpdateOrderStatus(ctx context.Context, query selector.Order, status, errorDescription string) error {

	filter, err := orderRepo.filterFrom(query)
	if err != nil {
		return errors.ErrorSwitch(err)
	}

	updateFields := bson.D{{Key: "status", Value: status}}

	if errorDescription != "" {
		updateFields = append(updateFields, bson.E{
			Key:   "creation_result.message",
			Value: errorDescription,
		})
	}

	if status == models.FAILED.String() ||
		status == models.CANCELLED_BY_POS_SYSTEM.String() ||
		status == models.CLOSED.String() ||
		status == string(models.STATUS_CANCELLED_BY_DELIVERY_SERVICE) ||
		status == string(models.STATUS_CANCELLED) ||
		status == string(models.STATUS_SKIPPED) {

		updateFields = append(updateFields, bson.E{
			Key:   "is_active",
			Value: false,
		})
	}

	statusHistory := models.OrderStatusUpdate{
		Name: status,
		Time: time.Now(),
	}

	update := bson.D{{Key: "$set", Value: updateFields}, {Key: "$push", Value: bson.M{"statuses_history": statusHistory}}}

	res, err := orderRepo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (orderRepo *OrderRepository) UpdateOrderStatusByID(ctx context.Context, id, posName, status string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.D{
		{Key: "_id", Value: objID},
		{Key: "pos_type", Value: posName},
	}

	updateFields := bson.D{{Key: "status", Value: status}}

	update := bson.D{{Key: "$set", Value: updateFields}}

	res, err := orderRepo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (orderRepo *OrderRepository) filterFrom(query selector.Order) (bson.D, error) {
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

	if query.HasIsPickedUpByCustomer() {
		result = append(result, bson.E{
			Key:   "is_picked_up_by_customer",
			Value: query.HasIsPickedUpByCustomer(),
		})
	}

	if query.HasIsParentOrder() {
		result = append(result, bson.E{
			Key:   "is_parent_order",
			Value: query.HasIsParentOrder(),
		})
	}

	if query.HasDeliveryArray() {
		result = append(result, bson.E{
			Key: "delivery_service",
			Value: bson.M{
				"$in": query.DeliveryArray,
			},
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
			Value: bson.M{"$regex": query.OrderCode, "$options": "i"},
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

	if query.HasExternalStoreID() {
		result = append(result, bson.E{
			Key:   "store_id",
			Value: query.ExternalStoreID,
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

	if query.HasPreorderPickupTimeFrom() {
		result = append(result, bson.E{
			Key: "preorder.time.value",
			Value: bson.M{
				"$gte": primitive.NewDateTimeFromTime(query.PreorderPickUpTimeFrom),
			},
		})
	}

	if query.HasPreorderPickupTimeTo() {
		result = append(result, bson.E{
			Key: "preorder.time.value",
			Value: bson.M{
				"$lte": primitive.NewDateTimeFromTime(query.PreorderPickUpTimeTo),
			},
		})
	}
	if query.HasCreatedAtTimeFrom() {
		result = append(result, bson.E{
			Key: "created_at",
			Value: bson.M{
				"$gte": primitive.NewDateTimeFromTime(query.CreatedAtTimeFrom),
			},
		})
	}
	if query.HasEstimatedPickupTimeTo() {
		result = append(result, bson.E{
			Key: "estimated_pickup_time.value",
			Value: bson.M{
				"$lte": primitive.NewDateTimeFromTime(query.EstimatedPickupTimeTo),
			},
		})
	}
	if query.HasCookingCompleteClosedStatus() {
		result = append(result, bson.E{
			Key: "statuses_history",
			Value: bson.D{
				{Key: "$not", Value: bson.D{
					{Key: "$elemMatch", Value: bson.D{
						{Key: "name", Value: bson.D{
							{Key: "$in", Value: []string{
								string(models.STATUS_COOKING_COMPLETE),
								string(models.STATUS_CLOSED),
							}},
						}},
					}},
				}},
			},
		})
	}

	return result, nil
}
