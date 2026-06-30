package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/models"
	"github.com/dbulyk/ai-for-developers-project-386/backend/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAdminEventTypesHandler() (*AdminEventTypesHandler, *store.MemoryStore) {
	s := store.NewMemoryStore().(*store.MemoryStore)
	h := NewAdminEventTypesHandler(s)
	return h, s
}

func TestAdminEventTypes_List(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		h, _ := setupAdminEventTypesHandler()
		r := chi.NewRouter()
		h.RegisterRoutes(r)

		req := httptest.NewRequest(http.MethodGet, "/admin/event-types", http.NoBody)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, "[]", rec.Body.String())
	})

	t.Run("returns event types", func(t *testing.T) {
		h, s := setupAdminEventTypesHandler()
		r := chi.NewRouter()
		h.RegisterRoutes(r)

		et := models.EventType{ID: "et-1", Name: "Consultation", Description: "Quick chat", DurationMinutes: 30}
		s.CreateEventType(et)

		req := httptest.NewRequest(http.MethodGet, "/admin/event-types", http.NoBody)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var got []models.EventType
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
		require.Len(t, got, 1)
		assert.Equal(t, et, got[0])
	})
}

func TestAdminEventTypes_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h, s := setupAdminEventTypesHandler()
		r := chi.NewRouter()
		h.RegisterRoutes(r)

		body := `{"name":"Consultation","description":"Quick chat","durationMinutes":30}`
		req := httptest.NewRequest(http.MethodPost, "/admin/event-types", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var got models.EventType
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
		assert.NotEmpty(t, got.ID)
		assert.Equal(t, "Consultation", got.Name)
		assert.Equal(t, "Quick chat", got.Description)
		assert.Equal(t, int32(30), got.DurationMinutes)

		stored, ok := s.GetEventType(got.ID)
		require.True(t, ok)
		assert.Equal(t, got, stored)
	})

	t.Run("invalid json", func(t *testing.T) {
		h, _ := setupAdminEventTypesHandler()
		r := chi.NewRouter()
		h.RegisterRoutes(r)

		req := httptest.NewRequest(http.MethodPost, "/admin/event-types", bytes.NewBufferString("not-json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid_json")
	})

	t.Run("validation error", func(t *testing.T) {
		h, _ := setupAdminEventTypesHandler()
		r := chi.NewRouter()
		h.RegisterRoutes(r)

		body := `{"name":"","description":"","durationMinutes":0}`
		req := httptest.NewRequest(http.MethodPost, "/admin/event-types", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "name_required")
	})
}

func TestAdminEventTypes_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h, s := setupAdminEventTypesHandler()
		r := chi.NewRouter()
		h.RegisterRoutes(r)

		et := models.EventType{ID: "et-1", Name: "Old", Description: "Old desc", DurationMinutes: 30}
		s.CreateEventType(et)

		body := `{"name":"New","description":"New desc","durationMinutes":60}`
		req := httptest.NewRequest(http.MethodPut, "/admin/event-types/et-1", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var got models.EventType
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
		assert.Equal(t, "et-1", got.ID)
		assert.Equal(t, "New", got.Name)
		assert.Equal(t, "New desc", got.Description)
		assert.Equal(t, int32(60), got.DurationMinutes)
	})

	t.Run("not found", func(t *testing.T) {
		h, _ := setupAdminEventTypesHandler()
		r := chi.NewRouter()
		h.RegisterRoutes(r)

		body := `{"name":"New","description":"New desc","durationMinutes":60}`
		req := httptest.NewRequest(http.MethodPut, "/admin/event-types/missing", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "not_found")
	})

	t.Run("invalid json", func(t *testing.T) {
		h, s := setupAdminEventTypesHandler()
		r := chi.NewRouter()
		h.RegisterRoutes(r)

		s.CreateEventType(models.EventType{ID: "et-1", Name: "Old", DurationMinutes: 30})

		req := httptest.NewRequest(http.MethodPut, "/admin/event-types/et-1", bytes.NewBufferString("not-json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid_json")
	})

	t.Run("validation error", func(t *testing.T) {
		h, s := setupAdminEventTypesHandler()
		r := chi.NewRouter()
		h.RegisterRoutes(r)

		s.CreateEventType(models.EventType{ID: "et-1", Name: "Old", DurationMinutes: 30})

		body := `{"name":"","description":"","durationMinutes":-5}`
		req := httptest.NewRequest(http.MethodPut, "/admin/event-types/et-1", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "name_required")
	})
}

func TestAdminEventTypes_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h, s := setupAdminEventTypesHandler()
		r := chi.NewRouter()
		h.RegisterRoutes(r)

		s.CreateEventType(models.EventType{ID: "et-1", Name: "Old", DurationMinutes: 30})

		req := httptest.NewRequest(http.MethodDelete, "/admin/event-types/et-1", http.NoBody)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Empty(t, rec.Body.String())

		_, ok := s.GetEventType("et-1")
		assert.False(t, ok)
	})

	t.Run("not found", func(t *testing.T) {
		h, _ := setupAdminEventTypesHandler()
		r := chi.NewRouter()
		h.RegisterRoutes(r)

		req := httptest.NewRequest(http.MethodDelete, "/admin/event-types/missing", http.NoBody)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "not_found")
	})
}
