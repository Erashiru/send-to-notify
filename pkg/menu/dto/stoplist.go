package dto

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
)

type StopListItem struct {
	ID          string
	Price       float64
	IsAvailable bool
}

type ProductStopList struct {
	ProductID string
	SetToStop bool
	Data      StoreProductStopLists
}

type StoreProductStopLists []StoreProductStopList

type StoreProductStopList struct {
	ID          string
	Aggregators models.Aggregators
}

func (s StoreProductStopLists) ToModels() []models.UpdateStoreData {
	res := make([]models.UpdateStoreData, 0, len(s))
	for _, store := range s {
		res = append(res, models.UpdateStoreData{
			ID:          store.ID,
			Aggregators: store.Aggregators,
		})
	}

	return res
}

func FromStopListTransactions(req []models.StopListTransaction) StoreProductStopLists {
	res := make(StoreProductStopLists, 0, len(req))
	for _, tr := range req {
		res = append(res, StoreProductStopList{
			ID:          tr.StoreID,
			Aggregators: fromTransactionsToAggregators(tr.Transactions),
		})
	}
	return res
}

func fromTransactionsToAggregators(req []models.TransactionData) models.Aggregators {
	res := make(models.Aggregators, 0, len(req))
	for _, aggr := range req {
		res = append(res, models.Aggregator{
			Name:    models.AggregatorName(aggr.Delivery),
			Success: aggr.Status == models.SUCCESS,
			Msg:     aggr.Message,
		})
	}
	return res
}

type StopList struct {
	Type          int    `json:"type,omitempty"`
	ElementID     int    `json:"element_id,omitempty"`
	StorageID     int    `json:"storage_id,omitempty"`
	ValueRelative int    `json:"value_relative,omitempty"`
	ValueAbsolute int    `json:"value_absolute,omitempty"`
	ProductID     string `json:"product_id,omitempty"`
}

func (s StopListItem) ToModel() models.ItemStopList {
	return models.ItemStopList{
		ID:          s.ID,
		Price:       s.Price,
		IsAvailable: s.IsAvailable,
	}
}
