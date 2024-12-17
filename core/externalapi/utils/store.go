package utils

import (
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
)

func ParseStores(req []coreStoreModels.Store, service string) models.GetStoreResponse {
	response := models.GetStoreResponse{
		Places: []models.Place{},
	}

	for _, store := range req {
		place := models.Place{
			Title:   store.Name,
			Address: store.Address.Street,
		}

		for _, externalStore := range store.ExternalConfig {
			if externalStore.Type == service || (service == "yandex" && externalStore.Type == "emenu") {
				if len(externalStore.StoreID) != 0 {
					place.Id = externalStore.StoreID[0]
				}
			}
		}

		response.Places = append(response.Places, place)
	}

	return response
}
