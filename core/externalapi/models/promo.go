package models

type Promo struct {
	PromoItems []PromoItem `json:"promoItems"`
}

type PromoItem struct {
	Id      string `json:"id"`
	PromoId string `json:"promoId"`
}
