package github_pkg

/*
get users and link them up to one of the active teams
*/
/*
func GetUsers(ctx context.Context, c *github.Client, owner string) ([]*CustomUser, error) {

	listMembersOpt := &github.ListMembersOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var allUsers []*github.User

	for {
		users, resp, err := c.Organizations.ListMembers(ctx, owner, listMembersOpt)
		if err != nil {
			return nil, err
		}

		allUsers = append(allUsers, users...)

		if resp.NextPage == 0 {
			break
		}
		listMembersOpt.Page = resp.NextPage
	}

	teams, err := readTeams()
	if err != nil {
		return nil, err
	}

	userTeams := make(map[string]*CustomTeam)

	// find team members of each team in org and add it to a map
	for _, team := range teams {

		if team.ReviewEnabled == nil || !*team.ReviewEnabled {
			continue
		}

		teamMembersOpt := &github.TeamListTeamMembersOptions{
			ListOptions: github.ListOptions{PerPage: 100},
		}

		var teamMembers []*github.User

		for {
			respMembers, resp, err := c.Teams.ListTeamMembersBySlug(ctx, owner, *team.Slug, teamMembersOpt)
			if err != nil {
				return nil, err
			}

			teamMembers = append(teamMembers, respMembers...)

			if resp.NextPage == 0 {
				break
			}
			teamMembersOpt.Page = resp.NextPage
		}

		for _, teamUser := range teamMembers {
			userTeams[*teamUser.Login] = team
		}
	}

	var wg sync.WaitGroup
	userChannel := make(chan *CustomUser)

	// go through all org members to get extended user info, also add team info
	for _, member := range allUsers {
		wg.Add(1)
		go processUser(userChannel, &wg, ctx, c, *member.Login, userTeams)
	}

	go func() {
		wg.Wait()
		close(userChannel)
	}()

	userMap := make(map[string]*CustomUser)
	customUsers := make([]*CustomUser, 0)

	for processedUser := range userChannel {
		userMap[*processedUser.Login] = processedUser
		customUsers = append(customUsers, processedUser)
	}

	writeUsers(userMap)

	return customUsers, nil
}
*/
/*
process a member into the member channel
*/
/*
func processUser(userChannel chan<- *CustomUser, wg *sync.WaitGroup, ctx context.Context, c *github.Client, login string, teams map[string]*CustomTeam) {
	defer wg.Done()

	user, _, err := c.Users.Get(ctx, login)
	if err != nil {
		log.Println("error fetching user", login, err)
	}

	customUser := new(CustomUser)
	customUser.User = user
	customUser.Team = teams[*user.Login]

	userChannel <- customUser
}

func writeUsers(users map[string]*CustomUser) error {
	jsonData, err := json.Marshal(users)
	if err != nil {
		return err
	}

	db_pkg.Write("users", jsonData)

	return nil
}
*/
