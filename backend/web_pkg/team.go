package web_pkg

import (
	"context"
	"database/sql"
	"encoding/json"
	"github-pull-request-dashboard/db_pkg"
	"github-pull-request-dashboard/github_pkg"
	"io"
	"net/http"

	"github.com/google/go-github/v71/github"
)

func GetTeams(ctx context.Context, db *sql.DB, c *github.Client, owner string, defaultRepo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(&w, "json")

		mu.Lock()
		defer mu.Unlock()

		refresh := r.URL.Query().Get("refresh")
		repositoryName := r.URL.Query().Get("repo")

		if repositoryName == "" || repositoryName == "null" || repositoryName == "undefined" {
			repositoryName = defaultRepo
		}

		if refresh == "y" {
			_, err := github_pkg.GetTeams(ctx, db, c, owner)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// fetch teams again as somebody might refresh to get new teams
		// but have existing team information that shouldn't be lost
		teams, err := db_pkg.GetTeams(ctx, db, repositoryName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(teams)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)
	}
}

func SetTeams(ctx context.Context, db *sql.DB) http.HandlerFunc {
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

		teams := make([]*db_pkg.Team, 0)

		err = json.Unmarshal(body, &teams)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = db_pkg.UpsertTeamReviews(ctx, db, teams)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("Team review data saved successfully"))
	}
}
