package dto

type Pagination struct {
	SyncTime  string `json:"sync_time"`
	PageCount int    `json:"page_count"`
	Page      int    `json:"page"`
	PerPage   int    `json:"per_page"`
}
