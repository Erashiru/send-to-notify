package models

type Store struct {
	ID               string
	Name             string
	PosType          string
	DeliveryServices []string
}
type ManageAggregatorStoreRequest struct {
	DeliveryService       string
	IsOpen                bool
	PosIntegrationStoreID string
}
