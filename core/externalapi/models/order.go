package models

import (
	coreOrderModels "github.com/kwaaka-team/orders-core/core/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type CancelOrderRequest struct {
	EatsId  string `json:"eatsId"`
	Comment string `json:"comment"`
}

type DeliveryInfo struct {
	ClientName  string `json:"clientName,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	//CourierArrivementDate string       `json:"courierArrivementDate"`
	MarketPlaceDeliveryDate string                      `json:"deliveryDate,omitempty"`
	DeliveryDate            string                      `json:"courierArrivementDate,omitempty"`
	DeliveryAddress         *MarketPlaceDeliveryAddress `json:"deliveryAddress,omitempty"`
	ClientArrivementDate    string                      `json:"clientArrivementDate,omitempty"`

	//DeliveryAddress DeliveryAddress `json:"deliveryAddress,omitempty"`
}

type MarketPlaceDeliveryAddress struct {
	Full      string `json:"full"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type DeliveryAddress struct {
	Full      string `json:"full"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type PaymentInfo struct {
	PaymentType    string `json:"paymentType"`
	ItemsCost      int    `json:"itemsCost"`
	NettingPayment *bool  `json:"netting_payment,omitempty"`
	//DeliveryFee int    `json:"deliveryFee"`
	//Total       int    `json:"total,omitempty"`
	//Change      int    `json:"change"`
}

type OrderModification struct {
	Id       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

type OrderItem struct {
	Id                   string              `json:"id"`
	Name                 string              `json:"name,omitempty"`
	Quantity             float64             `json:"quantity"`
	Price                int                 `json:"price"`
	PriceWithoutDiscount *int                `json:"price_without_discount,omitempty"`
	Modifications        []OrderModification `json:"modifications,omitempty"`
	Promos               []OrderPromo        `json:"promos,omitempty"`
}

type OrderPromo struct {
	Type     string `json:"type"`
	Discount int    `json:"discount"`
}

// Order model info
// @Description Discriminator Дискриминатор схемы обьекта. Для MarketplaceOrder равен
// @Description EatsId Сквозной идентификатор заказа на стороне Яндекс.Еды в формате DDDDDD-DDDDDD
type Order struct {
	Platform      string       `json:"platform,omitempty" description:"For Tillypad - Идентификатор платформы. YE - Yandex Eda, DC - Delivery club"`
	Discriminator string       `json:"discriminator" example:"marketplace" description:"Дискриминатор схемы обьекта. Для MarketplaceOrder равен"`
	EatsId        string       `json:"eatsId" example:"190330-123456" description:"Сквозной идентификатор заказа на стороне Яндекс.Еды в формате DDDDDD-DDDDDD"`
	RestaurantId  string       `json:"restaurantId,omitempty"`
	DeliveryInfo  DeliveryInfo `json:"deliveryInfo"`
	PaymentInfo   PaymentInfo  `json:"paymentInfo"`
	Items         []OrderItem  `json:"items"`
	Persons       int          `json:"persons,omitempty"`
	Comment       string       `json:"comment"`
	Promos        []OrderPromo `json:"promos,omitempty"`
}

type OrderStatusResponse struct {
	Status    string `json:"status"`
	Comment   string `json:"comment,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
}

// CreationResult - model
type CreationResult struct {
	Result  string `json:"result"`
	OrderId string `json:"orderId"`
}

func (o Order) ToModel(store coreStoreModels.Store, service string) (coreOrderModels.Order, error) {
	var (
		deliveryDate     time.Time
		err              error
		pickUpByCustomer bool
		isMarketplace    bool
	)

	switch o.Discriminator {
	case "pickup":
		pickUpByCustomer = true
		deliveryDate, err = time.Parse(time.RFC3339Nano, o.DeliveryInfo.ClientArrivementDate)
		if err != nil {
			return coreOrderModels.Order{}, errors.New("pickup invalid date type")
		}
	case "yandex":
		isMarketplace = true
		deliveryDate, err = time.Parse(time.RFC3339Nano, o.DeliveryInfo.DeliveryDate)
		if err != nil {
			return coreOrderModels.Order{}, errors.New("yandex invalid date type")
		}
	case "marketplace":
		deliveryDate, err = time.Parse(time.RFC3339Nano, o.DeliveryInfo.MarketPlaceDeliveryDate)
		if err != nil {
			return coreOrderModels.Order{}, errors.New("marketplace invalid date type")
		}
	}

	loc, _ := time.LoadLocation(store.Settings.TimeZone.TZ)

	deliveryDate = time.Date(deliveryDate.Year(), deliveryDate.Month(), deliveryDate.Day(), deliveryDate.Hour(), deliveryDate.Minute(), deliveryDate.Second(), deliveryDate.Nanosecond(), loc)

	deliveryAddress := coreOrderModels.DeliveryAddress{}
	if o.DeliveryInfo.DeliveryAddress != nil {
		deliveryAddress.Label = o.DeliveryInfo.DeliveryAddress.Full
	}

	order := coreOrderModels.Order{
		Type:            "INSTANT",
		Discriminator:   o.Discriminator,
		DeliveryService: service,
		PosType:         store.PosType,
		RestaurantID:    store.ID,
		OrderID:         o.EatsId,
		OrderCode:       o.EatsId,
		StoreID:         o.RestaurantId,
		Status:          "NEW",
		StatusesHistory: []coreOrderModels.OrderStatusUpdate{
			{
				Name: "NEW",
				Time: time.Now().UTC(),
			},
		},
		OrderTime: coreOrderModels.TransactionTime{
			Value:     coreOrderModels.Time{Time: time.Now().UTC()},
			TimeZone:  store.Settings.TimeZone.TZ,
			UTCOffset: strconv.Itoa(int(store.Settings.TimeZone.UTCOffset)),
		},
		EstimatedPickupTime: coreOrderModels.TransactionTime{
			Value:     coreOrderModels.Time{Time: deliveryDate.UTC()},
			TimeZone:  store.Settings.TimeZone.TZ,
			UTCOffset: strconv.Itoa(int(store.Settings.TimeZone.UTCOffset)),
		},
		Currency:            store.Settings.Currency,
		SpecialRequirements: o.Comment,
		AllergyInfo:         o.Comment,
		EstimatedTotalPrice: coreOrderModels.Price{
			Value:        float64(o.PaymentInfo.ItemsCost),
			CurrencyCode: store.Settings.Currency,
		},
		TotalCustomerToPay: coreOrderModels.Price{
			Value:        float64(o.PaymentInfo.ItemsCost),
			CurrencyCode: store.Settings.Currency,
		},
		Customer: coreOrderModels.Customer{
			Name:        o.DeliveryInfo.ClientName,
			PhoneNumber: o.DeliveryInfo.PhoneNumber,
		},
		Persons:              o.Persons,
		IsPickedUpByCustomer: pickUpByCustomer,
		CutleryRequested:     true,
		DeliveryAddress:      deliveryAddress,
		IsMarketplace:        isMarketplace,
		IsActive:             true,
	}

	for _, promo := range o.Promos {
		order.Promos = append(order.Promos, coreOrderModels.Promo{
			Discount: promo.Discount,
			Type:     promo.Type,
		})
	}

	switch o.PaymentInfo.PaymentType {
	case "CASH":
		order.PaymentMethod = "CASH"
	default:
		order.PaymentMethod = "DELAYED"
	}

	order.Products = make([]coreOrderModels.OrderProduct, 0, len(o.Items))
	for _, item := range o.Items {
		product := coreOrderModels.OrderProduct{
			ID:                 item.Id,
			PurchasedProductID: item.Id,
			Name:               item.Name,
			Quantity:           int(item.Quantity),
			Price: coreOrderModels.Price{
				Value:        float64(item.Price),
				CurrencyCode: store.Settings.Currency,
			},
		}

		for _, promo := range item.Promos {
			product.Promos = append(product.Promos, coreOrderModels.Promo{
				Discount: promo.Discount,
				Type:     promo.Type,
			})
		}

		if item.Modifications != nil {
			for _, modifier := range item.Modifications {
				product.Attributes = append(product.Attributes, coreOrderModels.ProductAttribute{
					ID:       modifier.Id,
					Name:     modifier.Name,
					Quantity: modifier.Quantity,
					Price: coreOrderModels.Price{
						Value:        float64(modifier.Price),
						CurrencyCode: store.Settings.Currency,
					},
				})
			}
		}

		order.Products = append(order.Products, product)
	}

	return order, nil
}
