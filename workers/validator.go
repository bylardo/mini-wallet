package workers

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func ValidateToken(r *http.Request, tokenID string) bool {
	reqToken := r.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Token ")
	reqToken = splitToken[1]

	if tokenID == reqToken {
		return true
	}

	return false
}

func ValidateRequestBody(r *http.Request) bool {

	err := r.ParseMultipartForm(0)

	if err != nil {
		return false
	}

	return true
}

func ValidateJSONFileReader(path string) (bool, []byte) {
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		return false, nil
	}
	return true, jsonData
}
