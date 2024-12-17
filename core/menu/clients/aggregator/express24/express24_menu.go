// When I wrote this, only God and I understood what I was doing
// Now, God only knows

package express24

import (
	"context"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kwaaka-team/orders-core/config/menu"
	errors2 "github.com/kwaaka-team/orders-core/core/errors"
	constErrors "github.com/kwaaka-team/orders-core/core/menu/clients/aggregator/errors"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	selector2 "github.com/kwaaka-team/orders-core/core/storecore/managers/selector"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	expressConf "github.com/kwaaka-team/orders-core/pkg/express24_v2"
	expressMenuCli "github.com/kwaaka-team/orders-core/pkg/express24_v2/clients"
	"github.com/kwaaka-team/orders-core/pkg/express24_v2/clients/dto"
	"github.com/kwaaka-team/orders-core/pkg/store"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"slices"
	"strconv"
	"time"
)

type Express24Menu interface {
	BulkUpdate(ctx context.Context, restaurantID, storeID string, products models.Products, attributes models.Attributes, store storeModels.Store) (string, error)
	ModifyProduct(ctx context.Context, storeID string, product models.Product) (models.ProductModifyResponse, error)
	UploadMenu(ctx context.Context, menuId, extStoreId string, menu models.Menu, store storeModels.Store, offers []models.BkOffers, sv3 *s3.S3, userRole string) (models.ExtTransaction, error)
	VerifyMenu(ctx context.Context, transaction models.ExtTransaction) (models.Status, error)
	BulkAttribute(ctx context.Context, storeID string, attributes models.Attributes) (string, error)
	GetMenu(ctx context.Context, extStoreId string) (models.Menu, error)
	ValidateMenu(ctx context.Context, request models.MenuValidateRequest) (models.MenuUploadTransaction, error)
}

type express24MenuImpl struct {
	cli               expressMenuCli.Express24V2
	s3Bucket          menu.S3_BUCKET
	restGroupMenuRepo drivers.RestaurantGroupMenuRepository
	storeCli          store.Client
}

func NewMenuManager(ctx context.Context, cfg menu.Configuration, repoMenu drivers.RestaurantGroupMenuRepository, client store.Client) (Express24Menu, error) {
	cli, err := expressConf.NewExpress24Client(&expressMenuCli.Config{
		Protocol: "http",
		BaseURL:  cfg.Express24Configuration.BaseURL,
		Token:    cfg.Express24Configuration.Token,
	})

	if err != nil {
		log.Trace().Err(err).Msg("can't initialize Express24 menu client ")
		return nil, err
	}
	return &express24MenuImpl{
		cli:               cli,
		s3Bucket:          cfg.S3_BUCKET,
		restGroupMenuRepo: repoMenu,
		storeCli:          client,
	}, nil
}

func (s express24MenuImpl) BulkUpdate(ctx context.Context, restaurantID, storeID string, products models.Products, attributes models.Attributes, store storeModels.Store) (string, error) {
	if len(products) == 0 && len(attributes) == 0 {
		return "", nil
	}

	branchID, err := strconv.Atoi(storeID)
	if err != nil {
		return "", err
	}
	req := dto.StopListBulkRequest{
		BranchIDs: []int{branchID},
	}

	if len(products) != 0 {
		req.Products = s.toProducts(products)
	}

	if len(attributes) != 0 {
		req.Modifiers = s.toAttributes(attributes)
	}

	err = s.cli.StopListBulk(ctx, req)
	if err != nil {
		return "", errors.Wrap(err, "bulk update error")
	}

	return "", nil
}

func (s express24MenuImpl) toProducts(req models.Products) dto.Products {
	items := make([]dto.ProductItem, 0, len(req))

	for i := range req {
		items = append(items, dto.ProductItem{
			ExternalID:  req[i].ExtID,
			Quantity:    models.BASEQUANTITY,
			IsAvailable: req[i].IsAvailable,
		})
	}
	return dto.Products{
		Items:                   items,
		MakeAvailableOtherItems: false,
	}
}

