package models

type Pos string

const (
	IIKO  Pos = "iiko"
	SYRVE Pos = "syrve"
	Yaros Pos = "yaros"
)

func (p Pos) String() string {
	return string(p)
}
