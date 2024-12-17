package dto

type ResponseListRestaurants struct {
	Status      int          `json:"status"`
	Restaurants []Restaurant `json:"restaurants"`
	ErrorResponse
}

type Restaurant struct {
	Id            string `json:"id"`
	Title         string `json:"title"`
	Timezone      string `json:"timezone"`
	Description   string `json:"description"`
	IsDispatching bool   `json:"is_dispatching"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}
