package get_order_request

type RK7Query struct {
	RK7Command RK7Command `xml:"RK7Command"`
}

type RK7Command struct {
	CMD   string `xml:"CMD,attr"`
	Order Order  `xml:"Order"`
}

type Order struct {
	Visit      string `xml:"visit,attr"`
	OrderIdent string `xml:"orderIdent,attr"`
}
