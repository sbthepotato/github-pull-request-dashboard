package github_pkg

/*
get list of github pull requests and process them with review information
*/
/*
func GetPullRequests(ctx context.Context, c *github.Client, owner string, repo string, prevResult *PullRequestInfo) (*PullRequestInfo, error) {

	if prevResult == nil {
		prevResult = new(PullRequestInfo)
	}

	if prevResult.Updated == nil {
		prevResult.Updated = &time.Time{}
	}

	log.Println("time ", *prevResult.Updated)

	opts := &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var ghPrs []*github.PullRequest

	for {
		respPrs, resp, err := c.PullRequests.List(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}

		ghPrs = append(ghPrs, respPrs...)

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	users, err := readUsers()
	if err != nil {
		return nil, err
	}

	teams, err := readTeams()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	idx := 0
	PrChannel := make(chan *CustomPullRequest)

	for _, pr := range ghPrs {
		if *pr.Draft {
			continue
		}

		wg.Add(1)
		go processPullRequest(PrChannel, &wg, ctx, c, owner, repo, pr, users, teams, idx)
		idx++ // manual index as we are skipping draft
	}

	result := new(PullRequestInfo)

	// sort the list of teams for the aggregation banner
	if teams != nil {
		slugs := make([]string, 0, len(teams))

		for _, slug := range teams {
			slugs = append(slugs, slug)
		}

		sort.SliceStable(slugs, func(i, j int) bool {
			return *teams[slugs[i]].ReviewOrder < *teams[slugs[j]].ReviewOrder
		})

		sorted_teams := make([]*CustomTeam, 0)
		for _, slug := range slugs {
			sorted_teams = append(sorted_teams, teams[slug])
		}

		result.ReviewTeams = sorted_teams
	}

	if users != nil {
		userList := make([]*CustomUser, 0)

		for _, user := range users {
			userList = append(userList, user)
		}

		result.Users = userList
	}

	go func() {
		wg.Wait()
		close(PrChannel)
	}()

	prs := make([]*CustomPullRequest, idx)

	for processed_pr := range PrChannel {
		prs[*processed_pr.Index] = processed_pr
	}

	result.PullRequests = prs
	currentTime := time.Now()
	result.Updated = &currentTime

	return result, nil

}

*/
/*
Process a pull request into the pull request channel
*/
/*
func processPullRequest(PrChannel chan<- *CustomPullRequest, wg *sync.WaitGroup, ctx context.Context, c *github.Client, owner string, repo string, pr *github.PullRequest, users map[string]*CustomUser, teams map[string]*CustomTeam, idx int) {

	defer wg.Done()

	customPr := new(CustomPullRequest)
	var ErrorMessage string
	var ErrorText string

	detailedPr, _, err := c.PullRequests.Get(ctx, owner, repo, *pr.Number)
	if err != nil {
		ErrorText = ErrorText + err.Error()
		ErrorMessage = ErrorMessage + "error fetching detailed pr info for pr " + strconv.Itoa(*pr.Number)
		log.Println(ErrorMessage, err.Error())
	}

	if *detailedPr.Draft {
		*detailedPr.State = "draft"
	}

	customPr.PullRequest = detailedPr

	review := new(Review)
	reviewOverview := make([]Review, 0)
	userReviewList := make([]string, 0)
	statusRequested := "REVIEW_REQUESTED"
	changesRequested := "Changes Requested"
	statusApproved := "APPROVED"
	teamOther := "other"

	if teams == nil {
		teamOther = "review"
	}

	// first populate requested teams and users. any previous state doesn't matter if you're requested
	if detailedPr.RequestedTeams != nil {
		for _, requestedTeam := range detailedPr.RequestedTeams {
			review.State = &statusRequested

			// if the team map isn't available, just use what we have
			if val, ok := teams[*requestedTeam.Slug]; ok {
				review.Team = val
			} else {
				custom_team := new(CustomTeam)
				custom_team.Team = requestedTeam
				review.Team = custom_team
			}

			review.State = &statusRequested

			reviewOverview = append(reviewOverview, *review)
		}
	}

	if detailedPr.RequestedReviewers != nil {
		for _, requestedUser := range detailedPr.RequestedReviewers {

			// if the user map isn't available, just use what we have
			if val, ok := users[*requestedUser.Login]; ok {
				review.User = val
			} else {
				customUser := new(CustomUser)
				customUser.User = requestedUser
				review.User = customUser
			}

			if review.User.Team != nil {
				review.Team = teams[*review.User.Team.Slug]
			}
			review.State = &statusRequested

			reviewOverview = append(reviewOverview, *review)
			userReviewList = append(userReviewList, *requestedUser.Login)
		}
	}

	var reviews []*github.PullRequestReview

	opt := &github.ListOptions{PerPage: 100}

	for {
		respReviews, resp, err := c.PullRequests.ListReviews(ctx, owner, repo, *detailedPr.Number, opt)
		if err != nil {
			ErrorText = ErrorText + err.Error()
			ErrorMessage = ErrorMessage + "error fetching pull request reviews for pr " + strconv.Itoa(*pr.Number)
			teamOther = "error"
			customPr.Awaiting = &teamOther

			log.Println(ErrorMessage, err.Error())
		}

		reviews = append(reviews, respReviews...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage

	}

	// loop in reverse because we're only interested in the most recent event
	for i := len(reviews) - 1; i >= 0; i-- {
		//for _, review := range reviews {
		ghReview := reviews[i]
		if (!slices.Contains(userReviewList, *ghReview.User.Login)) &&
			(*detailedPr.User.Login != *ghReview.User.Login) &&
			(*ghReview.State != "COMMENTED") {

			review := new(Review)
			userReviewList = append(userReviewList, *ghReview.User.Login)
			if val, ok := users[*ghReview.User.Login]; ok {
				review.User = val
			} else {
				customUser := new(CustomUser)
				customUser.User = ghReview.User
				review.User = customUser
			}

			if review.User != nil && review.User.Team != nil {
				review.Team = teams[*review.User.Team.Slug]
			}
			review.State = ghReview.State
			reviewOverview = append(reviewOverview, *review)

		}
	}

	customPr.CreatedBy = users[*detailedPr.User.Login]
	customPr.ReviewOverview = make([]*Review, 0)
	customPr.Index = &idx
	currentPriority := 100
	unassigned := false
	approvedCount := 0

	for _, customReview := range reviewOverview {
		if *customReview.State != "DISMISSED" {

			review := new(Review)
			review.User = customReview.User
			review.Team = customReview.Team
			review.State = customReview.State

			if *review.State == statusRequested {
				if review.Team != nil {

					if review.Team.ReviewOrder != nil &&
						*review.Team.ReviewOrder < currentPriority &&
						(customPr.Awaiting != nil &&
							*customPr.Awaiting != "error" ||
							customPr.Awaiting == nil) {

						currentPriority = *review.Team.ReviewOrder
						customPr.Awaiting = review.Team.Name

						if review.User == nil {
							unassigned = true
							customPr.Unassigned = &unassigned
						} else {
							unassigned = false
							customPr.Unassigned = &unassigned
						}

					} else if review.Team.ReviewOrder != nil &&
						*review.Team.ReviewOrder == currentPriority &&
						review.User != nil {

						unassigned = false
						customPr.Unassigned = &unassigned

					} else if customPr.Awaiting == nil {
						customPr.Awaiting = &teamOther
					}

				} else if review.Team == nil {
					customPr.Awaiting = &teamOther
				}

			} else if *review.State == "CHANGES_REQUESTED" {
				customPr.Awaiting = &changesRequested
				currentPriority = -1
				unassigned = false
				customPr.Unassigned = &unassigned

			} else if *review.State == statusApproved {
				approvedCount++
			}
			customPr.ReviewOverview = append(customPr.ReviewOverview, review)

		}

		if customPr.Awaiting == nil && approvedCount >= 1 {
			customPr.Awaiting = &statusApproved
			unassigned = false
			customPr.Unassigned = &unassigned
		}
	}

	if ErrorMessage != "" {
		ErrorMessage = ErrorMessage + ". See the console on the server or the errorText field in the network tab for more information"
		customPr.ErrorMessage = &ErrorMessage
		customPr.ErrorText = &ErrorText
	}

	PrChannel <- customPr

}
*/
