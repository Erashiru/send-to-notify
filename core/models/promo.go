package models

type PromoType string

const (
	Discount  PromoType = "discount"
)

func (a PromoType) String() string {
	return string(a)
}
