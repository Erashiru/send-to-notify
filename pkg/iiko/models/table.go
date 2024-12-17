package models

type TableRequest struct {
	TerminalGroupIds []string `json:"terminalGroupIds"`
}

type TableResponse struct {
	CorrelationId      string `json:"correlationId"`
	RestaurantSections []struct {
		Id              string      `json:"id"`
		TerminalGroupId string      `json:"terminalGroupId"`
		Name            string      `json:"name"`
		Tables          []Table     `json:"tables"`
		Schema          interface{} `json:"schema"`
	} `json:"restaurantSections"`
	Revision int64 `json:"revision"`
}

type Table struct {
	Id              string `json:"id"`
	Number          int    `json:"number"`
	Name            string `json:"name"`
	SeatingCapacity int    `json:"seatingCapacity"`
	Revision        int64  `json:"revision"`
	IsDeleted       bool   `json:"isDeleted"`
}

type OrdersByTablesRequest struct {
	//SourceKeys      []string `json:"sourceKeys"`
	OrganizationIds []string `json:"organizationIds"`
	TableIds        []string `json:"tableIds"`
	Statuses        []string `json:"statuses"`
	//DateFrom        string   `json:"dateFrom"`
	//DateTo          string   `json:"dateTo"`
}

type OrdersByTablesResponse struct {
	CorrelationID string   `json:"correlationId"`
	Orders        []Orders `json:"orders"`
}

type Waiter struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type Conception struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type RemovalType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DeletionMethod struct {
	ID          string      `json:"id"`
	Comment     string      `json:"comment"`
	RemovalType RemovalType `json:"removalType"`
}

type Deleted struct {
	DeletionMethod DeletionMethod `json:"deletionMethod"`
}

type Combos struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	Price    float64 `json:"price"`
	SourceID string  `json:"sourceId"`
	Size     Size    `json:"size"`
}

type PaymentType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type Payments struct {
	PaymentType            PaymentType `json:"paymentType"`
	Sum                    float64     `json:"sum"`
	IsPreliminary          bool        `json:"isPreliminary"`
	IsExternal             bool        `json:"isExternal"`
	IsProcessedExternally  bool        `json:"isProcessedExternally"`
	IsFiscalizedExternally bool        `json:"isFiscalizedExternally"`
	IsPrepay               bool        `json:"isPrepay"`
}

type TipsType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Tips struct {
	TipsType               TipsType    `json:"tipsType"`
	PaymentType            PaymentType `json:"paymentType"`
	Sum                    float64     `json:"sum"`
	IsPreliminary          bool        `json:"isPreliminary"`
	IsExternal             bool        `json:"isExternal"`
	IsProcessedExternally  bool        `json:"isProcessedExternally"`
	IsFiscalizedExternally bool        `json:"isFiscalizedExternally"`
	IsPrepay               bool        `json:"isPrepay"`
}

type DiscountType struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SelectivePositionsWithSum struct {
	PositionID string  `json:"positionId"`
	Sum        float64 `json:"sum"`
}

type Discounts struct {
	DiscountType              DiscountType                `json:"discountType"`
	Sum                       float64                     `json:"sum"`
	SelectivePositions        []string                    `json:"selectivePositions"`
	SelectivePositionsWithSum []SelectivePositionsWithSum `json:"selectivePositionsWithSum"`
}

type LoyaltyInfo struct {
	Coupon                  string   `json:"coupon"`
	AppliedManualConditions []string `json:"appliedManualConditions"`
}

type ExternalData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Orders struct {
	ID             string    `json:"id"`
	PosID          string    `json:"posId"`
	ExternalNumber string    `json:"externalNumber"`
	OrganizationID string    `json:"organizationId"`
	Timestamp      int       `json:"timestamp"`
	CreationStatus string    `json:"creationStatus"`
	ErrorInfo      ErrorInfo `json:"errorInfo"`
	Order          OrderResp `json:"order"`
}

type OrderResp struct {
	TableIds                       []string       `json:"tableIds"`
	Customer                       Customer       `json:"customer"`
	Phone                          string         `json:"phone"`
	Status                         string         `json:"status"`
	WhenCreated                    string         `json:"whenCreated"`
	Waiter                         Waiter         `json:"waiter"`
	TabName                        string         `json:"tabName"`
	SplitOrderBetweenCashRegisters string         `json:"splitOrderBetweenCashRegisters"`
	MenuID                         string         `json:"menuId"`
	Sum                            float64        `json:"sum"`
	Number                         float64        `json:"number"`
	SourceKey                      string         `json:"sourceKey"`
	WhenBillPrinted                string         `json:"whenBillPrinted"`
	WhenClosed                     string         `json:"whenClosed"`
	Conception                     Conception     `json:"conception"`
	GuestsInfo                     GuestsInfo     `json:"guestsInfo"`
	Items                          []ItemsResp    `json:"items"`
	Combos                         []Combos       `json:"combos"`
	Payments                       []Payments     `json:"payments"`
	Tips                           []Tips         `json:"tips"`
	Discounts                      []Discounts    `json:"discounts"`
	OrderType                      OrderType      `json:"orderType"`
	TerminalGroupID                string         `json:"terminalGroupId"`
	ProcessedPaymentsSum           float64        `json:"processedPaymentsSum"`
	LoyaltyInfo                    LoyaltyInfo    `json:"loyaltyInfo"`
	ExternalData                   []ExternalData `json:"externalData"`
}

type ItemsResp struct {
	Type    string `json:"type"`
	Product struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"product"`
	Modifiers        []Modifiers      `json:"modifiers"`
	Price            float64          `json:"price"`
	Cost             float64          `json:"cost"`
	PricePredefined  bool             `json:"pricePredefined"`
	PositionId       string           `json:"positionId"`
	ResultSum        float64          `json:"resultSum"`
	Status           string           `json:"status"`
	Deleted          Deleted          `json:"deleted"`
	Amount           float64          `json:"amount"`
	Comment          string           `json:"comment"`
	WhenPrinted      string           `json:"whenPrinted"`
	Size             Size             `json:"size"`
	ComboInformation ComboInformation `json:"comboInformation"`
}
type Modifiers struct {
	Product struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"product"`
	Amount       float64 `json:"amount"`
	ProductGroup struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"productGroup"`
	Price     float64 `json:"price"`
	ResultSum float64 `json:"resultSum"`
}

type GetOrdersByIDsRequest struct {
	OrganizationIds []string `json:"organizationIds"`
	OrderIds        []string `json:"orderIds"`
}

type ChangePaymentReq struct {
	OrganizationId string       `json:"organizationId"`
	OrderId        string       `json:"orderId"`
	Payments       []PaymentReq `json:"payments"`
}

type GetCommandStatusReq struct {
	OrganizationId string `json:"organizationId"`
	CorrelationId  string `json:"correlationId"`
}
type Exception struct {
	Message string `json:"message"`
}

type GetCommandStatusResp struct {
	State     string    `json:"state"`
	Exception Exception `json:"exception"`
}

type CloseTableOrderReq struct {
	OrganizationId string `json:"organizationId"`
	OrderId        string `json:"orderId"`
}

type PaymentReq struct {
	PaymentTypeKind        string  `json:"paymentTypeKind"`
	Sum                    float64 `json:"sum"`
	PaymentTypeId          string  `json:"paymentTypeId"`
	IsProcessedExternally  bool    `json:"isProcessedExternally"`
	IsFiscalizedExternally bool    `json:"isFiscalizedExternally"`
	IsPrepay               bool    `json:"isPrepay"`
}
