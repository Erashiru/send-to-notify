package dto

type RequestAuth struct {
	ApiKey string `form:"api_key"`
	Sig    string `form:"sig"`
}
