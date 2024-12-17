package dto

type UpdateOrderStatusRequest struct {
	Status       string `json:"status"`
	RejectReason string `json:"reject_reason"`
	Notes        string `json:"notes"`
}

type CreateSyncStatusRequest struct {
	Status     string `json:"status"`
	Reason     string `json:"reason"`
	Notes      string `json:"notes"`
	OccurredAt string `json:"occurred_at"`
}
