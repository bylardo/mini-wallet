package models

type Wallet []struct {
	ID         string `json:"id"`
	OwnedBy    string `json:"owned_by"`
	Status     string `json:"status"`
	EnabledAt  string `json:"enabled_at"`
	DisabledAt string `json:"disabled_at"`
	Balance    int    `json:"balance"`
}
