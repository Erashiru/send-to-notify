package selector

import (
	"github.com/kwaaka-team/orders-core/core/menu/models/pointer"
)

type BkOffers struct {
	ID       string
	IsActive *bool
}

func EmptyBkOffersSearch() BkOffers {
	return BkOffers{}
}

func (m BkOffers) HasID() bool {
	return m.ID != ""
}

func (m BkOffers) SetID(id string) BkOffers {
	m.ID = id
	return m
}

func (m BkOffers) HasIsActive() bool {
	return m.IsActive != nil
}

func (m BkOffers) SetIsActive(isActive bool) BkOffers {
	m.IsActive = pointer.OfBool(isActive)
	return m
}
