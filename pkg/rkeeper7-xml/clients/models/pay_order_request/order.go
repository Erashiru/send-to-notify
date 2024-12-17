package pay_order_request

type RK7Query struct {
	RK7CMD RK7CMD `xml:"RK7CMD"`
}

type RK7CMD struct {
	CMD     string  `xml:"CMD,attr"`
	Order   Order   `xml:"Order"`
	Cashier Cashier `xml:"Cashier"`
	Station Station `xml:"Station"`
	Payment Payment `xml:"Payment"`
}

type LicenseInfo struct {
	Anchor          string          `xml:"anchor,attr"`
	LicenseToken    string          `xml:"licenseToken,attr"`
	LicenseInstance LicenseInstance `xml:"LicenseInstance"`
}

type LicenseInstance struct {
	Guid      string `xml:"guid,attr"`
	SeqNumber string `xml:"seqNumber,attr"`
}

type Order struct {
	Guid string `xml:"guid,attr"`
}

type Cashier struct {
	Code string `xml:"code,attr"`
}

type Station struct {
	Code string `xml:"code,attr"`
}

type Payment struct {
	ID     string `xml:"id,attr"`
	Amount string `xml:"amount,attr"`
}
