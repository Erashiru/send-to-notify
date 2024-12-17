package models

type StoreGroup struct {
	ID                  string              `bson:"_id,omitempty" json:"id"`
	Name                string              `bson:"name" json:"name"`
	StoreIds            []string            `bson:"restaurant_ids" json:"store_ids"`
	Locations           []Location          `bson:"locations,omitempty" json:"locations"`
	IsTopPartner        bool                `bson:"is_top_partner" json:"is_top_partner"`
	RetryCount          int                 `bson:"retry_count" json:"retry_count"`
	Tags                string              `bson:"tags,omitempty" json:"tags,omitempty"`
	Logo                string              `bson:"logo,omitempty" json:"logo,omitempty"`
	ExtraLogo           string              `bson:"extra_logo,omitempty" json:"extra_logo,omitempty"`
	Country             string              `bson:"country,omitempty" json:"country,omitempty"`
	BrandType           string              `bson:"brand_type,omitempty" json:"brand_type,omitempty"`
	Category            string              `bson:"category,omitempty" json:"category,omitempty"`
	Chats               []Chat              `bson:"chats,omitempty" json:"chats,omitempty"`
	Contacts            []Contact           `bson:"contacts,omitempty" json:"contacts,omitempty"`
	SalesComments       string              `bson:"sales_comments,omitempty" json:"sales_comments,omitempty"`
	Status              string              `bson:"status" json:"status"`
	ColumnView          bool                `bson:"column_view" json:"column_view"`
	Description         string              `bson:"description,omitempty" json:"description,omitempty"`
	HeaderImage         string              `bson:"header_image,omitempty" json:"header_image,omitempty"`
	WorkSchedule        []Schedule          `bson:"work_schedule,omitempty" json:"work_schedule,omitempty"`
	SocialMediaLinks    []SocialMediaLink   `bson:"social_media_links,omitempty" json:"social_media_links,omitempty"`
	BrandInfo           BrandInfo           `bson:"brand_info,omitempty" json:"brand_info,omitempty"`
	DomainName          string              `bson:"domain_name" json:"domain_name"`
	DefaultRestaurantId string              `bson:"default_restaurant_id" json:"default_restaurant_id"`
	DefaultCity         string              `bson:"default_city" json:"default_city"`
	IsShowcase          bool                `bson:"is_showcase" json:"is_showcase"`
	DirectPromoBanners  []DirectPromoBanner `bson:"direct_promo_banners" json:"direct_promo_banners"`
	CancelOrderAllowed  bool                `bson:"cancel_order_allowed,omitempty" json:"cancel_order_allowed,omitempty"` //настройка для отмены заказа в поске
}

type SocialMediaLink struct {
	Name string `bson:"name" json:"name"`
	URL  string `bson:"url" json:"url"`
	Logo string `bson:"logo" json:"logo"`
}

type Chat struct {
	ChatName string `bson:"chat_name,omitempty" json:"chat_name,omitempty"`
	ChatLink string `bson:"chat_link,omitempty" json:"chat_link,omitempty"`
}

type Contact struct {
	FullName string `bson:"full_name,omitempty" json:"full_name,omitempty"`
	Position string `bson:"position,omitempty" json:"position,omitempty"`
	Phone    string `bson:"phone,omitempty" json:"phone,omitempty"`
	Comment  string `bson:"comment,omitempty" json:"comment,omitempty"`
	IsMain   bool   `bson:"is_main" json:"is_main"`
}

type Location struct {
	Address      string             `json:"address,omitempty"`
	MenuName     string             `json:"menu_name,omitempty"`
	MenuID       string             `json:"menu_id,omitempty"`
	RestaurantID string             `json:"restaurant_id"`
	Delivery     string             `json:"delivery,omitempty"`
	GlovoConfig  GlovoLimitedConfig `json:"glovo_config,omitempty"`
	WoltConfig   WoltLimitedConfig  `json:"wolt_config,omitempty"`

	RestaurantName string                 `json:"restaurant_name"`
	Deliveries     []AggregatorInLocation `json:"deliveries,omitempty"`
}

type GlovoLimitedConfig struct {
	StoreIds  []string `json:"store_ids"`
	IsOpen    bool     `json:"is_open"`
	SendToPos bool     `json:"send_to_pos"`
}
type WoltLimitedConfig struct {
	StoreIds  []string `json:"store_ids"`
	IsOpen    bool     `json:"is_open"`
	SendToPos bool     `json:"send_to_pos"`
}

type AggregatorInLocation struct {
	Type     string `json:"type"`
	MenuID   string `json:"menu_id"`
	Name     string `json:"name,omitempty"`
	IsActive bool   `json:"is_active"`
}

type StoreGroupIdAndName struct {
	ID   string `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
}

type StoreIdAndName struct {
	ID   string `bson:"_id" json:"id"`
	Name string `bson:"name" json:"name"`
}

type BrandInfo struct {
	Name        string `bson:"name,omitempty"`
	Description string `bson:"description,omitempty"`
	LogoURL     string `bson:"logo_url,omitempty"`
	HeaderURL   string `bson:"header_url,omitempty"`
}

type DirectPromoBanner struct {
	ID            string   `bson:"promo_id" json:"promo_id"`
	RestaurantIDs []string `bson:"restaurant_ids" json:"restaurant_ids"`
	Image         string   `bson:"image,omitempty" json:"image,omitempty"`
	IsActive      bool     `bson:"is_active" json:"is_active"`
}

type UpdateDirectPromoBanner struct {
	ID            *string   `bson:"promo_id"`
	RestaurantIDs *[]string `bson:"restaurant_ids"`
	Image         *string   `bson:"image"`
	IsActive      *bool     `bson:"is_active"`
}
