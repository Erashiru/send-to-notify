package models

type UpdateStoreGroup struct {
	ID                  *string           `bson:"_id,omitempty" json:"id"`
	Name                *string           `bson:"name" json:"name"`
	StoreIds            []string          `bson:"restaurant_ids" json:"store_ids"`
	Locations           []Location        `bson:"locations,omitempty" json:"locations"`
	IsTopPartner        *bool             `bson:"is_top_partner" json:"is_top_partner"`
	RetryCount          *int              `bson:"retry_count" json:"retry_count"`
	Logo                *string           `bson:"logo,omitempty" json:"logo,omitempty"`
	Country             *string           `bson:"country,omitempty" json:"country,omitempty"`
	BrandType           *string           `bson:"brand_type,omitempty" json:"brand_type,omitempty"`
	Category            *string           `bson:"category,omitempty" json:"category,omitempty"`
	Chats               []Chat            `bson:"chats,omitempty" json:"chats,omitempty"`
	Contacts            []Contact         `bson:"contacts,omitempty" json:"contacts,omitempty"`
	SalesComments       *string           `bson:"sales_comments,omitempty" json:"sales_comments,omitempty"`
	Status              *string           `bson:"status" json:"status"`
	ColumnView          *bool             `bson:"column_view" json:"column_view"`
	Description         *string           `bson:"description" json:"description"`
	HeaderImage         *string           `bson:"header_image" json:"header_image"`
	WorkSchedule        []Schedule        `bson:"work_schedule,omitempty" json:"work_schedule,omitempty"`
	SocialMediaLinks    []SocialMediaLink `bson:"social_media_links" json:"social_media_links"`
	DomainName          *string           `bson:"domain_name" json:"domain_name"`
	DefaultRestaurantId *string           `bson:"default_restaurant_id" json:"default_restaurant_id"`
}
