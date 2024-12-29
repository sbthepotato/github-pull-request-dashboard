package web_pkg

/*
func GetTeams(ctx context.Context, c *github.Client, owner string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setHeaders(&w, "json")

		mu.Lock()
		defer mu.Unlock()

		var err error

		refresh := r.URL.Query().Get("refresh")
		currentTime := time.Now()

		if refresh == "y" {
			cached_teams, err = github_pkg.GetTeams(ctx, c, owner)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			last_fetched_teams = time.Now()

		} else if (currentTime.Sub(last_fetched_teams).Hours() < 1) || (len(cached_teams) == 0) {

			teams, err := read_teams(false)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			for _, team := range teams {
				cached_teams = append(cached_teams, team)
			}
		}

		jsonData, err := json.Marshal(cached_teams)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonData)
	}
}

func SetTeams(w http.ResponseWriter, r *http.Request) {
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

	team_data := make([]SetTeam, 0)
	cached_teams = make([]*db_pkg.Team, 0)

	err = json.Unmarshal(body, &team_data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	team_map, err := read_teams(false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	active_team_map := make(map[string]*db_pkg.Team)

	for _, team := range team_data {
		*team_map[team.Slug].ReviewEnabled = team.ReviewEnabled
		*team_map[team.Slug].ReviewOrder = team.ReviewOrder

		if team.ReviewEnabled {
			active_team_map[team.Slug] = team_map[team.Slug]
		}

		updated_team := team_map[team.Slug]

		cached_teams = append(cached_teams, updated_team)
	}

	write_teams(active_team_map, true)
	write_teams(team_map, false)

	w.Write([]byte("Team data saved successfully"))
}
*/
