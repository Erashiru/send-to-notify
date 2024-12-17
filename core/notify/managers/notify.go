package managers

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/notify/clients/bitrix"
	"github.com/kwaaka-team/orders-core/core/notify/clients/clickup"
	"github.com/kwaaka-team/orders-core/core/notify/config"
	models2 "github.com/kwaaka-team/orders-core/core/notify/models"
	"github.com/rs/zerolog/log"
)

type NotifyManager interface {
	SendNotification(ctx context.Context, message models2.Message) ([]models2.Result, error)
}

type nnm struct {
	bitrixCli  bitrix.Bitrix
	clickupCli clickup.ClickUp
}

func NewNotifyManager(cfg config.Configuration) NotifyManager {
	return &nnm{
		bitrixCli:  bitrix.NewBitrix(cfg.BitrixConfiguration),
		clickupCli: clickup.NewClickUp(cfg.ClickUpConfiguration),
	}
}

func (n *nnm) SendNotification(ctx context.Context, message models2.Message) ([]models2.Result, error) {

	results := make([]models2.Result, 0, len(message.Services))

	for _, service := range message.Services {

		switch service {
		case models2.BITRIX:
			taskNumber, err := n.bitrixCli.CreateTask(ctx, message.Title, message.Description)
			if err != nil {
				results = append(results, models2.Result{
					Service: service,
					Status:  models2.ERROR,
					Message: "add task in bitrix",
				})
				continue
			}
			results = append(results, models2.Result{
				Service: service,
				Status:  models2.SUCCESS,
				Message: fmt.Sprintf("successfully created task %d", taskNumber.Result),
			})
			log.Info().Msgf("successfully created task %d", taskNumber.Result)

		case models2.CLICKUP:
			task, err := n.clickupCli.CreateTask(ctx, message)
			if err != nil {
				results = append(results, models2.Result{
					Service: service,
					Status:  models2.ERROR,
					Message: "add task in clickup",
				})
				continue
			}
			results = append(results, models2.Result{
				Service: service,
				Status:  models2.SUCCESS,
				Message: fmt.Sprintf("successfully created task %s with id %s", task.Name, task.Id),
			})
			log.Info().Msgf("successfully created task in clickup %s with id %s", task.Name, task.Id)
		default:
			results = append(results, models2.Result{
				Service: service,
				Status:  models2.ERROR,
				Message: models2.ErrNoService.Error(),
			})
			continue
		}
	}

	return results, nil
}
