package handler

import (
	"encoding/json"
	"errors"
	"github.com/leavemeal0ne/SolidgateTestTask/internal/domen"
	"net/http"
)

var ( //
	LuhnErrorCode       = "001"
	ExpErrorCode        = "002"
	IINErrorCode        = "003"
	UnknownIINErrorCode = "004"
)

type Validator interface {
	Validate(card domen.Card) error
}

type Handler struct {
	Validator Validator
}

func InitHandler(Validator Validator) *Handler {
	return &Handler{Validator: Validator}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /validate", h.ValidateCardHandler)
	return mux
}

func (h *Handler) ValidateCardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var card domen.Card
	err := json.NewDecoder(r.Body).Decode(&card)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	err = domen.EmptyFieldRaiseErr(&card)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := domen.Response{Valid: false}

	err = h.Validator.Validate(card)

	switch {
	case errors.Is(err, domen.LuhnError):
		response.Error = &domen.ErrorInfo{Code: LuhnErrorCode, Message: err.Error()}
	case errors.Is(err, domen.IINError):
		response.Error = &domen.ErrorInfo{Code: IINErrorCode, Message: err.Error()}
	case errors.Is(err, domen.UnknownIINError):
		response.Error = &domen.ErrorInfo{Code: UnknownIINErrorCode, Message: err.Error()}
	case errors.Is(err, domen.ExpError):
		response.Error = &domen.ErrorInfo{Code: ExpErrorCode, Message: err.Error()}
	case errors.Is(err, nil):
		response.Valid = true
	default:
		http.Error(w, "Failed to validate card", http.StatusInternalServerError)
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response data", http.StatusInternalServerError)
	}
}
