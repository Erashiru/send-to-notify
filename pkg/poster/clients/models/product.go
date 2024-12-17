package models

type Price struct {
	Field1 string `json:"1"`
}

type GetProductsResponseStop struct {
	SpotId      string `json:"spot_id"`
	Price       string `json:"price"`
	Profit      string `json:"profit"`
	ProfitNetto string `json:"profit_netto"`
	Visible     string `json:"visible"`
}

type GetProductsResponseSource struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Price   string `json:"price"`
	Visible int    `json:"visible"`
}

type GetProductsGroupModification struct {
	DishModificationGroupId int                       `json:"dish_modification_group_id"`
	Name                    string                    `json:"name"`
	NumMin                  int                       `json:"num_min"`
	NumMax                  int                       `json:"num_max"`
	Type                    int                       `json:"type"`
	IsDeleted               int                       `json:"is_deleted"`
	Modifications           []GetProductsModification `json:"modifications"`
}
type GetTovarModification struct {
	ModificatorID             string `json:"modificator_id"`
	ModificatorName           string `json:"modificator_name"`
	ModificatorSelfprice      string `json:"modificator_selfprice"`
	ModificatorSelfpriceNetto string `json:"modificator_selfprice_netto"`
	Order                     string `json:"order"`
	ModificatorBarcode        string `json:"modificator_barcode"`
	ModificatorProductCode    string `json:"modificator_product_code"`
	IngredientId              string `json:"ingredient_id"`
	FiscalCode                string `json:"fiscal_code"`
	MasterId                  string `json:"master_id"`
}
type GetProductsModification struct {
	DishModificationId int    `json:"dish_modification_id"`
	Name               string `json:"name"`
	IngredientId       int    `json:"ingredient_id"`
	Type               int    `json:"type"`
	Price              int    `json:"price"`
	PhotoOrig          string `json:"photo_orig"`
	PhotoLarge         string `json:"photo_large"`
	PhotoSmall         string `json:"photo_small"`
	LastModifiedTime   string `json:"last_modified_time"`
}

type GetProductsResponseBody struct {
	Barcode              string                         `json:"barcode"`
	CategoryName         string                         `json:"category_name"`
	Unit                 string                         `json:"unit"`
	Cost                 string                         `json:"cost"`
	CostNetto            string                         `json:"cost_netto"`
	Fiscal               string                         `json:"fiscal"`
	Hidden               string                         `json:"hidden"`
	MenuCategoryId       string                         `json:"menu_category_id"`
	Workshop             string                         `json:"workshop"`
	Nodiscount           string                         `json:"nodiscount"`
	Photo                string                         `json:"photo"`
	PhotoOrigin          string                         `json:"photo_origin"`
	Price                Price                          `json:"price"`
	ProductCode          string                         `json:"product_code"`
	ProductId            string                         `json:"product_id"`
	ProductName          string                         `json:"product_name"`
	Profit               Price                          `json:"profit"`
	SortOrder            string                         `json:"sort_order"`
	TaxId                string                         `json:"tax_id"`
	ProductTaxId         string                         `json:"product_tax_id"`
	Type                 string                         `json:"type"`
	WeightFlag           string                         `json:"weight_flag"`
	Color                string                         `json:"color"`
	Spots                []GetProductsResponseStop      `json:"spots"`
	IngredientId         string                         `json:"ingredient_id"`
	CookingTime          string                         `json:"cooking_time"`
	DifferentSpotsPrices string                         `json:"different_spots_prices"`
	Sources              []GetProductsResponseSource    `json:"sources"`
	GroupModifications   []GetProductsGroupModification `json:"group_modifications"`
	Modifications        []GetTovarModification         `json:"modifications"`
	//Out                          interface{}                    `json:"out"`
	ProductProductionDescription string `json:"product_production_description"`
	//Ingredients                  []interface{}                  `json:"ingredients"`
}

type GetIngridient struct {
	IngredientId           int    `json:"ingredient_id"`
	IngredientName         string `json:"ingredient_name"`
	IngredientBarcode      string `json:"ingredient_barcode"`
	CategoryId             int    `json:"category_id"`
	IngredientLeft         int    `json:"ingredient_left"`
	LimitValue             int    `json:"limit_value"`
	TimeNotif              int    `json:"time_notif"`
	IngredientUnit         string `json:"ingredient_unit"`
	IngredientWeight       int    `json:"ingredient_weight"`
	IngredientsLossesClear int    `json:"ingredients_losses_clear"`
	IngredientsLossesCook  int    `json:"ingredients_losses_cook"`
	IngredientsLossesFry   int    `json:"ingredients_losses_fry"`
	IngredientsLossesStew  int    `json:"ingredients_losses_stew"`
	IngredientsLossesBake  int    `json:"ingredients_losses_bake"`
	IngredientsType        int    `json:"ingredients_type"`
	PartialWriteOff        int    `json:"partial_write_off"`
}

type GetIngridientsResponse struct {
	Response []GetIngridient `json:"response"`
	ErrorResponse
}

type GetProductsResponse struct {
	Response []GetProductsResponseBody `json:"response"`
	ErrorResponse
}

type GetProductResponse struct {
	Response GetProductsResponseBody `json:"response"`
	ErrorResponse
}
