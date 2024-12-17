package models

type ApplePaySessionOpenRequest struct {
	PaymentSystem string
	OrderID       string
	Url           string `json:"url"`
	Platform      string `json:"platform"`
	DomainName    string `json:"domain_name"`
}

type ApplePaySessionOpenResponse struct {
	StatusCode                     string `json:"statusCode"`
	StatusMessage                  string `json:"statusMessage"`
	DisplayName                    string `json:"displayName"`
	DomainName                     string `json:"domainName"`
	EpochTimestamp                 int64  `json:"epochTimestamp"`
	ExpiresAt                      int64  `json:"expiresAt"`
	MerchantIdentifier             string `json:"merchantIdentifier"`
	MerchantSessionIdentifier      string `json:"merchantSessionIdentifier"`
	Nonce                          string `json:"nonce"`
	OperationalAnalyticsIdentifier string `json:"operationalAnalyticsIdentifier"`
	Retries                        int    `json:"retries"`
	Signature                      string `json:"signature"`
}

type ApplePayPayment struct {
	PaymentSystem string
	OrderID       string
	ToolType      string   `json:"tool_type"`
	ApplePay      ApplePay `json:"apple_pay"`
}

type Header struct {
	EphemeralPublicKey string `json:"ephemeralPublicKey"`
	PublicKeyHash      string `json:"publicKeyHash"`
	TransactionID      string `json:"transactionId"`
}

type PaymentData struct {
	Data      string `json:"data"`
	Header    Header `json:"header"`
	Signature string `json:"signature"`
	Version   string `json:"version"`
}

type PaymentMethod struct {
	DisplayName string `json:"displayName"`
	Network     string `json:"network"`
	Type        string `json:"type"`
}

type ApplePay struct {
	PaymentData           PaymentData   `json:"paymentData"`
	PaymentMethod         PaymentMethod `json:"paymentMethod"`
	TransactionIdentifier string        `json:"transactionIdentifier"`
}
