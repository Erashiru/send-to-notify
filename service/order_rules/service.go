package order_rules

import (
	"context"
	coreOrderModels "github.com/kwaaka-team/orders-core/core/models"
	"github.com/kwaaka-team/orders-core/service/order_rules/models"
	"github.com/rs/zerolog/log"
)

type Service interface {
	UseTheOrderRules(ctx context.Context, order coreOrderModels.Order) (coreOrderModels.Order, error)
}

type serviceImpl struct {
	repo Repository
}

func NewOrderRuleService(repo Repository) (*serviceImpl, error) {
	service := &serviceImpl{
		repo: repo,
	}

	return service, nil
}

func (svc *serviceImpl) UseTheOrderRules(ctx context.Context, order coreOrderModels.Order) (coreOrderModels.Order, error) {
	rules, err := svc.repo.FindOrderRulesByRestaurantId(ctx, order.RestaurantID)
	if err != nil {
		log.Err(err).Msgf("find order rules by restaurant id error")
		return order, nil
	}

	for _, rule := range rules {
		switch rule.Type {
		case models.OrderAddition:
			order = svc.useOrderAdditionRule(rule, order)
		}
	}

	return order, nil
}

func (svc *serviceImpl) useOrderAdditionRule(rule models.OrderRule, order coreOrderModels.Order) coreOrderModels.Order {
	switch rule.SupplementType {
	case models.ExceedOrderAmount:
		if rule.OrderAmount == 0 {
			order = svc.addSupplementProductsSeveralTimes(rule.SupplementaryProducts, 1, order)
		} else {
			count := int(order.EstimatedTotalPrice.Value) / rule.OrderAmount
			order = svc.addSupplementProductsSeveralTimes(rule.SupplementaryProducts, count, order)
		}
	}

	return order
}

func (svc *serviceImpl) addSupplementProductsSeveralTimes(supplementProducts []models.SupplementaryProduct, count int, order coreOrderModels.Order) coreOrderModels.Order {
	for _, supplementProduct := range supplementProducts {
		order.Products = append(order.Products, coreOrderModels.OrderProduct{
			ID:   supplementProduct.ProductId,
			Name: supplementProduct.Name,
			Price: coreOrderModels.Price{
				Value: supplementProduct.Price,
			},
			Quantity: supplementProduct.Quantity * count,
		})
	}

	return order
}
