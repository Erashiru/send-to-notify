package selector

type VirtualStore struct {
	ExternalStoreID   string
	RestaurantID      string
	VirtualStoreType  string
	Name              string
	ChildRestaurantID string
	DeliveryService   string
	ClientSecret      string
	Sorting
	Pagination
}

func EmptyVirtualStoreSearch() VirtualStore {
	return VirtualStore{}
}

func VirtualStoreSearch() VirtualStore {
	return VirtualStore{
		Pagination: Pagination{
			Limit: DefaultLimit,
		},
	}
}

func (vs VirtualStore) SetName(name string) VirtualStore {
	vs.Name = name
	return vs
}

func (vs VirtualStore) HasName() bool {
	return vs.Name != ""
}

func (vs VirtualStore) SetChildRestaurantID(id string) VirtualStore {
	vs.ChildRestaurantID = id
	return vs
}

func (vs VirtualStore) HasChildRestaurantID() bool {
	return vs.ChildRestaurantID != ""
}

func (vs VirtualStore) SetExternalStoreID(id string) VirtualStore {
	vs.ExternalStoreID = id
	return vs
}

func (vs VirtualStore) HasExternalStoreID() bool {
	return vs.ExternalStoreID != ""
}

func (vs VirtualStore) SetRestaurantID(id string) VirtualStore {
	vs.RestaurantID = id
	return vs
}

func (vs VirtualStore) HasRestaurantID() bool {
	return vs.RestaurantID != ""
}

func (vs VirtualStore) SetVirtualStoreType(storeType string) VirtualStore {
	vs.VirtualStoreType = storeType
	return vs
}

func (vs VirtualStore) HasVirtualStoreType() bool {
	return vs.VirtualStoreType != ""
}

func (vs VirtualStore) SetDeliveryService(delivery string) VirtualStore {
	vs.DeliveryService = delivery
	return vs
}

func (vs VirtualStore) HasDeliveryService() bool {
	return vs.DeliveryService != ""
}

func (vs VirtualStore) SetClientSecret(secret string) VirtualStore {
	vs.ClientSecret = secret
	return vs
}

func (vs VirtualStore) HasClientSecret() bool {
	return vs.ClientSecret != ""
}
