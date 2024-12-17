package dto

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"message,omitempty"`
}
