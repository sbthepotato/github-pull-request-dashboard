package web_pkg

/*
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
			cached_members, err = github_pkg.GetUsers(ctx, c, owner)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		} else if currentTime.Sub(last_fetched_members).Hours() < 1 || (len(cached_members) == 0) {
			cached_members = make([]*db_pkg.User, 0)
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

*/
