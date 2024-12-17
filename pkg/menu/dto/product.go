package dto

type DeleteProducts struct {
	MenuID     string   `json:"menu_id"`
	ProductIds []string `json:"product_ids"`
}
