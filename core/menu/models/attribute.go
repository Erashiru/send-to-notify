package models

import (
	"time"
)

type Attributes []Attribute

func (attr Attributes) ListIDs() []string {
	res := make([]string, 0, len(attr))

	for _, attribute := range attr {
		res = append(res, attribute.ExtID)
	}

	return res
}

type Attribute struct {
	ExtID                string                 `bson:"ext_id" json:"ext_id"`
	PosID                string                 `bson:"pos_id" json:"pos_id"`
	IngredientID         string                 `bson:"ingredient_id" json:"ingredient_id"`
	StarterAppID         string                 `bson:"starter_app_id" json:"starter_app_id"`
	StarterAppOfferID    string                 `bson:"starter_app_offer_id" json:"starter_app_offer_id"`
	ChocofoodFoodId      string                 `bson:"chocofood_food_id" json:"chocofood_food_id"`
	Code                 string                 `bson:"code" json:"code"`
	ExtName              string                 `bson:"ext_name" json:"ext_name"`
	Name                 string                 `bson:"name" json:"name"`
	Default              bool                   `bson:"selected_by_default" json:"selected_by_default"`
	Price                float64                `bson:"price_impact" json:"price_impact"`
	IsAvailable          bool                   `bson:"available" json:"available"`
	IsDeleted            bool                   `bson:"is_deleted" json:"is_deleted"`
	IsDisabled           bool                   `bson:"is_disabled" json:"is_disabled"`
	IsIgnored            bool                   `bson:"is_ignored" json:"is_ignored"`
	IsSync               bool                   `bson:"sync" json:"sync"`
	IncludedInMenu       bool                   `bson:"included_in_menu" json:"included_in_menu"`
	ByAdmin              bool                   `bson:"by_admin" json:"by_admin"`
	Balance              float64                `bson:"balance" json:"balance"`
	Min                  int                    `bson:"min" json:"min"`
	Max                  int                    `bson:"max" json:"max"`
	ParentAttributeGroup string                 `bson:"parent_attribute_group" json:"parent_attribute_group"`
	AttributeGroupName   string                 `bson:"attribute_group_name_iiko" json:"attribute_group_name_iiko"`
	AttributeGroupExtID  string                 `bson:"attribute_group_ext_id" json:"attribute_group_ext_id"`
	HasAttributeGroup    bool                   `bson:"has_attribute_group" json:"has_attribute_group"`
	AttributeGroupMin    int                    `bson:"attribute_group_min" json:"attribute_group_min"`
	AttributeGroupMax    int                    `bson:"attribute_group_max" json:"attribute_group_max"`
	IsComboAttribute     bool                   `bson:"is_combo_attribute" json:"is_combo_attribute"`
	UpdatedAt            time.Time              `bson:"-" json:"updated_at,omitempty"`
	Description          []LanguageDescription  `bson:"description" json:"description"`
	AttributeGroups      []AttributeGroupObject `bson:"attribute_groups" json:"attribute_groups"`
	NamesByLanguage      []LanguageDescription  `bson:"names_by_language" json:"names_by_language"`
	ProductInformation   ProductInformation     `bson:"product_information" json:"product_information"`
}
type AttributeGroupObject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (p Attributes) Unique() Attributes {

	existProducts := make(map[string]struct{}, len(p))
	results := make(Attributes, 0, len(p))

	for _, attribute := range p {

		if _, ok := existProducts[attribute.ExtID]; ok {
			continue
		}

		results = append(results, attribute)
		existProducts[attribute.ExtID] = struct{}{}
	}

	return results
}

type AttributeGroups []AttributeGroup

