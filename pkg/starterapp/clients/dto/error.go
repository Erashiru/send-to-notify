package dto

type ErrorResponse struct {
	Detail []Detail `json:"detail"`
}

type Detail struct {
	Loc  []string `json:"loc"`
	Msg  string   `json:"msg"`
	Type string   `json:"type"`
}
