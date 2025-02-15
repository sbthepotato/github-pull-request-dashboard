package web_pkg

import (
	"context"
	"database/sql"
	"encoding/json"
	"github-pull-request-dashboard/db_pkg"
	"github-pull-request-dashboard/github_pkg"
	"io"
	"net/http"

	"github.com/google/go-github/v69/github"
)

func GetDefaultRepository(ctx context.Context, repository string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(&w, "text")
		w.Write([]byte(repository))
	}
}

func GetRepositories(ctx context.Context, db *sql.DB, c *github.Client, owner string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(&w, "json")

		mu.Lock()
		defer mu.Unlock()

		refresh := r.URL.Query().Get("refresh")
		activeOnly := r.URL.Query().Get("active")

		repos := make([]*db_pkg.Repository, 0)
		var err error

		if refresh == "y" {
			repos, err = github_pkg.GetRepositories(ctx, db, c, owner)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			repos, err = db_pkg.GetRepositories(ctx, db, activeOnly == "y")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		jsonData, err := json.Marshal(repos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)

	}
}

func SetRepositories(ctx context.Context, db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		setHeaders(&w, "text")

		mu.Lock()
		defer mu.Unlock()

		if r.Method != http.MethodPost {
			http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		defer r.Body.Close()

		repositories := make([]*db_pkg.Repository, 0)

		err = json.Unmarshal(body, &repositories)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = db_pkg.SetRepositories(ctx, db, repositories)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Repo data saved successfully"))
	}
}
