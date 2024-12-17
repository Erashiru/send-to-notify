package dto

type CreatePaymentInvoiceRequest struct {
	Amount      int64       `json:"amount"`
	Currency    string      `json:"currency,omitempty"`
	Merchant    string      `json:"merchant,omitempty"`
	Service     string      `json:"service"`
	ExternalID  string      `json:"external_id,omitempty"`
	Regulations Regulations `json:"regulations"`
}

type ClientIdentification struct {
	Phone string `json:"phone"`
	Email string `json:"email,omitempty"`
}

type ResultURL struct {
	URL  string `json:"url,omitempty"`
	Type string `json:"type,omitempty"`
}

type Behavior struct {
	ResultURL    ResultURL `json:"result_url,omitempty"`
	AutoRedirect bool      `json:"auto_redirect,omitempty"`
	SuccessURL   string    `json:"success_url,omitempty"`
	FailureURL   string    `json:"failure_url,omitempty"`
}

type Description struct {
	Main            string `json:"main"`
	InternalComment string `json:"internal_comment"`
}

type Regulations struct {
	Language             string               `json:"language,omitempty"`
	UserPhoneRequired    bool                 `json:"user_phone_required,omitempty"`
	ExpirationDate       string               `json:"expiration_date,omitempty"`
	PaymentMethods       []string             `json:"payment_methods,omitempty"`
	ClientIdentification ClientIdentification `json:"client_identification"`
	Behavior             Behavior             `json:"behavior,omitempty"`
	Description          Description          `json:"description,omitempty"`
}

type CreatePaymentInvoiceResponse struct {
	Data Data `json:"data"`
}

type Attributes struct {
	OperationID int `json:"operation_id"`
}

type Data struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
	URL        string     `json:"url"`
}

type CreatePaymentInvoiceErrorResponse struct {
	Error Error `json:"error"`
}

type Detail struct {
	RegulationsExpirationDate []string `json:"regulations.expiration_date"`
	RegulationsPaymentMethods []string `json:"regulations.payment_methods"`
}

type Error struct {
	Code   int    `json:"code"`
	Title  string `json:"title"`
	Detail Detail `json:"detail"`
}
