package models

type SupplementList struct { // ModifierGroups
	Success bool            `json:"success"`
	Count   int             `json:"count"`
	Rows    []SupplementRow `json:"rows"`
	Error   interface{}     `json:"error"`
}

type SupplementRow struct {
	Guid             string               `json:"guid"`
	Name             string               `json:"name"`
	Priority         int                  `json:"priority"`
	AllowSeveralItem bool                 `json:"allowSeveralItem"`
	MaxSelectedCount int                  `json:"maxSelectedCount"`
	ItemList         []SupplementItemList `json:"itemList"`
}

type SupplementItemList struct {
	Guid             string             `json:"guid"`
	Name             string             `json:"name"`
	Priority         int                `json:"priority"`
	DefaultPrice     float64            `json:"defaultPrice"`
	DefaultTogoPrice float64            `json:"defaultTogoPrice"`
	MenuTypeToPrice  map[string]float64 `json:"menuTypeToPrice"`
}
