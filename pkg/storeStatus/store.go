package storeStatus

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/kwaaka-team/orders-core/service/aggregator"
	"github.com/kwaaka-team/orders-core/service/store"
	"github.com/kwaaka-team/orders-core/service/store/repository/storeclosedtime"
	storeGroupServicePkg "github.com/kwaaka-team/orders-core/service/storegroup"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

type StoreStatusService struct {
	aggFactory        aggregator.Factory
	storeService      store.Service
	datastoreClient   storeclosedtime.Repository
	storeGroupService storeGroupServicePkg.Service
	subject           *Subject
	StoreSchedule     *StoreSchedule
}

type Status interface {
	UpdateStoresSchedule(ctx context.Context, deliveryService string) error
	CheckStoreStatus(ctx context.Context, deliveryService string) error
	ReportStatuses(ctx context.Context) error
}

func NewStoreStatus(aggFactory aggregator.Factory, storeService store.Service, subject *Subject, datastoreClient storeclosedtime.Repository, storeGroupService storeGroupServicePkg.Service) (Status, error) {
	if aggFactory == nil {
		return nil, errors.New("aggregator factory is empty")
	}

	if storeService == nil {
		return nil, errors.New("store service service is empty")
	}

	if subject == nil {
		return nil, errors.New("store service service is empty")
	}

	return StoreStatusService{
		aggFactory:        aggFactory,
		storeService:      storeService,
		datastoreClient:   datastoreClient,
		subject:           subject,
		storeGroupService: storeGroupService,
	}, nil
}

type StoreSchedule struct {
	StoreAutoOpen bool
	Schedule      storeModels.AggregatorSchedule
}

func newStoreSchedule(store storeModels.Store, deliveryService string) *StoreSchedule {
	switch deliveryService {
	case models.WOLT.String():
		return &StoreSchedule{
			StoreAutoOpen: store.Wolt.StoreAutoOpen,
			Schedule:      store.StoreSchedule.WoltSchedule,
		}
	case models.GLOVO.String():
		return &StoreSchedule{
			StoreAutoOpen: store.Glovo.StoreAutoOpen,
			Schedule:      store.StoreSchedule.GlovoSchedule,
		}
	}

	return &StoreSchedule{}
}

func (ss StoreStatusService) ReportStatuses(ctx context.Context) error {
	restaurants, err := ss.storeService.FindAllStores(ctx)
	if err != nil {
		return err
	}

	aggregators := []string{models.GLOVO.String(), models.WOLT.String()}

	for _, restaurant := range restaurants {
		if restaurant.Telegram.StoreStatusChatId == "" {
			continue
		}

		log.Info().Msgf("Get open close status for %s store", restaurant.Name)
		var durations []storeModels.OpenTimeDuration

		for _, agg := range aggregators {
			aggSvc, err := ss.aggFactory.GetAggregator(agg, restaurant)
			if err != nil {
				log.Err(err)
				continue
			}

			storeIDs, err := ss.storeService.GetStoreExternalIds(restaurant, agg)
			if err != nil {
				log.Err(err)
				continue
			}

			var actualOpenDuration time.Duration
			var totalOpenDuration time.Duration
			for _, storeID := range storeIDs {
				schedule, err := aggSvc.GetStoreSchedule(ctx, storeID)
				if err != nil {
					log.Err(err) // Тут ошибка с сообщением You are not authorized. у Садыхана
					continue
				}

				if len(schedule.Schedule) == 0 {
					continue
				}

				weekday, err := ss.getWeekdayInTimezone(schedule.Timezone)
				if err != nil {
					return err
				}

				scheduleIndex, ok := ss.getScheduleIndex(schedule, weekday)
				if !ok {
					continue
				}

				actOpenDur, err := ss.actualOpenDuration(schedule, scheduleIndex)
				if err != nil {
					return err
				}

				totOpenDur, err := ss.totalOpenDuration(schedule, scheduleIndex)
				if err != nil {
					return err
				}

				actualOpenDuration += actOpenDur
				totalOpenDuration += totOpenDur
			}

			durations = append(durations, storeModels.OpenTimeDuration{
				DeliveryService:        agg,
				ActualOpenTimeDuration: actualOpenDuration,
				TotalOpenTimeDuration:  totalOpenDuration,
			})
		}

		if err := ss.subject.NotifyStatusReportObservers(ctx, restaurant, durations); err != nil {
			log.Err(err).Msgf("failed to notify status report for %s restaurant", restaurant.Name)
			continue
		}
	}

	return nil
}

