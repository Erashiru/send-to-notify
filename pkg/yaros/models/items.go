package models

type Item struct {
	Id           string `json:"ID"`
	CategoryId   string `json:"CATEGORY_ID"`
	Title        string `json:"TITLE"`
	Price        string `json:"PRICE"`
	Quantity     string `json:"QUANTITY"`
	ImageUrl     string `json:"IMAGE_URL"`
	Description  string `json:"DESCRIPTION"`
	Measure      string `json:"MEASURE"`
	SortPriority int    `json:"SORT_PRIORITY"`
}

type Category struct {
	ParentId     string `json:"PARENT_ID"`
	Id           string `json:"ID"`
	Title        string `json:"TITLE"`
	SortPriority int    `json:"SORT_PRIORITY"`
}

type GetItemsResponse struct {
	Status string `json:"status"`
	Items  []Item `json:"goods"`
}

type GetCategoriesResponse struct {
	Status     string     `json:"status"`
	Categories []Category `json:"categories"`
}
