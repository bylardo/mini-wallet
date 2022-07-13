package models

type Deposit struct {
	ID          string `json:"id"`
	DepositBy   string `json:"deposited_by"`
	Status      string `json:"status"`
	DepositAt   string `json:"deposited_at"`
	Amount      int    `json:"amount"`
	ReferenceID string `json:"reference_id"`
}
