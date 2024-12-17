package models

type OrderStat struct {
	TotalOrderNumber  float64            `json:"total_order_number"`
	TotalFailed       float64            `json:"total_failed"`
	TimeoutErrs       float64            `json:"timeout_err"`
	Errors            map[string]float64 `json:"errors"`
	ConstructedErrMsg string             `json:"constructed_err_msg"`
}
