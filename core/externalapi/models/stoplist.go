package models

type StopListResponse struct {
	Items     []StopListItem     `json:"items"`
	Modifiers []StopListModifier `json:"modifiers"`
}

type StopListItem struct {
	ItemId string  `json:"itemId"`
	Stock  float64 `json:"stock,omitempty"`
}

type StopListModifier struct {
	ModifierId string  `json:"modifierId"`
	Stock      float64 `json:"stock,omitempty"`
}
