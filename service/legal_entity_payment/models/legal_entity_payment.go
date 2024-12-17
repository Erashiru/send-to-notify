package models

import (
	"strings"
	"time"
)

//go:generate stringer -type=LegalEntityPaymentStatus
type LegalEntityPaymentStatus int

const (
	UNBILLED LegalEntityPaymentStatus = iota + 1
	BILLED
	PAID
	PAID_CONFIRMED
	FAILED
)

type LegalEntityPayment struct {
	ID               string    `bson:"_id,omitempty" json:"id"`
	Contract         string    `bson:"contract" json:"contract"`
	Name             string    `bson:"name" json:"name"`
	LegalEntityID    string    `bson:"legal_entity_id" json:"legal_entity_id"`
	LegalEntityName  string    `bson:"legal_entity_name" json:"legal_entity_name"`
	Amount           float64   `bson:"amount" json:"amount"`
	PaidAmount       float64   `bson:"paid_amount" json:"paid_amount"`
	StartDate        time.Time `bson:"start_date,omitempty" json:"start_date"`
	EndDate          time.Time `bson:"end_date,omitempty" json:"end_date"`
	PaymentType      string    `bson:"payment_type" json:"payment_type"`
	Status           string    `bson:"status" json:"status"`
	Bill             string    `bson:"bill,omitempty" json:"bill"`
	BillPayment      string    `bson:"bill_payment,omitempty" json:"bill_payment"`
	BillingAt        time.Time `bson:"billing_at,omitempty" json:"billing_at"`
	BillPaymentAt    time.Time `bson:"bill_payment_at,omitempty" json:"bill_payment_at"`
	ConfirmPaymentAt time.Time `bson:"confirm_payment_at,omitempty" json:"confirm_payment_at"`
	Comments         string    `bson:"comments,omitempty" json:"comments"`
	CreatedAt        time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time `bson:"updated_at,omitempty" json:"updated_at"`
}

type ListLegalEntityPaymentQuery struct {
	LegalEntityIDs []string  `json:"legal_entity_ids"`
	PaymentTypes   []string  `json:"payment_types"`
	Statuses       []string  `json:"statuses"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`
	Pagination
}

type ListLegalEntityPayment struct {
	LegalEntityId       string               `json:"legal_entity_id"`
	Name                string               `json:"name"`
	PaymentType         string               `json:"payment_type"`
	Status              string               `json:"status"`
	BilledAmount        float64              `json:"billed_amount"`
	PaidAmount          float64              `json:"paid_amount"`
	Brands              []string             `json:"brands"`
	LegalEntityPayments []LegalEntityPayment `json:"legal_entity_payments"`
}

type LegalEntityPaymentAnalyticsRequest struct {
	LegalEntityIDs []string  `json:"legal_entity_ids"`
	StartDate      time.Time `json:"start_date" bson:"start_date"`
	EndDate        time.Time `json:"end_date" bson:"end_date"`
	Statuses       []string  `json:"statuses"`
}

type LegalEntityPaymentAnalyticsResponse struct {
	Paid   LegalEntityPaymentAnalytics `json:"paid"`
	Unpaid LegalEntityPaymentAnalytics `json:"unpaid"`
	Total  LegalEntityPaymentAnalytics `json:"total"`
}

type LegalEntityPaymentAnalytics struct {
	Amount          float64 `json:"amount"`
	Quantity        int     `json:"quantity"`
	AmountPercent   float64 `json:"amount_percent"`
	QuantityPercent float64 `json:"quantity_percent"`
}

