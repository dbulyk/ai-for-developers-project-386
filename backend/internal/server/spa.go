package server

import (
	"io"
	"io/fs"
	"log/slog"
	"mime"
	"net/http"
	"path"
	"strings"
	"time"
)

// spaHandler serves a single-page application from the provided filesystem.
// It returns real files when they exist and falls back to index.html for
// any other path so that client-side routing can take over.
func spaHandler(dist fs.FS, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		name := strings.TrimPrefix(r.URL.Path, "/")

		f, err := dist.Open(name)
		if err == nil {
			stat, statErr := f.Stat()
			if statErr == nil && !stat.IsDir() {
				serveFile(w, r, f, name)
				return
			}
			_ = f.Close()
		}

		fallback, err := dist.Open("index.html")
		if err != nil {
			logger.Error("index.html not found in embedded assets", slog.Any("error", err))
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		defer func() { _ = fallback.Close() }()
		serveFile(w, r, fallback, "index.html")
	})
}

func serveFile(w http.ResponseWriter, r *http.Request, f fs.File, name string) {
	defer func() { _ = f.Close() }()

	rs, ok := f.(io.ReadSeeker)
	if !ok {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	contentType := mime.TypeByExtension(path.Ext(name))
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	http.ServeContent(w, r, name, time.Time{}, rs)
}
