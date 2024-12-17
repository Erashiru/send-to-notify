package models

type PointInfo struct {
	Success bool        `json:"success"`
	Count   int         `json:"count"`
	Rows    []Row       `json:"rows"`
	Error   interface{} `json:"error"`
}

type Row struct {
	Guid string `json:"guid"`
	Name string `json:"name"`
	Type string `json:"type"`
}
