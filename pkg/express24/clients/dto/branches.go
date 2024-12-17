package dto

type GetBranchesResponse struct {
	BranchId   int    `json:"branch_id"`
	Name       string `json:"name"`
	IsActive   string `json:"is_active"`
	StoreId    int    `json:"store_id"`
	ExternalId string `json:"external_id"`
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type BranchesError struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
	Code   int    `json:"code"`
}
type UpdateBranchesRequest struct {
	BranchIds []int  `json:"branch_ids"`
	IsActive  string `json:"is_active"`
}
type UpdateBranchesResponse struct {
	BranchId int    `json:"branch_id"`
	Name     string `json:"name"`
	IsActive string `json:"is_active"`
	StoreId  int    `json:"store_id"`
}

type UpdateBranchesError struct {
	Status int    `json:"status"`
	Data   string `json:"data,omitempty"`
	Error  string `json:"error,omitempty"`
}
