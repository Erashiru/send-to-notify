package models

type AwakeResponse struct {
	SuccessfullyProcessed []string `json:"successfullyProcessed"`
	FailedProcessed       []string `json:"failedProcessed"`
}