func (s express24MenuImpl) toAttributes(req models.Attributes) dto.Modifiers {
	items := make([]dto.AttributeItem, 0, len(req))

	for i := range req {
		items = append(items, dto.AttributeItem{
			ExternalID:  req[i].ExtID,
			IsAvailable: req[i].IsAvailable,
		})
	}
	return dto.Modifiers{
		Items:                   items,
		MakeAvailableOtherItems: false,
	}
}

func (s express24MenuImpl) ModifyProduct(ctx context.Context, storeID string, product models.Product) (models.ProductModifyResponse, error) {
	return models.ProductModifyResponse{}, constErrors.ErrNotImplemented
}

func (s express24MenuImpl) UploadMenu(ctx context.Context, menuId, extStoreId string, menu models.Menu, store storeModels.Store, offers []models.BkOffers, sv3 *s3.S3, userRole string) (models.ExtTransaction, error) {
	err := s.uploadExpress24Menu(ctx, store, menu)
	if err != nil {
		log.Err(err).Msgf("publication express24 menu error")
		return models.ExtTransaction{
			Status:     models.NOT_PROCESSED.String(),
			MenuID:     menuId,
			ExtStoreID: extStoreId,
			Details:    []string{err.Error()},
		}, err
	}

	return models.ExtTransaction{
		Status:     models.SUCCESS.String(),
		MenuID:     menuId,
		ExtStoreID: extStoreId,
		Details:    []string{},
	}, nil
}

func (s express24MenuImpl) VerifyMenu(ctx context.Context, transaction models.ExtTransaction) (models.Status, error) {
	return "", nil
}

func (s express24MenuImpl) BulkAttribute(ctx context.Context, storeID string, attributes models.Attributes) (string, error) {
	return "", constErrors.ErrNotImplemented
}

func (s express24MenuImpl) GetMenu(ctx context.Context, extStoreId string) (models.Menu, error) {
	return models.Menu{}, constErrors.ErrNotImplemented
}

func (s express24MenuImpl) ValidateMenu(ctx context.Context, request models.MenuValidateRequest) (models.MenuUploadTransaction, error) {
	return models.MenuUploadTransaction{}, nil
}

func (s express24MenuImpl) uploadExpress24Menu(ctx context.Context, store storeModels.Store, menu models.Menu) error {

	restGroup, err := s.storeCli.FindStoreGroup(ctx, selector2.NewEmptyStoreGroupSearch().SetStoreIDs([]string{store.ID}))
	if err != nil {
		return err
	}

	oldMenu, err := s.getOldMenu(ctx, restGroup.ID)
	if err != nil {
		return err
	}

	newMenu, err := s.createSyncRequest(ctx, menu, store.Express24.Vat, oldMenu)
	if err != nil {
		return err
	}

	time.Sleep(10 * time.Second)

	// alter products which are required to us (attached branch)
	if err = s.alterProducts(ctx, menu, store); err != nil {
		return err
	}

	// alter attributes which are required to us (attached branch)
	if err = s.alterAttributes(ctx, menu, store); err != nil {
		return err
	}

	if err = s.createOrUpdateMenu(ctx, restGroup.ID, oldMenu, *newMenu); err != nil {
		return err
	}

	return nil
}

func (s express24MenuImpl) createOrUpdateMenu(ctx context.Context, restGroupId string, oldMenu *dto.MenuSyncReq, newMenu dto.MenuSyncReq) error {
	if oldMenu == nil {
		if err := s.restGroupMenuRepo.UpdateOrCreateMenu(ctx, restGroupId, s.toRestGroupMenu(restGroupId, newMenu)); err != nil {
			return err
		}

		return nil
	}

	generalMenu := s.concatMenu(oldMenu, &newMenu)

	if err := s.restGroupMenuRepo.UpdateOrCreateMenu(ctx, restGroupId, s.toRestGroupMenu(restGroupId, *generalMenu)); err != nil {
		return err
	}

	return nil
}

