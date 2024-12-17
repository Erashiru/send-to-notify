package create_order_request

type RK7CMD struct {
	CMD   string             `xml:"CMD,attr"`
	Order CreateOrderRequest `xml:"Order"`
}

type Table struct {
	Code string `xml:"code,attr"`
}

type Station struct {
	ID string `xml:"id,attr"`
}

type GuestType struct {
	ID string `xml:"id,attr"`
}

type Guests struct {
	Guest Guest `xml:"Guest"`
}

type Guest struct {
	GuestLabel string `xml:"GuestLabel,attr"`
}

type OrderType struct {
	Code string `xml:"code,attr"`
}

type CreateOrderRequest struct {
	PersistentComment string    `xml:"persistentComment,attr"`
	Table             Table     `xml:"Table"`
	OrderType         OrderType `xml:"OrderType"`
	Station           Station   `xml:"Station"`
	GuestType         GuestType `xml:"GuestType"`
	Guests            Guests    `xml:"Guests"`
}

type RK7Query struct {
	RK7CMD RK7CMD `xml:"RK7CMD"`
}
