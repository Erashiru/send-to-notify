package dto

type ErrorResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int64  `json:"id"`
	Error   Error  `json:"error"`
}
type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    string `json:"data"`
}
