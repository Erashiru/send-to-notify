package manager

import (
	"context"
	"github.com/kwaaka-team/orders-core/domain/logger"
	"github.com/kwaaka-team/orders-core/core/talabat/models"
	"github.com/kwaaka-team/orders-core/pkg/menu"
	menuModels "github.com/kwaaka-team/orders-core/pkg/menu/dto"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Menu interface {
	UpdateMenuUploadTransaction(ctx context.Context, req models.MenuUploadCallbackRequest) error
}

type menuImplementation struct {
	logger  *zap.SugaredLogger
	menuCli menu.Client
}

func NewMenu(menuCli menu.Client, logger *zap.SugaredLogger) Menu {
	return &menuImplementation{
		logger:  logger,
		menuCli: menuCli,
	}
}

func (man *menuImplementation) UpdateMenuUploadTransaction(ctx context.Context, req models.MenuUploadCallbackRequest) error {
	if req.RequestId == "" {
		man.logger.Error(logger.LoggerInfo{
			System:   "talabat update menu upload transaction error",
			Response: errors.New("invalid request id"),
		})
		return errors.New("invalid request id")
	}
	mut, err := man.menuCli.GetMenuUploadTransaction(ctx, menuModels.MenuUploadTransaction{
		Service: "talabat",
		ExtTransactions: []menuModels.ExtTransaction{
			{
				ID: req.RequestId,
			},
		},
	})

	if err != nil {
		man.logger.Error(logger.LoggerInfo{
			System:   "talabat update menu upload transaction error",
			Response: err,
		})
		return err
	}

	switch req.Status {
	case "COMPLETED":
		mut.Status = "SUCCESS"
	case "COMPLETED_WITH_ERRORS":
		mut.Status = "ERROR"
	default:
		man.logger.Error(logger.LoggerInfo{
			System:   "talabat update menu upload transaction error",
			Response: errors.New("invalid status"),
		})
		return errors.New("invalid status")
	}

	for i, ext := range mut.ExtTransactions {
		if ext.ID == req.RequestId {
			mut.ExtTransactions[i].Status = mut.Status
			if req.Description != "" {
				mut.ExtTransactions[i].Details = append(mut.ExtTransactions[i].Details, req.Description)
			}
		}
	}
	if req.Description != "" {
		mut.Details = append(mut.Details, req.Description)
	}

	err = man.menuCli.UpdateMenuUploadTransaction(ctx, mut)
	if err != nil {
		man.logger.Error(logger.LoggerInfo{
			System:   "talabat update menu upload transaction error",
			Response: err,
		})
		return err
	}

	return nil
}
