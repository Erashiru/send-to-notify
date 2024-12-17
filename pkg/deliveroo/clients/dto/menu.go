package dto

type Menu struct {
	Name     string   `json:"name"`
	MenuData MenuData `json:"menu"`
	SiteIDs  []string `json:"site_ids"`
}

type MenuData struct {
	Categories []Category  `json:"categories"`
	Items      []MenuItem  `json:"items"`
	Mealtimes  []Mealtimes `json:"mealtimes"`
	Modifiers  []Modifier  `json:"modifiers"`
}

type Category struct {
	Description Description `json:"description"`
	ID          string      `json:"id"`
	ItemIDs     []string    `json:"item_ids"`
	Name        Name        `json:"name"`
}

type Description struct {
	EN string `json:"en"`
}

type Name struct {
	EN string `json:"en"`
}

type MenuItem struct {
	Allergies                 []string        `json:"allergies"`
	Classifications           []string        `json:"classifications"`
	ContainsAlcohol           bool            `json:"contains_alcohol"`
	Description               Description     `json:"description"`
	Diets                     []string        `json:"diets"`
	ExternalData              string          `json:"external_data"`
	Highlights                []string        `json:"highlights"`
	IAN                       string          `json:"ian"`
	ID                        string          `json:"id"`
	Image                     ItemImage       `json:"image"`
	IsEligibleAsReplacement   bool            `json:"is_eligible_as_replacement"`
	IsEligibleForSubstitution bool            `json:"is_eligible_for_substitution"`
	MaxQuantity               int             `json:"max_quantity"`
	ModifierIDs               []string        `json:"modifier_ids"`
	Name                      Name            `json:"name"`
	NutritionalInfo           NutritionalInfo `json:"nutritional_info"`
	OperationalName           string          `json:"operational_name"`
	PLU                       string          `json:"plu"`
	PriceInfo                 PriceInfo       `json:"price_info"`
	TaxRate                   string          `json:"tax_rate"`
	Type                      string          `json:"type"`
}

type ItemImage struct {
	URL string `json:"url"`
}

type NutritionalInfo struct {
	EnergyKcal NutritionalInfoKcal `json:"energy_kcal"`
	HFSS       bool                `json:"hfss"`
}

type NutritionalInfoKcal struct {
	High int `json:"high"`
	Low  int `json:"low"`
}

type PriceInfo struct {
	Overrides []PriceOverride `json:"overrides"`
	Price     int             `json:"price"`
}

type PriceOverride struct {
	ID    string `json:"id"`
	Price int    `json:"price"`
	Type  string `json:"type"`
}

type Mealtimes struct {
	CategoryIDs    []string    `json:"category_ids"`
	Description    Description `json:"description"`
	ID             string      `json:"id"`
	Image          ItemImage   `json:"image"`
	Name           Name        `json:"name"`
	Schedule       []Schedule  `json:"schedule"`
	SEODescription interface{} `json:"seo_description"`
}

type Schedule struct {
	DayOfWeek   int           `json:"day_of_week"`
	TimePeriods []TimePeriods `json:"time_periods"`
}

type TimePeriods struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type Modifier struct {
	Description  Description `json:"description"`
	ID           string      `json:"id"`
	ItemIDs      []string    `json:"item_ids"`
	MaxSelection int         `json:"max_selection"`
	MinSelection int         `json:"min_selection"`
	Name         Name        `json:"name"`
	Repeatable   bool        `json:"repeatable"`
}

type UploadMenuResponse struct {
	Status string `json:"status"`
	Result string `json:"result"`
}