func (s express24MenuImpl) toRestGroupMenu(restGroupId string, menu dto.MenuSyncReq) models.RestGroupMenu {
	result := models.RestGroupMenu{
		RestGroupId: restGroupId,
	}

	for _, category := range menu.Categories {
		result.Category = append(result.Category, models.Category{
			ExternalID:    category.ExternalID,
			Name:          category.Name,
			IsActive:      category.IsActive,
			Sort:          category.Sort,
			SubCategories: category.SubCategories,
		})
	}

	for _, product := range menu.Products {
		result.MenuProduct = append(result.MenuProduct, models.MenuProduct{
			ExternalID:  product.ExternalID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			CategoryID:  product.CategoryID,
			Modifiers:   product.Modifiers,
			Fiscalization: models.Fiscalization{
				SpicID:      product.Fiscalization.SpicID,
				PackageCode: product.Fiscalization.PackageCode,
			},
			Vat:    product.Vat,
			Images: s.toProductImage(product.Images),
		})
	}

	for _, subcategory := range menu.SubCategories {
		result.SubCategory = append(result.SubCategory, models.SubCategory{
			ExternalID: subcategory.ExternalID,
			Name:       subcategory.Name,
			IsActive:   subcategory.IsActive,
			Sort:       subcategory.Sort,
		})
	}

	for _, modItem := range menu.ModifierItems {
		result.ModifierItems = append(result.ModifierItems, models.ModifierItems{
			Name:       modItem.Name,
			ExternalID: modItem.ExternalID,
			Price:      modItem.Price,
		})
	}

	for _, mod := range menu.Modifiers {
		result.Modifier = append(result.Modifier, models.Modifier{
			Name:       mod.Name,
			ExternalID: mod.ExternalID,
			Items:      mod.Items,
		})
	}

	return result
}

func (s express24MenuImpl) toProductImage(images []dto.Images) []models.Images {
	res := []models.Images{}
	for _, image := range images {
		res = append(res, models.Images{URL: image.URL, IsPreview: image.IsPreview})
	}

	return res
}

func (s express24MenuImpl) concatMenu(menu1, menu2 *dto.MenuSyncReq) *dto.MenuSyncReq {
	resultMenu := dto.MenuSyncReq{
		Categories:    append(menu1.Categories, menu2.Categories...),
		SubCategories: append(menu1.SubCategories, menu2.SubCategories...),
		Products:      append(menu1.Products, menu2.Products...),
		Modifiers:     append(menu1.Modifiers, menu2.Modifiers...),
		ModifierItems: append(menu1.ModifierItems, menu2.ModifierItems...),
	}
	resultMenu.Categories = s.uniqueCategories(resultMenu.Categories)
	resultMenu.SubCategories = s.uniqueSubCategories(resultMenu.SubCategories)
	resultMenu.Products = s.uniqueProducts(resultMenu.Products)
	resultMenu.Modifiers = s.uniqueModifiers(resultMenu.Modifiers)
	resultMenu.ModifierItems = s.uniqueModifierItems(resultMenu.ModifierItems)

	return &resultMenu
}

func (s express24MenuImpl) uniqueCategories(categories []dto.MenuSyncCategoryReq) []dto.MenuSyncCategoryReq {
	seen := make(map[string]bool)
	result := []dto.MenuSyncCategoryReq{}

	for _, category := range categories {
		if _, ok := seen[category.ExternalID]; !ok {
			seen[category.ExternalID] = true
			result = append(result, category)
		}
	}
	return result
}
func (s express24MenuImpl) uniqueSubCategories(subCategories []dto.MenuSyncSubCategoryReq) []dto.MenuSyncSubCategoryReq {
	seen := make(map[string]bool)
	result := []dto.MenuSyncSubCategoryReq{}

	for _, subCategory := range subCategories {
		if _, ok := seen[subCategory.ExternalID]; !ok {
			seen[subCategory.ExternalID] = true
			result = append(result, subCategory)
		}
	}
	return result
}

