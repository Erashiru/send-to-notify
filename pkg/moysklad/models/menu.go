package models

type GetMenuRequest struct {
	Limit  string `json:"limit"`
	Offset string `json:"offset"`
}

type Menu struct {
	Context Context    `json:"context"`
	Meta    Meta       `json:"meta"`
	Rows    []MenuRows `json:"rows"`
}

type MenuRows struct {
	Meta                Meta        `json:"meta"`
	ID                  string      `json:"id"`
	AccountID           string      `json:"accountId"`
	Owner               MetaData    `json:"owner"`
	Shared              bool        `json:"shared,omitempty"`
	Group               MetaData    `json:"group"`
	Updated             string      `json:"updated"`
	Name                string      `json:"name"`
	Code                string      `json:"code"`
	ExternalCode        string      `json:"externalCode"`
	Archived            bool        `json:"archived,omitempty"`
	PathName            string      `json:"pathName,omitempty"`
	UseParentVat        bool        `json:"useParentVat,omitempty"`
	Vat                 int         `json:"vat,omitempty"`
	VatEnabled          bool        `json:"vatEnabled,omitempty"`
	EffectiveVat        int         `json:"effectiveVat,omitempty"`
	EffectiveVatEnabled bool        `json:"effectiveVatEnabled,omitempty"`
	Uom                 MetaData    `json:"uom,omitempty"`
	Images              MetaData    `json:"images,omitempty"`
	MinPrice            Price       `json:"minPrice,omitempty"`
	SalePrices          []SalePrice `json:"salePrices,omitempty"`
	Supplier            MetaData    `json:"supplier,omitempty"`
	BuyPrice            Price       `json:"buyPrice,omitempty"`
	Article             string      `json:"article,omitempty"`
	Weight              float64     `json:"weight,omitempty"`
	Volume              float64     `json:"volume,omitempty"`
	Barcodes            []Barcode   `json:"barcodes,omitempty"`
	VariantsCount       int         `json:"variantsCount,omitempty"`
	IsSerialTrackable   bool        `json:"isSerialTrackable,omitempty"`
	Stock               float64     `json:"stock,omitempty"`
	Reserve             float64     `json:"reserve,omitempty"`
	InTransit           float64     `json:"inTransit,omitempty"`
	Quantity            float64     `json:"quantity,omitempty"`
	Label               string      `json:"label,omitempty"`
	Assortment          MetaData    `json:"assortment,omitempty"`
	Components          MetaData    `json:"components,omitempty"`
	Characteristics     []PriceType `json:"characteristics,omitempty"`
	Product             MetaData    `json:"product,omitempty"`
}

type Barcode struct {
	Ean13 string `json:"ean13"`
}
type Price struct {
	Value    float64  `json:"value"`
	Currency MetaData `json:"currency"`
}

type SalePrice struct {
	Value     float64   `json:"value"`
	Currency  Price     `json:"currency"`
	PriceType PriceType `json:"priceType"`
}

type PriceType struct {
	Meta         Meta   `json:"meta"`
	ID           string `json:"id"`
	Name         string `json:"name"`
	ExternalCode string `json:"externalCode"`
}
