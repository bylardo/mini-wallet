package models

type Withdraw struct {
	ID          string `json:"id"`
	WithdrawnBy string `json:"withdrawn_by"`
	Status      string `json:"status"`
	WithdrawnAt string `json:"withdrawn_at"`
	Amount      int    `json:"amount"`
	ReferenceID string `json:"reference_id"`
}
