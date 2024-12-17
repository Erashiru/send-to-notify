package models

type MenuUploadResult struct {
	HttpStatus int         `json:"http_status"`
	BrandId    string      `json:"brand_id"`
	MenuId     string      `json:"menu_id"`
	SiteIds    []string    `json:"site_ids"`
	Errors     interface{} `json:"errors"`
}
