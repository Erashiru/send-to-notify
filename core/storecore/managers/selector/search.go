package selector

const (
	DefaultLimit  = 20
	DirectionASC  = 1
	DirectionDESC = -1
)

type Pagination struct {
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"`
}

func (p Pagination) HasPagination() bool {
	return p.Page >= 0 && p.Limit > 0
}

func (p Pagination) Skip() int64 {
	return p.Page * p.Limit
}

type Sorting struct {
	Param     string
	Direction int8
}

func (s Sorting) HasSorting() bool {
	return s.Param != "" && (s.Direction == DirectionASC || s.Direction == DirectionDESC)
}
