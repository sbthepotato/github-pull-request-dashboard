package web_pkg

import (
	"context"
	"database/sql"
	"encoding/json"
	"github-pull-request-dashboard/db_pkg"
	"github-pull-request-dashboard/github_pkg"
	"net/http"

	"github.com/google/go-github/v71/github"
)

/*
get list of users part of owner organization
*/
func GetUsers(ctx context.Context, db *sql.DB, c *github.Client, owner string, defaultRepoName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(&w, "json")

		mu.Lock()
		defer mu.Unlock()

		var err error

		refresh := r.URL.Query().Get("refresh")
		syncType := r.URL.Query().Get("type")
		repositoryName := r.URL.Query().Get("repo")

		if repositoryName == "" || repositoryName == "null" || repositoryName == "undefined" {
			repositoryName = defaultRepoName
		}

		if refresh == "y" {
			if syncType == "users" {
				_, err = github_pkg.GetUsers(ctx, db, c, owner)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				_, err = github_pkg.GetUserTeams(ctx, db, c, owner, repositoryName)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		users, err := db_pkg.GetUsersAsTeamMap(ctx, db, repositoryName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonData, err := json.Marshal(users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)
	}
}
