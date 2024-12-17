package menu

import coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"

/*
MapPosProductsToAggregatorProducts

Multiple products in POS can have the same ID.
Therefore, one product in POS can be represented as multiple products in kwaaka:
with same ProductID but with different ExtIDs
*/
func (s *Service) MapPosProductsToAggregatorProducts(posMenu coreMenuModels.Menu, aggregatorMenu coreMenuModels.Menu, posMenuProduct []coreMenuModels.Product) []coreMenuModels.Product {
	extIDs := s.mapPosIDToExtIDsFromPosMenuProduct(posMenu, posMenuProduct)
	return s.getAggregatorMenuProductsByExtIDs(aggregatorMenu, extIDs)
}

/*
MapPosProductToAggregatorProduct

Multiple products in POS can have the same ID.
Therefore, one product in POS can be represented as multiple products in kwaaka:
with same ProductID but with different ExtIDs
*/
func (s *Service) MapPosProductToAggregatorProduct(posMenu coreMenuModels.Menu, aggregatorMenu coreMenuModels.Menu, posMenuProduct coreMenuModels.Product) []coreMenuModels.Product {
	extIDs := s.mapPosIDToExtIDsFromPosMenuProduct(posMenu, []coreMenuModels.Product{posMenuProduct})
	return s.getAggregatorMenuProductsByExtIDs(aggregatorMenu, extIDs)
}

func (s *Service) mapPosIDToExtIDsFromPosMenuProduct(posMenu coreMenuModels.Menu, posMenuProducts []coreMenuModels.Product) []string {

	productIDs := make(map[string]struct{})
	for i := range posMenuProducts {
		posMenuProduct := posMenuProducts[i]
		productIDs[posMenuProduct.ProductID] = struct{}{}
	}

	result := make([]string, 0)
	for i := range posMenu.Products {
		posProduct := posMenu.Products[i]
		if _, ok := productIDs[posProduct.ProductID]; !ok {
			continue
		}
		extID := s.getExtIDFromPosProduct(posProduct)
		result = append(result, extID)
	}

	return result
}

func (s *Service) getExtIDFromPosProduct(posProduct coreMenuModels.Product) string {
	if posProduct.ExtID != "" {
		return posProduct.ExtID
	}
	return posProduct.ProductID
}

func (s *Service) GetSystemIDFromPosProduct(posProduct coreMenuModels.Product) string {
	return s.getExtIDFromPosProduct(posProduct)
}

func (s *Service) GetPosIDFromPosProduct(posProduct coreMenuModels.Product) string {
	return posProduct.ProductID
}

func (s *Service) GetSystemIDFromAggregatorProduct(aggProduct coreMenuModels.Product) string {
	return aggProduct.PosID
}

func (s *Service) getAggregatorMenuProductsByExtIDs(aggregatorMenu coreMenuModels.Menu, extIDs []string) []coreMenuModels.Product {
	extIDsMap := s.toMap(extIDs)
	result := make([]coreMenuModels.Product, 0)
	for i := range aggregatorMenu.Products {
		aggProduct := aggregatorMenu.Products[i]
		posID := aggProduct.PosID
		if _, ok := extIDsMap[posID]; !ok {
			continue
		}
		result = append(result, aggProduct)
	}
	return result
}

func (s *Service) toMap(items []string) map[string]struct{} {
	result := make(map[string]struct{})
	for i := range items {
		extID := items[i]
		result[extID] = struct{}{}
	}
	return result
}

func (s *Service) GetAggrMenuProductsByIDs(aggregatorMenu coreMenuModels.Menu, extIDs []string) []coreMenuModels.Product {
	return s.getAggregatorMenuProductsByExtIDs(aggregatorMenu, extIDs)
}
