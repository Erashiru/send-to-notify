package models

import (
	"strconv"
	"time"
)

type DeliveryService string

const (
	STARTERAPP DeliveryService = "starter_app"
)

func (d DeliveryService) String() string {
	return string(d)
}

type OrderDto struct {
	StarterId         int             `json:"starterId"` // +
	GlobalId          string          `json:"globalId"`  // +
	OrderItems        []OrderItemDto  `json:"orderItems"`
	Bonuses           int             `json:"bonuses"`
	Price             float64         `json:"price"`
	DiscountPrice     float64         `json:"discountPrice"`
	DeliveryPrice     float64         `json:"deliveryPrice"` // +
	ChangeFrom        float64         `json:"changeFrom"`
	TotalPrice        float64         `json:"totalPrice"` // +
	Address           Address         `json:"address"`
	FlatwareAmount    int             `json:"flatwareAmount"`
	DeliveryType      string          `json:"deliveryType"` // +
	PaymentType       string          `json:"paymentType"`  // +
	PaymentStatus     string          `json:"paymentStatus"`
	SubmittedDatetime time.Time       `json:"submittedDatetime"` // +
	DeliveryDatetime  time.Time       `json:"deliveryDatetime"`  // +
	DeliveryDuration  int             `json:"deliveryDuration"`  // +
	UserId            int             `json:"userId"`
	Username          string          `json:"username"`  // +
	UserPhone         string          `json:"userPhone"` // +
	UserLang          string          `json:"userLang"`
	Comment           string          `json:"comment"` // +
	Status            string          `json:"status"`
	ShopId            int             `json:"shopId"` // +
	NotCall           bool            `json:"notCall"`
	IsPreorder        bool            `json:"isPreorder"` // +
	Source            string          `json:"source"`
	Discounts         []Discount      `json:"discounts"`
	DeliveryProduct   DeliveryProduct `json:"deliveryProduct"`
	Timezone          string          `json:"timezone"`
	TerminalId        string          `json:"terminalId"`
}

type OrderItemDto struct {
	OrderItemId   int           `json:"orderItemId"`
	MealId        int           `json:"mealId"`
	Quantity      int           `json:"quantity"`
	Price         float64       `json:"price"`
	TotalPrice    float64       `json:"totalPrice"`
	DiscountPrice float64       `json:"discountPrice"`
	Modifiers     []ModifierDto `json:"modifiers"`
}

type ModifierDto struct {
	ModifierId         int    `json:"modifierId"`
	Amount             int    `json:"amount"`
	Price              int    `json:"price"`
	Title              string `json:"title"`
	ModifiersGroupId   int    `json:"modifiersGroupId"`
	ModifiersGroupName string `json:"modifiersGroupName"`
}

type Address struct {
	Street    string  `json:"street"`   // +
	Flat      string  `json:"flat"`     // +
	Floor     string  `json:"floor"`    // +
	Entrance  string  `json:"entrance"` // +
	Comment   string  `json:"comment"`  // +
	City      string  `json:"city"`     // +
	Doorphone string  `json:"doorphone"`
	House     string  `json:"house"`
	Longitude float64 `json:"longitude"` // +
	Latitude  float64 `json:"latitude"`  // +
}

type Discount struct {
	DiscountId  string  `json:"discountId"`
	Title       string  `json:"title"`
	Type        string  `json:"type"`
	Sum         int     `json:"sum"`
	SumWithCent float64 `json:"sumWithCent"`
	Promocode   string  `json:"promocode"`
}

type DeliveryProduct struct {
	Id    int     `json:"id"`
	Price float64 `json:"price"`
}

type Order struct {
	StarterId         string          `json:"starterId"` // +
	GlobalId          string          `json:"globalId"`  // +
	OrderItems        []OrderItem     `json:"orderItems"`
	Bonuses           int             `json:"bonuses"`
	Price             float64         `json:"price"`
	DiscountPrice     float64         `json:"discountPrice"`
	DeliveryPrice     float64         `json:"deliveryPrice"` // +
	ChangeFrom        float64         `json:"changeFrom"`
	TotalPrice        float64         `json:"totalPrice"` // +
	Address           Address         `json:"address"`
	FlatwareAmount    int             `json:"flatwareAmount"`
	DeliveryType      string          `json:"deliveryType"` // +
	PaymentType       string          `json:"paymentType"`  // +
	PaymentStatus     string          `json:"paymentStatus"`
	SubmittedDatetime time.Time       `json:"submittedDatetime"` // +
	DeliveryDatetime  time.Time       `json:"deliveryDatetime"`  // +
	DeliveryDuration  int             `json:"deliveryDuration"`  // +
	UserId            string          `json:"userId"`
	Username          string          `json:"username"`  // +
	UserPhone         string          `json:"userPhone"` // +
	UserLang          string          `json:"userLang"`
	Comment           string          `json:"comment"` // +
	Status            string          `json:"status"`
	ShopId            string          `json:"shopId"` // +
	NotCall           bool            `json:"notCall"`
	IsPreorder        bool            `json:"isPreorder"` // +
	Source            string          `json:"source"`
	Discounts         []Discount      `json:"discounts"`
	DeliveryProduct   DeliveryProduct `json:"deliveryProduct"`
	Timezone          string          `json:"timezone"`
	TerminalId        string          `json:"terminalId"`
	CookingTime       int32           `json:"cookingTime,omitempty"`
}

