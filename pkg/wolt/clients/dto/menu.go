package dto

type Menu struct {
	Currency        string     `json:"currency"`
	PrimaryLanguage string     `json:"primary_language"`
	Categories      []Category `json:"categories"`
}

type Category struct {
	ID          string     `json:"id"`
	Name        []ItemName `json:"name"`
	Description []ItemName `json:"description"`
	Items       []MenuItem `json:"items"`
}

type ItemName struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

type MenuItem struct {
	Name               []ItemName         `json:"name"`
	Description        []ItemName         `json:"description"`
	ImageUrl           string             `json:"image_url"`
	Price              float32            `json:"price"`
	SalesTaxPercentage int                `json:"sales_tax_percentage"`
	AlcoholPercentage  float32            `json:"alcohol_percentage"`
	Enabled            bool               `json:"enabled"`
	ExternalData       string             `json:"external_data"`
	DeliveryMethods    []string           `json:"delivery_methods"`
	Options            []OptionItem       `json:"options"`
	ProductInformation ProductInformation `json:"product_information"`
}

type OptionItem struct {
	Name           []ItemName      `json:"name"`
	Type           string          `json:"type"`
	ExternalData   string          `json:"external_data"`
	Values         []ValueItem     `json:"values"`
	SelectionRange *SelectionRange `json:"selection_range,omitempty"`
}

type SelectionRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type ValueItem struct {
	Name               []ItemName         `json:"name"`
	Price              float32            `json:"price"`
	Enabled            bool               `json:"enabled"`
	ExternalData       string             `json:"external_data"`
	Default            bool               `json:"default"`
	SubOptionValues    []SubOptionValue   `json:"sub_option_values,omitempty"`
	SelectionRange     *SelectionRange    `json:"selection_range,omitempty"`
	ProductInformation ProductInformation `json:"product_information"`
}

type SubOptionValue struct {
	Enabled        bool           `json:"enabled"`
	Name           []ItemName     `json:"name"`
	Price          int            `json:"price"`
	Default        bool           `json:"default,omitempty"`
	SelectionRange SelectionRange `json:"selection_range,omitempty"`
}

type ProductInformation struct {
	RegulatoryInformation []RegulatoryInformationValues `json:"regulatory_information"`
}

type RegulatoryInformationValues struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
