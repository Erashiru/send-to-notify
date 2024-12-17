package pos

import (
	"github.com/kwaaka-team/orders-core/config/menu"
	"github.com/kwaaka-team/orders-core/core/menu/clients/pos/base"
	"github.com/kwaaka-team/orders-core/core/menu/clients/pos/iiko"
	"github.com/kwaaka-team/orders-core/core/menu/clients/pos/jowi"
	"github.com/kwaaka-team/orders-core/core/menu/clients/pos/paloma"
	"github.com/kwaaka-team/orders-core/core/menu/clients/pos/poster"
	"github.com/kwaaka-team/orders-core/core/menu/clients/pos/rkeeper"
	rkeeper_xml "github.com/kwaaka-team/orders-core/core/menu/clients/pos/rkeeper7xml"
	"github.com/kwaaka-team/orders-core/core/menu/clients/pos/tillypad"
	"github.com/kwaaka-team/orders-core/core/menu/clients/pos/yaros"
	"github.com/kwaaka-team/orders-core/core/menu/clients/pos/ytimes"
	"github.com/kwaaka-team/orders-core/core/menu/database/drivers"
	"github.com/kwaaka-team/orders-core/core/menu/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
)

func NewPosManager(globalConfig menu.Configuration, menuRepo drivers.MenuRepository, store storeModels.Store) (base.Manager, error) {
	switch models.PosName(store.PosType) {
	case models.IIKO, models.SYRVE:
		return iiko.NewIIKOManager(globalConfig, menuRepo, store)
	case models.RKEEPER:
		return rkeeper.NewManager(globalConfig, menuRepo, store)
	case models.PALOMA:
		return paloma.NewManager(globalConfig, menuRepo, store)
	case models.POSTER:
		return poster.NewManager(globalConfig, menuRepo, store)
	case models.JOWI:
		return jowi.NewJowiManager(globalConfig, menuRepo, store)
	case models.YAROS:
		return yaros.NewManager(globalConfig, menuRepo, store)
	case models.RKEEPER7XML:
		return rkeeper_xml.NewManager(globalConfig, menuRepo, store)
	case models.Tillypad:
		return tillypad.NewTillypadManager(globalConfig, menuRepo, store)
	case models.Ytimes:
		return ytimes.NewYtimesManager(globalConfig.Ytimes.BaseUrl, store)

	}
	return nil, ErrPosNotFound
}
