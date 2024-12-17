package dto

type StoreSelector struct {
	ID                       string
	ClientSecret             string
	DeliveryService          string
	ExternalStoreID          string
	PosType                  string
	Hash                     string
	PosOrganizationID        string
	Token                    string
	AggregatorMenuID         string
	AggregatorMenuIDs        []string
	IsActiveMenu             *bool
	IDs                      []string
	Express24StoreId         []string
	PosterAccountNumber      string
	HasVirtualStore          *bool
	GroupID                  string
	TalabatRemoteBranchId    string
	HasScheduledStatusChange bool
	YarosStoreId             string
	City                     string
	IsChildStore             *bool
	OrderAutoClose           *bool
}