type UpdateLegalEntityPayment struct {
	ID               *string    `bson:"id" json:"id"`
	Name             *string    `bson:"name" json:"name"`
	LegalEntityID    *string    `bson:"legal_entity_id" json:"legal_entity_id"`
	LegalEntityName  *string    `bson:"legal_entity_name" json:"legal_entity_name"`
	Amount           *float64   `bson:"amount" json:"amount"`
	StartDate        *time.Time `bson:"start_date,omitempty" json:"start_date"`
	EndDate          *time.Time `bson:"end_date,omitempty" json:"end_date"`
	PaymentType      *string    `bson:"payment_type" json:"payment_type"`
	Status           *string    `bson:"status" json:"status"`
	Bill             *string    `bson:"bill,omitempty" json:"bill"`
	BillPayment      *string    `bson:"bill_payment,omitempty" json:"bill_payment"`
	BillingAt        *time.Time `bson:"billing_at,omitempty" json:"billing_at"`
	BillPaymentAt    *time.Time `bson:"bill_payment_at,omitempty" json:"bill_payment_at"`
	ConfirmPaymentAt *time.Time `bson:"confirm_payment_at,omitempty" json:"confirm_payment_at"`
}

type LegalEntityPaymentCreateBillRequest struct {
	LegalEntityPaymentID string    `json:"legal_entity_payment_id"`
	Name                 string    `json:"name"`
	Amount               float64   `json:"amount"`
	StartDate            time.Time `json:"start_date"`
	EndDate              time.Time `json:"end_date"`
	BillLink             string    `json:"bill_link"`
}

type LegalEntityPaymentConfirmPaymentRequest struct {
	LegalEntityPaymentID string `json:"legal_entity_payment_id"`
	BillLink             string `json:"bill_link"`
}

type LegalEntityPaymentDownloadPDFRequest struct {
	LegalEntityPaymentID string `json:"legal_entity_payment_id"`
	File                 []byte `json:"file"`
}

type LegalEntityPaymentFilePDFResponse struct {
	PDFUrl string `json:"pdf_url"`
}

type Pagination struct {
	Page   int64 `json:"page"`
	Limit  int64 `json:"limit"`
	Offset int64
}

func (p Pagination) AddOffset() Pagination {
	p.Offset = (p.Page - 1) * p.Limit
	return p
}

func (s LegalEntityPayment) ToUpdateModel() UpdateLegalEntityPayment {
	var update UpdateLegalEntityPayment
	update.ID = &s.ID

	if s.LegalEntityID != "" {
		update.LegalEntityID = &s.LegalEntityID
	}

	if s.Name != "" {
		update.Name = &s.Name
	}

	if s.LegalEntityName != "" {
		update.LegalEntityName = &s.LegalEntityName
	}

	if s.Amount != 0 {
		update.Amount = &s.Amount
	}

	if !s.StartDate.IsZero() {
		update.StartDate = &s.StartDate
	}

	if !s.EndDate.IsZero() {
		update.EndDate = &s.EndDate
	}

	if s.PaymentType != "" {
		update.PaymentType = &s.PaymentType
	}

	if s.Status != "" {
		update.Status = &s.Status
	}

	if s.Bill != "" {
		update.Bill = &s.Bill
	}

	if s.BillPayment != "" {
		update.BillPayment = &s.BillPayment
	}

	if !s.BillingAt.IsZero() {
		update.BillingAt = &s.BillingAt
	}

	if !s.BillPaymentAt.IsZero() {
		update.BillPaymentAt = &s.BillPaymentAt
	}

	if !s.ConfirmPaymentAt.IsZero() {
		update.ConfirmPaymentAt = &s.ConfirmPaymentAt
	}

	return update
}

func StringToLegalEntityPaymentStatus(s string) int {
	switch strings.ToUpper(s) {
	case "UNBILLED":
		return int(UNBILLED)
	case "BILLED":
		return int(BILLED)
	case "PAID":
		return int(PAID)
	case "PAID_CONFIRMED":
		return int(PAID_CONFIRMED)
	case "FAILED":
		return int(FAILED)
	default:
		return 0
	}
}

type PaymentXlsx struct {
	Number      string `json:"number"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Bank        string `json:"bank"`
	BIK         string `json:"BIK"`
	Code        string `json:"code"`
	Month       string `json:"month"`
	BillingDate string `json:"billing_date"`
	Buyer       string `json:"buyer"`
	Contract    string `json:"contract"`
	Status      string `json:"status"`
}

type IntegrationXlsx struct {
	Number   string `json:"number"`
	Name     string `json:"name"`
	Amount   string `json:"amount"`
	Price    string `json:"price"`
	SumPrice string `json:"sum_price"`
}
