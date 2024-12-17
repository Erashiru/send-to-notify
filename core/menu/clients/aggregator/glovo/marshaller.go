package glovo

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
	glovoModels "github.com/kwaaka-team/orders-core/pkg/glovo/clients/dto"
	"github.com/rs/zerolog/log"
	"strings"
)

func toProducts(storeId string, req models.Products) []glovoModels.Product {
	products := make([]glovoModels.Product, 0, len(req))

	for i := range req {
		log.Info().Msgf("[STORE %v BULK] PRODUCT_ID: %v, AVAILABLE: %v", storeId, req[i].ExtID, req[i].IsAvailable)
		products = append(products, toProduct(req[i]))
	}

	return products
}

func toProduct(req models.Product) glovoModels.Product {

	product := glovoModels.Product{
		ID:               req.ExtID,
		Available:        &req.IsAvailable,
		AttributesGroups: req.AttributesGroups,
		Description:      "",
		ExtraImageUrls:   []string{},
		Restrictions:     &glovoModels.Restriction{IsTobacco: req.IsTobacco, IsAlcoholic: req.IsAlcohol},
	}
	if len(product.AttributesGroups) == 0 {
		product.AttributesGroups = []string{}
	}

	if len(req.Name) != 0 {
		product.Name = req.Name[0].Value
	}

	if len(req.Price) != 0 {
		product.Price = req.Price[0].Value
	}

	if len(req.ImageURLs) != 0 {
		product.ImageURL = strings.TrimSpace(req.ImageURLs[0])
	}

	if len(req.ImageURLs) > 1 {
		product.ExtraImageUrls = req.ImageURLs[1:]
	}

	if len(req.Description) != 0 {
		product.Description = req.Description[0].Value
	}

	return product
}

func toAttributes(storeId string, req models.Attributes) []glovoModels.Attribute {

	attributes := make([]glovoModels.Attribute, 0, len(req))
	for i := range req {
		log.Info().Msgf("[STORE %v BULK] ATTRIBUTE_ID: %v, AVAILABLE: %v", storeId, req[i].ExtID, req[i].IsAvailable)
		attributes = append(attributes, toAttribute(req[i]))
	}

	return attributes
}

func toAttribute(req models.Attribute) glovoModels.Attribute {
	return glovoModels.Attribute{
		ID:                req.ExtID,
		Name:              req.Name,
		PriceImpact:       req.Price,
		Available:         req.IsAvailable,
		SelectedByDefault: req.Default,
	}
}

