package models

import "time"

type PaymentOrder struct {
	ExternalID                string          `bson:"_id,omitempty" json:"external_id,omitempty"`
	PaymentOrderID            string          `bson:"payment_order_id" json:"payment_order_id,omitempty"`
	ShopID                    string          `bson:"shop_id" json:"shop_id,omitempty"`
	PaymentOrderStatus        string          `bson:"status" json:"payment_order_status,omitempty"`
	PaymentOrderStatusHistory []StatusHistory `bson:"payment_order_status_history" json:"payment_order_status_history,omitempty"`
	PaymentStatusHistory      []StatusHistory `bson:"payment_status_history" json:"payment_status_history,omitempty"`
	CreatedAtPaymentSystem    time.Time       `bson:"created_at_payment_system" json:"created_at_payment_system"`
	Amount                    int             `bson:"amount" json:"amount,omitempty"`
	RefundAmount              int64           `bson:"refund_amount,omitempty" json:"refund_amount,omitempty"`
	RefundReason              string          `bson:"refund_reason,omitempty" json:"refund_reason,omitempty"`
	RefundAuthor              string          `bson:"refund_author,omitempty" json:"refund_author,omitempty"`
	Currency                  string          `bson:"currency" json:"currency,omitempty"`
	CaptureMethod             string          `bson:"capture_method" json:"capture_method,omitempty"`
	Description               string          `bson:"description" json:"description,omitempty"`
	ExtraInfo                 string          `bson:"extra_info" json:"extra_info,omitempty"`
	Attempts                  int             `bson:"attempts" json:"attempts,omitempty"`
	DueDate                   string          `bson:"due_date" json:"due_date,omitempty"`
	CustomerID                string          `bson:"customer_id" json:"customer_id,omitempty"`
	CardID                    string          `bson:"card_id" json:"card_id,omitempty"`
	BackURL                   string          `bson:"back_url" json:"back_url,omitempty"`
	SuccessURL                string          `bson:"success_url" json:"success_url,omitempty"`
	FailureURL                string          `bson:"failure_url" json:"failure_url,omitempty"`
	Template                  string          `bson:"template" json:"template,omitempty"`
	CheckoutURL               string          `bson:"checkout_url" json:"checkout_url,omitempty"`
	AccessToken               string          `bson:"access_token" json:"access_token,omitempty"`
	Mcc                       string          `bson:"mcc" json:"mcc,omitempty"`
	PaymentSystem             string          `bson:"payment_system" json:"payment_system,omitempty"`
	CreatedAt                 time.Time       `bson:"created_at" json:"created_at"`
	UpdatedAt                 time.Time       `bson:"updated_at" json:"updated_at"`
	CartID                    string          `bson:"cart_id" json:"cart_id,omitempty"`
	OrderSource               string          `bson:"order_source" json:"order_source,omitempty"`
	PaymentID                 string          `bson:"payment_id" json:"payment_id,omitempty"`
	Payer                     Payer           `bson:"payer" json:"payer"`
	Acquirer                  Acquirer        `bson:"acquirer" json:"acquirer"`
	Action                    Action          `bson:"action" json:"action"`
	Error                     Error           `bson:"error" json:"error"`
	OrderID                   string          `bson:"order_id" json:"order_id,omitempty"`
	CustomerPhoneNumber       string          `bson:"customer_phone_number" json:"customer_phone_number"`
	RestaurantID              string          `bson:"restaurant_id" json:"restaurant_id"`
	PaymentInvoiceID          string          `bson:"payment_order_operation_id" json:"payment_order_operation_id"`
	PaymentProcessingDate     time.Time       `json:"payment_processing_date" bson:"payment_processing_date"`
	RestaurantName            string          `json:"restaurant_name" bson:"restaurant_name"`
	RestaurantGroupName       string          `bson:"restaurant_group_name" json:"restaurant_group_name"`
	PaymentTypeID             string          `json:"payment_type_id" bson:"payment_type_id"`
	NotificationCount         int             `bson:"notification_count" json:"notification_count,omitempty"`
	WhatsappPaymentChatId     string          `bson:"whatsapp_payment_chat_id" json:"whatsapp_payment_chat_id,omitempty"`
	CustomerName              string          `bson:"customer_name,omitempty" json:"customer_name,omitempty"`
	MulticardRefundUuid       string          `bson:"multicard_refund_uuid,omitempty" json:"multicard_refund_uuid,omitempty"`
}

type StatusHistory struct {
	Status string    `bson:"status"`
	Time   time.Time `bson:"time"`
}

type UpdatePaymentOrderStatus struct {
	CartID string `json:"cart_id" bson:"cart_id"`
	Status string `json:"status" bson:"status"`
}

type RefundResponse struct {
	ID        string         `json:"id,omitempty"`
	PaymentID string         `json:"payment_id"`
	OrderID   string         `json:"order_id,omitempty"`
	Status    string         `json:"status"`
	CreatedAt string         `json:"created_at"`
	Error     ErrorRefund    `json:"error,omitempty"`
	Acquirer  AcquirerRefund `json:"acquirer,omitempty"`
}

type ErrorRefund struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type AcquirerRefund struct {
	Name      string `json:"name"`
	Reference string `json:"reference"`
}
