package models

type CreateOrderRequest struct {
	SpotID          string                       `json:"spot_id"` // required
	ClientID        int                          `json:"client_id,omitempty"`
	FirstName       string                       `json:"first_name,omitempty"`
	LastName        string                       `json:"last_name,omitempty"`
	Phone           string                       `json:"phone"` // required (обязательно 11 цифр с + или без
	Email           string                       `json:"email,omitempty"`
	Sex             string                       `json:"sex,omitempty"`
	Birthday        string                       `json:"birthday,omitempty"`
	ClientAddressID int                          `json:"client_address_id,omitempty"`
	ClientAddress   CreateOrderAddressRequest    `json:"client_address"`
	ServiceMode     int                          `json:"service_mode,omitempty"`   // 1 — в заведении, 2 — навынос, 3 — доставка
	DeliveryPrice   int                          `json:"delivery_price,omitempty"` // указывается только для заказа с доставкой - 3
	Comment         string                       `json:"comment,omitempty"`
	Products        []CreateOrderProductRequest  `json:"products"` // required
	Payment         *CreateOrderPaymentRequest   `json:"payment,omitempty"`
	Promotion       *CreateOrderPromotionRequest `json:"promotion,omitempty"`
	DeliveryTime    string                       `json:"delivery_time,omitempty"` // время доставки
}

type CreateOrderProductRequest struct {
	ProductID     int                              `json:"product_id"`
	ModificatorID int                              `json:"modificator_id,omitempty"`
	Modifications []CreateOrderModificationRequest `json:"modification,omitempty"` // ??? модификации тех-карты
	Count         string                           `json:"count"`                  // required
	Price         int                              `json:"price"`
}

type CreateOrderAddressRequest struct {
	Address1  string `json:"address1"`
	Address2  string `json:"address2"`
	Comment   string `json:"comment"`
	Latitude  string `json:"lat,omitempty"`
	Longitude string `json:"lng,omitempty"`
}

type CreateOrderPaymentRequest struct {
	Type     int    `json:"type"`     // required - редварительная оплата (0 - не была, 1 - была) -> если оплаты нету, не нужно передавать
	Sum      int    `json:"sum"`      // required - сумма
	Currency string `json:"currency"` // required - payment ISO (RUB - рубль, USD - доллар) - тестовый EUR (получить информации об ISO POS можно в getAllSettings
}

type CreateOrderPromotionRequest struct {
	ID               int                           `json:"id"`                // Id акции которую нужно применить
	InvolvedProducts []CreateOrderInvolvedProducts `json:"involved_products"` //	Массив товаров которые участвуют в акции.
	ResultProducts   []CreateOrderResultProducts   `json:"result_products"`   // Массив товаров которые являются результатом акции. Нужно передавать только в бонусных акциях.
}

type CreateOrderInvolvedProducts struct {
	ID           int                            `json:"id"`
	Modification CreateOrderModificationRequest `json:"modification,omitempty"`
	Count        string                         `json:"count"`
}

type CreateOrderResultProducts struct {
	ID           int                            `json:"id"`
	Modification CreateOrderModificationRequest `json:"modification,omitempty"`
	Count        string                         `json:"count"`
}

type CreateOrderModificationRequest struct {
	M int `json:"m,omitempty"`
	A int `json:"a,omitempty"`
}

type CreateOrderProductResponse struct {
	IoProductID     int    `json:"io_product_id"`
	ProductId       int    `json:"product_id"`
	ModificatorID   int    `json:"modificator_id"`
	IncomingOrderID int    `json:"incoming_order_id"`
	Count           string `json:"count"`
	Price           int    `json:"price"`
	CreatedAt       string `json:"created_at"`
}

type CreateOrderBodyResponse struct {
	IncomingOrderID int `json:"incoming_order_id"`
	Type            int `json:"type"`
	SpotId          int `json:"spot_id"`
	Status          int `json:"status"`
	ClientId        int `json:"client_id"`
	ClientAddressID int `json:"client_address_id"`
	//TableId         interface{}          `json:"table_id"`
	//Comment         interface{}          `json:"comment"`
	CreatedAt     string      `json:"created_at"`
	UpdatedAt     string      `json:"updated_at"`
	TransactionId interface{} `json:"transaction_id"`
	ServiceMode   int         `json:"service_mode"`
	//DeliveryPrice   interface{}          `json:"delivery_price"`
	//FiscalSpreading int    `json:"fiscal_spreading"`
	FiscalMethod string `json:"fiscal_method"`
	//Promotion       interface{}          `json:"promotion"`
	DeliveryTime string `json:"delivery_time"`
	//PaymentMethodId interface{}          `json:"payment_method_id"`
	//FirstName       interface{}          `json:"first_name"`
	//LastName        interface{}          `json:"last_name"`
	Phone string `json:"phone"`
	//Email           interface{}          `json:"email"`
	//Sex             interface{}          `json:"sex"`
	//Birthday        interface{}          `json:"birthday"`
	//Address         interface{}          `json:"address"`
	Products []CreateOrderProductResponse `json:"products"`
}

type CreateOrderResponse struct {
	Response CreateOrderBodyResponse `json:"response"`
	ErrorResponse
}

type GetOrdersRequest struct {
	Status   string `json:"status"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}

type GetOrdersResponse struct {
	Response []Response `json:"response"`
}

type Response struct {
	IncomingOrderId int       `json:"incoming_order_id"`
	SpotId          int       `json:"spot_id"`
	Status          int       `json:"status"`
	ClientId        int       `json:"client_id"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Phone           string    `json:"phone"`
	Email           string    `json:"email"`
	Sex             int       `json:"sex"`
	Birthday        string    `json:"birthday"`
	Address         string    `json:"address"`
	Comment         string    `json:"comment"`
	CreatedAt       string    `json:"created_at"`
	UpdatedAt       string    `json:"updated_at"`
	TransactionId   int       `json:"transaction_id"`
	ServiceMode     int       `json:"service_mode"`
	Products        []Product `json:"products"`
}

type Product struct {
	IoProductId     int         `json:"io_product_id"`
	ProductId       int         `json:"product_id"`
	ModificatorId   interface{} `json:"modificator_id"`
	IncomingOrderId int         `json:"incoming_order_id"`
	Count           string      `json:"count"`
	CreatedAt       string      `json:"created_at"`
}
