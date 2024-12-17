package models

type TapRestaurant struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Image       string `json:"img" bson:"img"`
	QRMenuLink  string `json:"qr_menu_link" bson:"qr_menu_link"`
	Tel         string `json:"tel" bson:"tel"`
	Instagram   string `json:"instagram" bson:"instagram"`
	Website     string `json:"website" bson:"website"`
}

type UpdateTapRestaurant struct {
	ID          *string `json:"id" bson:"_id"`
	Name        *string `json:"name" bson:"name"`
	Description *string `json:"description" bson:"description"`
	Image       *string `json:"img" bson:"img"`
	QRMenuLink  *string `json:"qr_menu_link" bson:"qr_menu_link"`
	Tel         *string `json:"tel" bson:"tel"`
	Instagram   *string `json:"instagram" bson:"instagram"`
	Website     *string `json:"website" bson:"website"`
}
