package models

type ErrorResponse struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type GetStoreResponse struct {
	Places []Place `json:"places"`
}

type Place struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Address string `json:"address"`
}

type Store struct {
	Token     string          `json:"token"`
	Name      string          `json:"name"`
	MenuID    string          `json:"menu_id"`
	PosType   string          `json:"pos_type"`
	Address   StoreAddress    `json:"address"`
	IIKOCloud StoreIIKOConfig `json:"iiko_cloud"`
	Delivery  []StoreDelivery `json:"delivery"`
	Menus     []StoreDSMenu   `json:"menus"`
}

type StoreAddress struct {
	City   string `json:"city"`
	Street string `json:"street"`
}

type StoreIIKOConfig struct {
	OrganizationID string `json:"organization_id"`
	TerminalID     string `json:"terminal_id"`
	Key            string `json:"key"`
}

type StoreDelivery struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	Price    int    `json:"price"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

type StoreDSMenu struct {
	MenuID    string `json:"menu_id"`
	Name      string `json:"name"`
	IsActive  bool   `json:"is_active"`
	IsDeleted bool   `json:"is_deleted"`
}
