package save_order_response

type RK7QueryResult struct {
	ServerVersion   string      `xml:"ServerVersion,attr"`
	XmlVersion      string      `xml:"XmlVersion,attr"`
	NetName         string      `xml:"NetName,attr"`
	Status          string      `xml:"Status,attr"`
	CMD             string      `xml:"CMD,attr"`
	ErrorText       string      `xml:"ErrorText,attr"`
	DateTime        string      `xml:"DateTime,attr"`
	WorkTime        string      `xml:"WorkTime,attr"`
	Processed       string      `xml:"Processed,attr"`
	ArrivalDateTime string      `xml:"ArrivalDateTime,attr"`
	Order           Order       `xml:"Order"`
	Session         MainSession `xml:"Session"`
}

type MainSession struct {
	LineGuid  string `xml:"line_guid,attr"`
	SessionID string `xml:"sessionID,attr"`
}

type Creator struct {
	ID   string `xml:"id,attr"`
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
	Role Role   `xml:"Role"`
}

type Role struct {
	ID   string `xml:"id,attr"`
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
}

type Waiter struct {
	ID   string `xml:"id,attr"`
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
	Role Role   `xml:"Role"`
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
	Guest Guest  `xml:"Guest"`
}

type Guest struct {
	GuestLabel string `xml:"guestLabel,attr"`
}

type ExternalProps struct {
	Prop Prop `xml:"Prop"`
}

type Prop struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type Author struct {
	ID   string `xml:"id,attr"`
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
	Role Role   `xml:"Role"`
}

type Dish struct {
	ID          string `xml:"id,attr"`
	Code        string `xml:"code,attr"`
	Name        string `xml:"name,attr"`
	Uni         string `xml:"uni,attr"`
	LineGuid    string `xml:"line_guid,attr"`
	State       string `xml:"state,attr"`
	Guid        string `xml:"guid,attr"`
	Price       string `xml:"price,attr"`
	Amount      string `xml:"amount,attr"`
	Quantity    string `xml:"quantity,attr"`
	SrcQuantity string `xml:"srcQuantity,attr"`
}

type PriceScale struct {
	ID   string `xml:"id,attr"`
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
}

type TradeGroup struct {
	ID   string `xml:"id,attr"`
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
}

type Session struct {
	Uni          string     `xml:"uni,attr"`
	LineGuid     string     `xml:"line_guid,attr"`
	State        string     `xml:"state,attr"`
	SessionID    string     `xml:"sessionID,attr"`
	IsDraft      string     `xml:"isDraft,attr"`
	RemindTime   string     `xml:"remindTime,attr"`
	StartService string     `xml:"startService,attr"`
	Printed      string     `xml:"printed,attr"`
	CookMins     string     `xml:"cookMins,attr"`
	Station      Station    `xml:"Station"`
	Author       Author     `xml:"Author"`
	Creator      Creator    `xml:"Creator"`
	Dish         []Dish     `xml:"Dish"`
	PriceScale   PriceScale `xml:"PriceScale"`
	TradeGroup   TradeGroup `xml:"TradeGroup"`
}

type Order struct {
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
	BasicSum             string        `xml:"basicSum,attr"`
	NationalSum          string        `xml:"nationalSum,attr"`
	Creator              Creator       `xml:"Creator"`
	Waiter               Waiter        `xml:"Waiter"`
	OrderCategory        OrderCategory `xml:"OrderCategory"`
	OrderType            OrderType     `xml:"OrderType"`
	Table                Table         `xml:"Table"`
	Station              Station       `xml:"Station"`
	Guests               Guests        `xml:"Guests"`
	ExternalProps        ExternalProps `xml:"ExternalProps"`
	Session              []Session     `xml:"Session"`
}
