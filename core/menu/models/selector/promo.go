package selector

type Promo struct {
	ID              string
	StoreID         string
	MenuID          string
	DeliveryService string
	POS             string
	ExternalStoreID string
	IsActive        bool
	ProductIDs      []string
	Sorting
	Pagination
}

func EmptyPromoSearch() Promo {
	return Promo{}
}

func PromoSearch() Promo {
	return Promo{
		Pagination: Pagination{
			Limit: DefaultLimit,
		},
	}
}

func (p Promo) SetIsActive(active bool) Promo {
	p.IsActive = active
	return p
}

func (p Promo) HasIsActive() bool {
	return p.IsActive
}

func (p Promo) SetID(id string) Promo {
	p.ID = id
	return p
}

func (p Promo) HasID() bool {
	return p.ID != ""
}

func (p Promo) SetStoreID(id string) Promo {
	p.StoreID = id
	return p
}

func (p Promo) HasStoreID() bool {
	return p.StoreID != ""
}

func (p Promo) SetMenuID(id string) Promo {
	p.MenuID = id
	return p
}

func (p Promo) HasMenuID() bool {
	return p.MenuID != ""
}

func (p Promo) SetDeliveryService(service string) Promo {
	p.DeliveryService = service
	return p
}

func (p Promo) HasDeliveryService() bool {
	return p.DeliveryService != ""
}

func (p Promo) SetPOS(pos string) Promo {
	p.POS = pos
	return p
}

func (p Promo) HasPOS() bool {
	return p.POS != ""
}

func (p Promo) SetExternalStoreID(id string) Promo {
	p.ExternalStoreID = id
	return p
}

func (p Promo) HasExternalStoreID() bool {
	return p.ExternalStoreID != ""
}

func (p Promo) SetProductIDs(ids []string) Promo {
	p.ProductIDs = ids
	return p
}

func (p Promo) HasProductIDs() bool {
	return len(p.ProductIDs) > 0
}
