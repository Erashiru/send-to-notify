package models

type AuthRequest struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
	ApiLogin string `json:"apiLogin,omitempty"`
}

type AuthResponse struct {
	Status int    `json:"status"`
	Token  string `json:"data"`
}
