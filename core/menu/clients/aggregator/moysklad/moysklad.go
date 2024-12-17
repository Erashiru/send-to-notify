package moysklad

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kwaaka-team/orders-core/config/menu"
	constErrors "github.com/kwaaka-team/orders-core/core/menu/clients/aggregator/errors"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/managers/validator"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	"github.com/kwaaka-team/orders-core/core/menu/models/utils"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/pkg/errors"
	"strconv"
	"time"

	moyskladCli "github.com/kwaaka-team/orders-core/pkg/moysklad"
	moysklad "github.com/kwaaka-team/orders-core/pkg/moysklad/clients"
	msModels "github.com/kwaaka-team/orders-core/pkg/moysklad/models"
	storeCore "github.com/kwaaka-team/orders-core/pkg/store"
	storeDto "github.com/kwaaka-team/orders-core/pkg/store/dto"
	"github.com/rs/zerolog/log"
)

type Config struct {
	BaseURL     string
	ProductHref string
	ProductType string
	Protocol    string
	UserName    string
	Password    string
	Quantity    string
}

type Moysklad interface {
	ModifyProduct(ctx context.Context, storeID string, product models.Product) (models.ProductModifyResponse, error)
	UploadMenu(ctx context.Context, menuId, extStoreId string, menu models.Menu, store storeModels.Store, offers []models.BkOffers, sv3 *s3.S3, userRole string) (models.ExtTransaction, error)
	BulkUpdate(ctx context.Context, restaurantID, storeID string, products models.Products, attributes models.Attributes, store storeModels.Store) (string, error)
	VerifyMenu(ctx context.Context, transaction models.ExtTransaction) (models.Status, error)
	BulkAttribute(ctx context.Context, storeID string, attributes models.Attributes) (string, error)
	GetMenu(ctx context.Context, extStoreId string) (models.Menu, error)
	ValidateMenu(ctx context.Context, request models.MenuValidateRequest) (models.MenuUploadTransaction, error)
}

type Manager struct {
	cli          moysklad.MoySklad
	cfg          *Config
	storeCli     storeCore.Client
	mspRepo      drivers.MSPositionsRepository
	stRepo       drivers.StopListTransactionRepository
	orderID      string
	restaurantID string
}

func (m Manager) ModifyProduct(ctx context.Context, storeID string, product models.Product) (models.ProductModifyResponse, error) {
	return models.ProductModifyResponse{}, constErrors.ErrNotImplemented
}
func (m Manager) BulkUpdate(ctx context.Context, restaurantID, storeID string, products models.Products, attributes models.Attributes, store storeModels.Store) (string, error) {
	m.orderID = store.MoySklad.OrderID
	m.restaurantID = restaurantID
	var product models.ProductRequest
	reqProd := models.ProductRequest{}
	reqAttr := models.AttributeRequest{}

	productsRequest := reqProd.ToModel(products)
	attributesRequest := reqAttr.ToModel(attributes)

	if len(productsRequest) == 0 {
		return "", errors.New("empty request products")
	}

	log.Info().Msgf("before prod len Product %v len Attr %v ", len(productsRequest), len(attributesRequest))
	productsRequest = append(productsRequest, product.FromAttribute(attributesRequest)...) //case: attr like product
	log.Info().Msgf("after merge with attributes %v", len(productsRequest))
	productsRequest = product.RemoveDuplicate(productsRequest)

	if !store.IikoCloud.StopListByBalance {
		log.Info().Msgf("StopListByBalance state false, lessZeroCase")
		productsRequest = product.LessZeroList(productsRequest) // for ms case, add soplist only <= 0 balanceproduct
	}
	log.Info().Msgf("after prod check less then Balance len %v", len(productsRequest))

	positions, err := m.mspRepo.GetPositions(
		ctx,
		selector.EmptyMoySkladSearch().
			SetRestaurantID(store.ID).
			SetIsDeleted(false),
	)
	if err != nil {
		return "", err
	}

	if len(productsRequest) == 0 {
		return "", errors.New("empty request products")
	}

	positionIDs, err := m.updatePositions(ctx, positions, productsRequest)
	if err != nil {
		return "", err
	}

	var transaction models.StopListTransaction
	for _, positionID := range positionIDs {
		transaction.Append(positionID, models.MOYSKLAD.String(), store.MoySklad.OrganizationID, "", models.SUCCESS)
	}
	reqProd = models.ProductRequest{}
	reqAttr = models.AttributeRequest{}

	toModelProducts := reqProd.FromModel(productsRequest)
	toModelAttributes := reqAttr.FromModel(attributesRequest)

	transaction.Fill(store.ID, toModelProducts, toModelAttributes)
	id, err := m.stRepo.Insert(ctx, transaction)
	if err != nil {
		return fmt.Sprintf("%s from delivery %s", validator.ErrAddTransaction, "moysklad"), nil
	}
	return id, nil
}

