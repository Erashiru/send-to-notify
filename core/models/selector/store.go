package selector

type Store struct {
	ID string

	ExternalStoreID string
	DeliveryService string
}

func EmptyStore() Store {
	return Store{}
}

func (s Store) HasID() bool {
	return s.ID != ""
}

func (s Store) HasExternalStoreID() bool {
	return s.ExternalStoreID != ""
}

func (s Store) SetID(storeID string) Store {
	s.ID = storeID
	return s
}

func (s Store) HasDeliveryService() bool {
	return s.DeliveryService != ""
}

func (s Store) SetDeliveryService(deliveryService string) Store {
	s.DeliveryService = deliveryService
	return s
}

func (s Store) SetExternalStoreID(storeID string) Store {
	s.ExternalStoreID = storeID
	return s
}
