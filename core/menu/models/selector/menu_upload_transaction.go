package selector

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"time"
)

type MenuUploadTransaction struct {
	ID               string
	StoreID          string
	MenuID           string
	Status           string
	ExtTransactionID string
	Service          models.AggregatorName

	Sorting
	Pagination

	CreatedTo   time.Time
	CreatedFrom time.Time
}

func EmptyMenuUploadTransactionSearch() MenuUploadTransaction {
	return MenuUploadTransaction{}
}

func MenuUploadTransactionSearch() MenuUploadTransaction {
	return MenuUploadTransaction{
		Pagination: Pagination{
			Limit: DefaultLimit,
		},
	}
}

func (m MenuUploadTransaction) HasID() bool {
	return m.ID != ""
}

func (m MenuUploadTransaction) HasMenuID() bool {
	return m.MenuID != ""
}

func (m MenuUploadTransaction) HasStoreID() bool {
	return m.StoreID != ""
}

func (m MenuUploadTransaction) HasExtTransactionID() bool {
	return m.ExtTransactionID != ""
}

func (m MenuUploadTransaction) HasStatus() bool {
	return m.Status != ""
}

func (m MenuUploadTransaction) HasCreatedTo() bool {
	return m.CreatedTo != time.Time{}
}

func (m MenuUploadTransaction) HasCreatedFrom() bool {
	return m.CreatedFrom != time.Time{}
}

func (m MenuUploadTransaction) HasService() bool {
	return m.Service != ""
}

func (m MenuUploadTransaction) SetID(id string) MenuUploadTransaction {
	m.ID = id
	return m
}

func (m MenuUploadTransaction) SetCreatedTo(date time.Time) MenuUploadTransaction {
	m.CreatedTo = date
	return m
}

func (m MenuUploadTransaction) SetCreatedFrom(date time.Time) MenuUploadTransaction {
	m.CreatedFrom = date
	return m
}

func (m MenuUploadTransaction) SetStoreID(id string) MenuUploadTransaction {
	m.StoreID = id
	return m
}

func (m MenuUploadTransaction) SetStatus(status string) MenuUploadTransaction {
	m.Status = status
	return m
}

func (m MenuUploadTransaction) SetMenuID(id string) MenuUploadTransaction {
	m.MenuID = id
	return m
}

func (m MenuUploadTransaction) SetService(service models.AggregatorName) MenuUploadTransaction {
	m.Service = service
	return m
}

func (m MenuUploadTransaction) SetPage(page int64) MenuUploadTransaction {
	if page > 0 {
		m.Pagination.Page = page - 1
	}
	return m
}

func (m MenuUploadTransaction) SetLimit(limit int64) MenuUploadTransaction {
	if limit > 0 {
		m.Pagination.Limit = limit
	}
	return m
}

func (m MenuUploadTransaction) SetSorting(key string, dir int8) MenuUploadTransaction {
	m.Sorting.Param = key
	m.Sorting.Direction = dir
	return m
}
