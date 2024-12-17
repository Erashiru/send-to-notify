package order

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/rs/zerolog/log"
)

type Subscriber interface {
	SendOrder(ctx context.Context, order models.Order, store coreStoreModels.Store, posStatus models.PosStatus) error
}

type Publisher struct {
	Subscribers []Subscriber
}

func (s *Publisher) AddSubscriber(subscriber Subscriber) {
	s.Subscribers = append(s.Subscribers, subscriber)
}

func (s *Publisher) NotifySubscribers(ctx context.Context, order models.Order, store coreStoreModels.Store, posStatus models.PosStatus) error {
	for i := range s.Subscribers {
		subscriber := s.Subscribers[i]
		if err := subscriber.SendOrder(ctx, order, store, posStatus); err != nil {
			log.Info().Msgf("can not create order on 3pl service by order id: %s", order.ID)
		}
	}
	return nil
}
