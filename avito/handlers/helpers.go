package handlers

import (
	"encoding/json"
	"net/http"
)

func getErrorStatus(isInternal bool) int {
	if !isInternal {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}

func sendResponse(httpStatus int, response interface{}, w http.ResponseWriter) {
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(response)
}