func (ss StoreStatusService) UpdateStoresSchedule(ctx context.Context, deliveryService string) error {
	restaurants, err := ss.storeService.FindStoresByDeliveryService(ctx, deliveryService)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, restaurant := range restaurants {
		wg.Add(1)

		go func(rest storeModels.Store) {
			defer wg.Done()

			externalStoreIds, err := ss.storeService.GetStoreExternalIds(rest, deliveryService)
			if err != nil {
				log.Err(err).Msgf("error occurred during getting store external ids for store ID %s", rest.ID)
				return
			}

			for _, externalStoreId := range externalStoreIds {
				if err = ss.updateStoreSchedule(ctx, rest, externalStoreId, deliveryService); err != nil {
					log.Err(err).Msgf("error occurred during updating schedule for restaurant ID %s and external store ID %s: ", rest.ID, externalStoreId)
				}
			}
		}(restaurant)
	}

	wg.Wait()

	return nil
}

func (ss StoreStatusService) updateStoreSchedule(ctx context.Context, store storeModels.Store, externalStoreId, deliveryService string) error {
	aggSvc, err := ss.aggFactory.GetAggregator(deliveryService, store)
	if err != nil {
		return err
	}

	schedule, err := aggSvc.GetStoreSchedule(ctx, externalStoreId)
	if err != nil {
		return err
	}

	if schedule.Timezone == "" {
		schedule.Timezone = store.Settings.TimeZone.TZ
	}

	if err = ss.storeService.UpdateStoreSchedule(ctx, storeModels.UpdateStoreSchedule{
		RestaurantID:    store.ID,
		DeliveryService: deliveryService,
		StoreSchedule:   schedule,
	}); err != nil {
		return err
	}

	log.Printf("restaurant by external store id %s was succesfully updated schedule in restaurant by id %s", externalStoreId, store.ID)

	return nil
}

func (ss StoreStatusService) CheckStoreStatus(ctx context.Context, deliveryService string) error {
	restaurants, err := ss.storeService.FindStoresByDeliveryService(ctx, deliveryService)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, restaurant := range restaurants {
		wg.Add(1)

		go func(store storeModels.Store) {
			defer wg.Done()

			externalStoreIDs, err := ss.storeService.GetStoreExternalIds(store, deliveryService)
			if err != nil {
				log.Err(err).Msgf("error occurred during getting store external ids for store ID %s:", store.ID)
				return
			}

			for _, externalStoreId := range externalStoreIDs {
				if err := ss.sendNotify(ctx, store, externalStoreId, deliveryService); err != nil {
					log.Err(err).Msgf("error occurred during sending notification for restaurant ID %s and store ID %s", store.ID, externalStoreId)
				}
			}
		}(restaurant)

	}

	wg.Wait()

	return nil
}

func (ss StoreStatusService) calculateOpenCloseTimeDifference(timeSlot storeModels.TimeSlot) (time.Duration, error) {
	if timeSlot.Opening == "24:00" {
		timeSlot.Opening = "23:59"
	}

	if timeSlot.Closing == "24:00" {
		timeSlot.Closing = "23:59"
	}

	openTime, err := time.Parse("15:04", timeSlot.Opening)
	if err != nil {
		return time.Duration(0), err
	}
	closeTime, err := time.Parse("15:04", timeSlot.Closing)
	if err != nil {
		return time.Duration(0), err
	}

	return openTime.Sub(closeTime), nil
}

func (ss StoreStatusService) getWeekdayInTimezone(timezone string) (int, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return 0, err
	}

	now := time.Now().In(loc)
	return int(now.Weekday()), nil
}

func (ss StoreStatusService) getScheduleIndex(schedule storeModels.AggregatorSchedule, weekday int) (int, bool) {
	scheduleIndex := -1
	for i, sch := range schedule.Schedule {
		if sch.DayOfWeek == weekday {
			scheduleIndex = i
			break
		}
	}

	if scheduleIndex == -1 {
		return -1, false
	}

	return scheduleIndex, true
}

func (ss StoreStatusService) actualOpenDuration(schedule storeModels.AggregatorSchedule, scheduleIndex int) (time.Duration, error) {
	result := time.Duration(0)
	for _, timeSlot := range schedule.Schedule[scheduleIndex].TimeSlots {
		openTime, err := ss.calculateOpenCloseTimeDifference(timeSlot)
		if err != nil {
			return time.Duration(0), err
		}

		result += openTime
	}

	return result, nil
}

func (ss StoreStatusService) totalOpenDuration(schedule storeModels.AggregatorSchedule, scheduleIndex int) (time.Duration, error) {
	return ss.calculateOpenCloseTimeDifference(storeModels.TimeSlot{
		Opening: schedule.Schedule[scheduleIndex].TimeSlots[0].Opening,
		Closing: schedule.Schedule[scheduleIndex].TimeSlots[len(schedule.Schedule[scheduleIndex].TimeSlots)-1].Closing,
	})
}

func (ss StoreStatusService) sendNotify(ctx context.Context, store storeModels.Store, externalStoreId, deliveryService string) error {
	needToNotify, storeIsOpened, err := ss.needToSendNotify(ctx, store, externalStoreId, deliveryService)
	if err != nil {
		return err
	}

	if needToNotify {
		log.Printf("restaurant by id: %s and store id: %s  was closed in working hour and it will send notify", store.ID, externalStoreId)
		if err = ss.subject.NotifyObservers(ctx, store, externalStoreId, deliveryService, storeIsOpened); err != nil {
			return err
		}
	}

	return nil
}

