package server

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func TestSpaHandler_ServesExistingFile(t *testing.T) {
	fsys := fstest.MapFS{
		"index.html":     {Data: []byte("<div id=\"root\">app</div>")},
		"assets/main.js": {Data: []byte("console.log('hello')")},
	}

	h := spaHandler(fsys, slog.Default())
	req := httptest.NewRequest(http.MethodGet, "/assets/main.js", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "console.log('hello')") {
		t.Fatalf("expected JS content, got %q", body)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "javascript") {
		t.Fatalf("expected javascript content type, got %q", contentType)
	}
}

func TestSpaHandler_FallsBackToIndexForUnknownPath(t *testing.T) {
	fsys := fstest.MapFS{
		"index.html":     {Data: []byte("<div id=\"root\">app</div>")},
		"assets/main.js": {Data: []byte("console.log('hello')")},
	}

	h := spaHandler(fsys, slog.Default())
	req := httptest.NewRequest(http.MethodGet, "/admin/event-types", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, `<div id="root">app</div>`) {
		t.Fatalf("expected index.html fallback, got %q", body)
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Fatalf("expected text/html content type, got %q", contentType)
	}
}

func TestSpaHandler_FallsBackToIndexForDirectoryPath(t *testing.T) {
	fsys := fstest.MapFS{
		"index.html": {Data: []byte("<div id=\"root\">app</div>")},
	}

	h := spaHandler(fsys, slog.Default())
	req := httptest.NewRequest(http.MethodGet, "/admin/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, `<div id="root">app</div>`) {
		t.Fatalf("expected index.html fallback, got %q", body)
	}
}

func TestSpaHandler_RootReturnsIndex(t *testing.T) {
	fsys := fstest.MapFS{
		"index.html": {Data: []byte("<div id=\"root\">app</div>")},
	}

	h := spaHandler(fsys, slog.Default())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	body, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	if !strings.Contains(string(body), `<div id="root">app</div>`) {
		t.Fatalf("expected index.html at root, got %q", string(body))
	}
}
