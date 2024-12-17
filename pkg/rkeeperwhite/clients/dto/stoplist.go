package dto

type StopListResponse struct {
	TaskResponse   StopListTaskResponse `json:"taskResponse"`
	ResponseCommon ResponseCommon       `json:"responseCommon"`
	ErrResponse    ErrResponse          `json:"error,omitempty"`
}

type StopListTaskResponse struct {
	StopList StopList `json:"stopList"`
}

type StopList struct {
	Dishes []Dish `json:"dishes"`
}

type Dish struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}
