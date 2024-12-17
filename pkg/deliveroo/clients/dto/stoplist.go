package dto

type GetUnavailabilitiesRequest struct {
	BrandID string `json:"brand_id"`
	MenuID  string `json:"menu_id"`
	SiteID  string `json:"site_id"`
}

type GetUnavailabilitiesResponse struct {
	UnavailableIDs []string `json:"unavailable_ids"`
	HiddenIDs      []string `json:"hidden_ids"`
}

type UpdateUnavailabilitesRequest struct {
	BrandID string `json:"brand_id"`
	MenuID  string `json:"menu_id"`
	SiteID  string `json:"site_id"`

	UpdateUnavailabilitesRequestBody UpdateUnavailabilitesRequestBody
}

type UpdateUnavailabilitesRequestBody struct {
	UnavailableIDs []string `json:"unavailable_ids"`
	HiddenIDs      []string `json:"hidden_ids"`
}
