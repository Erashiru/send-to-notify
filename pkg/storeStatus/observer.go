package storeStatus

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/rs/zerolog/log"
)

type Observer interface {
	Notify(ctx context.Context, restaurant models.Store, externalStoreID, deliveryService string, storeIsOpened bool) error
	NotifyStatusReport(ctx context.Context, restaurant models.Store, durations []models.OpenTimeDuration) error
	NotifyStatusChange(ctx context.Context, status string, phone string) error
}

type Subject struct {
	observers []Observer
}

func (s *Subject) AddObserver(observer Observer) {
	s.observers = append(s.observers, observer)
}

func (s *Subject) NotifyObservers(ctx context.Context, restaurant models.Store, externalStoreID, deliveryService string, storeIsOpened bool) error {
	for _, observer := range s.observers {
		if err := observer.Notify(ctx, restaurant, externalStoreID, deliveryService, storeIsOpened); err != nil {
			log.Info().Msgf("error was occured during sending message notification service by restaurant id: %s", restaurant.ID)
		}
	}
	return nil
}

func (s *Subject) NotifyStatusReportObservers(ctx context.Context, restaurant models.Store, durations []models.OpenTimeDuration) error {
	for _, observer := range s.observers {
		if err := observer.NotifyStatusReport(ctx, restaurant, durations); err != nil {
			log.Info().Msgf("error was occured during sending message notification service: %v", durations)
		}
	}
	return nil
}

func (s *Subject) NotifyStatusChangeObservers(ctx context.Context, status string, phone string) error {
	for _, observer := range s.observers {
		if err := observer.NotifyStatusChange(ctx, status, phone); err != nil {
			log.Info().Msgf("error was occured during sending message notification service by phone: %v", phone)
		}
	}
	return nil
}
