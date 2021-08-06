package helpers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func JSONError(err error, w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")

	var response = Response{
		Message: err.Error(),
		Status:  status,
	}

	w.WriteHeader(response.Status)
	json.NewEncoder(w).Encode(response)
}

func JSONResult(result interface{}, status int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	var response = Response{
		Status: status,
		Result: result,
	}

	json.NewEncoder(w).Encode(response)
}
