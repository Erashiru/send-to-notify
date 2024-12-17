package dto

type AuthenticateData struct {
	ClientID     string `json:"client_id" validate:"required" form:"client_id"`
	ClientSecret string `json:"client_secret" validate:"required" form:"client_secret"`
	GrantType    string `json:"grant_type" validate:"required" form:"grant_type"`
	Scope        string `json:"scope" validate:"required" form:"scope"`
}

type SuccessResponse struct {
	AccessToken string `json:"access_token"`
}

type UnauthorizedResponse struct {
	Reason string `json:"reason"`
}
