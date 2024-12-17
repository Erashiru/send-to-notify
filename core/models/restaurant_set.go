package models

type RestaurantSet struct {
	ID                 string   `bson:"_id,omitempty" json:"id"`
	Name               string   `json:"name" bson:"name"`
	Logo               string   `json:"logo" bson:"logo"`
	HeaderImage        string   `json:"header_image" bson:"header_image"`
	DomainName         string   `json:"domain_name" json:"domain_name"`
	RestaurantGroupIds []string `json:"restaurant_group_ids" bson:"restaurant_group_ids"`
}
