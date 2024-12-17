package models

type SubmitCatalogRequest struct {
	ChainCode   string   `json:"-"`
	Vendors     []string `json:"vendors"`
	Catalog     Catalog  `json:"catalog"`
	CallbackUrl string   `json:"callbackUrl,omitempty"`
}

type SubmitCatalogResponse struct {
	Status          string `json:"status"`
	CatalogImportId string `json:"catalogImportId"`
}

type Catalog struct {
	Items map[string]CatalogItem `json:"items"`
}

type CatalogItem struct {
	Id          string             `json:"id"`
	Type        string             `json:"type"`
	MenuType    string             `json:"menuType,omitempty"`
	Title       *Title             `json:"title,omitempty"`
	Products    map[string]SubItem `json:"products,omitempty"`
	Schedule    map[string]SubItem `json:"schedule,omitempty"`
	Images      map[string]SubItem `json:"images,omitempty"`
	IsActive    *bool              `json:"active,omitempty"`
	Price       string             `json:"price,omitempty"`
	Toppings    map[string]SubItem `json:"toppings,omitempty"`
	Description *Title             `json:"description,omitempty"`
	Quantity    *Quantity          `json:"quantity,omitempty"`
	URL         string             `json:"url,omitempty"`
	Alt         *Title             `json:"alt,omitempty"`
	StartTime   string             `json:"startTime,omitempty"`
	EndTime     string             `json:"endTime,omitempty"`
	StartDate   string             `json:"startDate,omitempty"`
	EndDate     string             `json:"endDate,omitempty"`
	WeekDays    []string           `json:"weekDays,omitempty"`
	Order       int                `json:"order,omitempty"`
}

type Title struct {
	Default string `json:"default"`
}

type SubItem struct {
	Id    string `json:"id"`
	Order int    `json:"order,omitempty"`
	Type  string `json:"type"`
	Price string `json:"price,omitempty"`
	Title *Title `json:"title,omitempty"`
}

type Quantity struct {
	Min int `json:"minimum"`
	Max int `json:"maximum"`
}
