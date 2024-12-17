package jowi

import (
	"context"
	"github.com/kwaaka-team/orders-core/config/menu"
	"github.com/kwaaka-team/orders-core/core/menu/clients/pos/base"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	pkg "github.com/kwaaka-team/orders-core/pkg/jowi"
	jowiClient "github.com/kwaaka-team/orders-core/pkg/jowi/client"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type manager struct {
	cli          pkg.Jowi
	restaurantID string
	globalConfig menu.Configuration
	menuRepo     drivers.MenuRepository
}

func NewJowiManager(
	conf menu.Configuration,
	menuRepo drivers.MenuRepository,
	store storeModels.Store) (base.Manager, error) {

	cli, err := jowiClient.New(pkg.Config{
		ApiKey:    store.Jowi.ApiKey,
		ApiSecret: conf.JowiConfiguration.ApiSecret,
		BaseURL:   conf.JowiConfiguration.BaseURL,
		Protocol:  "http",
	})

	if err != nil {
		log.Trace().Err(err).Msg("can't initialize Jowi Client")
		return nil, err
	}

	return &manager{
		cli:          cli,
		restaurantID: store.Jowi.RestaurantID,
		globalConfig: conf,
		menuRepo:     menuRepo,
	}, nil
}

func (man manager) GetAggMenu(ctx context.Context, store storeModels.Store) ([]models.Menu, error) {
	return nil, errors.New("method not implemented")
}

func (m manager) GetMenu(ctx context.Context, store storeModels.Store) (models.Menu, error) {
	courses, err := m.cli.GetCourses(ctx, m.restaurantID)
	if err != nil {
		log.Trace().Err(err).Msgf("Get courses error")
		return models.Menu{}, err
	}

	sections, collections, err := m.getCategories(ctx)
	if err != nil {
		return models.Menu{}, err
	}

	stopList, err := m.getStopList(ctx)
	if err != nil {
		return models.Menu{}, err
	}

	return menuFromClient(courses, store.Settings, sections, collections, stopList), nil
}

func (m manager) getStopList(ctx context.Context) (models.StopListItems, error) {
	items, err := m.cli.GetStopList(ctx, m.restaurantID)
	if err != nil {
		log.Trace().Err(err).Msgf("Get stop list error")
		return models.StopListItems{}, err
	}

	return stopListFromClient(items)
}

func (m manager) getCategories(ctx context.Context) ([]models.Section, []models.MenuCollection, error) {
	courseCategories, err := m.cli.GetCourseCategories(ctx, m.restaurantID)
	if err != nil {
		log.Trace().Err(err).Msgf("Get course categories error")
		return nil, nil, err
	}

	return categoryFromClient(courseCategories)
}