type OrderItem struct {
	ExtID         string     `json:"extId,omitempty"`
	Name          string     `json:"name,omitempty"`
	OrderItemId   string     `json:"orderItemId"`
	MealId        string     `json:"mealId"`
	Quantity      int        `json:"quantity"`
	Price         float64    `json:"price"`
	TotalPrice    float64    `json:"totalPrice"`
	DiscountPrice float64    `json:"discountPrice"`
	Modifiers     []Modifier `json:"modifiers"`
}

type Modifier struct {
	ExtId              string `json:"extId,omitempty"`
	Name               string `json:"name,omitempty"`
	ModifierId         string `json:"modifierId"`
	Amount             int    `json:"amount"`
	Price              int    `json:"price"`
	Title              string `json:"title"`
	ModifiersGroupId   string `json:"modifiersGroupId"`
	ModifiersGroupName string `json:"modifiersGroupName"`
}

func (o *OrderDto) FromDto() Order {
	return Order{
		StarterId:         strconv.Itoa(o.StarterId),
		GlobalId:          o.GlobalId,
		OrderItems:        convertOrderItemsFromDto(o.OrderItems),
		Bonuses:           o.Bonuses,
		Price:             o.Price,
		DiscountPrice:     o.DiscountPrice,
		DeliveryPrice:     o.DeliveryPrice,
		ChangeFrom:        o.ChangeFrom,
		TotalPrice:        o.TotalPrice,
		Address:           o.Address,
		FlatwareAmount:    o.FlatwareAmount,
		DeliveryType:      o.DeliveryType,
		PaymentType:       o.PaymentType,
		PaymentStatus:     o.PaymentStatus,
		SubmittedDatetime: o.SubmittedDatetime,
		DeliveryDatetime:  o.DeliveryDatetime,
		DeliveryDuration:  o.DeliveryDuration,
		UserId:            strconv.Itoa(o.UserId),
		Username:          o.Username,
		UserPhone:         o.UserPhone,
		UserLang:          o.UserLang,
		Comment:           o.Comment,
		Status:            o.Status,
		ShopId:            strconv.Itoa(o.ShopId),
		NotCall:           o.NotCall,
		IsPreorder:        o.IsPreorder,
		Source:            o.Source,
		Discounts:         o.Discounts,
		DeliveryProduct:   o.DeliveryProduct,
		Timezone:          o.Timezone,
		TerminalId:        o.TerminalId,
	}
}

func convertOrderItemsFromDto(itemsDto []OrderItemDto) []OrderItem {
	var items []OrderItem
	for _, itemDto := range itemsDto {
		items = append(items, OrderItem{
			OrderItemId:   strconv.Itoa(itemDto.OrderItemId),
			MealId:        strconv.Itoa(itemDto.MealId),
			Quantity:      itemDto.Quantity,
			Price:         itemDto.Price,
			TotalPrice:    itemDto.TotalPrice,
			DiscountPrice: itemDto.DiscountPrice,
			Modifiers:     convertModifiersFromDto(itemDto.Modifiers),
		})
	}
	return items
}

func convertModifiersFromDto(modifiersDto []ModifierDto) []Modifier {
	var modifiers []Modifier
	for _, modDto := range modifiersDto {
		modifiers = append(modifiers, Modifier{
			ModifierId:         strconv.Itoa(modDto.ModifierId),
			Amount:             modDto.Amount,
			Price:              modDto.Price,
			Title:              modDto.Title,
			ModifiersGroupId:   strconv.Itoa(modDto.ModifiersGroupId),
			ModifiersGroupName: modDto.ModifiersGroupName,
		})
	}
	return modifiers
}
