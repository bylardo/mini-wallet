package models

type ErrorResponse struct {
	Error string `json:"error"`
}

type APIResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
