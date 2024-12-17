package dto

type RestaurantGroupSetResponse struct {
	ID               string              `json:"id"`
	Name             string              `json:"name"`
	Logo             string              `json:"logo"`
	DomainName       string              `json:"domain_name"`
	HeaderImage      string              `json:"header_image"`
	RestaurantGroups []RestGroupResponse `json:"restaurant_groups"`
}

type RestGroupResponse struct {
	Id                  string `json:"id"`
	Name                string `json:"name"`
	ColumnView          bool   `json:"column_view"`
	Description         string `json:"description,omitempty"`
	HeaderImage         string `json:"header_image,omitempty"`
	DomainName          string `json:"domain_name"`
	DefaultRestaurantId string `json:"default_restaurant_id"`
	DefaultCity         string `json:"default_city"`
	Logo                string `json:"logo"`
	ExtraLogo           string `json:"extra_logo"`
	Tags                string `json:"tags"`
}
