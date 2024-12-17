package dto

import "time"

type RequestCreateOrder struct {
	ApiKey       string                 `json:"api_key"`
	Sig          string                 `json:"sig"`
	RestaurantID string                 `json:"restaurant_id"`
	Order        RequestCreateOrderBody `json:"order"`
}

type RequestCreateOrderBody struct {
	RestaurantId     string                     `json:"restaurant_id"`
	ToRestaurantId   string                     `json:"to_restaurant_id"`
	Address          string                     `json:"address"`
	Phone            string                     `json:"phone"`
	Contact          string                     `json:"contact"`
	Description      string                     `json:"description"`
	PeopleCount      int                        `json:"people_count"`
	OrderType        int                        `json:"order_type"`
	AmountOrder      string                     `json:"amount_order"`
	PaymentMethod    int                        `json:"payment_method"`
	PaymentType      int                        `json:"payment_type"`
	DeliveryTime     string                     `json:"delivery_time"`
	DeliveryTimeType int                        `json:"delivery_time_type"`
	DeliveryPrice    int                        `json:"delivery_price"`
	Discount         int                        `json:"discount"`
	DiscountSum      int                        `json:"discount_sum"`
	Courses          []RequestCreateOrderCourse `json:"courses"`
}

type RequestCreateOrderCourse struct {
	CourseId    string `json:"course_id"`
	Count       int    `json:"count"`
	Price       int    `json:"price"`
	Description string `json:"description"`
}

type ResponseOrder struct {
	Status int               `json:"status"`
	Order  ResponseOrderBody `json:"order"`
	ErrorResponse
}

type ResponseOrderBody struct {
	Id                      string                `json:"id"`
	RestaurantId            string                `json:"restaurant_id"`
	ToRestaurantId          string                `json:"to_restaurant_id"`
	Status                  int                   `json:"status"`
	Number                  int                   `json:"number"`
	AmountOrder             string                `json:"amount_order"`
	DateTime                time.Time             `json:"date_time"`
	WorkDate                string                `json:"work_date"`
	Address                 string                `json:"address"`
	Phone                   string                `json:"phone"`
	Description             string                `json:"description"`
	DeliveryPrice           string                `json:"delivery_price"`
	IsDeliveryInCash        bool                  `json:"is_delivery_in_cash"`
	CourierName             string                `json:"courier_name"`
	CourierPhone            string                `json:"courier_phone"`
	CancellationReason      string                `json:"cancellation_reason"`
	CancellationType        int                   `json:"cancellation_type"`
	IsCancellationConfirmed bool                  `json:"is_cancellation_confirmed"`
	PeopleCount             int                   `json:"people_count"`
	OrderType               int                   `json:"order_type"`
	Discount                int                   `json:"discount"`
	DiscountAmount          string                `json:"discount_amount"`
	DiscountSum             string                `json:"discount_sum"`
	IsPayed                 bool                  `json:"is_payed"`
	PaymentMethod           int                   `json:"payment_method"`
	PaymentStatus           int                   `json:"payment_status"`
	PaymentType             int                   `json:"payment_type"`
	ClientId                string                `json:"client_id"`
	ClientCardId            string                `json:"client_card_id"`
	AccumulationAccount     string                `json:"accumulation_account"`
	BillId                  string                `json:"bill_id"`
	Contact                 string                `json:"contact"`
	History                 string                `json:"history"`
	Courses                 []ResponseOrderCourse `json:"courses"`
}

type ResponseOrderCourse struct {
	Id           string `json:"id"`
	CourseId     string `json:"course_id"`
	CourseTitle  string `json:"course_title"`
	Count        string `json:"count"`
	CoursePrice  string `json:"course_price"`
	CourseAmount string `json:"course_amount"`
	Description  string `json:"description"`
	IsException  bool   `json:"is_exception"`
}

type RequestCancelOrder struct {
	ApiKey             string `json:"api_key"`
	Sig                string `json:"sig"`
	CancellationReason string `json:"cancellation_reason"`
}
