package rkeeper

import (
	"github.com/kwaaka-team/orders-core/core/menu/models"
	coreModels "github.com/kwaaka-team/orders-core/core/models"
	storeModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"strconv"
	"time"

	rkeeperModels "github.com/kwaaka-team/orders-core/pkg/rkeeperwhite/clients/dto"
	"github.com/rs/zerolog/log"
)

func menuFromClient(req rkeeperModels.Menu, productsExist map[string]string, posProducts map[string]models.Product, stopList rkeeperModels.StopListResponse, store storeModels.Store) models.Menu {

	menu := models.Menu{
		Name:             models.RKEEPER.String(),
		ExtName:          models.MAIN.String(),
		Description:      "rkeeper pos menu",
		AttributesGroups: attributeGroupsToModel(req),
		Products:         productsToModel(req, productsExist, store),
		Groups:           groupsToModel(req.Categories),
		CreatedAt:        coreModels.TimeNow(),
		UpdatedAt:        coreModels.TimeNow(),
	}
	menu.Products = append(menu.Products, existProductsByStopList(posProducts, stopList)...).Unique()
	menu.Attributes = attributesToModel(req.Ingredients, menu.AttributesGroups)
	menu.Sections = addSectionOrderForSections(activeSections(menu.Products, splitGroups(menu.Groups)).Unique())

	return menu
}

func existProductsByStopList(posProducts map[string]models.Product, stopList rkeeperModels.StopListResponse) models.Products {
	products := make(models.Products, 0, len(stopList.TaskResponse.StopList.Dishes))

	for _, dish := range stopList.TaskResponse.StopList.Dishes {
		if val, ok := posProducts[dish.ID]; ok {
			val.IsAvailable = false
			products = append(products, val)
		}
	}

	return products
}

func attributesToModel(ingredients rkeeperModels.Ingredients, attributeGroups models.AttributeGroups) models.Attributes {

	attrGroupsExist := make(map[string]models.AttributeGroup, len(attributeGroups))

	for _, attrGroup := range attributeGroups {
		for _, ingredient := range attrGroup.Attributes {
			attrGroupsExist[ingredient] = attrGroup
		}
	}

	res := make(models.Attributes, 0, len(ingredients))

	for _, attribute := range ingredients {
		attr, err := attributeToModel(attribute)
		if err != nil {
			log.Err(err).Msg("rkeeper cli err: get attribute")
			continue
		}

		if attrGroup, ok := attrGroupsExist[attribute.ID]; ok {
			attr.ParentAttributeGroup = attrGroup.ExtID
			attr.HasAttributeGroup = true
			attr.AttributeGroupName = attrGroup.Name
			attr.AttributeGroupMin = attrGroup.Min
			attr.AttributeGroupMax = attrGroup.Max
		}

		res = append(res, attr)
	}

	return res
}

func attributeToModel(attributes rkeeperModels.Ingredient) (models.Attribute, error) {

	res := models.Attribute{
		ExtID:       attributes.ID,
		Name:        attributes.Name,
		ExtName:     attributes.Description,
		IsAvailable: true,
		UpdatedAt:   time.Now(),
	}

	price, err := strconv.ParseFloat(attributes.Price, 64)
	if err != nil {
		return models.Attribute{}, err
	}

	res.Price = price

	return res, nil
}

func attributeGroupsToModel(req rkeeperModels.Menu) models.AttributeGroups {

	groupsInSchemes := make(map[string]rkeeperModels.IngredientsSchemeGroup, len(req.IngredientsSchemes))
	for _, scheme := range req.IngredientsSchemes {
		if scheme.IngredientsGroups != nil {
			for _, attributeGroups := range scheme.IngredientsGroups {
				groupsInSchemes[attributeGroups.ID] = attributeGroups
			}
		}
	}

	res := make(models.AttributeGroups, 0, len(req.IngredientsGroups))
	for _, ingredientGroup := range req.IngredientsGroups {

		attributeGroup := attributeGroupToModel(ingredientGroup)

		if data, ok := groupsInSchemes[attributeGroup.ExtID]; ok {
			attributeGroup.Min = data.MinCount
			attributeGroup.Max = data.MaxCount
		}

		res = append(res, attributeGroup)
	}
	return res
}

func attributeGroupToModel(req rkeeperModels.IngredientsGroup) models.AttributeGroup {
	return models.AttributeGroup{
		ExtID:      req.ID,
		Name:       req.Name,
		Attributes: req.Ingredients,
	}
}

func groupsToModel(req rkeeperModels.Categories) models.Groups {

	res := make(models.Groups, 0, len(req))
	for _, group := range req {
		res = append(res, groupToModel(group))
	}
	return res
}

func groupToModel(req rkeeperModels.Category) models.Group {
	return models.Group{
		ID:          req.ID,
		Name:        req.Name,
		ParentGroup: req.ParentId,
	}
}

func getSections(req rkeeperModels.Products) (models.Sections, map[string]struct{}) {

	sectionExist := make(map[string]struct{}, len(req))
	sections := make(models.Sections, 0, len(req))

	for _, product := range req {
		if product.SchemeId == "" {
			continue
		}
		sectionExist[product.SchemeId] = struct{}{}
		sections = append(sections, models.Section{
			ExtID: product.ID,
		})

	}

	return sections, sectionExist
}

func getSection(groups models.Groups) models.Sections {

	var sections models.Sections
	for _, group := range groups {
		section := models.Section{
			ExtID:      group.ID,
			Name:       group.Name,
			Collection: group.ParentGroup,
		}
		sections = append(sections, section)
	}
	return sections
}

func splitGroups(groups models.Groups) models.Sections {

	var groupsForSection models.Groups

	for _, group := range groups {
		if group.ParentGroup != "" {
			groupsForSection = append(groupsForSection, group)
		}
	}

	sections := getSection(groupsForSection)

	return sections
}

func activeSections(products models.Products, sections models.Sections) models.Sections {

	var (
		activeSections []string
		resultSections models.Sections
		sectionMap     = make(map[string]models.Section, len(sections))
	)

	for _, product := range products {
		activeSections = append(activeSections, product.Section)
	}

	for _, section := range sections {
		sectionMap[section.ExtID] = section
	}

	for extID := range activeSections {
		resultSections = append(resultSections, sectionMap[activeSections[extID]])
	}

	return resultSections
}

func addSectionOrderForSections(sections models.Sections) models.Sections {

	for i := range sections {
		sections[i].SectionOrder = i + 1
	}
	return sections
}
