package models

import (
	"encoding/json"
	"github.com/pkg/errors"
)

var ErrNoEvent = errors.New("no such event")

//go:generate stringer -type=Event
type Event int

const (
	DeliveryOrderUpdate Event = iota + 1
	DeliveryOrderError
	ReserveUpdate
	ReserveError
	TableOrderUpdate
	TableOrderError
	StopListUpdate
)

//go:generate stringer -type=Status
type Status int

const (
	Success = iota + 1
	InProgress
	Error
)

type WebhookEvents []WebhookEvent

type WebhookEvent struct {
	EventType      Event      `json:"eventType"`
	EventTime      DateTime   `json:"eventTime"`
	OrganizationID string     `json:"organizationId"`
	CorrelationID  string     `json:"correlationId"`
	EventInfo      *EventInfo `json:"eventInfo"`
}

type EventInfo struct {
	ID              string     `json:"id"`
	OrganizationID  string     `json:"organizationId"`
	PosID           string     `json:"posId"`
	Timestamp       int64      `json:"timestamp"`
	CreationsStatus string     `json:"creationStatus"`
	Error           ErrorInfo  `json:"errorInfo,omitempty"`
	Order           OrderEvent `json:"order,omitempty"`
	Problem         *Problem   `json:"problem,omitempty"`
	IsDeleted       bool       `json:"isDeleted"`
	UpdateType      string     `json:"stopListType"`
}

type Problem struct {
	HasProblem  bool   `json:"hasProblem"`
	Description string `json:"description,omitempty"`
}

type GuestsInfo struct {
	Count               int  `json:"count"`
	SplitBetweenPersons bool `json:"splitBetweenPersons"`
}

func (e Event) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

func (e *Event) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	for i := DeliveryOrderUpdate; i <= StopListUpdate; i++ {
		if i.String() == str {
			*e = i
			return nil
		}
	}
	return ErrNoEvent
}

type OrderEvent struct {
	Customer            Customer              `json:"customer"`
	TerminalID          string                `json:"terminalGroupId"`
	Phone               string                `json:"phone"`
	Status              string                `json:"status"`
	CompleteBefore      DateTime              `json:"completeBefore"`
	WhenCreated         DateTime              `json:"whenCreated"`
	CookingStartTime    DateTime              `json:"cookingStartTime"`
	IsDeleted           bool                  `json:"isDeleted"`
	Sum                 json.Number           `json:"sum"`
	Number              int                   `json:"number"`
	GuestInfo           GuestsInfo            `json:"guestsInfo"`
	Items               []Item                `json:"items"`
	OfflineOrderPayment []OfflineOrderPayment `json:"payments"`
}

type ExtendedOrderEvent struct {
	Order               OrderEvent              `json:"order"`
	RegularCustomerInfo GetCustomerInfoResponse `json:"regularCustomerInfo"`
	StoreId             string                  `json:"storeId"`
	StoreName           string                  `json:"storeName"`
	EventId             string                  `json:"eventId"`
}
type OfflineOrderPayment struct {
	PaymentTypes          OfflineOrderPaymentType `json:"paymentType"`
	Sum                   float64                 `json:"sum"`
	IsProcessedExternally bool                    `json:"isProcessedExternally"`
}
type OfflineOrderPaymentType struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type ItemEvent struct {
	Product         ProductEvent `json:"product"`
	Price           json.Number  `json:"price"`
	Cost            json.Number  `json:"cost"`
	PricePredefined bool         `json:"pricePredefined"`
	Type            string       `json:"type"`
	Status          string       `json:"string"`
	Amount          json.Number  `json:"amount"`
}

type ProductEvent struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type WebhookEventResponse struct {
	Details []string `json:"details"`
}
