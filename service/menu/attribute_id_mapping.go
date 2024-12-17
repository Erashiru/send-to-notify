package menu

import coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"

func (s *Service) MapPosAttributesToAggregatorAttributes(posMenu coreMenuModels.Menu, aggregatorMenu coreMenuModels.Menu, posMenuAttributes []coreMenuModels.Attribute) []coreMenuModels.Attribute {
	extIDs := s.mapPosIDToSystemIDsFromPosMenuAttribute(posMenu, posMenuAttributes)
	return s.getAggregatorMenuAttributesBySystemIDs(aggregatorMenu, extIDs)
}

func (s *Service) GetSystemIDFromPosAttribute(posAttribute coreMenuModels.Attribute) string {
	return posAttribute.ExtID
}

func (s *Service) GetPosIDFromPosAttribute(posAttribute coreMenuModels.Attribute) string {
	if posAttribute.PosID != "" {
		return posAttribute.PosID
	}
	return posAttribute.ExtID
}

func (s *Service) GetSystemIDFromAggregatorAttribute(aggAttribute coreMenuModels.Attribute) string {
	if aggAttribute.PosID != "" {
		return aggAttribute.PosID
	}
	return aggAttribute.ExtID
}

func (s *Service) GetAggregatorIDFromAggregatorAttribute(aggAttribute coreMenuModels.Attribute) string {
	return aggAttribute.ExtID
}

func (s *Service) mapPosIDToSystemIDsFromPosMenuAttribute(posMenu coreMenuModels.Menu, posMenuAttributes []coreMenuModels.Attribute) []string {
	attributeIDs := make(map[string]struct{})
	for i := range posMenuAttributes {
		posMenuAttribute := posMenuAttributes[i]
		posID := s.GetPosIDFromPosAttribute(posMenuAttribute)
		attributeIDs[posID] = struct{}{}
	}

	result := make([]string, 0)
	for i := range posMenu.Attributes {
		posAttribute := posMenu.Attributes[i]
		posID := s.GetPosIDFromPosAttribute(posAttribute)
		if _, ok := attributeIDs[posID]; !ok {
			continue
		}
		systemID := s.GetSystemIDFromPosAttribute(posAttribute)
		result = append(result, systemID)
	}

	return result
}

func (s *Service) getAggregatorMenuAttributesBySystemIDs(aggregatorMenu coreMenuModels.Menu, systemIDs []string) []coreMenuModels.Attribute {
	systemIDsMap := s.toMap(systemIDs)
	result := make([]coreMenuModels.Attribute, 0)
	for i := range aggregatorMenu.Attributes {
		aggAttribute := aggregatorMenu.Attributes[i]
		systemID := s.GetSystemIDFromAggregatorAttribute(aggAttribute)
		if _, ok := systemIDsMap[systemID]; !ok {
			continue
		}
		result = append(result, aggAttribute)
	}
	return result
}
