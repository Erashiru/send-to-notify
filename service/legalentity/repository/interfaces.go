package repository

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	"github.com/kwaaka-team/orders-core/service/legalentity/models"
	leSel "github.com/kwaaka-team/orders-core/service/legalentity/models/selector"
)

type LegalEntityRepository interface {
	Insert(ctx context.Context, req models.LegalEntityForm) (string, error)
	GetByID(ctx context.Context, id string) (models.LegalEntityView, error)
	Update(ctx context.Context, updatedProfile leSel.LegalEntityForm, id string) error
	Disable(ctx context.Context, id string) error
	List(ctx context.Context, pagination selector.Pagination, filter models.Filter) ([]models.GetListOfLegalEntitiesDB, error)
	GetStores(ctx context.Context, pagination selector.Pagination, id string) (models.GetListOfStoresDB, error)
	InsertDocument(ctx context.Context, id string, newDoc models.Document) error
	DisableDocument(ctx context.Context, legalEntityID, documentID string) error
	GetAllDocumentsByLegalEntityID(ctx context.Context, pagination selector.Pagination, filter models.DocumentFilter, legalEntityID string) ([]models.Document, error)
	GetDocumentDownloadLink(ctx context.Context, legalEntityID, documentID string) (string, error)
}
