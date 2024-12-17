package managers

import (
	"context"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/foodband/models"
	"github.com/kwaaka-team/orders-core/domain/logger"
	orderCore "github.com/kwaaka-team/orders-core/pkg/order"
	"go.uber.org/zap"
)

type Order interface {
	UpdateOrderStatus(ctx context.Context, req models.UpdateOrderStatusReq) error
}

type orderImplementation struct {
	orderCoreCli orderCore.Client
	logger       *zap.SugaredLogger
}

func NewOrderManager(orderCli orderCore.Client, logger *zap.SugaredLogger) Order {
	return &orderImplementation{
		orderCoreCli: orderCli,
		logger:       logger,
	}
}

func (man *orderImplementation) UpdateOrderStatus(ctx context.Context, req models.UpdateOrderStatusReq) error {
	man.logger.Info(logger.LoggerInfo{
		System:  "foodband request",
		Request: req,
	})

	if req.OrderID == "" || req.Status == "" {
		man.logger.Error(logger.LoggerInfo{
			System:   "foodband response error",
			Response: fmt.Sprintf("invalid inpus, orderID %s, status %s", req.OrderID, req.Status),
		})
		return fmt.Errorf("invalid order id or status")
	}

	if err := man.orderCoreCli.UpdateOrderStatus(ctx, req.OrderID, "foodband", req.Status, ""); err != nil {
		man.logger.Error(logger.LoggerInfo{
			System:   "foodband response error",
			Response: err,
		})
		return err
	}

	return nil
}
