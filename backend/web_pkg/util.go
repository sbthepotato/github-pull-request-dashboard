package web_pkg

import (
	"context"
	"encoding/json"
	"github-pull-request-dashboard/github_pkg"
	"net/http"
	"sync"
	"time"

	"github.com/google/go-github/v80/github"
)

var mu sync.Mutex
var lastRateLimit time.Time
var cachedRateLimits *github.RateLimits

func setHeaders(w *http.ResponseWriter, content_type string) {

	switch content_type {
	case "text":
		(*w).Header().Set("Content-Type", "text/plain")
	case "json":
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

/*
Hello World from the backend
*/
func HelloGo(w http.ResponseWriter, r *http.Request) {
	setHeaders(&w, "text")
	w.Write([]byte("Hello, from the golang backend " + time.Now().String()))
}

func GetRateLimit(ctx context.Context, c *github.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(&w, "json")

		mu.Lock()
		defer mu.Unlock()

		currentTime := time.Now()

		if cachedRateLimits == nil ||
			currentTime.Sub(lastRateLimit).Minutes() > 1 {
			lastRateLimit = time.Now()
			rateLimit, err := github_pkg.GetApiLimit(ctx, c)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			cachedRateLimits = rateLimit
		}

		jsonData, err := json.Marshal(cachedRateLimits)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)
	}
}
