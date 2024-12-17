package dto

type Pos string

func (p Pos) String() string {
	return string(p)
}

const (
	IIKO    Pos = "iiko"
	RKeeper Pos = "rkeeper"
)
