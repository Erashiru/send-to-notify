package dto

type WebhookEvent struct {
	Event            string           `json:"event"`
	OrderEvent       OrderEvent       `json:"order"`
	PaymentEvent     PaymentEvent     `json:"payment"`
	RefundEvent      RefundEvent      `json:"refund"`
	InstallmentEvent InstallmentEvent `json:"installment"`
	CustomerEvent    CustomerEvent    `json:"customer"`
	CardEvent        CardEvent        `json:"card"`
}

type OrderEvent struct {
	ID             string `json:"id"`
	ShopID         string `json:"shop_id"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
	Amount         int    `json:"amount"`
	Currency       string `json:"currency"`
	CaptureMethod  string `json:"capture_method"`
	ExternalID     string `json:"external_id"`
	Description    string `json:"description"`
	DueDate        string `json:"due_date"`
	SubscriptionID string `json:"subscription_id"`
}

type PayerEvent struct {
	Type          string `json:"type"`
	PanMasked     string `json:"pan_masked"`
	ExpiryDate    string `json:"expiry_date"`
	Holder        string `json:"holder"`
	PaymentSystem string `json:"payment_system"`
	Emitter       string `json:"emitter"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	CustomerID    string `json:"customer_id"`
	CardID        string `json:"card_id"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type PaymentEvent struct {
	ID             string   `json:"id"`
	OrderID        string   `json:"order_id"`
	Status         string   `json:"status"`
	CreatedAt      string   `json:"created_at"`
	ApprovedAmount int      `json:"approved_amount"`
	CapturedAmount int      `json:"captured_amount"`
	RefundedAmount int      `json:"refunded_amount"`
	ProcessingFee  int      `json:"processing_fee"`
	Payer          Payer    `json:"payer"`
	Error          Error    `json:"error"`
	Acquirer       Acquirer `json:"acquirer"`
	Action         Action   `json:"action"`
}

type RefundEvent struct {
	ID        string   `json:"id"`
	PaymentID string   `json:"payment_id"`
	OrderID   string   `json:"order_id"`
	Status    string   `json:"status"`
	CreatedAt string   `json:"created_at"`
	Error     Error    `json:"error"`
	Acquirer  Acquirer `json:"acquirer"`
	Amount    int64    `json:"amount"`
	Reason    string   `json:"reason"`
	Author    string   `json:"author"`
}

type InstallmentEvent struct {
	ID               string  `json:"id"`
	OrderID          string  `json:"order_id"`
	RedirectURL      string  `json:"redirect_url"`
	Status           string  `json:"status"`
	CreatedAt        string  `json:"created_at"`
	ProcessingFee    int     `json:"processing_fee"`
	MonthlyPayment   float64 `json:"monthly_payment"`
	InterestRate     float64 `json:"interest_rate"`
	EffectiveRate    float64 `json:"effective_rate"`
	Iin              string  `json:"iin"`
	Phone            string  `json:"phone"`
	Period           string  `json:"period"`
	ContractNumber   string  `json:"contract_number"`
	ContractSignedAt string  `json:"contract_signed_at"`
	Error            Error   `json:"error"`
}

type CustomerEvent struct {
	ID         string `json:"id"`
	ShopID     string `json:"shop_id"`
	CreatedAt  string `json:"created_at"`
	Status     string `json:"status"`
	ExternalID string `json:"external_id"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
}

type CardEvent struct {
	ID            string `json:"id"`
	CustomerID    string `json:"customer_id"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
	PanMasked     string `json:"pan_masked"`
	ExpiryDate    string `json:"expiry_date"`
	Holder        string `json:"holder"`
	PaymentSystem string `json:"payment_system"`
	Emitter       string `json:"emitter"`
	CvcRequired   bool   `json:"cvc_required"`
	Error         Error  `json:"error"`
	Action        Action `json:"action"`
}
