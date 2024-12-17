package dto

type SendPaymentOrderToCustomerRequest struct {
	Id     int64                                   `json:"id"`
	Method string                                  `json:"method"`
	Params SendPaymentOrderToCustomerRequestParams `json:"params"`
}

type SendPaymentOrderToCustomerRequestParams struct {
	Id    string `json:"id"`
	Phone string `json:"phone"`
}

type SendPaymentOrderToCustomerResponse struct {
	Result SendPaymentOrderToCustomerResponseResult `json:"result"`
}

type SendPaymentOrderToCustomerResponseResult struct {
	Success bool `json:"success"`
}
