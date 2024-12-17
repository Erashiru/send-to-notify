package models

type GetCustomerInfoRequest struct {
	Phone          string `json:"phone"`
	Type           string `json:"type"`
	OrganizationId string `json:"organizationId"`
}

type GetCustomerInfoResponse struct {
	Id                            string             `json:"id"`
	ReferrerId                    string             `json:"referrerId"`
	Name                          string             `json:"name"`
	Surname                       string             `json:"surname"`
	MiddleName                    string             `json:"middleName"`
	Comment                       string             `json:"comment"`
	Phone                         string             `json:"phone"`
	CultureName                   string             `json:"cultureName"`
	Birthday                      string             `json:"birthday"`
	Email                         string             `json:"email"`
	Sex                           int                `json:"sex"`
	ConsentStatus                 int                `json:"consentStatus"`
	Anonymized                    bool               `json:"anonymized"`
	Cards                         []CustomerCard     `json:"cards"`
	Categories                    []CustomerCategory `json:"categories"`
	WalletBalances                []WalletBalance    `json:"walletBalances"`
	UserData                      string             `json:"userData"`
	ShouldReceivePromoActionsInfo bool               `json:"shouldReceivePromoActionsInfo"`
	ShouldReceiveLoyaltyInfo      bool               `json:"shouldReceiveLoyaltyInfo"`
	ShouldReceiveOrderStatusInfo  bool               `json:"shouldReceiveOrderStatusInfo"`
	PersonalDataConsentFrom       string             `json:"personalDataConsentFrom"`
	PersonalDataConsentTo         string             `json:"personalDataConsentTo"`
	PersonalDataProcessingFrom    string             `json:"personalDataProcessingFrom"`
	PersonalDataProcessingTo      string             `json:"personalDataProcessingTo"`
	IsDeleted                     bool               `json:"isDeleted"`
}

type GetTransactionInfoReq struct {
	CustomerId     string `json:"customerId"`
	PageSize       int    `json:"pageSize"`
	OrganizationId string `json:"organizationId"`
}

type GetTransactionInfoResp struct {
	Transactions      []Transaction `json:"transactions"`
	LastRevision      int           `json:"lastRevision"`
	LastTransactionId string        `json:"lastTransactionId,omitempty"`
	PageSize          int           `json:"pageSize"`
}

type Transaction struct {
	ApiClientLogin       string      `json:"apiClientLogin"`
	BalanceAfter         float64     `json:"balanceAfter"`
	BalanceBefore        float64     `json:"balanceBefore"`
	BlockReason          string      `json:"blockReason"`
	Certificate          Certificate `json:"certificate"`
	Comment              string      `json:"comment"`
	Counteragent         string      `json:"counteragent"`
	CounteragentType     int         `json:"counteragentType"`
	CounteragentTypeName string      `json:"counteragentTypeName"`
	Coupon               Coupon      `json:"coupon"`
	EmitentName          string      `json:"emitentName"`
	LoyaltyUser          string      `json:"loyaltyUser"`
	MarketingCampaignId  string      `json:"marketingCampaignId"`
	Nominal              float64     `json:"nominal"`
	OrderNumber          int         `json:"orderNumber"`
	OrderSum             float64     `json:"orderSum"`
	OrganizationId       string      `json:"organizationId"`
	PosBalanceBefore     float64     `json:"posBalanceBefore"`
	ProgramId            string      `json:"programId"`
	Sum                  float64     `json:"sum"`
	Type                 int         `json:"type"`
	TypeName             string      `json:"typeName"`
	WalletId             string      `json:"walletId"`
	WhenCreated          string      `json:"whenCreated"`
	WhenCreatedOrder     string      `json:"whenCreatedOrder"`
	Id                   string      `json:"id"`
	IsDelivery           bool        `json:"isDelivery"`
	IsIgnored            bool        `json:"isIgnored"`
	PosOrderId           string      `json:"posOrderId"`
	Revision             int         `json:"revision"`
	TerminalGroupId      string      `json:"terminalGroupId"`
}

type Certificate struct {
	Number     string `json:"number"`
	Series     string `json:"series"`
	StatusName string `json:"statusName"`
	TypeName   string `json:"typeName"`
}

type Coupon struct {
	Number string `json:"number"`
	Series string `json:"series"`
}

type CustomerCard struct {
	Id          string `json:"id"`
	Track       string `json:"track"`
	Number      string `json:"number"`
	ValidToDate string `json:"validToDate"`
}

type CustomerCategory struct {
	Id                    string `json:"id"`
	Name                  string `json:"name"`
	IsActive              bool   `json:"isActive"`
	IsDefaultForNewGuests bool   `json:"isDefaultForNewGuests"`
}

type WalletBalance struct {
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	Type    int     `json:"type"`
	Balance float64 `json:"balance"`
}
