package get_order_response

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
	OrderCategory        OrderCategory `xml:"OrderCategory"`
	OrderType            OrderType     `xml:"OrderType"`
	Table                Table         `xml:"Table"`
}

type CommandResult struct {
	CMD       string `xml:"CMD,attr"`
	Status    string `xml:"Status,attr"`
	ErrorText string `xml:"ErrorText,attr"`
	DateTime  string `xml:"DateTime,attr"`
	WorkTime  string `xml:"WorkTime,attr"`
	Order     Order  `xml:"Order"`
}

type RK7QueryResult struct {
	ServerVersion   string        `xml:"ServerVersion,attr"`
	XmlVersion      string        `xml:"XmlVersion,attr"`
	NetName         string        `xml:"NetName,attr"`
	Status          string        `xml:"Status,attr"`
	Processed       string        `xml:"Processed,attr"`
	ArrivalDateTime string        `xml:"ArrivalDateTime,attr"`
	CommandResult   CommandResult `xml:"CommandResult"`
}
