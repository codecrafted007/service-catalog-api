package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    int64       `json:"code"`
	Data    interface{} `json:"data"`
	Error   string      `json:"error"`
	Success bool        `json:"success"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}, errStr string) {
	resp := Response{
		Code:    int64(statusCode),
		Data:    data,
		Error:   errStr,
		Success: errStr == "",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
