package web_pkg

import (
	"context"
	"encoding/json"
	"github-pull-request-dashboard/db_pkg"
	"github-pull-request-dashboard/github_pkg"
	"io"
	"net/http"

	"github.com/google/go-github/v68/github"
)

func GetRepositories(ctx context.Context, c *github.Client, owner string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(&w, "json")

		mu.Lock()
		defer mu.Unlock()

		repos, err := github_pkg.GetRepositories(ctx, c, owner)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(repos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)

	}
}

func SetRepos(ctx context.Context) http.HandlerFunc {
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

		setRepos := make([]setRepo, 0)

		err = json.Unmarshal(body, &setRepos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, setRepo := range setRepos {
			repository := new(db_pkg.Repository)
			repository.Enabled = &setRepo.Enabled
			repository.Name = &setRepo.Name
			err = db_pkg.SetRepository(ctx, repository)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		w.Write([]byte("Repo data saved successfully"))
	}
}
