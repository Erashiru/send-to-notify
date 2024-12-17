package models

const (
	SERVICE_TYPE_DELIVERY_BY_COURIER = "DeliveryByCourier"
	SERVICE_TYPE_DELIVERY_BY_CLIENT  = "DeliveryByClient"
	SERVICE_TYPE_DELIVERY_PICKUP     = "DeliveryPickUp"
	LegacyAddressType                = "legacy"
	CityAddressType                  = "city"
)

type CreateOrderSettings struct {
	TransportToFrontTimeout int `json:"transportToFrontTimeout,omitempty"`
}

type Street struct {
	ClassifierID string `json:"classifierId,omitempty"`
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	City         string `json:"city,omitempty"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

type Address struct {
	Street    *Street `json:"street,omitempty"`
	Index     string  `json:"index,omitempty"`
	House     string  `json:"house,omitempty"`
	Building  string  `json:"building,omitempty"`
	Flat      string  `json:"flat,omitempty"`
	Entrance  string  `json:"entrance,omitempty"`
	Floor     string  `json:"floor,omitempty"`
	Doorphone string  `json:"doorphone,omitempty"`
	RegionID  string  `json:"regionId,omitempty"`
	Type      string  `json:"type,omitempty"`
}

type DeliveryPoint struct {
	Coordinates           *Coordinates `json:"coordinates,omitempty"`
	Address               *Address     `json:"address,omitempty"`
	ExternalCartographyID string       `json:"externalCartographyId,omitempty"`
	Comment               string       `json:"comment,omitempty"`
}

type Customer struct {
	ID                                    string `json:"id,omitempty"`
	Name                                  string `json:"name,omitempty"`
	Surname                               string `json:"surname,omitempty"`
	Comment                               string `json:"comment,omitempty"`
	Birthdate                             string `json:"birthdate,omitempty"`
	Email                                 string `json:"email,omitempty"`
	ShouldReceiveOrderStatusNotifications bool   `json:"shouldReceiveOrderStatusNotifications,omitempty"`
	Gender                                string `json:"gender,omitempty"`
	InBlacklist                           bool   `json:"inBlacklist,omitempty"`
	BlacklistReason                       string `json:"blacklistReason,omitempty"`
	Type                                  string `json:"type,omitempty"`
}

type Guests struct {
	Count               int  `json:"count,omitempty"`
	SplitBetweenPersons bool `json:"splitBetweenPersons,omitempty"`
}

type PaymentAdditionalData struct {
	Credential  string `json:"credential,omitempty"`
	SearchScope string `json:"searchScope,omitempty"`
	Type        string `json:"type,omitempty"`
}

type Payment struct {
	PaymentTypeKind        string                 `json:"paymentTypeKind,omitempty"`
	Sum                    int                    `json:"sum,omitempty"`
	PaymentTypeID          string                 `json:"paymentTypeId,omitempty"`
	IsProcessedExternally  bool                   `json:"isProcessedExternally"`
	PaymentAdditionalData  *PaymentAdditionalData `json:"paymentAdditionalData,omitempty"`
	IsFiscalizedExternally bool                   `json:"isFiscalizedExternally,omitempty"`
}

type Tip struct {
	PaymentTypeKind        string                 `json:"paymentTypeKind,omitempty"`
	TipsTypeID             string                 `json:"tipsTypeId,omitempty"`
	Sum                    int                    `json:"sum,omitempty"`
	PaymentTypeID          string                 `json:"paymentTypeId,omitempty"`
	IsProcessedExternally  bool                   `json:"isProcessedExternally"`
	PaymentAdditionalData  *PaymentAdditionalData `json:"paymentAdditionalData,omitempty"`
	IsFiscalizedExternally bool                   `json:"isFiscalizedExternally,omitempty"`
}

type Card struct {
	Track string `json:"track,omitempty"`
}

type Discount struct {
	DiscountTypeId     string   `json:"discountTypeId,omitempty"`
	Sum                float64  `json:"sum,omitempty"`
	SelectivePositions []string `json:"selectivePositions,omitempty"`
	Type               string   `json:"type,omitempty"`
}

type DiscountsInfo struct {
	Card      *Card      `json:"card,omitempty"`
	Discounts []Discount `json:"discounts,omitempty"`
}

type IikoCard5Info struct {
	Coupon                     string   `json:"coupon,omitempty"`
	ApplicableManualConditions []string `json:"applicableManualConditions,omitempty"`
}

type Combo struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Amount    int    `json:"amount,omitempty"`
	Price     int    `json:"price,omitempty"`
	SourceID  string `json:"sourceId,omitempty"`
	ProgramID string `json:"programId,omitempty"`
}

type ComboInformation struct {
	ComboID       string `json:"comboId,omitempty"`
	ComboSourceID string `json:"comboSourceId,omitempty"`
	ComboGroupID  string `json:"comboGroupId,omitempty"`
}

type ItemModifier struct {
	ProductId      string  `json:"productId,omitempty"`
	Amount         float64 `json:"amount,omitempty"`
	ProductGroupId string  `json:"productGroupId,omitempty"`
	Price          float64 `json:"price"`
	PositionId     string  `json:"positionId,omitempty"`
}

func (im ItemModifier) IsProductGroupID(productGroupID string) {

}

type Item struct {
	ProductId        string            `json:"productId,omitempty"`
	Modifiers        []ItemModifier    `json:"modifiers,omitempty"`
	Price            *float64          `json:"price,omitempty"` // TODO: logic for GIFT promo when price equal 0 (yandex)???
	PositionId       string            `json:"positionId,omitempty"`
	Type             string            `json:"type,omitempty"`
	Amount           float64           `json:"amount,omitempty"`
	ProductSizeID    string            `json:"productSizeId,omitempty"`
	ComboInformation *ComboInformation `json:"comboInformation,omitempty"`
	Comment          string            `json:"comment,omitempty"`
	TableProduct     TableProduct      `json:"product,omitempty"`
}

type TableProduct struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Order struct {
	ID                string         `json:"id,omitempty"`
	ExternalNumber    string         `json:"externalNumber,omitempty"`
	CompleteBefore    string         `json:"completeBefore,omitempty"`
	Phone             string         `json:"phone,omitempty"`
	OrderTypeID       string         `json:"orderTypeId,omitempty"`
	OrderServiceType  string         `json:"orderServiceType,omitempty"`
	DeliveryPoint     *DeliveryPoint `json:"deliveryPoint,omitempty"`
	Comment           string         `json:"comment,omitempty"`
	Customer          *Customer      `json:"customer,omitempty"`
	Guests            *Guests        `json:"guests,omitempty"`
	MarketingSourceID string         `json:"marketingSourceId,omitempty"`
	OperatorID        string         `json:"operatorId,omitempty"`
	Items             []Item         `json:"items,omitempty"`
	Combos            []Combo        `json:"combos,omitempty"`
	Payments          []Payment      `json:"payments,omitempty"`
	Tips              []Tip          `json:"tips,omitempty"`
	SourceKey         string         `json:"sourceKey,omitempty"`
	DiscountsInfo     *DiscountsInfo `json:"discountsInfo,omitempty"`
	IikoCard5Info     *IikoCard5Info `json:"iikoCard5Info,omitempty"`
}

type ErrorInfo struct {
	Code           string      `json:"code,omitempty"`
	Message        string      `json:"message,omitempty"`
	Description    string      `json:"description,omitempty"`
	AdditionalData interface{} `json:"additionalData,omitempty"`
}

type OrderInfo struct {
	ID             string     `json:"id,omitempty"`
	ExternalNumber string     `json:"externalNumber,omitempty"`
	OrganizationID string     `json:"organizationId,omitempty"`
	Timestamp      int        `json:"timestamp,omitempty"`
	CreationStatus string     `json:"creationStatus,omitempty"`
	ErrorInfo      *ErrorInfo `json:"errorInfo,omitempty"`
}

type CreateDeliveryRequest struct {
	OrganizationID      string               `json:"organizationId,omitempty"`
	TerminalGroupID     string               `json:"terminalGroupId,omitempty"`
	CreateOrderSettings *CreateOrderSettings `json:"createOrderSettings,omitempty"`
	Order               *Order               `json:"order,omitempty"`
}

type CreateDeliveryResponse struct {
	CorrelationID string     `json:"correlationId,omitempty"`
	OrderInfo     *OrderInfo `json:"orderInfo,omitempty"`
}

type RetrieveDeliveryRequest struct {
	OrganizationId string   `json:"organizationId"`
	OrderIds       []string `json:"orderIds"`
}

type RetrieveOrder struct {
	Id             string    `json:"id"`
	PosId          string    `json:"posId"`
	ExternalNumber string    `json:"externalNumber"`
	OrganizationId string    `json:"organizationId"`
	Timestamp      int       `json:"timestamp"`
	CreationStatus string    `json:"creationStatus"`
	ErrorInfo      ErrorInfo `json:"errorInfo"`
}

type RetrieveDeliveryResponse struct {
	CorrelationId string          `json:"correlationId"`
	Orders        []RetrieveOrder `json:"orders"`
}

type OrderItem struct {
	OrganizationId string  `json:"organizationId"`
	OrderId        string  `json:"orderId"`
	Items          []Item  `json:"items"`
	Combos         []Combo `json:"combos,omitempty"`
}
type OrderItemResponse struct {
	CorrelationId    string `json:"correlationId"`
	ErrorDescription string `json:"errorDescription"`
	Error            string `json:"error"`
}

type CancelDeliveryResponse struct {
	OrganizationId string `json:"organizationId"`
	RemovalTypeId  string `json:"removalTypeId"`
	OrderId        string `json:"orderId"`
}

type SendNotificationRequest struct {
	OrderSource    string `json:"orderSource"`
	OrderId        string `json:"orderId"`
	AdditionalInfo string `json:"additionalInfo"`
	MessageType    string `json:"messageType"`
	OrganizationId string `json:"organizationId"`
}

type UpdateOrderProblem struct {
	OrganizationId string `json:"organizationId"`
	OrderId        string `json:"orderId"`
	HasProblem     bool   `json:"hasProblem"`
	Problem        string `json:"problem"`
}

type CloseOrderRequest struct {
	OrganizationId string `json:"organizationId"`
	OrderId        string `json:"orderId"`
}
