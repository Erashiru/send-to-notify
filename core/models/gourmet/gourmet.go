package gourmet

import (
	iikoModels "github.com/kwaaka-team/orders-core/pkg/iiko/models"
	"time"
)

type GourmetGetTablesResponse struct {
	Tables []GourmetTabels `json:"tables"`
}

type GourmetTabels struct {
	Id              string `json:"id"`
	Number          int    `json:"number"`
	Name            string `json:"name"`
	SeatingCapacity int    `json:"seatingCapacity"`
	SectionId       string `json:"sectionId"`
	SectionName     string `json:"sectionName"`
}
type GourmetGetOrdersResponse struct {
	Orders []GourmetOrder `json:"orders"`
}
type GourmetOrder struct {
	Id              string                 `json:"id"`
	TableID         string                 `json:"tableId"`
	Customer        iikoModels.Customer    `json:"customer"`
	Waiter          iikoModels.Waiter      `json:"waiter"`
	Phone           string                 `json:"phone"`
	Status          string                 `json:"status"`
	Sum             float64                `json:"sum"`
	TotalToPay      float64                `json:"totalToPay"`
	WhenCreated     string                 `json:"whenCreated"`
	WhenBillPrinted string                 `json:"whenBillPrinted"`
	WhenClosed      string                 `json:"whenClosed"`
	Items           []iikoModels.ItemsResp `json:"items"`
	Discounts       []iikoModels.Discounts `json:"discounts"`
	Payments        []iikoModels.Payments  `json:"payments"`
}

type Size struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type ComboInformation struct {
	GroupName string `json:"groupName"`
}
type Product struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type ProductGroup struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type Modifiers struct {
	Product      Product      `json:"product"`
	Amount       int          `json:"amount"`
	ProductGroup ProductGroup `json:"productGroup"`
	Price        int          `json:"price"`
	ResultSum    int          `json:"resultSum"`
}
type Items struct {
	Type             string           `json:"type"`
	Status           string           `json:"status"`
	Deleted          bool             `json:"deleted"`
	Amount           int              `json:"amount"`
	Comment          string           `json:"comment"`
	WhenPrinted      time.Time        `json:"whenPrinted"`
	Size             Size             `json:"size"`
	ComboInformation ComboInformation `json:"comboInformation"`
	Modifiers        []Modifiers      `json:"modifiers"`
}
type DiscountType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type SelectivePositionsWithSum struct {
	PositionID string `json:"positionId"`
	Sum        int    `json:"sum"`
}
type Discounts struct {
	DiscountType              DiscountType                `json:"discountType"`
	Sum                       int                         `json:"sum"`
	SelectivePositions        []string                    `json:"selectivePositions"`
	SelectivePositionsWithSum []SelectivePositionsWithSum `json:"selectivePositionsWithSum"`
}
type PaymentType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}
type Payments struct {
	PaymentType            PaymentType `json:"paymentType"`
	Sum                    int         `json:"sum"`
	IsPreliminary          bool        `json:"isPreliminary"`
	IsExternal             bool        `json:"isExternal"`
	IsProcessedExternally  bool        `json:"isProcessedExternally"`
	IsFiscalizedExternally bool        `json:"isFiscalizedExternally"`
	IsPrepay               bool        `json:"isPrepay"`
}

type PaymentChangeResponse struct {
	OrderId string `json:"orderId"`
	TableId string `json:"tableId"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type PaymentChangeRequest struct {
	PaymentTypeId   string `json:"paymentTypeId"`
	PaymentTypeKind string `json:"paymentTypeKind"`
	IsPaid          bool   `json:"isPaid"`
}
