package event

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /v1/events", h.handleCreate)
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var req EvenCreationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Payload", http.StatusBadRequest)
		return
	}

	es, err := h.service.CreateEvent(req)

	if err != nil {
		http.Error(w, "Not processable", http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(es)
}
