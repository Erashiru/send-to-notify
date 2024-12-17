package managers

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
)

type StopListTransaction interface {
	Insert(ctx context.Context, transactions []models.StopListTransaction) ([]string, error)
}

type stm struct {
	stRepo drivers.StopListTransactionRepository
}

func NewStopListTransactionManager(
	stRepo drivers.StopListTransactionRepository,
) StopListTransaction {
	return &stm{
		stRepo: stRepo,
	}
}

// Insert temporary method to save transaction in DB
func (s stm) Insert(ctx context.Context, transactions []models.StopListTransaction) ([]string, error) {
	ids := make([]string, len(transactions))
	for _, transaction := range transactions {
		res, err := s.stRepo.Insert(ctx, transaction)
		if err != nil {
			return nil, err
		}
		ids = append(ids, res)
	}

	return ids, nil
}