type AttributeGroup struct {
	ExtID           string                `bson:"ext_id" json:"ext_id"`
	PosID           string                `bson:"pos_id" json:"pos_id"`
	StarterAppID    string                `bson:"starter_app_id" json:"starter_app_id"`
	Name            string                `bson:"name,omitempty" json:"name"`
	Max             int                   `bson:"max" json:"max"`
	Min             int                   `bson:"min" json:"min"`
	Collapse        bool                  `bson:"collapse" json:"collapse"`
	MultiSelection  bool                  `bson:"multiple_selection" json:"multi_selection"`
	Attributes      []string              `bson:"attributes" json:"attributes"`
	IsSync          bool                  `bson:"sync" json:"is_sync"`
	IsComboGroup    bool                  `bson:"is_combo_group" json:"is_combo_group"`
	Description     []LanguageDescription `bson:"description" json:"description"`
	AttributeObject []AttributeIdAndName  `bson:"attributeobject" json:"attribute_object,omitempty"`
	AttributeMinMax []AttributeIdMinMax   `bson:"attribute_min_max" json:"attribute_min_max"`
	NamesByLanguage []LanguageDescription `bson:"names_by_language" json:"names_by_language"`
}

type AttributeIdAndName struct {
	ExtId string `bson:"ext_id" json:"ext_id"`
	Name  string `bson:"name" json:"name,omitempty"`
}

type AttributeIdMinMax struct {
	ExtId string `bson:"ext_id" json:"ext_id"`
	Min   int    `bson:"min" json:"min"`
	Max   int    `bson:"max" json:"max"`
}

func (attrGroup AttributeGroups) Unique() AttributeGroups {

	existProducts := make(map[string]struct{}, len(attrGroup))
	results := make(AttributeGroups, 0, len(attrGroup))

	for _, attributeGroup := range attrGroup {

		if _, ok := existProducts[attributeGroup.ExtID]; ok {
			continue
		}

		results = append(results, attributeGroup)
		existProducts[attributeGroup.ExtID] = struct{}{}
	}

	return results
}

func (attr *AttributeGroup) AddAttributeNames(attributes []Attribute, attributeGroups []AttributeGroup) []AttributeGroup {
	var result []AttributeGroup
	mapAttributes := make(map[string]Attribute)
	for _, attribute := range attributes {
		mapAttributes[attribute.ExtID] = attribute
	}

	for index, group := range attributeGroups {
		var attrName []AttributeIdAndName
		for _, attributeId := range group.Attributes {
			if attribute, ok := mapAttributes[attributeId]; ok {
				attrTmp := AttributeIdAndName{
					Name:  attribute.Name,
					ExtId: attribute.ExtID,
				}
				attrName = append(attrName, attrTmp)
			}
		}
		attributeGroups[index].AttributeObject = attrName
		result = append(result, attributeGroups[index])
	}
	return result
}

func (attrGroup AttributeGroups) ListIDs() []string {
	res := make([]string, 0, len(attrGroup))

	for _, attributeGroup := range attrGroup {
		res = append(res, attributeGroup.ExtID)
	}

	return res
}

func (groups AttributeGroups) GetAttributeGroup(id string) AttributeGroup {
	for _, attributeGroup := range groups {
		if attributeGroup.ExtID == id {
			return attributeGroup
		}
	}
	return AttributeGroup{}
}

func (p *Attribute) RemoveDuplicate(attributes []Attribute) []Attribute {
	var unique []Attribute
Loop:
	for _, v := range attributes {
		for i, u := range unique {
			if v.ExtID == u.ExtID {
				unique[i] = v
				continue Loop
			}
		}
		unique = append(unique, v)
	}
	return unique
}

func (attr *Attribute) AddAttributeGroupNames(attributes []Attribute, attributeGroups []AttributeGroup) []Attribute {

	attributeGroupMap := make(map[string][]AttributeGroupObject)
	for _, attributeGroup := range attributeGroups {
		for _, attributeID := range attributeGroup.Attributes {
			attributeGroupMap[attributeID] = append(attributeGroupMap[attributeID], AttributeGroupObject{
				ID:   attributeGroup.ExtID,
				Name: attributeGroup.Name,
			})
		}
	}
	for index, attribute := range attributes {
		if attrGroups, ok := attributeGroupMap[attribute.ExtID]; ok {
			attributes[index].AttributeGroups = attrGroups
		}
	}
	return attributes
}
