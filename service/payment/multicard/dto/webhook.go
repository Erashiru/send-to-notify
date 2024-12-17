package dto

import (
	"encoding/json"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"

type Webhook struct {
	StoreId     string      `json:"store_id"`
	Amount      int         `json:"amount"`
	InvoiceId   string      `json:"invoice_id"`
	InvoiceUuid string      `json:"invoice_uuid"`
	BillingId   interface{} `json:"billing_id"`
	PaymentTime Time        `json:"payment_time"`
	Phone       string      `json:"phone"`
	CardPan     string      `json:"card_pan"`
	Uuid        string      `json:"uuid"`
	ReceiptUrl  string      `json:"receipt_url"`
	Sign        string      `json:"sign"`
}

type Time struct {
	time.Time
}

func (ct *Time) UnmarshalJSON(data []byte) error {
	strTime := string(data)
	strTime = strTime[1 : len(strTime)-1]

	parsedTime, err := time.Parse(timeFormat, strTime)
	if err != nil {
		return err
	}

	ct.Time = parsedTime
	return nil
}

func (ct Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.Format(timeFormat))
}
