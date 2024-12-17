package models

import coreModels "github.com/kwaaka-team/orders-core/core/models"

type MenuUploadTransaction struct {
	ID               string          `bson:"_id,omitempty" json:"id"`
	StoreID          string          `bson:"restaurant_id" json:"store_id"`
	ExtTransactions  ExtTransactions `bson:"ext_transactions" json:"ext_transactions"`
	Status           string          `bson:"status,omitempty" json:"status"`
	MenuURL          string          `bson:"menu_url" json:"menu_url"`
	Service          string          `bson:"service" json:"service"`
	Details          []string        `bson:"details,omitempty" json:"details"`
	UserName         string          `bson:"user_name" json:"user_name"`
	MenuDBVersionUrl string          `bson:"menu_db_version_url" json:"menu_db_version_url"`

	CreatedAt coreModels.TransactionTime `bson:"created_at" json:"created_at"`
	UpdatedAt coreModels.TransactionTime `bson:"updated_at" json:"updated_at"`
}

type UpdateMenuUploadTransaction struct {
	RestaurantID     *string         `bson:"restaurant_id" validate:"omitempty" json:"restaurant_id"`
	MenuID           *string         `bson:"menu_id" json:"menu_id"`
	ExtTransactions  ExtTransactions `bson:"ext_transactions" validate:"omitempty" json:"ext_transactions"`
	Status           *string         `bson:"status" validate:"omitempty" json:"status"`
	Service          *string         `bson:"service" validate:"omitempty" json:"service"`
	Details          []string        `bson:"details" validate:"omitempty" json:"details"`
	UserName         *string         `bson:"user_name" validate:"omitempty" json:"user_name"`
	MenuDBVersionUrl *string         `bson:"menu_db_version_url" validate:"omitempty" json:"menu_db_version_url"`

	CreatedAt coreModels.TransactionTime `bson:"created_at" validate:"omitempty" json:"created_at"`
	UpdatedAt coreModels.TransactionTime `bson:"updated_at" validate:"omitempty" json:"updated_at"`
}

type ExtTransactions []ExtTransaction

type ExtTransaction struct {
	ID         string   `bson:"id"  validate:"omitempty" json:"id"`
	ExtStoreID string   `bson:"store_id"  validate:"omitempty" json:"ext_store_id"`
	MenuID     string   `bson:"menu_id" json:"menu_id"`
	Status     string   `bson:"status"  validate:"omitempty" json:"status"`
	Details    []string `bson:"details"  validate:"omitempty" json:"details"`
	MenuUrl    string   `bson:"menu_url" validate:"omitempty" json:"menu_url"`
}

func (trx ExtTransactions) GetByMenu(menuId string) (ExtTransaction, error) {
	for _, tr := range trx {
		if tr.MenuID == menuId {
			return tr, nil
		}
	}
	return ExtTransaction{}, ErrNotFound
}

func (trx ExtTransactions) HasProcessingStatus() bool {

	for _, tr := range trx {
		if tr.Status == PROCESSING.String() {
			return true
		}
	}

	return false
}

func ToUpdateMenuTransactions(req MenuUploadTransaction) UpdateMenuUploadTransaction {

	return UpdateMenuUploadTransaction{
		RestaurantID:     &req.StoreID,
		ExtTransactions:  req.ExtTransactions,
		Status:           &req.Status,
		Service:          &req.Service,
		Details:          req.Details,
		CreatedAt:        req.CreatedAt,
		UpdatedAt:        req.UpdatedAt,
		UserName:         &req.UserName,
		MenuDBVersionUrl: &req.MenuDBVersionUrl,
	}
}
func (trx ExtTransactions) HasNotSuccessProcessingStatus() bool {

	for _, tr := range trx {
		if tr.Status != SUCCESS.String() && tr.Status != PROCESSING.String() {
			return true
		}
	}

	return false
}
