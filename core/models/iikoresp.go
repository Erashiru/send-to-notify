package models

type Balance struct {
	Balance float64 `json:"balance"`
}

type Transactions struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	BalanceBefore float64 `json:"balance_before"`
	BalanceAfter  float64 `json:"balance_after"`
	TypeName      string  `json:"type_name"`
}
