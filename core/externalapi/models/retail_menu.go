package models

type Categories struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	ParentId *string `json:"parentId,omitempty"`
}

type NomenclatureItem struct {
	Barcode       Barcode         `json:"barcode"`
	CategoryId    string          `json:"categoryId"`
	Description   ItemDescription `json:"description"`
	Id            string          `json:"id"`
	Images        []Images        `json:"images"`
	IsCatchWeight bool            `json:"isCatchWeight"`
	Measure       Measure         `json:"measure"`
	Name          string          `json:"name"`
	Price         float64         `json:"price"`
	VendorCode    string          `json:"vendorCode"`
	ExciseValue   string          `json:"exciseValue,omitempty"`
	Labels        []string        `json:"labels,omitempty"`
	LimitPerOrder int             `json:"limitPerOrder,omitempty"`
	MarkingType   string          `json:"markingType,omitempty"`
	OldPrice      float64         `json:"oldPrice,omitempty"`
	Vat           int             `json:"vat,omitempty"`
	VendorInn     string          `json:"vendorInn,omitempty"`
	Volume        *Volume         `json:"volume,omitempty"`
}

type Barcode struct {
	Type           string `json:"type"`
	Value          string `json:"value"`
	WeightEncoding string `json:"weightEncoding"`
	Values         string `json:"values,omitempty"`
}

type ItemDescription struct {
	Composition         string `json:"composition,omitempty"`
	ExpiresIn           string `json:"expiresIn,omitempty"`
	General             string `json:"general,omitempty"`
	NutritionalValue    string `json:"nutritionalValue,omitempty"`
	PackageInfo         string `json:"packageInfo,omitempty"`
	Purpose             string `json:"purpose,omitempty"`
	StorageRequirements string `json:"storageRequirements,omitempty"`
	VendorCountry       string `json:"vendorCountry,omitempty"`
	VendorName          string `json:"vendorName,omitempty"`
}

type Images struct {
	Url   string `json:"url"`
	Hash  string `json:"hash,omitempty"`
	Order int    `json:"order,omitempty"`
}

type Measure struct {
	Unit    string  `json:"unit"`
	Value   int     `json:"value"`
	Quantum float64 `json:"quantum,omitempty"`
}

type Volume struct {
	Unit  string `json:"unit,omitempty"`
	Value int    `json:"value,omitempty"`
}

type RetailMenu struct {
	Categories        []Categories       `json:"categories"`
	NomenclatureItems []NomenclatureItem `json:"items"`
}
