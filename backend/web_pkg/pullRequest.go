package web_pkg

import (
	"context"
	"database/sql"
	"encoding/json"
	"github-pull-request-dashboard/db_pkg"
	"github-pull-request-dashboard/github_pkg"
	"net/http"
	"time"

	"github.com/google/go-github/v80/github"
)

var cachedPrListResults map[string]*db_pkg.PullRequestInfo

func GetPullRequests(ctx context.Context, db *sql.DB, c *github.Client, owner string, defaultRepo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		mu.Lock()
		defer mu.Unlock()

		setHeaders(&w, "json")
		refresh := r.URL.Query().Get("refresh")

		repo := r.URL.Query().Get("repo")
		if repo == "" || repo == "null" || repo == "undefined" {
			repo = defaultRepo
		}

		currentTime := time.Now()

		if cachedPrListResults == nil ||
			cachedPrListResults[repo] == nil ||
			(refresh == "y" && currentTime.Sub(*cachedPrListResults[repo].Updated).Minutes() > 1) ||
			currentTime.Sub(*cachedPrListResults[repo].Updated).Minutes() > 2 {

			if cachedPrListResults == nil {
				cachedPrListResults = make(map[string]*db_pkg.PullRequestInfo)
			}

			pullRequestResult, err := github_pkg.GetPullRequests(ctx, db, c, owner, repo, cachedPrListResults[repo])

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			cachedPrListResults[repo] = pullRequestResult
		}

		jsonData, err := json.Marshal(cachedPrListResults[repo])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)
	}
}
