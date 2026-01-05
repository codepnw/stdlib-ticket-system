package helper

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func SuccessResponse(w http.ResponseWriter, code int, msg string, data any) {
	if msg == "" {
		msg = "successful"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

func ErrorResponse(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response{
		Success: false,
		Message: msg,
		Data:    "-",
	})
}

func ParseInt64(id string) (int64, error) {
	return strconv.ParseInt(id, 10, 64)
}
