package dto

import menuCoreModel "github.com/kwaaka-team/orders-core/pkg/menu/dto"

type MenuUploadRequest struct {
	MenuURL string `json:"menu_url" binding:"required" example:"https://kwaaka-menu-files.s3.eu-west-1.amazonaws.com/unit_tests_menu/unit_test_menu_5.json"`
}

type MenuUploadResponse struct {
	TransactionID string `json:"transaction_id" example:"64fef088b0ea9e7d5d401637"`
}

type MenuUploadStatusResponse struct {
	TransactionID string   `json:"transaction_id" example:"64fef088b0ea9e7d5d401637"`
	Status        string   `json:"status" example:"SUCCESS"`
	Details       []string `json:"details"`
}

func FromMenuUploadTransaction(req menuCoreModel.MenuUploadTransaction) MenuUploadStatusResponse {
	details := req.Details
	for _, extTr := range req.ExtTransactions {
		details = append(details, extTr.Details...)
	}

	return MenuUploadStatusResponse{
		TransactionID: req.ID,
		Status:        req.Status,
		Details:       details,
	}
}