func (m Manager) UploadMenu(ctx context.Context, menuId, extStoreId string, menu models.Menu, store storeModels.Store, offers []models.BkOffers, sv3 *s3.S3, userRole string) (models.ExtTransaction, error) {
	return models.ExtTransaction{}, constErrors.ErrNotImplemented
}

func (m Manager) VerifyMenu(ctx context.Context, transaction models.ExtTransaction) (models.Status, error) {
	return "", constErrors.ErrNotImplemented
}

func (m Manager) BulkAttribute(ctx context.Context, storeID string, attributes models.Attributes) (string, error) {
	return "", constErrors.ErrNotImplemented
}

func NewManager(ctx context.Context, cfg menu.Configuration, mspRepo drivers.MSPositionsRepository,
	stRepo drivers.StopListTransactionRepository) (Moysklad, error) {

	cli, err := moyskladCli.NewMoySkladClient(&moysklad.Config{
		Protocol: cfg.MoySkladConfiguration.Protocol,
		BaseURL:  cfg.MoySkladConfiguration.BaseURL,
		Username: cfg.MoySkladConfiguration.Username,
		Password: cfg.MoySkladConfiguration.Password,
	})

	if err != nil {
		log.Trace().Err(err).Msg("can't initialize moysklad client.")
		return nil, err
	}

	storeCli, err := storeCore.NewClient(storeDto.Config{
		Region:    cfg.AwsConfig.Region,
		SecretEnv: cfg.SecretEnvironments,
	})
	if err != nil {
		return nil, fmt.Errorf("cannot initialize Store Client")
	}

	return &Manager{
		cli: cli,
		cfg: &Config{
			BaseURL:     cfg.MoySkladConfiguration.BaseURL,
			ProductHref: cfg.MoySkladConfiguration.ProductHref,
			ProductType: cfg.MoySkladConfiguration.ProductType,
			Protocol:    cfg.MoySkladConfiguration.Protocol,
			Quantity:    cfg.MoySkladConfiguration.Quantity,
		},
		storeCli: storeCli,
		mspRepo:  mspRepo,
		stRepo:   stRepo,
	}, nil
}

func (m Manager) GetMenu(ctx context.Context, extStoreId string) (models.Menu, error) {
	return models.Menu{}, constErrors.ErrNotImplemented
}

