package models

type Category struct {
	Id        string          `json:"id"`
	ParentId  string          `json:"parentId,omitempty"`
	Name      string          `json:"name"`
	SortOrder int             `json:"sortOrder,omitempty"`
	Images    []CategoryImage `json:"images,omitempty"`
}

type ItemImage struct {
	Url  string `json:"url"`
	Hash string `json:"hash"`
}

type Nutrients struct {
	Calories      float64 `json:"calories,omitempty"`
	Carbohydrates float64 `json:"carbohydrates,omitempty"`
	Fat           float64 `json:"fats,omitempty"`
	IsDeactivated bool    `json:"is_deactivated,omitempty"`
	Proteins      float64 `json:"proteins,omitempty"`
}

type CategoryImage struct {
	Url       string `json:"url"`
	UpdatedAt string `json:"updatedAt"`
}

type Item struct {
	Id                     string                                    `json:"id"`
	CategoryId             string                                    `json:"categoryId"`
	Name                   string                                    `json:"name"`
	Description            string                                    `json:"description,omitempty"`
	Price                  float64                                   `json:"price"`
	Vat                    int                                       `json:"vat,omitempty"`
	IsCatchWeight          bool                                      `json:"isCatchweight,omitempty"`
	Measure                int                                       `json:"measure"`
	WeightQuantum          float64                                   `json:"weightQuantum,omitempty"`
	MeasureUnit            string                                    `json:"measureUnit"`
	SortOrder              int                                       `json:"sortOrder,omitempty"`
	Nutrients              Nutrients                                 `json:"nutrients,omitempty"`
	ModifierGroups         []ModifierGroup                           `json:"modifierGroups,omitempty"`
	Images                 []ItemImage                               `json:"images,omitempty"`
	AdditionalDescriptions MenuCompositionItemAdditionalDescriptions `json:"additional_descriptions,omitempty"`
}

type ModifierGroup struct {
	Id                   string     `json:"id"`
	Name                 string     `json:"name"`
	Modifiers            []Modifier `json:"modifiers,omitempty"`
	MinSelectedModifiers int        `json:"minSelectedModifiers"`
	MaxSelectedModifiers int        `json:"maxSelectedModifiers"`
	SortOrder            int        `json:"sortOrder,omitempty"`
}

type Modifier struct {
	Id        string  `json:"id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Vat       int     `json:"vat,omitempty"`
	MinAmount int     `json:"minAmount"`
	MaxAmount int     `json:"maxAmount"`
}

type Menu struct {
	Categories []Category `json:"categories"`
	Items      []Item     `json:"items"`
	LastChange string     `json:"lastChange"`
}

type MenuCompositionItemAdditionalDescriptions struct {
	Badges []FoodSpecifics `json:"badges,omitempty"`
}

type FoodSpecifics struct {
	Category string `json:"category,omitempty"`
	Value    string `json:"value,omitempty"`
}
