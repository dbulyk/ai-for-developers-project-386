package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/models"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/store"
	"github.com/go-chi/chi/v5"
)

// AdminEventTypesHandler exposes CRUD operations for event types under /admin.
type AdminEventTypesHandler struct {
	store store.Store
}

// NewAdminEventTypesHandler creates a new AdminEventTypesHandler.
func NewAdminEventTypesHandler(s store.Store) *AdminEventTypesHandler {
	return &AdminEventTypesHandler{store: s}
}

// eventTypeRequest is the request body for creating or updating an event type.
type eventTypeRequest struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	DurationMinutes int32  `json:"durationMinutes"`
}

// RegisterRoutes mounts the admin event types routes on the given router.
func (h *AdminEventTypesHandler) RegisterRoutes(r chi.Router) {
	r.Get("/admin/event-types", h.List)
	r.Post("/admin/event-types", h.Create)
	r.Put("/admin/event-types/{id}", h.Update)
	r.Delete("/admin/event-types/{id}", h.Delete)
}

// List returns all event types.
func (h *AdminEventTypesHandler) List(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, h.store.ListEventTypes())
}

// Create creates a new event type.
func (h *AdminEventTypesHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req eventTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}

	id, err := generateID()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "failed to generate id")
		return
	}

	et, err := models.NewEventType(id, req.Name, req.Description, req.DurationMinutes)
	if err != nil {
		var vErr *models.ValidationError
		if errors.As(err, &vErr) {
			writeError(w, http.StatusBadRequest, vErr.Code, vErr.Message)
			return
		}
		writeError(w, http.StatusBadRequest, "validation_failed", err.Error())
		return
	}

	h.store.CreateEventType(et)
	writeJSON(w, http.StatusCreated, et)
}

// Update updates an existing event type.
func (h *AdminEventTypesHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id_required", "id is required")
		return
	}

	et, ok := h.store.GetEventType(id)
	if !ok {
		writeError(w, http.StatusNotFound, "not_found", "event type not found")
		return
	}

	var req eventTypeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "invalid request body")
		return
	}

	if err := et.Update(req.Name, req.Description, req.DurationMinutes); err != nil {
		var vErr *models.ValidationError
		if errors.As(err, &vErr) {
			writeError(w, http.StatusBadRequest, vErr.Code, vErr.Message)
			return
		}
		writeError(w, http.StatusBadRequest, "validation_failed", err.Error())
		return
	}

	if !h.store.UpdateEventType(et) {
		writeError(w, http.StatusNotFound, "not_found", "event type not found")
		return
	}

	writeJSON(w, http.StatusOK, et)
}

// Delete removes an event type by id.
func (h *AdminEventTypesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "id_required", "id is required")
		return
	}

	if !h.store.DeleteEventType(id) {
		writeError(w, http.StatusNotFound, "not_found", "event type not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// generateID returns a cryptographically secure random 16-byte hex string.
func generateID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
