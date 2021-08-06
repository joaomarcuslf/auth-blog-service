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
	var response = Response{
		Status:  status,
		Message: err.Error(),
	}

	JSONResponse(response, w, status)
}

func JSONSuccess(result interface{}, w http.ResponseWriter, status int) {
	var response = Response{
		Status: status,
		Result: result,
	}

	JSONResponse(response, w, status)
}

func JSONResponse(response Response, w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	json.NewEncoder(w).Encode(response)
}
