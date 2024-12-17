package clients

import (
	"context"
	"github.com/kwaaka-team/orders-core/pkg/express24_v2/clients/dto"
)

type Config struct {
	Protocol string
	BaseURL  string
	Token    string
}
type Express24V2 interface {
	GetCategories(ctx context.Context) ([]dto.Category, error)
	CreateCategory(ctx context.Context, req dto.CreateCategoryRequest) (dto.CreateCategoryResponse, error)
	GetCategoryProducts(ctx context.Context, categoryID string) ([]dto.GetCategoryProductsResponse, error)

	GetSubCategories(ctx context.Context, categoryID string) ([]dto.Category, error)
	CreateSubCategory(ctx context.Context, req dto.CreateCategoryRequest) (dto.CreateCategoryResponse, error)
	GetSubCategoryProducts(ctx context.Context, subCategoryID string) ([]dto.GetCategoryProductsResponse, error)

	CreateProduct(ctx context.Context, req dto.CreateProductRequest) (dto.CreateProductResponse, error)
	UpdateProduct(ctx context.Context, req dto.UpdateProductRequest) (dto.UpdateProductResponse, error)

	CreateProductsAttributeGroup(ctx context.Context, req dto.CreateProductsAttributeGroupRequest) (dto.AttributeGroup, error)
	GetProductsAttributeGroups(ctx context.Context, productID string) (dto.GetAttributeGroupResponse, error)

	CreateAttributeGroup(ctx context.Context, req dto.CreateProductsAttributeGroupRequest) (dto.AttributeGroup, error)
	GetAttributeGroups(ctx context.Context) ([]dto.AttributeGroup, error)

	CreateAttributeGroupsItem(ctx context.Context, req dto.CreateAttributeGroupsItemRequest) (dto.AttributeGroupItem, error)
	GetAttributeGroupsItems(ctx context.Context, attributeGroupID string) ([]dto.GetAttributeGroupsItemResponse, error)

	StopListByAttributes(ctx context.Context, req dto.StopListByAttributesRequest) error
	StopListByProducts(ctx context.Context, req dto.StopListByProductsRequest) error
	StopListBulk(ctx context.Context, req dto.StopListBulkRequest) error

	SyncMenu(ctx context.Context, req dto.MenuSyncReq) error
}
