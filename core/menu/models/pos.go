package models

type PosName string

const (
	IIKO        PosName = "iiko"
	IIKOWEB     PosName = "iikoweb"
	RKEEPER     PosName = "rkeeper"
	MAIN        PosName = "POS Menu"
	PALOMA      PosName = "paloma"
	POSTER      PosName = "poster"
	SYRVE       PosName = "syrve"
	JOWI        PosName = "jowi"
	YAROS       PosName = "yaros"
	RKEEPER7XML PosName = "rkeeper7_xml"
	Tillypad    PosName = "tillypad"
	Ytimes      PosName = "ytimes"
)

func (p PosName) String() string {
	return string(p)
}
