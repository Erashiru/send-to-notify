package models

type GetStopListResponse struct {
	StopLists []StopList `json:"response"`
}
type StopList struct {
	IngredientId              string  `json:"ingredient_id"`
	IngredientName            string  `json:"ingredient_name"`
	IngredientLeft            string  `json:"ingredient_left"`
	LimitValue                string  `json:"limit_value"`
	IngredientUnit            string  `json:"ingredient_unit"`
	IngredientsType           string  `json:"ingredients_type"`
	StorageIngredientSum      string  `json:"storage_ingredient_sum"`
	StorageIngredientSumNetto string  `json:"storage_ingredient_sum_netto"`
	PrimeCost                 float64 `json:"prime_cost"`
	PrimeCostNetto            string  `json:"prime_cost_netto"`
	Hidden                    string  `json:"hidden"`
}
