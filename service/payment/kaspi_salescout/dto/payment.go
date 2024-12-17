package dto

type CreatePaymentOrderRequest struct {
	Amount        float64 `json:"amount"`
	MerchantID    string  `json:"merchantId"`
	TransactionID string  `json:"transactionId"`
}

type CreatePaymentOrderResponse struct {
	PaymentLink            string                 `json:"paymentLink"`
	PaymentId              int                    `json:"paymentId"`
	PaymentTypes           []string               `json:"paymentTypes"`
	ExternalId             string                 `json:"externalId"`
	Status                 string                 `json:"status"`
	PaymentBehaviorOptions PaymentBehaviorOptions `json:"paymentBehaviorOptions"`
}

type CreatePaymentTokenResponse struct {
	QrToken                string                      `json:"qrToken"`
	PaymentId              int                         `json:"paymentId"`
	PaymentTypes           []string                    `json:"paymentTypes"`
	ExternalId             string                      `json:"externalId"`
	Status                 string                      `json:"status"`
	PaymentBehaviorOptions PaymentBehaviorOptionsToken `json:"paymentBehaviorOptions"`
}

type PaymentBehaviorOptionsToken struct {
	QrCodeScanWaitTimeout      int `json:"qrCodeScanWaitTimeout"`
	PaymentConfirmationTimeout int `json:"paymentConfirmationTimeout"`
	StatusPollingInterval      int `json:"statusPollingInterval"`
}

type PaymentBehaviorOptions struct {
	LinkActivationWaitTimeout  int `json:"linkActivationWaitTimeout"`
	PaymentConfirmationTimeout int `json:"paymentConfirmationTimeout"`
	StatusPollingInterval      int `json:"statusPollingInterval"`
}

type PaymentStatusResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

type RefundRequest struct {
	Amount     float64 `json:"amount"`
	PaymentID  int     `json:"paymentId"`
	MerchantID string  `json:"merchantId"`
}

type RefundResponse struct {
	ReturnOperationID int    `json:"returnOperationId"`
	Status            string `json:"status"`
}
