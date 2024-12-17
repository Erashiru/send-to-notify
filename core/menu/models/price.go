package models

type Price struct {
	Value        float64 `bson:"value" json:"value"`
	CurrencyCode string  `bson:"currency_code" json:"currency_code"`
}

func (lp Price) Get(price []Price) float64 {
	if len(price) != 0 {
		return price[0].Value
	}
	return 0
}
