package web_pkg

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/go-github/github"
)

func GetRepos(ctx context.Context, c *github.Client, owner string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(&w, "json")

		mu.Lock()
		defer mu.Unlock()

		repos, err := gh_get_repos(ctx, c, owner)

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

func SetRepos(w http.ResponseWriter, r *http.Request) {
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

	repos, err := read_repos(false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// probably not the fastest way to do this but the list should never be huge...
	for _, setRepo := range setRepos {
		if setRepo.Enabled {
			for _, repo := range repos {
				if *repo.Name == setRepo.Name {
					repo.Enabled = &setRepo.Enabled
				}
			}
		}
	}

	err = write_repos(repos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Repo data saved successfully"))
}
