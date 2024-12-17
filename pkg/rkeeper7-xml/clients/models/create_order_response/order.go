package create_order_response

type Creator struct {
	ID   string `xml:"id,attr"`
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
	Role struct {
		ID   string `xml:"id,attr"`
		Code string `xml:"code,attr"`
		Name string `xml:"name,attr"`
	} `xml:"Role"`
}

type Waiter struct {
	ID   string `xml:"id,attr"`
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
	Role struct {
		ID   string `xml:"id,attr"`
		Code string `xml:"code,attr"`
		Name string `xml:"name,attr"`
	} `xml:"Role"`
}

type OrderCategory struct {
	ID   string `xml:"id,attr"`
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
}

type OrderType struct {
	ID   string `xml:"id,attr"`
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
}

type Table struct {
	ID   string `xml:"id,attr"`
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
}

type Station struct {
	ID   string `xml:"id,attr"`
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
}

type Guests struct {
	Count string `xml:"count,attr"`
	Guest struct {
		GuestLabel string `xml:"guestLabel,attr"`
	} `xml:"Guest"`
}

type ExternalProp struct {
	Prop struct {
		Name  string `xml:"name,attr"`
		Value string `xml:"value,attr"`
	} `xml:"Prop"`
}

type CreateOrderBodyResponse struct {
	Visit                string        `xml:"visit,attr"`
	OrderIdent           string        `xml:"orderIdent,attr"`
	Guid                 string        `xml:"guid,attr"`
	URL                  string        `xml:"url,attr"`
	OrderName            string        `xml:"orderName,attr"`
	Version              string        `xml:"version,attr"`
	Crc32                string        `xml:"crc32,attr"`
	OrderSum             string        `xml:"orderSum,attr"`
	UnpaidSum            string        `xml:"unpaidSum,attr"`
	DiscountSum          string        `xml:"discountSum,attr"`
	TotalPieces          string        `xml:"totalPieces,attr"`
	SeqNumber            string        `xml:"seqNumber,attr"`
	Paid                 string        `xml:"paid,attr"`
	Finished             string        `xml:"finished,attr"`
	PersistentComment    string        `xml:"persistentComment,attr"`
	NonPersistentComment string        `xml:"nonPersistentComment,attr"`
	OpenTime             string        `xml:"openTime,attr"`
	CookMins             string        `xml:"cookMins,attr"`
	Creator              Creator       `xml:"Creator"`
	Waiter               Waiter        `xml:"Waiter"`
	OrderCategory        OrderCategory `xml:"OrderCategory"`
	OrderType            OrderType     `xml:"OrderType"`
	Table                Table         `xml:"Table"`
	Station              Station       `xml:"Station"`
	Guests               Guests        `xml:"Guests"`
	ExternalProps        ExternalProp  `xml:"ExternalProps"`
}

type CreateOrderResponse struct {
	ServerVersion   string                  `xml:"ServerVersion,attr"`
	XmlVersion      string                  `xml:"XmlVersion,attr"`
	NetName         string                  `xml:"NetName,attr"`
	Status          string                  `xml:"Status,attr"`
	CMD             string                  `xml:"CMD,attr"`
	VisitID         string                  `xml:"VisitID,attr"`
	OrderID         string                  `xml:"OrderID,attr"`
	Guid            string                  `xml:"guid,attr"`
	ErrorText       string                  `xml:"ErrorText,attr"`
	DateTime        string                  `xml:"DateTime,attr"`
	WorkTime        string                  `xml:"WorkTime,attr"`
	Processed       string                  `xml:"Processed,attr"`
	ArrivalDateTime string                  `xml:"ArrivalDateTime,attr"`
	Order           CreateOrderBodyResponse `xml:"Order"`
}