func ToGlovoMenu(req models.Menu, offers []models.BkOffers) glovoModels.ValidateMenuRequest {
	superCollectionsMap := make(map[string]models.MenuSuperCollection, len(req.SuperCollections))
	for _, sc := range req.SuperCollections {
		superCollectionsMap[sc.ExtID] = sc
	}
	collectionsMap := make(map[string]models.MenuCollection, len(req.Collections))
	for _, c := range req.Collections {
		collectionsMap[c.ExtID] = c
	}
	sectionMap := make(map[string]models.Section, len(req.Sections))
	for _, s := range req.Sections {
		sectionMap[s.ExtID] = s
	}
	offersMap := make(map[string]models.BkOffers, len(offers))
	for _, o := range offers {
		offersMap[o.ProductID] = o
	}

	resAttributes := make([]glovoModels.Attribute, 0, len(req.Attributes))
	for _, a := range req.Attributes {
		resAttributes = append(resAttributes, toAttribute(a))
	}

	resAttributeGroups := make([]glovoModels.AttributeGroup, 0, len(req.AttributesGroups))
	for _, ag := range req.AttributesGroups {
		resAttributeGroups = append(resAttributeGroups, toAttributeGroup(ag))
	}

	tempSections := make(map[string]glovoModels.Section)
	tempCollections := make(map[string]glovoModels.Collection)
	tempSuperCollections := make(map[string]glovoModels.SuperCollection)

	resProducts := make([]glovoModels.Product, 0, len(req.Products))
	for _, p := range req.Products {
		if !isValidProduct(p) {
			continue
		}
		if _, ok := sectionMap[p.Section]; !ok {
			continue
		}
		if sectionMap[p.Section].IsDeleted {
			continue
		}

		tempProduct := toProduct(p)
		if val, ok := offersMap[tempProduct.ID]; ok {
			tempProduct.Price = float64(val.GlovoPrice)
		}

		resProducts = append(resProducts, tempProduct)
		tempSection, ok := tempSections[p.Section]
		if !ok {
			tempSections[p.Section] = glovoModels.Section{
				ID:       p.Section,
				Name:     sectionMap[p.Section].Name,
				Position: sectionMap[p.Section].SectionOrder,
				Products: []string{p.ExtID},
			}
			continue
		}
		tempSection.Products = append(tempSections[p.Section].Products, p.ExtID)
		tempSections[tempSection.ID] = tempSection
	}

	for _, s := range tempSections {
		sectionMapVal, ok := sectionMap[s.ID]
		if !ok {
			continue
		}
		if sectionMapVal.Collection == "" {
			continue
		}

		collectionsMapVal, ok := collectionsMap[sectionMapVal.Collection]
		if !ok {
			continue
		}
		tempCollection, ok := tempCollections[sectionMapVal.Collection]
		if ok {
			tempCollection.Sections = append(tempCollection.Sections, s)
			tempCollections[tempCollection.ID] = tempCollection
			continue
		}

		tempCol := glovoModels.Collection{
			ID:       sectionMapVal.Collection,
			Name:     collectionsMapVal.Name,
			Position: collectionsMapVal.CollectionOrder,
			ImageUrl: collectionsMapVal.ImageURL,
			Sections: []glovoModels.Section{s},
		}

		scheduleAvailabilities := toCollectionAvailabilities(collectionsMapVal.Schedule.Availabilities)
		if len(scheduleAvailabilities) > 0 {
			tempCol.Schedule = &glovoModels.Schedule{
				ID:             collectionsMapVal.Schedule.ID,
				Name:           collectionsMapVal.Schedule.Name,
				Availabilities: toCollectionAvailabilities(collectionsMapVal.Schedule.Availabilities),
			}
		}

		tempCollections[sectionMapVal.Collection] = tempCol
	}

	for _, c := range tempCollections {
		collectionsMapVal, ok := collectionsMap[c.ID]
		if !ok {
			continue
		}
		if collectionsMapVal.SuperCollection == "" {
			continue
		}

		superCollectionsMapVal, ok := superCollectionsMap[collectionsMapVal.SuperCollection]
		if !ok {
			continue
		}
		tempSuperCollection, ok := tempSuperCollections[collectionsMapVal.SuperCollection]
		if !ok {
			tempSuperCollections[collectionsMapVal.SuperCollection] = glovoModels.SuperCollection{
				ID:          collectionsMapVal.SuperCollection,
				Name:        superCollectionsMapVal.Name,
				Position:    superCollectionsMapVal.SuperCollectionOrder,
				ImageUrl:    superCollectionsMapVal.ImageUrl,
				Collections: []string{c.Name},
			}
			continue
		}
		tempSuperCollection.Collections = append(tempSuperCollection.Collections, c.Name)
		tempSuperCollections[collectionsMapVal.SuperCollection] = tempSuperCollection
	}

	resCollections := make([]glovoModels.Collection, 0, len(req.Collections))
	for _, c := range req.Collections {
		if c.ExtID == "" {
			continue
		}
		if val, ok := tempCollections[c.ExtID]; ok {
			resCollections = append(resCollections, val)
		}
	}

	resSuperColletions := make([]glovoModels.SuperCollection, 0, len(req.SuperCollections))
	for _, sc := range req.SuperCollections {
		if sc.ExtID == "" {
			continue
		}
		if val, ok := tempSuperCollections[sc.ExtID]; ok {
			tempRes := glovoModels.SuperCollection{
				ID:          val.ID,
				Name:        val.Name,
				Position:    val.Position,
				Collections: val.Collections,
			}
			if val.ImageUrl != "" {
				tempRes.ImageUrl = val.ImageUrl
			}
			resSuperColletions = append(resSuperColletions, tempRes)
		}
	}

	return glovoModels.ValidateMenuRequest{
		Attributes:       resAttributes,
		AttributeGroups:  resAttributeGroups,
		Products:         resProducts,
		Collections:      resCollections,
		SuperCollections: resSuperColletions,
	}
}

func toCollectionAvailabilities(req []models.Availability) []glovoModels.Availability {
	res := make([]glovoModels.Availability, 0, len(req))

	for _, v := range req {
		res = append(res, glovoModels.Availability{
			Day:       v.Day,
			TimeSlots: toAvalabilityTimeSlots(v.TimeSlots),
		})
	}
	return res
}

func toAvalabilityTimeSlots(req []models.TimeSlot) []glovoModels.TimeSlot {
	res := make([]glovoModels.TimeSlot, 0, len(req))

	for _, v := range req {
		res = append(res, glovoModels.TimeSlot{
			Start: v.Start,
			End:   v.End,
		})
	}
	return res
}

func toAttributeGroup(req models.AttributeGroup) glovoModels.AttributeGroup {
	return glovoModels.AttributeGroup{
		ID:                req.ExtID,
		Name:              req.Name,
		Min:               req.Min,
		Max:               req.Max,
		Collapse:          req.Collapse,
		MultipleSelection: req.MultiSelection,
		Attributes:        req.Attributes,
	}
}

func isValidProduct(p models.Product) bool {
	if p.IsDeleted {
		return false
	}
	if p.ExtID == "" {
		return false
	}
	if len(p.Name) == 0 {
		return false
	}
	if len(p.Price) == 0 {
		return false
	}
	return true
}
