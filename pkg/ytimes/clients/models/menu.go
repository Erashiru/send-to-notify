package models

type Menu struct {
	Success bool        `json:"success"`
	Count   int         `json:"count"`
	Rows    []MenuRow   `json:"rows"`
	Error   interface{} `json:"error"`
}

type ItemList struct {
	Guid                          string         `json:"guid"`
	Name                          string         `json:"name"`
	Priority                      int            `json:"priority"`
	ImageLink                     string         `json:"imageLink"`
	Description                   string         `json:"description"`
	Recipe                        string         `json:"recipe"`
	TypeList                      []TypeList     `json:"typeList"`
	SupplementCategoryToFreeCount map[string]int `json:"supplementCategoryToFreeCount"`
	DefaultSupplements            []interface{}  `json:"defaultSupplements"`
}

type TypeList struct {
	Guid   string  `json:"guid"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
	IsTogo bool    `json:"isTogo"`
}

type GoodsList struct {
	Guid        string  `json:"guid"`
	Name        string  `json:"name"`
	Priority    int     `json:"priority"`
	Price       float64 `json:"price"`
	ImageLink   string  `json:"imageLink"`
	Description string  `json:"description"`
	Recipe      string  `json:"recipe"`
}

type MenuRow struct {
	Guid         string         `json:"guid"`
	Name         string         `json:"name"`
	Priority     int            `json:"priority"`
	ImageLink    string         `json:"imageLink"`
	CategoryList []CategoryList `json:"categoryList"`
	ItemList     []ItemList     `json:"itemList"`
	GoodsList    []GoodsList    `json:"goodsList"`
}

type CategoryList struct {
	Guid      string `json:"guid"`
	Name      string `json:"name"`
	Priority  int    `json:"priority"`
	ImageLink string `json:"imageLink"`
	//CategoryList []interface{} `json:"categoryList"`
	ItemList  []ItemList  `json:"itemList"`
	GoodsList []GoodsList `json:"goodsList"`
	ComboList interface{} `json:"comboList"`
}
