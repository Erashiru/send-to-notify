package dto

type CreateReceiptRequest struct {
	ID     int64                      `json:"id"`
	Method string                     `json:"method"`
	Params CreateReceiptRequestParams `json:"params"`
}

type Account struct {
	OrderID string `json:"order_id"`
}

type Shipping struct {
	Title string `json:"title"`
	Price int    `json:"price"`
}

type Items struct {
	Discount    string `json:"discount"`
	Title       string `json:"title"`
	Price       int    `json:"price"`
	Count       int    `json:"count"`
	Code        string `json:"code"`
	Units       int    `json:"units"`
	VatPercent  int    `json:"vat_percent"`
	PackageCode string `json:"package_code"`
}

type Detail struct {
	ReceiptType int      `json:"receipt_type"`
	Shipping    Shipping `json:"shipping"`
	Items       []Items  `json:"items"`
}

type CreateReceiptRequestParams struct {
	Amount  int     `json:"amount"`
	Account Account `json:"account"`
	Detail  Detail  `json:"detail"`
}

type CreateReceiptResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  Result `json:"result"`
}
type Discount struct {
	Title string `json:"title"`
	Price int    `json:"price"`
}

type Epos struct {
	MerchantID string `json:"merchantId"`
	TerminalID string `json:"terminalId"`
}

type Merchant struct {
	ID           string `json:"_id"`
	Name         string `json:"name"`
	Organization string `json:"organization"`
	Address      string `json:"address"`
	BusinessID   string `json:"business_id"`
	Epos         Epos   `json:"epos"`
	Date         int64  `json:"date"`
	Logo         any    `json:"logo"`
	Type         string `json:"type"`
	Terms        any    `json:"terms"`
}

type Meta struct {
	Source string `json:"source"`
	Owner  string `json:"owner"`
}

type Receipt struct {
	ID           string    `json:"_id"`
	CreateTime   int64     `json:"create_time"`
	PayTime      int       `json:"pay_time"`
	CancelTime   int       `json:"cancel_time"`
	State        int       `json:"state"`
	Type         int       `json:"type"`
	External     bool      `json:"external"`
	Operation    int       `json:"operation"`
	Category     any       `json:"category"`
	Error        any       `json:"error"`
	Description  string    `json:"description"`
	Detail       Detail    `json:"detail"`
	Amount       int       `json:"amount"`
	Currency     int       `json:"currency"`
	Commission   int       `json:"commission"`
	Account      []Account `json:"account"`
	Card         any       `json:"card"`
	Merchant     Merchant  `json:"merchant"`
	Meta         Meta      `json:"meta"`
	ProcessingID any       `json:"processing_id"`
}

type Result struct {
	Receipt Receipt `json:"receipt"`
}
