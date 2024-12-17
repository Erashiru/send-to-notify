package selector

import (
	"github.com/kwaaka-team/orders-core/service/error_solutions/models"
	"time"
)

type Order struct {
	ID                          string
	PosOrderID                  string
	Type                        string
	OrderID                     string
	DeliveryService             string
	DeliveryOrderID             string
	ExternalStoreID             string
	PosType                     string
	StoreID                     string
	SentToPos                   *bool
	IgnoreStatus                string
	Status                      string
	OrderCode                   string
	OnlyActive                  bool
	Restaurants                 []string
	OrderTimeTo                 time.Time
	OrderTimeFrom               time.Time
	PickupTimeFrom              time.Time
	PickupTimeTo                time.Time
	PreorderPickUpTimeFrom      time.Time
	PreorderPickUpTimeTo        time.Time
	DeliveryArray               []string
	IsDeferSubmission           bool
	IsParentOrder               bool
	CustomerNumber              string
	FailReason                  string
	FailedReasonCode            string
	FailedReasonTimeoutCodes    []string
	IsActive                    bool
	IsPickedUpByCustomer        bool
	DeliveryDispatcher          string
	SearchForReport             string
	DeliveryServices            []string
	CreatedAtTimeFrom           time.Time
	EstimatedPickupTimeTo       time.Time
	CookingCompleteClosedStatus *bool

	Sorting
	Pagination
}

func EmptyOrderSearch() Order {
	return Order{}
}

func OrderSearch() Order {
	return Order{
		Pagination: Pagination{
			Limit: DefaultLimit,
		},
	}
}

func (o Order) SetDeliveryOrderId(id string) Order {
	o.DeliveryOrderID = id
	return o
}

func (o Order) HasDeliveryOrderId() bool {
	return o.DeliveryOrderID != ""
}

func (o Order) SetIsPickedUpByCustomer(value bool) Order {
	o.IsPickedUpByCustomer = value
	return o
}

func (o Order) HasIsPickedUpByCustomer() bool {
	return o.IsPickedUpByCustomer
}

func (o Order) SetExternalStoreID(externalStoreID string) Order {
	o.ExternalStoreID = externalStoreID
	return o
}

func (o Order) HasExternalStoreID() bool {
	return o.ExternalStoreID != ""
}

func (o Order) SetIsDeferSubmission(isDeferSubmisson bool) Order {
	o.IsDeferSubmission = isDeferSubmisson
	return o
}

func (o Order) HasIsDeferSubmission() bool {
	return o.IsDeferSubmission
}

func (o Order) SetIsParentOrder(isParentOrder bool) Order {
	o.IsParentOrder = isParentOrder
	return o
}

func (o Order) HasIsParentOrder() bool {
	return o.IsParentOrder
}

func (o Order) SetOrderCode(orderCode string) Order {
	o.OrderCode = orderCode
	return o
}

func (o Order) SetID(id string) Order {
	o.ID = id
	return o
}

func (o Order) HasID() bool {
	return o.ID != ""
}

func (o Order) SetPosType(posType string) Order {
	o.PosType = posType
	return o
}

func (o Order) HasPosType() bool {
	return o.PosType != ""
}

func (o Order) HasCustomerNumber() bool {
	return o.CustomerNumber != ""
}

func (o Order) SetCustomerNumber(customerNumber string) Order {
	o.CustomerNumber = customerNumber
	return o
}

func (o Order) SetDeliveryService(deliveryService string) Order {
	o.DeliveryService = deliveryService
	return o
}

func (o Order) HasDeliveryService() bool {
	return o.DeliveryService != ""
}

func (o Order) SetOrderID(orderID string) Order {
	o.OrderID = orderID
	return o
}

func (o Order) HasOrderID() bool {
	return o.OrderID != ""
}

func (o Order) SetStoreID(StoreID string) Order {
	o.StoreID = StoreID
	return o
}

func (o Order) HasStoreID() bool {
	return o.StoreID != ""
}

func (o Order) SetPosOrderID(id string) Order {
	o.PosOrderID = id
	return o
}

func (o Order) HasPosOrderID() bool {
	return o.PosOrderID != ""
}

func (o Order) SetIgnoreStatus(IgnoreStatus string) Order {
	o.IgnoreStatus = IgnoreStatus
	return o
}

func (o Order) HasIgnoreStatus() bool {
	return o.IgnoreStatus != "" && !o.HasStatus()
}

func (o Order) SetStatus(Status string) Order {
	o.Status = Status
	return o
}

func (o Order) HasStatus() bool {
	return o.Status != ""
}

func (o Order) SetOrderTimeFrom(orderTimeFrom time.Time) Order {
	o.OrderTimeFrom = orderTimeFrom
	return o
}

func (o Order) SetOrderTimeTo(orderTimeTo time.Time) Order {
	o.OrderTimeTo = orderTimeTo
	return o
}

func (o Order) HasOrderTimeFrom() bool {
	return o.OrderTimeFrom != time.Time{}
}

func (o Order) HasOrderTimeTo() bool {
	return o.OrderTimeTo != time.Time{}
}

func (o Order) HasOnlyActive() bool {
	return o.OnlyActive
}

