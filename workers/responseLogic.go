package workers

import (
	"encoding/json"
	"net/http"

	"miniwallet.co.id/models"
)

func PrepareAPIResponse(statusResponse bool, message string, data interface{}) models.APIResponse {

	var varResponseMessage models.APIResponse
	var varErrorResponse models.ErrorResponse
	if !statusResponse {
		varErrorResponse.Error = message
		varResponseMessage.Status = "fail"
		varResponseMessage.Data = varErrorResponse
	}

	if statusResponse {
		varResponseMessage.Status = "success"
		varResponseMessage.Data = data
	}

	return varResponseMessage
}

func sendAPIResponse(responseMessage models.APIResponse, w http.ResponseWriter) {
	json.NewEncoder(w).Encode(responseMessage)
	return
}