func (ss StoreStatusService) needToSendNotify(ctx context.Context, store storeModels.Store, externalStoreID, deliveryService string) (bool, bool, error) {
	needToNotify, storeIsOpened, err := ss.compareCurrentVsExpectedStatus(ctx, store, externalStoreID, deliveryService)
	if err != nil {
		return false, storeIsOpened, err
	}

	res, hasRecord, err := ss.datastoreClient.GetByFilter(ctx, storeModels.FilterStoreActiveTime{
		RestaurantID:    store.ID,
		StoreID:         externalStoreID,
		DeliveryService: deliveryService,
	})
	if err != nil {
		return false, storeIsOpened, err
	}

	if needToNotify && hasRecord {
		return false, storeIsOpened, nil
	}

	if !needToNotify && hasRecord {
		if err = ss.datastoreClient.UpdateEndTime(ctx, res.ID); err != nil {
			return false, storeIsOpened, err
		}
	}

	return needToNotify, storeIsOpened, nil
}

func (ss StoreStatusService) compareCurrentVsExpectedStatus(ctx context.Context, store storeModels.Store, externalStoreId, deliveryService string) (bool, bool, error) {
	aggSvc, err := ss.aggFactory.GetAggregator(deliveryService, store)
	if err != nil {
		return false, false, err
	}

	storeStatus, err := aggSvc.GetStoreStatus(ctx, externalStoreId)
	if err != nil {
		log.Err(err).Msgf("error was occured during getting store status by and external store id: %s", externalStoreId)
		return false, false, err
	}

	ss.StoreSchedule = newStoreSchedule(store, deliveryService)

	location, err := time.LoadLocation(ss.StoreSchedule.Schedule.Timezone)
	if err != nil {
		return false, false, err
	}

	currentTime := time.Now().In(location)

	expectedStatus := checkSchedule(currentTime, store.StoreSchedule.GlovoSchedule)

	storeIsOpened, err := ss.openStoreStatusIfNeeded(ctx, store, externalStoreId, deliveryService, storeStatus, expectedStatus)
	if err != nil {
		return false, false, err
	}

	return storeStatus != expectedStatus && expectedStatus, storeIsOpened, nil
}

func checkSchedule(currentTime time.Time, schedule storeModels.AggregatorSchedule) bool {
	currentDayOfWeek := int(currentTime.Weekday())

	if currentDayOfWeek == 0 {
		currentDayOfWeek = 7
	}

	for _, daySchedule := range schedule.Schedule {
		if daySchedule.DayOfWeek == currentDayOfWeek {
			for _, timeSlot := range daySchedule.TimeSlots {
				openingTime, _ := time.Parse("15:04", timeSlot.Opening)
				closingTime, _ := time.Parse("15:04", timeSlot.Closing)
				currentTimeParsed, _ := time.Parse("15:04", currentTime.Format("15:04"))
				if openingTime == closingTime {
					return true
				}

				if closingTime.Before(openingTime) {
					if currentTimeParsed.After(openingTime) || currentTimeParsed.Before(closingTime) {
						return true
					}
				}
				if currentTimeParsed.After(openingTime) && currentTimeParsed.Before(closingTime) {
					return true
				}

			}
		}
	}
	return false
}

func checkStartTime(currentTime time.Time, schedule storeModels.AggregatorSchedule) bool {
	currentDayOfWeek := int(currentTime.Weekday())
	currentDayOfWeek = currentDayOfWeek + 1

	for _, daySchedule := range schedule.Schedule {
		if daySchedule.DayOfWeek == currentDayOfWeek {
			for _, timeSlot := range daySchedule.TimeSlots {
				openingTime, _ := time.Parse("15:04", timeSlot.Opening)
				currentTimeParsed, _ := time.Parse("15:04", currentTime.Format("15:04"))

				sub := currentTimeParsed.Sub(openingTime).Minutes()

				if sub < 15 {
					return true
				}
			}
		}
	}
	return false
}

func (ss StoreStatusService) openStoreStatusIfNeeded(ctx context.Context, store storeModels.Store, externalStoreId, deliveryService string, currentStatus, expectedStatus bool) (bool, error) {
	if expectedStatus == false {
		return false, nil
	}

	if currentStatus == expectedStatus {
		return false, nil
	}

	if !ss.StoreSchedule.StoreAutoOpen {
		return false, nil
	}

	aggSvc, err := ss.aggFactory.GetAggregator(deliveryService, store)
	if err != nil {
		return false, err
	}

	if err := aggSvc.OpenStore(ctx, externalStoreId); err != nil {
		log.Info().Msgf("The restaurant: %s, with store ID %s will be opened because it's time for opening.", store.ID, externalStoreId)
		return false, err
	}

	return true, nil
}
