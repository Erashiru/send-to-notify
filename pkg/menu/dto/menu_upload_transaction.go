package dto

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
)

type MenuUploadTransaction struct {
	ID               string           `bson:"_id,omitempty"`
	StoreID          string           `bson:"restaurant_id"`
	ExtTransactions  []ExtTransaction `bson:"ext_transactions"`
	Status           string           `bson:"status,omitempty"`
	Service          string           `bson:"service"`
	MenuURL          string           `bson:"menu_url"`
	Details          []string         `bson:"details,omitempty"`
	UserName         string           `bson:"user_name"`
	MenuDBVersionUrl string           `bson:"menu_db_version_url"`

	CreatedAt coreModels.TransactionTime `bson:"created_at"`
	UpdatedAt coreModels.TransactionTime `bson:"updated_at"`
}

type ExtTransaction struct {
	ID         string   `bson:"id"  validate:"omitempty"`
	ExtStoreID string   `bson:"store_id"  validate:"omitempty"`
	MenuID     string   `bson:"menu_id"`
	Status     string   `bson:"status"  validate:"omitempty"`
	Details    []string `bson:"details"  validate:"omitempty"`
	MenuUrl    string   `bson:"menu_url" validate:"omitempty"`
}

func ToMenuUploadTransaction(req models.MenuUploadTransaction) MenuUploadTransaction {
	extTransactions := make([]ExtTransaction, 0, len(req.ExtTransactions))

	for _, extTransaction := range req.ExtTransactions {
		extTransactions = append(extTransactions, ExtTransaction{
			ID:         extTransaction.ID,
			ExtStoreID: extTransaction.ExtStoreID,
			MenuID:     extTransaction.MenuID,
			Status:     extTransaction.Status,
			Details:    extTransaction.Details,
			MenuUrl:    extTransaction.MenuUrl,
		})
	}

	return MenuUploadTransaction{
		ID:               req.ID,
		StoreID:          req.StoreID,
		ExtTransactions:  extTransactions,
		Status:           req.Status,
		Service:          req.Service,
		MenuURL:          req.MenuURL,
		Details:          req.Details,
		UserName:         req.UserName,
		MenuDBVersionUrl: req.MenuDBVersionUrl,
		CreatedAt: coreModels.TransactionTime{
			Value:     req.CreatedAt.Value,
			TimeZone:  req.CreatedAt.TimeZone,
			UTCOffset: req.CreatedAt.UTCOffset,
		},
		UpdatedAt: coreModels.TransactionTime{
			Value:     req.UpdatedAt.Value,
			TimeZone:  req.UpdatedAt.TimeZone,
			UTCOffset: req.UpdatedAt.UTCOffset,
		},
	}
}

func FromMenuUploadTransaction(req MenuUploadTransaction) models.MenuUploadTransaction {
	extTransactions := make([]models.ExtTransaction, 0, len(req.ExtTransactions))

	for _, extTransaction := range req.ExtTransactions {
		extTransactions = append(extTransactions, models.ExtTransaction{
			ID:         extTransaction.ID,
			ExtStoreID: extTransaction.ExtStoreID,
			MenuID:     extTransaction.MenuID,
			Status:     extTransaction.Status,
			Details:    extTransaction.Details,
		})
	}

	return models.MenuUploadTransaction{
		ID:               req.ID,
		StoreID:          req.StoreID,
		ExtTransactions:  extTransactions,
		Status:           req.Status,
		Service:          req.Service,
		MenuURL:          req.MenuURL,
		Details:          req.Details,
		UserName:         req.UserName,
		MenuDBVersionUrl: req.MenuDBVersionUrl,
		CreatedAt: coreModels.TransactionTime{
			Value:     req.CreatedAt.Value,
			TimeZone:  req.CreatedAt.TimeZone,
			UTCOffset: req.CreatedAt.UTCOffset,
		},
		UpdatedAt: coreModels.TransactionTime{
			Value:     req.UpdatedAt.Value,
			TimeZone:  req.UpdatedAt.TimeZone,
			UTCOffset: req.UpdatedAt.UTCOffset,
		},
	}
}

func ToMenuUploadTransactions(req []models.MenuUploadTransaction) []MenuUploadTransaction {
	transactions := make([]MenuUploadTransaction, 0, len(req))

	for _, tx := range req {
		transactions = append(transactions, ToMenuUploadTransaction(tx))
	}

	return transactions
}
