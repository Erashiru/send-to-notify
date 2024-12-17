package models

type AuthRequest struct {
	ApiLogin string `json:"apiLogin,omitempty"`
}

type AuthResponse struct {
	AccessToken   string `json:"token"`
	CorrelationID string `json:"correlationId"`
}
