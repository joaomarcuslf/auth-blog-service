package helpers

import (
	"encoding/json"
	"net/http"

	types "auth_blog_service/types"
)

func JSONError(err error, w http.ResponseWriter, status int) {
	var response = types.ResponseBody{
		Status:  status,
		Message: err.Error(),
	}

	JSONResponse(response, w, status)
}

func JSONSuccess(result interface{}, w http.ResponseWriter, status int) {
	var response = types.ResponseBody{
		Status: status,
		Result: result,
	}

	JSONResponse(response, w, status)
}

func JSONResponse(response types.ResponseBody, w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Status)
	json.NewEncoder(w).Encode(response)
}
