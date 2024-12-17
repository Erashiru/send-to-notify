package legalentity

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	"github.com/kwaaka-team/orders-core/service/legalentity/models"
	leSel "github.com/kwaaka-team/orders-core/service/legalentity/models/selector"
)

type LegalEntityService interface {
	Create(ctx context.Context, profile models.LegalEntityForm) (string, error)
	Get(ctx context.Context, id string) (models.LegalEntityView, error)
	List(ctx context.Context, pagination selector.Pagination, filter models.Filter) ([]models.GetListOfLegalEntities, error)
	Update(ctx context.Context, updatedProfile leSel.LegalEntityForm, id string) error
	Disable(ctx context.Context, id string) error
	GetStores(ctx context.Context, pagination selector.Pagination, id string) (models.GetListOfStores, error)
	UploadDocument(ctx context.Context, request models.UploadDocumentRequest) (string, error)
	DisableDocument(ctx context.Context, legalEntityID, documentID string) error
	GetAllDocumentsByLegalEntityID(ctx context.Context, pagination selector.Pagination, filter models.DocumentFilter, legalEntityID string) ([]models.Document, error)
	GetDocumentDownloadLink(ctx context.Context, legalEntityID, documentID string) (string, error)
	GenerateContract(contract models.ContractRequest) (models.ContractResponse, error)
}