func (m Manager) updatePositions(ctx context.Context, positions []models.MoySkladPosition, products []models.ProductRequest) ([]string, error) {
	var positionIDs []string

	if len(positions) == 0 {
		log.Info().Msgf("step4 list position == 0, first init %v", len(positions))
		for _, product := range products {
			temp := models.MoySkladPosition{
				MsID:      product.MSID,
				ProductID: product.ID,
				Available: *product.IsAvailable,
			}
			//add to pos to  moysklad
			positionID, err := m.add(ctx, temp)
			if err != nil {
				continue
			}
			positionIDs = append(positionIDs, positionID)
		}
	} else {

		mapPositions := make(map[string]models.MoySkladPosition, len(positions))
		for _, position := range positions {
			mapPositions[position.ProductID] = position
		}
		log.Info().Msgf("step5 have position in DB len map %v", len(mapPositions))
		//compare current positions if no has DB, create, else update is_deleted
		for _, product := range products {
			if position, ok := mapPositions[product.ID]; ok {
				//log.Info().Msgf("step6 have position in DB positID %v %v %v", position.ID, *product.IsAvailable, position.Available)

				if product.IsAvailable != nil && position.Available != *product.IsAvailable {
					log.Info().Msgf("step6.2 different state product %v", product.IsAvailable)

					if *product.IsAvailable {
						log.Info().Msgf("step7 available = true, removePosition and softDeleteDB %v", product.IsAvailable)
						if err := m.delete(ctx, position); err != nil {
							continue
						}
						positionIDs = append(positionIDs, position.ID)
						//{else } log.Info().Msgf("step8 not available, add position and insertDB ", product.IsAvailable)
					}
				} //else  skip
			} else {
				if product.IsAvailable != nil && !*product.IsAvailable {
					//not in map -> not in DB -> create new
					log.Info().Msgf("step9, not found in MAP,  add new position, woopay & DB %v %v", product.ID, *product.IsAvailable)
					temp := models.MoySkladPosition{
						RestaurantID: m.restaurantID,
						OrderID:      m.orderID,
						MsID:         product.MSID,
						ProductID:    product.ID,
					}
					positionID, err := m.add(ctx, temp)
					if err != nil {
						continue
					}
					positionIDs = append(positionIDs, positionID)
				}
			}
		}
	}
	return positionIDs, nil
}

func (m Manager) delete(ctx context.Context, position models.MoySkladPosition) error {

	if err := m.cli.DeleteProductSupplier(ctx, msModels.Position{
		ProductID: position.ID,
		OrderID:   m.orderID,
	}); err != nil {
		return err
	}
	log.Info().Msgf("delete  position in client")

	//soft delete
	if err := m.deleteDB(ctx, position); err != nil {
		return err
	}
	log.Info().Msgf("delete  position in DB")

	return nil
}

func (m Manager) deleteDB(ctx context.Context, position models.MoySkladPosition) error {

	position.IsDeleted = true
	position.Available = true
	position.UpdatedAt = time.Now().UTC()

	if err := m.mspRepo.RemovePosition(
		ctx,
		selector.MoySklad{
			OrderID:      position.OrderID,
			RestaurantID: position.RestaurantID,
			ProductID:    position.ProductID,
			MsID:         position.MsID,
			Code:         position.Code,
			ID:           position.ID,
			Available:    position.Available,
			IsDeleted:    &position.IsDeleted,
			CreatedAt:    position.CreatedAt,
			UpdatedAt:    position.UpdatedAt,
		},
	); err != nil {
		return err
	}
	return nil
}

func (m Manager) add(ctx context.Context, position models.MoySkladPosition) (string, error) {
	num, err := strconv.Atoi(m.cfg.Quantity)
	if err != nil {
		return "", err
	}
	request := msModels.Position{
		ProductID: position.MsID,
		OrderID:   m.orderID,
		Quantity:  num,
		Assortment: msModels.Assortment{
			Meta: msModels.Meta{
				Type: m.cfg.ProductType,
				Href: m.cfg.ProductHref + position.MsID,
			},
		},
	}

	id, err := m.cli.AddProductSupplier(ctx, request)
	if err != nil {
		return "", err
	}
	utils.Beautify("add position in client", request)
	position.ID = id
	position.CreatedAt = time.Now().UTC()
	position.UpdatedAt = time.Now().UTC()
	if err = m.mspRepo.CreatePosition(ctx, position); err != nil {
		return "", err
	}
	log.Info().Msgf("add position in DB")

	return id, nil
}

func (m Manager) ValidateMenu(ctx context.Context, request models.MenuValidateRequest) (models.MenuUploadTransaction, error) {
	return models.MenuUploadTransaction{}, nil
}
