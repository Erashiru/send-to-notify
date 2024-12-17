package models

type ModifierGroup struct {
	ObjectId  int        `json:"object_id"`
	Name      string     `json:"name"`
	Modifiers []Modifier `json:"modifiers"`
}

type Modifier struct {
	ObjectId    int     `json:"object_id"`
	MarkDeleted int     `json:"mark_deleted"`
	IUseInMenu  int     `json:"i_useInMenu"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	MinCount    int     `json:"min_count"`
	MaxCount    int     `json:"max_count"`
	Image       string  `json:"image"`
}

type ComplexGroup struct {
	ObjectId      int           `json:"object_id"`
	Name          string        `json:"name"`
	MinCount      int           `json:"min_count"`
	MaxCount      int           `json:"max_count"`
	DefaultItemId int           `json:"default_item_id"`
	ComplexItems  []ComplexItem `json:"complex_items"`
}

type ComplexItem struct {
	ObjectId    int     `json:"object_id"`
	MarkDeleted string  `json:"mark_deleted"`
	IUseInMenu  string  `json:"i_useInMenu"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Image       string  `json:"image"`
}

type OtherPrices struct {
	Menu       []PriceMenu `json:"menu"`
	PriceTypes []PriceType `json:"price_types"`
}

type PriceMenu struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
}

type PriceType struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
}

type Item struct {
	ObjectId    int         `json:"object_id"`
	MarkDeleted int         `json:"mark_deleted"`
	IUseInMenu  int         `json:"i_useInMenu"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Price       float64     `json:"price"`
	Quantity    float64     `json:"quantity"` // TODO: ??? in docs that's number (int), but in api it's float64
	Image       string      `json:"image"`
	EditDate    string      `json:"edit_date"`
	OtherPrices OtherPrices `json:"other_prices"`

	ModifierGroups []ModifierGroup `json:"modifier_groups"`
	ComplexGroups  []ComplexGroup  `json:"complex_groups"`
}

type ItemGroup struct {
	ObjectId int    `json:"object_id"`
	Name     string `json:"name"`
	Image    string `json:"image"`
	Items    []Item `json:"items"`
}

type Menu struct {
	ItemGroups []ItemGroup `json:"item_groups"`
}