func (s express24MenuImpl) uniqueModifiers(modifiers []dto.MenuSyncModifierReq) []dto.MenuSyncModifierReq {
	seen := make(map[string]bool)
	result := []dto.MenuSyncModifierReq{}

	for _, modifier := range modifiers {
		if _, ok := seen[modifier.ExternalID]; !ok {
			seen[modifier.ExternalID] = true
			result = append(result, modifier)
		}
	}
	return result
}

func (s express24MenuImpl) uniqueProducts(products []dto.MenuSyncProductReq) []dto.MenuSyncProductReq {
	seen := make(map[string]bool)
	result := []dto.MenuSyncProductReq{}

	for _, product := range products {
		if _, ok := seen[product.ExternalID]; !ok {
			seen[product.ExternalID] = true
			result = append(result, product)
		}
	}
	return result
}

func (s express24MenuImpl) uniqueModifierItems(modifierItems []dto.MenuSyncModifierItemsReq) []dto.MenuSyncModifierItemsReq {
	seen := make(map[string]bool)
	result := []dto.MenuSyncModifierItemsReq{}

	for _, modifierItem := range modifierItems {
		if _, ok := seen[modifierItem.ExternalID]; !ok {
			seen[modifierItem.ExternalID] = true
			result = append(result, modifierItem)
		}
	}
	return result
}

