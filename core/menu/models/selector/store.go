package selector

type Store struct {
	ID                string
	Name              string
	PosMenuID         string
	AggregatorMenuIDs []string
	AggregatorMenuID  string

	Token string

	ExternalStoreID string
	DeliveryService string

	IsActiveMenu *bool

	Sorting
	Pagination
}

func EmptyStoreSearch() Store {
	return Store{}
}

func StoreSearch() Store {
	return Store{
		Pagination: Pagination{
			Limit: DefaultLimit,
		},
	}
}

func (s Store) SetID(id string) Store {
	s.ID = id
	return s
}

func (s Store) HasID() bool {
	return s.ID != ""
}

func (s Store) HasToken() bool {
	return s.Token != ""
}

func (s Store) HasPosMenuID() bool {
	return s.PosMenuID != ""
}

func (s Store) HasAggregatorMenuID() bool {

	return s.AggregatorMenuID != ""
}

func (s Store) HasAggregatorMenuIDs() bool {
	return len(s.AggregatorMenuIDs) > 0
}

func (s Store) HasExternalStoreID() bool {
	return s.ExternalStoreID != ""
}

func (s Store) HasDeliveryService() bool {
	return s.DeliveryService != ""
}

func (s Store) ActiveMenu() bool {
	if s.IsActiveMenu != nil && *s.IsActiveMenu {
		return true
	}
	return false
}

func (s Store) HasIsActiveMenu() bool {
	return s.IsActiveMenu != nil
}

func (s Store) SetPosMenuID(menuID string) Store {
	s.PosMenuID = menuID
	return s
}

func (s Store) SetToken(token string) Store {
	s.Token = token
	return s
}

func (s Store) SetAggregatorMenuID(menuID string) Store {
	s.AggregatorMenuID = menuID
	return s
}

func (s Store) SetAggregatorMenuIDs(menuIDs []string) Store {
	s.AggregatorMenuIDs = menuIDs
	return s
}

func (s Store) SetDeliveryService(deliveryService string) Store {
	s.DeliveryService = deliveryService
	return s
}

func (s Store) SetExternalStoreID(storeID string) Store {
	s.ExternalStoreID = storeID
	return s
}

func (s Store) SetIsActiveMenu(isActive *bool) Store {
	s.IsActiveMenu = isActive
	return s
}
