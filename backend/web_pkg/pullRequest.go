package web_pkg

import (
	"github-pull-request-dashboard/db_pkg"
)

var cachedPrListResults map[string]*db_pkg.PullRequestInfo

/*
func GetPrList(ctx context.Context, c *github.Client, owner string, defaultRepo string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		mu.Lock()
		defer mu.Unlock()

		setHeaders(&w, "json")
		refresh := r.URL.Query().Get("refresh")
		currentTime := time.Now()

		repo := r.URL.Query().Get("repo")
		if repo == "" {
			repo = defaultRepo
		}

		log.Println("repo: ", repo)

		if currentTime.Sub(last_fetched_prs).Minutes() > 5 || (refresh == "y" && currentTime.Sub(last_fetched_prs).Seconds() > 2) {
			cached_prs = new(github_pkg.PullRequestInfo)

			prs, err := github_pkg.GetPullRequests(ctx, c, owner, repo, cachedPrListResults[repo])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			cachedPrListResults[repo] = prs
			last_fetched_prs = time.Now()
		}

		jsonData, err := json.Marshal(cachedPrListResults[repo])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)
	}
}

*/
