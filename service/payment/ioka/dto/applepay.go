package dto

type ApplePaySessionOpenRequest struct {
	OrderID    string `json:"order_id"`
	Url        string `json:"url"`
	Platform   string `json:"platform"`
	DomainName string `json:"domain_name"`
}
