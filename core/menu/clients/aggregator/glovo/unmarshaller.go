package glovo

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/pkg/glovo/clients/dto"
)

func productModifierFromClient(req dto.ProductModifyResponse) models.ProductModifyResponse {
	return models.ProductModifyResponse{
		ExtID:       req.ID,
		ExtName:     req.Name,
		IsAvailable: req.IsAvailable,
		Price:       req.Price,
	}
}

func uploadMenuToTransaction(req dto.UploadMenuResponse, menuId, storeId, menuUrl string) models.ExtTransaction {
	return models.ExtTransaction{
		ID:         req.TransactionID,
		ExtStoreID: storeId,
		Status:     req.Status.String(),
		Details:    req.Details,
		MenuID:     menuId,
		MenuUrl:    menuUrl,
	}
}