func (s express24MenuImpl) getOldMenu(ctx context.Context, restGroupId string) (*dto.MenuSyncReq, error) {

	oldMenu, err := s.restGroupMenuRepo.GetMenuByRestGroupId(ctx, restGroupId)
	if err != nil {
		if errors.Is(err, errors2.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var res dto.MenuSyncReq

	for _, cat := range oldMenu.Category {
		tmp := dto.MenuSyncCategoryReq{
			ExternalID:    cat.ExternalID,
			Name:          cat.Name,
			IsActive:      cat.IsActive,
			Sort:          cat.Sort,
			SubCategories: cat.SubCategories,
		}

		res.Categories = append(res.Categories, tmp)
	}

	for _, subcat := range oldMenu.SubCategory {
		tmp := dto.MenuSyncSubCategoryReq{
			ExternalID: subcat.ExternalID,
			Name:       subcat.Name,
			IsActive:   subcat.IsActive,
			Sort:       subcat.Sort,
		}

		res.SubCategories = append(res.SubCategories, tmp)
	}

	for _, pr := range oldMenu.MenuProduct {
		tmp := dto.MenuSyncProductReq{
			ExternalID:  pr.ExternalID,
			Name:        pr.Name,
			Description: pr.Description,
			Price:       pr.Price,
			CategoryID:  pr.CategoryID,
			Modifiers:   pr.Modifiers,
			Fiscalization: dto.Fiscalization{
				SpicID:      pr.Fiscalization.SpicID,
				PackageCode: pr.Fiscalization.PackageCode,
			},
			Vat:    pr.Vat,
			Images: s.toDtoImages(pr.Images),
		}

		res.Products = append(res.Products, tmp)
	}

	for _, modItem := range oldMenu.ModifierItems {
		tmp := dto.MenuSyncModifierItemsReq{
			ExternalID: modItem.ExternalID,
			Name:       modItem.Name,
			Price:      modItem.Price,
		}

		res.ModifierItems = append(res.ModifierItems, tmp)
	}

	for _, mod := range oldMenu.Modifier {
		tmp := dto.MenuSyncModifierReq{
			ExternalID: mod.ExternalID,
			Name:       mod.Name,
			Items:      mod.Items,
		}

		res.Modifiers = append(res.Modifiers, tmp)
	}

	return &res, nil
}

func (s express24MenuImpl) toDtoImages(images []models.Images) []dto.Images {
	res := []dto.Images{}
	for _, image := range images {
		res = append(res, dto.Images{URL: image.URL, IsPreview: image.IsPreview})
	}

	return res
}

func (s express24MenuImpl) alterAttributes(ctx context.Context, menu models.Menu, store storeModels.Store) error {
	if true { // right now we are not required in logic to change branches for attributes, however we must investigate this moment
		return nil
	}

	modsList, err := s.cli.GetAttributeGroups(ctx)
	if err != nil {
		return err
	}

	mpExtIdAttr := map[string]dto.GetAttributeGroupsItemResponse{}
	for _, attGroup := range modsList {
		attrs, err := s.cli.GetAttributeGroupsItems(ctx, strconv.Itoa(attGroup.ID))
		if err != nil {
			return err
		}
		for _, attr := range attrs {
			mpExtIdAttr[attr.ExternalID] = attr
		}
	}

	for _, attr := range menu.Attributes {
		exprAttr, ok := mpExtIdAttr[attr.ExtID]
		if !ok {
			return errors.New("could not find product")
		}

		if err := s.alterBranchAvailabilityAttr(ctx, attr, exprAttr, store.Express24.StoreID[0]); err != nil {
			return err
		}
	}

	return nil
}

func (s express24MenuImpl) alterBranchAvailabilityAttr(ctx context.Context, attribute models.Attribute, exprAttr dto.GetAttributeGroupsItemResponse, currStoreId string) error {

	stId, err := strconv.Atoi(currStoreId)
	if err != nil {
		return err
	}

	newAttBranch := make([]dto.AttachedBranch, 0, len(exprAttr.Branches))

	for _, branch := range exprAttr.Branches {
		if branch.ID == stId && branch.IsAvailable {
			newAttBranch = append(newAttBranch, dto.AttachedBranch{
				ID:          branch.ID,
				ExternalID:  branch.ExternalID,
				IsActive:    attribute.IsAvailable,
				IsAvailable: true,
				Qty:         branch.Qty,
			})
		}
		newAttBranch = append(newAttBranch, branch)
	}

	// TODO: add method to update attribute

	return nil
}

func (s express24MenuImpl) alterProducts(ctx context.Context, menu models.Menu, store storeModels.Store) error {

	categories, err := s.cli.GetCategories(ctx)
	if err != nil {
		return err
	}

	mpExtIdProduct := map[string]dto.GetCategoryProductsResponse{}
	for _, category := range categories {
		products, err := s.cli.GetCategoryProducts(ctx, strconv.Itoa(category.ID))
		if err != nil {
			return err
		}

		for _, product := range products {
			mpExtIdProduct[product.ExternalID] = product
		}
	}

	for _, product := range menu.Products {
		exprProduct, ok := mpExtIdProduct[product.ExtID]
		if !ok {
			return errors.New("could not find product")
		}

		if err := s.alterBranchAvailability(ctx, product, exprProduct, store.Express24.StoreID[0]); err != nil {
			return err
		}
	}

	return nil
}

func (s express24MenuImpl) alterBranchAvailability(ctx context.Context, product models.Product, exprProduct dto.GetCategoryProductsResponse, currStoreId string) error {

	stId, err := strconv.Atoi(currStoreId)
	if err != nil {
		return err
	}

	newAttBranch := make([]dto.AttachedBranch, 0, len(exprProduct.AttachedBranches))

	for _, branch := range exprProduct.AttachedBranches {
		if branch.ID == stId {
			newAttBranch = append(newAttBranch, dto.AttachedBranch{
				ID:          branch.ID,
				ExternalID:  branch.ExternalID,
				IsActive:    product.IsAvailable,
				IsAvailable: true,
				Qty:         branch.Qty,
			})
		} else {
			newAttBranch = append(newAttBranch, branch)
		}
	}

	mods := make([]int, 0, len(exprProduct.AttributeGroups))
	for _, mod := range exprProduct.AttributeGroups {
		mods = append(mods, mod.ID)
	}

	if !slices.Equal(newAttBranch, exprProduct.AttachedBranches) {
		if _, err := s.cli.UpdateProduct(ctx, dto.UpdateProductRequest{
			ID:               strconv.Itoa(exprProduct.ID),
			Name:             &exprProduct.Name,
			ExternalID:       &exprProduct.ExternalID,
			Description:      &exprProduct.Description,
			Price:            &exprProduct.Price,
			CategoryID:       &exprProduct.CategoryID,
			Fiscalization:    &exprProduct.Fiscalization,
			Vat:              &exprProduct.Vat,
			AttachedBranches: &newAttBranch,
			Images:           &exprProduct.Images,
			AttributeGroups:  &mods,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (s express24MenuImpl) createSyncRequest(ctx context.Context, menu models.Menu, vatValue int, oldMenu *dto.MenuSyncReq) (*dto.MenuSyncReq, error) {
	var (
		categoriesReq    []dto.MenuSyncCategoryReq
		subCategoriesReq []dto.MenuSyncSubCategoryReq
		modifierItemsReq []dto.MenuSyncModifierItemsReq
		modifiersReq     []dto.MenuSyncModifierReq
		productsReq      []dto.MenuSyncProductReq
	)

	for _, section := range menu.Sections {
		if section.IsDeleted {
			continue
		}

		categoriesReq = append(categoriesReq, s.toCategories(section))
	}

	for _, attGroup := range menu.AttributesGroups {
		modifiersReq = append(modifiersReq, s.toModifier(attGroup))
	}

	for _, att := range menu.Attributes {
		if att.IsDeleted {
			continue
		}

		modifierItemsReq = append(modifierItemsReq, s.toModifierItem(att))
	}

	for _, product := range menu.Products {
		if product.IsDeleted {
			continue
		}

		productsReq = append(productsReq, s.toSyncProducts(product, vatValue))
	}

	newMenu := dto.MenuSyncReq{
		Categories:    categoriesReq,
		SubCategories: subCategoriesReq,
		Products:      productsReq,
		Modifiers:     modifiersReq,
		ModifierItems: modifierItemsReq,
	}

	generalMenu := dto.MenuSyncReq{}
	if oldMenu != nil {
		concatMenu := s.concatMenu(oldMenu, &newMenu)
		generalMenu = *concatMenu
	} else {
		generalMenu = newMenu
	}

	if err := s.cli.SyncMenu(ctx, generalMenu); err != nil {
		return nil, err
	}

	return &newMenu, nil
}

func (s express24MenuImpl) toModifierItem(att models.Attribute) dto.MenuSyncModifierItemsReq {
	return dto.MenuSyncModifierItemsReq{
		Name:       att.Name,
		Price:      int(att.Price),
		ExternalID: att.ExtID,
	}
}

func (s express24MenuImpl) toSyncProducts(product models.Product, vatValue int) dto.MenuSyncProductReq {
	pkgcode, _ := strconv.Atoi(product.PackageCode)

	return dto.MenuSyncProductReq{
		ExternalID:  product.ExtID,
		Name:        product.Name[0].Value,
		Description: product.Description[0].Value,
		Price:       int(product.Price[0].Value),
		CategoryID:  product.Section,
		Modifiers:   product.AttributesGroups,
		Fiscalization: dto.Fiscalization{
			SpicID:      product.SpicID,
			PackageCode: pkgcode,
		},
		Vat:    vatValue,
		Images: s.toImageUrls(product.ImageURLs),
	}
}

func (s express24MenuImpl) toCategories(section models.Section) dto.MenuSyncCategoryReq {
	return dto.MenuSyncCategoryReq{
		ExternalID:    section.ExtID,
		Name:          section.Name,
		IsActive:      true,
		Sort:          section.SectionOrder,
		SubCategories: make([]string, 0),
	}
}

func (s express24MenuImpl) toModifier(attGrop models.AttributeGroup) dto.MenuSyncModifierReq {
	return dto.MenuSyncModifierReq{
		Name:       attGrop.Name,
		ExternalID: attGrop.ExtID,
		Items:      attGrop.Attributes,
	}
}

func (s express24MenuImpl) toImageUrls(urls []string) []dto.Images {
	if len(urls) == 0 {
		return nil
	}

	return []dto.Images{
		{
			URL:       urls[0],
			IsPreview: true,
		},
	}
}
