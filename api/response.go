package api

import (
	"encoding/json"
	"net/http"
)

var (
	ErrInvalidJson    = Error{StatusCode: http.StatusBadRequest, Type: "invalid_json", Message: "Invalid or malformed JSON"}
	ErrMissingData    = Error{StatusCode: http.StatusBadRequest, Type: "missing_data", Message: "Missing required Data"}
	ErrInternalServer = Error{StatusCode: http.StatusInternalServerError, Type: "server_error", Message: "Internal server Error"}
	ErrUnauthorized   = Error{StatusCode: http.StatusUnauthorized, Type: "unauthorized", Message: "Unauthorized"}
	ErrForbidden      = Error{StatusCode: http.StatusForbidden, Type: "forbidden", Message: "Forbidden"}
)

// Alert
type Alert struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"status,omitempty"`
	Type       string `json:"type"`
	Message    string `json:"message,omitempty"`
}

func (a Alert) Send(w http.ResponseWriter) error {
	statusCode := a.StatusCode

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(a)
}

type Response struct {
	Success    bool        `json:"success"`
	StatusCode int         `json:"status,omitempty"`
	Type       string      `json:"type,omitempty"`
	Result     interface{} `json:"result,omitempty"`
}

func Success(result interface{}, status int) *Response {
	return &Response{
		Success:    true,
		StatusCode: status,
		Type:       "ok",
		Result:     result,
	}
}

func (r *Response) Send(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.StatusCode)
	return json.NewEncoder(w).Encode(r)
}

// Error
type Error struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"status,omitempty"`
	Type       string `json:"type"`
	Message    string `json:"message,omitempty"`
}

func (e Error) Send(w http.ResponseWriter) error {
	statusCode := e.StatusCode

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(e)
}
