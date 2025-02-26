package api

import (
	"encoding/json"
	"net/http"
)

type AddContestantResponse struct {
	Code      int
	Id        int
	Name      string
	Created   string
	ExtRef    string
	InitPoint string
	BackURL   string
	PrefId    string
}

type PaymentResponse struct {
	Code    int
	BackURL string
	Message string
}

type Error struct {
	Code    int
	Message string
}

func writeError(w http.ResponseWriter, message string, code int) {
	resp := Error{
		Code:    code,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(resp)
}

var (
	RequestErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusBadRequest)
	}
	InternalErrorHandler = func(w http.ResponseWriter) {
		writeError(w, "An unexpected Error Ocurred", http.StatusInternalServerError)
	}
	UnauthorizedErrorHandler = func(w http.ResponseWriter) {
		writeError(w, "Invalid token", http.StatusUnauthorized)
	}
	PaymentError = func(w http.ResponseWriter, backURL string) {
		writeError(w, backURL, http.StatusBadRequest)
	}
)
