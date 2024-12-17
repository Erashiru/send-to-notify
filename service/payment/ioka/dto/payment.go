package dto

import "time"

type CreatePaymentOrderRequest struct {
	Amount        int       `json:"amount"`
	Currency      string    `json:"currency,omitempty"`
	CaptureMethod string    `json:"capture_method,omitempty"`
	ExternalID    string    `json:"external_id,omitempty"`
	Description   string    `json:"description,omitempty"`
	Mcc           string    `json:"mcc,omitempty"`
	ExtraInfo     ExtraInfo `json:"extra_info,omitempty"`
	Attempts      int       `json:"attempts,omitempty"`
	DueDate       string    `json:"due_date,omitempty"`
	CustomerID    string    `json:"customer_id,omitempty"`
	CardID        string    `json:"card_id,omitempty"`
	BackURL       string    `json:"back_url,omitempty"`
	SuccessURL    string    `json:"success_url,omitempty"`
	FailureURL    string    `json:"failure_url,omitempty"`
	Template      string    `json:"template,omitempty"`
}

type ExtraInfo struct {
	RestaurantName      string `json:"Ресторан:"`
	RestaurantGroupName string `json:"Группа Ресторанов:"`
	CustomerName        string `json:"Имя Клиента:"`
	CustomerPhoneNumber string `json:"Номер Клиента:"`
}

type CreatePaymentOrderResponse struct {
	Order            Order  `json:"order"`
	OrderAccessToken string `json:"order_access_token"`
}

type Order struct {
	ID            string    `json:"id"`
	ShopID        string    `json:"shop_id"`
	Status        string    `json:"status"`
	CreatedAt     string    `json:"created_at"`
	Amount        int       `json:"amount"`
	Currency      string    `json:"currency"`
	CaptureMethod string    `json:"capture_method"`
	ExternalID    string    `json:"external_id"`
	Description   string    `json:"description"`
	ExtraInfo     ExtraInfo `json:"extra_info"`
	Attempts      int       `json:"attempts"`
	DueDate       string    `json:"due_date"`
	CustomerID    string    `json:"customer_id"`
	CardID        string    `json:"card_id"`
	BackURL       string    `json:"back_url"`
	SuccessURL    string    `json:"success_url"`
	FailureURL    string    `json:"failure_url"`
	Template      string    `json:"template"`
	CheckoutURL   string    `json:"checkout_url"`
	AccessToken   string    `json:"access_token"`
	Mcc           string    `json:"mcc"`
}

type GetSubscriptionPaymentsResponse struct {
	ID             string        `json:"id"`
	OrderID        string        `json:"order_id"`
	Status         string        `json:"status"`
	CreatedAt      time.Time     `json:"created_at"`
	ApprovedAmount int           `json:"approved_amount"`
	CapturedAmount int           `json:"captured_amount"`
	RefundedAmount int           `json:"refunded_amount"`
	ProcessingFee  int           `json:"processing_fee"`
	Payer          Payer         `json:"payer"`
	Error          ErrorResponse `json:"error"`
	Acquirer       Acquirer      `json:"acquirer"`
	Action         Action        `json:"action"`
}

type Acquirer struct {
	Name      string `json:"name"`
	Reference string `json:"reference"`
}
type Action struct {
	URL string `json:"url"`
}

type RefundRequest struct {
	Amount    int        `json:"amount"`
	Reason    string     `json:"reason"`
	Rules     []Rule     `json:"rules"`
	Positions []Position `json:"positions"`
}

type Rule struct {
	AccountID string `json:"account_id"`
	Amount    int    `json:"amount"`
}

type Position struct {
	Name       string `json:"name"`
	Amount     int    `json:"amount"`
	Count      int    `json:"count"`
	Section    int    `json:"section"`
	TaxPercent int    `json:"tax_percent"`
	TaxType    int    `json:"tax_type"`
	TaxAmount  int    `json:"tax_amount"`
	UnitCode   int    `json:"unit_code"`
}

type RefundResponse struct {
	ID        string         `json:"id"`
	PaymentID string         `json:"payment_id"`
	OrderID   string         `json:"order_id"`
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
