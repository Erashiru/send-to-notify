package models

type LicenseResponse struct {
	Id             string  `json:"id"`
	ExpirationDate string  `json:"expirationDate"`
	Qty            float64 `json:"qty"`
}
