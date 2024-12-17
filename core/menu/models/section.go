package models

import (
	"time"
)

type Sections []Section

type Section struct {
	ExtID           string                `bson:"ext_id" json:"id"`
	StarterAppID    string                `bson:"starter_app_id" json:"starter_app_id"`
	ChocofoodFoodId string                `bson:"chocofood_food_id" json:"chocofood_food_id"`
	LastID          string                `bson:"last_id,omitempty" json:"last_id,omitempty"` // fixme: for what?
	Name            string                `bson:"name" json:"name"`
	SectionOrder    int                   `bson:"section_order" json:"section_order"`
	Collection      string                `bson:"collection" json:"collection"`
	Description     []LanguageDescription `bson:"description" json:"description"` // fixme: why?
	ImageUrl        string                `bson:"image_url" json:"image_url"`
	ImageUpdatedAt  time.Time             `bson:"image_updated_at" json:"image_updated_at"`
	IsDeleted       bool                  `bson:"is_deleted" json:"is_deleted"`
	Amount          int                   `bson:"amount" json:"amount"`
	NamesByLanguage []LanguageDescription `bson:"names_by_language" json:"names_by_language"`
}

type LanguageDescription struct {
	LanguageCode string `bson:"language_code" json:"language_code"`
	Value        string `bson:"value" json:"value"`
}

func (lp LanguageDescription) Get(name []LanguageDescription) string {
	if len(name) != 0 {
		return name[0].Value
	}
	return ""
}

func (s Sections) GetIndex(id string) (int, bool) {
	for i := range s {
		if s[i].ExtID == id {
			return i, true
		}
	}
	return 0, false
}

func (s Sections) Len() int           { return len(s) }
func (s Sections) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Sections) Less(i, j int) bool { return s[i].SectionOrder < s[j].SectionOrder }

func (s Sections) Unique() Sections {

	existSections := make(map[string]struct{}, len(s))
	result := make(Sections, 0, len(s))

	for _, section := range s {
		if _, ok := existSections[section.ExtID]; ok {
			continue
		}
		if section.ExtID == "" {
			continue
		}
		result = append(result, section)
		existSections[section.ExtID] = struct{}{}
	}
	return result
}
