package managers

import (
	"context"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
)

type MSPositions interface {
	GetPositions(ctx context.Context, query selector.MoySklad) ([]models.MoySkladPosition, error)
	RemovePosition(ctx context.Context, query selector.MoySklad) error
	InsertDB(ctx context.Context, req models.MoySkladPosition) error
}

type msp struct {
	mspRepo drivers.MSPositionsRepository
}

func (m msp) GetPositions(ctx context.Context, query selector.MoySklad) ([]models.MoySkladPosition, error) {
	return m.mspRepo.GetPositions(ctx, query)
}

func (m msp) RemovePosition(ctx context.Context, query selector.MoySklad) error {
	return m.mspRepo.RemovePosition(ctx, query)
}

func (m msp) InsertDB(ctx context.Context, req models.MoySkladPosition) error {
	err := m.mspRepo.CreatePosition(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func NewMSPositionsManager(
	mspRepo drivers.MSPositionsRepository,
) MSPositions {
	return &msp{
		mspRepo: mspRepo,
	}
}
