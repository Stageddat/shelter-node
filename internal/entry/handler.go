package entry

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Routes(r chi.Router) {
	r.Get("/", h.ListEntries)
	r.Post("/", h.CreateEntry)
	r.Get("/{id}", h.GetEntry)
	r.Put("/{id}", h.UpdateEntry)
	r.Delete("/{id}", h.DeleteEntry)
}

func (h *Handler) ListEntries(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		writeError(w, http.StatusBadRequest, "userId query param required")
		return
	}

	entries, err := h.svc.GetByUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, entries)
}

func (h *Handler) CreateEntry(w http.ResponseWriter, r *http.Request) {
	var req CreateEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	e, err := h.svc.Create(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, e)
}

func (h *Handler) GetEntry(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	e, err := h.svc.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "entry not found")
		return
	}

	writeJSON(w, http.StatusOK, e)
}

func (h *Handler) UpdateEntry(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	e, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, e)
}

func (h *Handler) DeleteEntry(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.svc.Delete(r.Context(), id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
