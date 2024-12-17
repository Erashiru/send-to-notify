package dto

import (
	"github.com/goccy/go-json"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

type AuthRequest struct {
	ApplicationId string `json:"application_id"`
	Secret        string `json:"secret"`
}

type AuthResponse struct {
	Token  string `json:"token"`
	Role   string `json:"role"`
	Expiry Time   `json:"expiry"`
}

type Time struct {
	time.Time
}

func (ct *Time) UnmarshalJSON(b []byte) (err error) {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	parsedTime, err := time.Parse(timeFormat, s)
	if err != nil {
		return err
	}

	ct.Time = parsedTime
	return nil
}

func (ct Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.Format(timeFormat))
}

type CreatePaymentInvoiceResponse struct {
	Success bool `json:"success"`
	Data    struct {
		StoreId     int                       `json:"store_id"`
		Amount      int                       `json:"amount"`
		InvoiceId   string                    `json:"invoice_id"`
		ReturnUrl   string                    `json:"return_url"`
		Ofds        []Ofd                     `json:"ofd"`
		Products    any                       `json:"products"`
		Split       any                       `json:"split"`
		Uuid        string                    `json:"uuid"`
		ShortLink   string                    `json:"short_link"`
		CallBackUrl string                    `json:"callback_url"`
		AddedOn     CreatePaymentInvoiceTime  `json:"added_on"`
		UpdatedOn   *CreatePaymentInvoiceTime `json:"updated_on"`
		Payment     any                       `json:"payment"`
		CheckoutUrl string                    `json:"checkout_url"`
	} `json:"data"`
}

type CreatePaymentInvoiceRequest struct {
	StoreId     int    `json:"store_id"`
	Amount      int    `json:"amount"`
	InvoiceId   string `json:"invoice_id"`
	ReturnUrl   string `json:"return_url"`
	CallbackUrl string `json:"callback_url"`
	Ofds        []Ofd  `json:"ofd"`
}

type Ofd struct {
	Vat         string `json:"vat"`
	Price       int    `json:"price"`
	Name        string `json:"name"`
	PackageCode string `json:"package_code"`
	Mxik        string `json:"mxik"`
	Total       string `json:"total"`
}

type Payment struct {
	Uuid             string `json:"uuid"`
	Status           string `json:"status"`
	Ps               string `json:"ps"`
	StoreInvoiceId   string `json:"store_invoice_id"`
	PaymentAmount    int    `json:"payment_amount"`
	CommissionAmount int    `json:"commission_amount"`
	TotalAmount      int    `json:"total_amount"`
}

type CreatePaymentInvoiceTime struct {
	time.Time
}

func (t *CreatePaymentInvoiceTime) UnmarshalJSON(b []byte) error {
	str := string(b)
	str = str[1 : len(str)-1]

	parsedTime, err := time.Parse(timeFormat, str)
	if err != nil {
		return err
	}
	t.Time = parsedTime
	return nil
}

type ErrResponse struct {
	Success bool `json:"success"`
	Error   struct {
		Code    string `json:"code"`
		Details string `json:"details"`
	} `json:"error"`
}
