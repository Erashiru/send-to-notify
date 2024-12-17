package models

type Order struct {
	Guid                  string          `json:"guid,omitempty"`
	ShopGuid              string          `json:"shopGuid"`
	Type                  string          `json:"type"`
	Client                *Client         `json:"client,omitempty"`
	ItemList              []OrderItemList `json:"itemList"`
	Comment               string          `json:"comment"`
	PaidValue             float64         `json:"paidValue"`
	PrintFiscalCheck      bool            `json:"printFiscalCheck"`
	PrintFiscalCheckEmail interface{}     `json:"printFiscalCheckEmail"`
}

type CreateOrderResponse struct {
	Success bool `json:"success"`
	Count   int  `json:"count"`
	Rows    []struct {
		Guid   string `json:"guid"`
		Status string `json:"status"`
	} `json:"rows"`
	Error interface{} `json:"error"`
}

type OrderItemList struct {
	Guid              string         `json:"guid,omitempty"`
	MenuItemGuid      *string        `json:"menuItemGuid"`
	MenuTypeGuid      *string        `json:"menuTypeGuid"`
	SupplementList    map[string]int `json:"supplementList"`
	GoodsItemGuid     *string        `json:"goodsItemGuid"`
	PriceWithDiscount float64        `json:"priceWithDiscount"`
	Quantity          int            `json:"quantity"`
}

type Client struct {
	Name       string      `json:"name"`
	CardNumber interface{} `json:"cardNumber"`
	PhoneCode  string      `json:"phoneCode"`
	Phone      string      `json:"phone"`
	Email      string      `json:"email"`
}
