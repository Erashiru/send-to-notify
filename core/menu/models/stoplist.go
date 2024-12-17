package models

import (
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	"strconv"
	"time"
)

type StopListProduct struct {
	ExtID       string  `bson:"ext_id" json:"ext_id"`
	PosID       string  `bson:"pos_id" json:"pos_id"`
	ProductID   string  `bson:"product_id" json:"product_id"`
	MsID        string  `bson:"ms_id" json:"ms_id,omitempty"`
	Name        string  `bson:"product_name" json:"product_name"`
	Price       float64 `bson:"price" json:"price"`
	IsAvailable bool    `bson:"is_available" json:"is_available"`
	Stock       string  `bson:"stock" json:"stock"`
}

type StopListProducts []StopListProduct

func ToStopListProducts(products Products) StopListProducts {
	res := make(StopListProducts, 0, len(products))

	for _, product := range products {
		stoplistProduct := StopListProduct{
			ExtID:       product.ExtID,
			IsAvailable: product.IsAvailable,
			Stock:       strconv.Itoa(int(product.Balance)),
		}

		if len(product.Name) != 0 {
			stoplistProduct.Name = product.Name[0].Value
		}

		if len(product.Price) != 0 {
			stoplistProduct.Price = product.Price[0].Value
		}

		res = append(res, stoplistProduct)
	}

	return res
}

func (s StopListProducts) GetNames() []string {

	res := make([]string, 0, len(s))
	for _, product := range s {
		res = append(res, product.Name)
	}

	return res
}

type StopListAttribute struct {
	AttributeID   string  `bson:"attribute_id" json:"attribute_id"`
	AttributeName string  `bson:"attribute_name" json:"attribute_name"`
	IsAvailable   bool    `bson:"is_available" json:"is_available"`
	Stock         string  `bson:"stock" json:"stock"`
	Price         float64 `bson:"price" json:"price"`
}

type StopListAttributes []StopListAttribute

func ToStopListAttributes(attributes Attributes) StopListAttributes {
	res := make(StopListAttributes, 0, len(attributes))

	for _, attribute := range attributes {
		res = append(res, StopListAttribute{
			AttributeID:   attribute.ExtID,
			AttributeName: attribute.Name,
			IsAvailable:   attribute.IsAvailable,
			Price:         attribute.Price,
			// Stock: product. //fixme: what is Stock?
		})
	}

	return res
}

func (s StopListAttributes) GetNames() []string {

	res := make([]string, 0, len(s))
	for _, attr := range s {
		res = append(res, attr.AttributeName)
	}

	return res
}

type StopListResponse struct {
	Products   StopListProducts   `json:"products"`
	Attributes StopListAttributes `json:"attributes"`
}

type StopListTransaction struct {
	ID               string            `bson:"_id,omitempty" json:"_id,omitempty"`
	StoreID          string            `bson:"restaurant_id" json:"restaurant_id"`
	PosStopListItems StopListItems     `bson:"pos_stoplist_items" json:"pos_stoplist_items"`
	Transactions     []TransactionData `bson:"transactions" json:"transactions"`

	Products   StopListProducts   `bson:"products,omitempty" json:"products,omitempty"`
	Attributes StopListAttributes `bson:"attributes,omitempty" json:"attributes,omitempty"`

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

func (st *StopListTransaction) Append(id, delivery, storeId, message string, status Status) {
	trn := TransactionData{
		ID:       id,
		Delivery: delivery,
		StoreID:  storeId,
		Status:   status,
		Message:  message,
	}
	st.Transactions = append(st.Transactions, trn)
}

type Transaction struct {
	ID                 string `bson:"id" json:"id"`
	StoreID            string `bson:"store_id" json:"store_id"`
	Delivery           string `bson:"delivery" json:"delivery"`
	MoySkladPositionID string `bson:"moy_sklad_position_id" json:"moy_sklad_position_id"`
}

func (st *StopListTransaction) Fill(restID string, products StopListProducts, attributes StopListAttributes) {
	st.StoreID = restID
	st.Products = products
	st.Attributes = attributes
	st.CreatedAt = time.Now()
}

type TransactionData struct {
	ID      string `bson:"id" json:"id"`
	StoreID string `bson:"store_id" json:"store_id"`

	Delivery           string `bson:"delivery" json:"delivery"`
	MoySkladPositionID string `bson:"moy_sklad_position_id,omitempty" json:"moy_sklad_position_id,omitempty"`

	Products   StopListProducts   `bson:"products" json:"products"`
	Attributes StopListAttributes `bson:"attributes" json:"attributes"`

	Status  Status `bson:"status,omitempty" json:"status,omitempty"`
	Message string `bson:"message,omitempty" json:"message,omitempty"`
}

type StopListTerminalResponse struct {
	TerminalGroupID string
	Items           StopListItems
}

type StopListItems []StopListItem

type StopListItem struct {
	Balance   float64
	ProductID string
}

func (s StopListItems) Products() []string {

	products := make([]string, 0, len(s))

	for _, req := range s {
		products = append(products, req.ProductID)
	}

	return products
}

type StopList struct {
	Type          int    `json:"type,omitempty"`
	ElementID     int    `json:"element_id,omitempty"`
	StorageID     int    `json:"storage_id,omitempty"`
	ValueRelative int    `json:"value_relative,omitempty"`
	ValueAbsolute int    `json:"value_absolute,omitempty"`
	ProductID     string `json:"product_id,omitempty"`
}

type StopListScheduler struct {
	ID                 string               `json:"_id"`
	Name               string               `json:"name"`
	RstGroupID         string               `json:"rst_group_id"`
	RstID              string               `json:"rst_id"`
	MenuID             string               `json:"menu_id"`
	Products           []ProductScheduler   `json:"products"`
	IsActive           bool                 `json:"is_active"`
	Available          bool                 `json:"available"`
	StartDate          time.Time            `json:"start_date"`
	EndDate            time.Time            `json:"end_date"`
	WeeklyAvailability []WeeklyAvailability `json:"weekly_availability"`
}
type WeeklyAvailability struct {
	WeeklyDayName string        `json:"weekly_day_name" bson:"weekly_day_name"`
	WeeklyDayNum  int           `json:"weekly_day_num" bson:"weekly_day_num"`
	StartTime     TimeScheduler `json:"start_time" bson:"start_time"`
	EndTime       TimeScheduler `json:"end_time" bson:"end_time"`
}

func (s *StopListScheduler) DefineAvailability() bool {
	now := coreModels.TimeNow().UTC()
	dayOfWeek := int(now.Weekday())

	// availability logic
	for _, day := range s.WeeklyAvailability {
		if day.WeeklyDayNum != dayOfWeek {
			continue
		}
		if day.EndTime.Hour == now.Hour() && day.EndTime.Minute <= now.Minute() {
			return false
		}
		if day.StartTime.Hour == now.Hour() && day.StartTime.Minute <= now.Minute() {
			return true
		}

	}
	return false //?
}

type TimeScheduler struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

type ProductScheduler struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type StopListEvent struct {
	GroupID     string   `json:"group_id"`
	GlovoIDs    []string `json:"glovo_ids"`
	WoltIDs     []string `json:"wolt_ids"`
	IsAvailable bool     `json:"is_available"`
}

type StopListByIDsEvent struct {
	IsAvailable         *bool               `json:"is_available"`
	StoreToProductIDs   map[string][]string `json:"store_to_product_ids"`
	StoreToAttributeIDs map[string][]string `json:"store_to_attribute_ids"`
}
