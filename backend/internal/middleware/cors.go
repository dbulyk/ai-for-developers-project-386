package middleware

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// CORS returns a chi middleware that validates the Origin header against
// a comma-separated allowlist. Preflight OPTIONS requests are answered directly.
func CORS(allowedOrigins string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{})
	for _, o := range strings.Split(allowedOrigins, ",") {
		o = strings.TrimSpace(o)
		if o != "" {
			allowed[o] = struct{}{}
		}
	}

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				if _, ok := allowed[origin]; ok {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
			}

			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// RegisterCORS mounts CORS middleware on a chi router.
func RegisterCORS(r chi.Router, allowedOrigins string) {
	r.Use(CORS(allowedOrigins))
}
