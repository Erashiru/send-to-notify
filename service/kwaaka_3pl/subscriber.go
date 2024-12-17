package kwaaka_3pl

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	models2 "github.com/kwaaka-team/orders-core/service/kwaaka_3pl/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"time"
)

type Subscriber3PL struct {
	kwaaka3plService Service
}

func NewSubscriber3PL(kwaaka3pl Service) (*Subscriber3PL, error) {
	if kwaaka3pl == nil {
		return nil, errors.New("kwaaka3pl is nil")
	}
	return &Subscriber3PL{
		kwaaka3plService: kwaaka3pl,
	}, nil
}

func (s Subscriber3PL) SendOrder(ctx context.Context, order models.Order, store coreStoreModels.Store, posStatus models.PosStatus) error {
	if !store.Kwaaka3PL.Is3pl {
		return nil
	}
	if order.IsMarketplace || !order.SendCourier {
		return nil
	}
	if order.DeliveryOrderID != "" || order.DeliveryDispatcher == "" || order.DeliveryDispatcher == coreStoreModels.SELFDELIVERY.String() {
		return nil
	}

	switch posStatus {
	case models.NEW, models.ACCEPTED, models.CANCELLED_BY_POS_SYSTEM, models.WAIT_SENDING, models.FAILED:
		return nil
	default:
		log.Info().Msgf("going to create 3pl order, posStatus: %s, orderID: %s", posStatus.String(), order.ID)
	}

	items := make([]models2.Item, 0, len(order.Products))
	for _, product := range order.Products {
		items = append(items, models2.Item{
			Name:     product.Name,
			ID:       product.ID,
			Quantity: product.Quantity,
			Price:    product.Price.Value,
		})
	}

	if order.DeliveryDispatcher == models2.IndriveDelivery {
		return nil
	}

	if err := s.kwaaka3plService.Create3plOrder(ctx, models2.CreateDeliveryRequest{
		ID:                order.ID,
		FullDeliveryPrice: order.FullDeliveryPrice,
		Provider:          order.DeliveryDispatcher,
		PickUpTime:        time.Now().Add(time.Duration(order.DispatcherDeliveryTime) * time.Minute),
		DeliveryAddress: models2.Address{
			Label:        order.DeliveryAddress.Label,
			Lat:          order.DeliveryAddress.Latitude,
			Lon:          order.DeliveryAddress.Longitude,
			Comment:      order.DeliveryAddress.Comment,
			BuildingName: order.DeliveryAddress.BuildingName,
			Street:       order.DeliveryAddress.Street,
			Flat:         order.DeliveryAddress.Flat,
			Porch:        order.DeliveryAddress.Porch,
			Floor:        order.DeliveryAddress.Floor,
		},
		StoreAddress: models2.Address{
			Label:   store.Address.City + ", " + store.Address.Street,
			Lon:     store.Address.Coordinates.Longitude,
			Lat:     store.Address.Coordinates.Latitude,
			Comment: store.Address.Entrance,
		},
		CustomerInfo: models2.CustomerInfo{
			Name:  s.setCustomerName(order.Customer.Name),
			Phone: s.setCustomerPhoneNumber(order.Customer.PhoneNumber),
			Email: s.setCustomerEmail(order.Customer.Email),
		},
		StoreInfo: models2.StoreInfo{
			Name:  store.Name,
			Phone: store.StorePhoneNumber,
			Email: store.Settings.Email,
		},
		PickUpCode:      order.PickUpCode,
		Currency:        order.Currency,
		Comment:         order.SpecialRequirements,
		Items:           items,
		ExternalStoreID: store.Kwaaka3PL.IndriveStoreID,
		TaxiClass:       store.Kwaaka3PL.TaxiClass,
	}); err != nil {
		return err
	}

	return nil
}

func (s *Subscriber3PL) setCustomerEmail(email string) string {
	customerEmail := email
	if customerEmail == "" {
		customerEmail = models.Default3plCustomerEmail
	}

	return customerEmail
}

func (s *Subscriber3PL) setCustomerPhoneNumber(phoneNumber string) string {
	customerPhone := phoneNumber
	if customerPhone == "" {
		customerPhone = models.Default3plCustomerPhone
	}

	return customerPhone
}

func (s *Subscriber3PL) setCustomerName(name string) string {
	customerName := name
	if customerName == "" {
		customerName = models.Default3plCustomerName
	}

	return customerName
}
