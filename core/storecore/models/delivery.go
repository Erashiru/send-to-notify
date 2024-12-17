package models

type DeliveryServicePaymentType struct {
	CASH    PaymentType `bson:"CASH" json:"cash"`
	DELAYED PaymentType `bson:"DELAYED" json:"delayed"`
}

type StoreDelivery struct {
	ID       string `bson:"id" json:"id"`
	Code     string `bson:"code" json:"code"`
	Price    int    `bson:"price" json:"price"`
	Name     string `bson:"name" json:"name"`
	IsActive bool   `bson:"is_active" json:"is_active"`
	Service  string `bson:"service" json:"service"`
}

type StoreAddress struct {
	ExtId          string      `bson:"ext_id" json:"ext_id"`
	City           string      `bson:"city" json:"city"`
	Street         string      `bson:"street" json:"street"`
	Entrance       string      `bson:"entrance,omitempty" json:"entrance,omitempty"`
	DeliveryRadius int         `bson:"delivery_radius" json:"delivery_radius"`
	Coordinates    Coordinates `bson:"coordinates" json:"coordinates"`
	Geometry       Geometry    `bson:"geometry,omitempty" json:"geometry,omitempty"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

type Geometry struct {
	Type               string        `bson:"type,omitempty" json:"type,omitempty"`
	PercentageModifier int           `bson:"percentage_modifier,omitempty" json:"percentage_modifier,omitempty"` //orders total sum part that restaurant pays for deliveries
	Coordinates        []Coordinates `bson:"coordinates,omitempty" json:"coordinates,omitempty"`
}

type Polygon struct {
	ID                                      string        `bson:"id" json:"id"`
	Priority                                int           `json:"priority" bson:"priority"`
	CPO                                     float64       `bson:"cpo" json:"cpo"`
	MinBasketForOrder                       float64       `bson:"min_basket_for_order" json:"min_basket_for_order"`
	MinBasketForFreeDelivery                float64       `bson:"min_basket_for_free_delivery" json:"min_basket_for_free_delivery"`
	MinBasketBasedDeliveryPrice             bool          `bson:"min_basket_based_delivery_price" json:"min_basket_based_delivery_price"`
	FixedDeliveryPriceForClient             float64       `bson:"fixed_delivery_price_for_client" json:"fixed_delivery_price_for_client"`
	PolygonBasedFixedDeliveryPriceForClient bool          `bson:"polygon_based_fixed_delivery_price_for_client" json:"polygon_based_fixed_delivery_price_for_client"`
	Coordinates                             []Coordinates `bson:"coordinates" json:"coordinates"`
}
