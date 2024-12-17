package models

type OrderTypesRequest struct {
	Organizations []string `json:"organizationIds"`
}

type OrderTypesResponse struct {
	OrganizationOrderTypes []OrganizationOrderType `json:"orderTypes"`
}

type OrganizationOrderType struct {
	OrganizationID string      `json:"organizationId"`
	OrderTypes     []OrderType `json:"items"`
}

type OrderType struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	OrderServiceType string `json:"orderServiceType"`
	IsDeleted        bool   `json:"isDeleted"`
	ExternalRevision int64  `json:"externalRevision"`
}
