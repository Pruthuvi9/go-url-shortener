package handler

import (
	"encoding/json"
	"net/http"

	"github.com/pruthuvifernando/url-shortner/internal/usecase"
)

type URLHandler struct {
	uc *usecase.URLUsecase
}

func NewURLHandler(uc *usecase.URLUsecase) *URLHandler {
	return &URLHandler{uc: uc}
}

type shortenRequest struct {
	LongURL        string  `json:"long_url"`
	CustomAlias    string  `json:"custom_alias,omitempty"`
	ExpirationDate *string `json:"expiration_date,omitempty"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
}

func (h *URLHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// TODO: call h.uc.Shorten, map domain errors to HTTP status codes
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *URLHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")

	// TODO: call h.uc.Redirect, map ErrNotFound->404, ErrExpired->410
	_ = shortCode
	http.Redirect(w, r, "", http.StatusFound)
}

func (h *URLHandler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	shortCode := r.PathValue("shortCode")

	// TODO: call h.uc.Delete, map ErrNotFound->404
	_ = shortCode
	w.WriteHeader(http.StatusNoContent)
}
