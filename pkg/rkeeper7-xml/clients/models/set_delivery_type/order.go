package set_delivery_type

type RK7Query struct {
	RK7CMD RK7CMD `xml:"RK7CMD"`
}

type RK7CMD struct {
	CMD       string    `xml:"CMD,attr"`
	Order     Order     `xml:"Order"`
	OrderType OrderType `xml:"OrderType"`
	ExtSource ExtSource `xml:"ExtSource"`
}

type Order struct {
	OrderIdent string `xml:"orderIdent,attr"`
	Visit      string `xml:"visit,attr"`
}

type DeliveryBlock struct {
	DeliveryState string `xml:"deliveryState,attr"`
}

type OrderType struct {
	ID string `xml:"id,attr"`
}

type ExtSource struct {
	Source string `xml:"source,attr"`
}
