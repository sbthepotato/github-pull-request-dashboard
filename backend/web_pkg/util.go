package web_pkg

import (
	"net/http"
	"sync"
)

var mu sync.Mutex

func setHeaders(w *http.ResponseWriter, content_type string) {

	if content_type == "text" {
		(*w).Header().Set("Content-Type", "text/plain")
	} else if content_type == "json" {
		(*w).Header().Set("Content-Type", "application/json")
	}
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func EnableCors(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Serve the actual request
		handler.ServeHTTP(w, r)
	})
}
