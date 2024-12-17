package models

type Pos string

const (
	IIKO        Pos = "iiko"
	JOWI        Pos = "jowi"
	Yaros       Pos = "yaros"
	RKeeper     Pos = "rkeeper"
	BurgerKing  Pos = "burger_king"
	Paloma      Pos = "paloma"
	Syrve       Pos = "syrve"
	WaitSending Pos = "wait_sending"
	FoodBand    Pos = "foodband"
	Poster      Pos = "poster"
	RKeeper7XML Pos = "rkeeper7_xml"
	CTMax       Pos = "CTMAX"
	TillyPad    Pos = "tillypad"
	Kwaaka      Pos = "kwaaka_pos"
	Ytimes      Pos = "ytimes"
	Posist      Pos = "posist"
)

func (p Pos) String() string {
	return string(p)
}