func (o Order) HasRestaurants() bool {
	return o.Restaurants != nil
}

func (o Order) SetRestaurants(restaurants []string) Order {
	o.Restaurants = restaurants
	return o
}

func (o Order) IsSentToPos() bool {
	return o.SentToPos != nil
}

func (o Order) SetSentToPos(SentToPos *bool) Order {
	o.SentToPos = SentToPos
	return o
}

func (o Order) SetType(Type string) Order {
	o.Type = Type
	return o
}

func (o Order) HasType() bool {
	return o.Type != ""
}

func (o Order) SetPickupTimeFrom(pickupTimeFrom time.Time) Order {
	o.PickupTimeFrom = pickupTimeFrom
	return o
}

func (o Order) HasPickupTimeFrom() bool {
	return o.PickupTimeFrom != time.Time{}
}

func (o Order) SetPickupTimeTo(pickupTimeTo time.Time) Order {
	o.PickupTimeTo = pickupTimeTo
	return o
}

func (o Order) HasPickupTimeTo() bool {
	return o.PickupTimeTo != time.Time{}
}

func (o Order) SetPreorderPickupTimeTo(preorderPickupTimeTo time.Time) Order {
	o.PreorderPickUpTimeTo = preorderPickupTimeTo
	return o
}

func (o Order) HasPreorderPickupTimeTo() bool {
	return o.PreorderPickUpTimeTo != time.Time{}
}

func (o Order) SetPreorderPickupTimeFrom(preorderPickupTimeFrom time.Time) Order {
	o.PreorderPickUpTimeFrom = preorderPickupTimeFrom
	return o
}

func (o Order) HasPreorderPickupTimeFrom() bool {
	return o.PreorderPickUpTimeFrom != time.Time{}
}

func (o Order) HasOrderCode() bool {
	return o.OrderCode != ""
}

func (o Order) HasFailedReasonTimeoutCodes() bool {
	return o.FailedReasonTimeoutCodes != nil
}

func (o Order) SetFailedReasonTimeoutCodes(errorSolutions []models.ErrorSolution) Order {
	errSolutionsCode := make([]string, 0, len(errorSolutions))
	for _, errSol := range errorSolutions {
		errSolutionsCode = append(errSolutionsCode, errSol.Code)
	}
	o.FailedReasonTimeoutCodes = errSolutionsCode
	return o
}

func (o Order) SetDeliveryOrderID(deliveryID string) Order {
	o.DeliveryOrderID = deliveryID
	return o
}

func (o Order) SetPage(page int64) Order {
	if page > 0 {
		o.Pagination.Page = page - 1
	}
	return o
}

func (o Order) SetLimit(limit int64) Order {
	if limit > 0 {
		o.Pagination.Limit = limit
	}
	return o
}

func (o Order) SetSorting(param string, dir int8) Order {
	o.Sorting.Param = param
	o.Sorting.Direction = dir
	return o
}

func (o Order) HasDeliveryArray() bool {
	return o.DeliveryArray != nil && len(o.DeliveryArray) != 0
}

func (o Order) SetDeliveryArray(value []string) Order {
	o.DeliveryArray = value
	return o
}

func (o Order) HasDeliveryDispatcher() bool {
	return o.DeliveryDispatcher != ""
}
func (o Order) SetDeliveryDispatcher(value string) Order {
	o.DeliveryDispatcher = value
	return o
}

func (o Order) HasSearchForReport() bool {
	return o.SearchForReport != ""
}
func (o Order) SetSearchReport(value string) Order {
	o.SearchForReport = value
	return o
}

func (o Order) HasDeliveryServices() bool {
	return o.DeliveryServices != nil
}
func (o Order) SetDeliveryServices(deliveryServices []string) Order {
	o.DeliveryServices = deliveryServices
	return o
}

func (o Order) HasFailReason() bool { return o.FailReason != "" && o.FailReason != notEmpty }
func (o Order) SetFailReason(failReason string) Order {
	o.FailReason = failReason
	return o
}

const notEmpty = "not empty"

func (o Order) HasStatusFailedOrFailReasonNotEmpty() bool {
	return o.FailReason == notEmpty && o.FailedReasonCode == notEmpty
}
func (o Order) SetStatusFailedOrFailReasonNotEmpty() Order {
	o.FailedReasonCode = notEmpty
	o.FailReason = notEmpty
	return o
}

func (o Order) SetEstimatedPickupTimeTo(ept time.Time) Order {
	o.EstimatedPickupTimeTo = ept
	return o
}

func (o Order) HasEstimatedPickupTimeTo() bool {
	return o.EstimatedPickupTimeTo != time.Time{}
}

func (o Order) SetCreatedAtTimeFrom(createdAt time.Time) Order {
	o.CreatedAtTimeFrom = createdAt
	return o
}

func (o Order) HasCreatedAtTimeFrom() bool {
	return o.CreatedAtTimeFrom != time.Time{}
}

func (o Order) SetCookingCompleteClosedStatus(value *bool) Order {
	o.CookingCompleteClosedStatus = value
	return o
}

func (o Order) HasCookingCompleteClosedStatus() bool {
	return o.CookingCompleteClosedStatus != nil
}
