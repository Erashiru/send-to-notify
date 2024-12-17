package order

import (
	"context"
	"errors"
	ocErr "github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/core/models"
)

type CartService interface {
	GetQRMenuCartByID(ctx context.Context, cartID string) (models.Cart, error)
	GetKwaakaAdminCartByOrderID(ctx context.Context, cartID string) (models.Cart, error)
	GetOldQRMenuCartByID(ctx context.Context, cartID string) (models.OldCart, error)
	GetCartById(ctx context.Context, cartID string) (models.Cart, error)
}

type CartServiceImpl struct {
	cartRepo CartRepository
}

func NewCartService(cartRepo CartRepository) *CartServiceImpl {
	return &CartServiceImpl{
		cartRepo: cartRepo,
	}
}

func (s *CartServiceImpl) GetQRMenuCartByID(ctx context.Context, cartID string) (models.Cart, error) {
	return s.cartRepo.GetQRMenuCartByID(ctx, cartID)
}

func (s *CartServiceImpl) GetKwaakaAdminCartByOrderID(ctx context.Context, cartID string) (models.Cart, error) {
	return s.cartRepo.GetKwaakaAdminCartByCartID(ctx, cartID)
}

func (s *CartServiceImpl) GetOldQRMenuCartByID(ctx context.Context, cartID string) (models.OldCart, error) {
	return s.cartRepo.GetOldQRMenuCartByID(ctx, cartID)
}

func (s *CartServiceImpl) GetCartById(ctx context.Context, cartID string) (models.Cart, error) {
	cart, err := s.cartRepo.GetQRMenuCartByID(ctx, cartID)
	if err != nil {
		if errors.Is(err, ocErr.ErrNotFound) {
			cart, err = s.cartRepo.GetKwaakaAdminCartByCartID(ctx, cartID)
			if err != nil {
				return models.Cart{}, err
			}
			return cart, nil
		}
		return models.Cart{}, err
	}

	return cart, nil
}
