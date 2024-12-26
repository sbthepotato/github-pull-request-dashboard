package web_pkg

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/go-github/github"
)

func GetMembers(ctx context.Context, c *github.Client, owner string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(&w, "json")

		mu.Lock()
		defer mu.Unlock()

		var err error
		refresh := r.URL.Query().Get("refresh")
		currentTime := time.Now()

		if refresh == "y" {
			last_fetched_members = time.Now()
			cached_members, err = gh_get_members(ctx, c, owner)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		} else if currentTime.Sub(last_fetched_members).Hours() < 1 || (len(cached_members) == 0) {
			cached_members = make([]*CustomUser, 0)
			users, err := read_users()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			for _, user := range users {
				cached_members = append(cached_members, user)
			}
		}

		jsonData, err := json.Marshal(cached_members)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)
	}
}
