package github_pkg

import (
	"context"
	"database/sql"
	"github-pull-request-dashboard/db_pkg"
	"log"
	"slices"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/google/go-github/v68/github"
)

/**** private ****/

/*
Process a pull request into the pull request channel
*/
func processPullRequest(prChannel chan<- *db_pkg.PullRequest,
	wg *sync.WaitGroup,
	ctx context.Context,
	c *github.Client,
	owner string,
	repositoryName string,
	pr *github.PullRequest,
	users map[string]*db_pkg.User,
	teams map[string]*db_pkg.Team,
	idx int) {

	defer wg.Done()

	var ErrorMessage string
	var ErrorText string

	detailedPr, _, err := c.PullRequests.Get(ctx, owner, repositoryName, *pr.Number)
	if err != nil {
		ErrorText = ErrorText + err.Error()
		ErrorMessage = ErrorMessage + "error fetching detailed pr info for pr " + strconv.Itoa(*pr.Number)
		log.Println(ErrorMessage, err.Error())
	}

	if *detailedPr.Draft {
		*detailedPr.State = "draft"
	}

	resultPr := new(db_pkg.PullRequest)
	resultPr.PullRequest = detailedPr

	reviewOverview := make([]db_pkg.Review, 0)
	userReviewList := make([]string, 0)
	statusRequested := "REVIEW_REQUESTED"
	changesRequested := "Changes Requested"
	statusApproved := "APPROVED"
	teamOther := "other"

	if teams == nil || len(teams) == 0 {
		teamOther = "review"
	}

	// first populate requested teams
	if detailedPr.RequestedTeams != nil {
		for _, requestedTeam := range detailedPr.RequestedTeams {
			review := new(db_pkg.Review)

			if team, ok := teams[*requestedTeam.Slug]; ok {
				review.Team = team
			} else {
				team := new(db_pkg.Team)
				team.Team = requestedTeam
				review.Team = team
			}

			review.State = &statusRequested
			reviewOverview = append(reviewOverview, *review)
		}
	}

	// then requested users, add them to the list as previous state doesn't matter if you're requested
	if detailedPr.RequestedReviewers != nil {
		for _, requestedUser := range detailedPr.RequestedReviewers {
			review := new(db_pkg.Review)

			if user, ok := users[*requestedUser.Login]; ok {
				review.User = user
				review.Team = user.Team
			} else {
				user := new(db_pkg.User)
				user.User = requestedUser
				review.User = user
			}

			review.State = &statusRequested

			userReviewList = append(userReviewList, *review.User.Login)
			reviewOverview = append(reviewOverview, *review)
		}
	}

	reviews := make([]*github.PullRequestReview, 0)
	opt := &github.ListOptions{PerPage: 100}

	for {
		respReviews, resp, err := c.PullRequests.ListReviews(ctx, owner, repositoryName, *detailedPr.Number, opt)
		if err != nil {
			ErrorText = ErrorText + err.Error()
			ErrorMessage = ErrorMessage + "error fetching pull request reviews for pr " + strconv.Itoa(*pr.Number)
			teamOther = "error"
			resultPr.Awaiting = &teamOther

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
		ghReview := reviews[i]
		if (!slices.Contains(userReviewList, *ghReview.User.Login)) &&
			(*detailedPr.User.Login != *ghReview.User.Login) &&
			(*ghReview.State != "COMMENTED") {

			review := new(db_pkg.Review)
			if user, ok := users[*ghReview.User.Login]; ok {
				review.User = user
				review.Team = user.Team
			} else {
				user := new(db_pkg.User)
				user.User = ghReview.User
				review.User = user
			}

			review.State = ghReview.State
			userReviewList = append(userReviewList, *ghReview.User.Login)
			reviewOverview = append(reviewOverview, *review)
		}
	}

	if val, ok := users[*detailedPr.User.Login]; ok {
		resultPr.CreatedBy = val
	} else {
		user := new(db_pkg.User)
		user.User = detailedPr.User
		resultPr.CreatedBy = user
	}

	resultPr.ReviewOverview = make([]*db_pkg.Review, 0)
	resultPr.Index = &idx
	currentPriority := 100
	unassigned := false
	approvedCount := 0

	for _, review := range reviewOverview {
		if *review.State != "DISMISSED" {

			finalReview := new(db_pkg.Review)
			finalReview.User = review.User
			finalReview.Team = review.Team
			finalReview.State = review.State

			if *finalReview.State == statusRequested {
				if finalReview.Team != nil {

					if finalReview.Team.ReviewOrder != nil &&
						*finalReview.Team.ReviewOrder < currentPriority &&
						(resultPr.Awaiting != nil &&
							*resultPr.Awaiting != "error" ||
							resultPr.Awaiting == nil) {

						currentPriority = *finalReview.Team.ReviewOrder
						resultPr.Awaiting = finalReview.Team.Name

						if finalReview.User == nil {
							unassigned = true
							resultPr.Unassigned = &unassigned
						} else {
							unassigned = false
							resultPr.Unassigned = &unassigned
						}

					} else if finalReview.Team.ReviewOrder != nil &&
						*finalReview.Team.ReviewOrder == currentPriority &&
						finalReview.User != nil {

						unassigned = false
						resultPr.Unassigned = &unassigned

					} else if resultPr.Awaiting == nil {
						resultPr.Awaiting = &teamOther
					}

				} else if finalReview.Team == nil {
					resultPr.Awaiting = &teamOther
				}

			} else if *review.State == "CHANGES_REQUESTED" {
				resultPr.Awaiting = &changesRequested
				currentPriority = -1
				unassigned = false
				resultPr.Unassigned = &unassigned

			} else if *review.State == statusApproved {
				approvedCount++
			}
			resultPr.ReviewOverview = append(resultPr.ReviewOverview, finalReview)

		}

		if resultPr.Awaiting == nil && approvedCount >= 1 {
			resultPr.Awaiting = &statusApproved
			unassigned = false
			resultPr.Unassigned = &unassigned
		}
	}

	if ErrorMessage != "" {
		ErrorMessage = ErrorMessage + ". See the console on the server or the errorText field in the network tab for more information"
		resultPr.ErrorMessage = &ErrorMessage
		resultPr.ErrorText = &ErrorText
	}

	prChannel <- resultPr

}

/**** public ****/

/*
get list of github pull requests and process them with review information
*/
func GetPullRequests(ctx context.Context, db *sql.DB, c *github.Client, owner string, RepositoryName string, prevResult *db_pkg.PullRequestInfo) (*db_pkg.PullRequestInfo, error) {

	if prevResult == nil {
		prevResult = new(db_pkg.PullRequestInfo)
	}

	if prevResult.Updated == nil {
		prevResult.Updated = &time.Time{}
	}

	previousPrMap := make(map[int]*db_pkg.PullRequest)

	if prevResult.PullRequests != nil {
		for _, pullRequest := range prevResult.PullRequests {
			previousPrMap[*pullRequest.Number] = pullRequest
		}
	}

	opts := &github.PullRequestListOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var ghPrs []*github.PullRequest

	for {
		respPrs, resp, err := c.PullRequests.List(ctx, owner, RepositoryName, opts)
		if err != nil {
			return nil, err
		}

		ghPrs = append(ghPrs, respPrs...)

		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	users, err := db_pkg.GetUsersAsLoginMap(ctx, db, RepositoryName)
	if err != nil {
		return nil, err
	}

	teams, err := db_pkg.GetTeamsAsMap(ctx, db, RepositoryName)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	idx := 0
	prChannel := make(chan *db_pkg.PullRequest)
	var cachedPrList []*db_pkg.PullRequest

	for _, pr := range ghPrs {
		if *pr.Draft {
			continue
		}

		// if the pr was updated longer ago than the last fetch then we don't need to re-fetch it
		if prevResult.Updated != nil && pr.UpdatedAt.GetTime().Before(*prevResult.Updated) {
			savedPr := previousPrMap[*pr.Number]
			tempIdx := idx
			savedPr.Index = &tempIdx
			cachedPrList = append(cachedPrList, savedPr)
		} else {
			wg.Add(1)
			go processPullRequest(prChannel, &wg, ctx, c, owner, RepositoryName, pr, users, teams, idx)
		}

		idx++ // manual index as we are skipping draft
	}

	result := new(db_pkg.PullRequestInfo)

	// sort the list of teams for the aggregation banner
	if teams != nil {
		slugs := make([]string, 0, len(teams))

		for slug := range teams {
			slugs = append(slugs, slug)
		}

		sort.SliceStable(slugs, func(i, j int) bool {
			return *teams[slugs[i]].ReviewOrder < *teams[slugs[j]].ReviewOrder
		})

		sorted_teams := make([]*db_pkg.Team, 0)
		for _, slug := range slugs {
			sorted_teams = append(sorted_teams, teams[slug])
		}

		result.ReviewTeams = sorted_teams
	}

	if users != nil {
		userList := make([]*db_pkg.User, 0)

		for _, user := range users {
			userList = append(userList, user)
		}

		result.Users = userList
	}

	go func() {
		wg.Wait()
		close(prChannel)
	}()

	prs := make([]*db_pkg.PullRequest, idx)

	for processed_pr := range prChannel {
		prs[*processed_pr.Index] = processed_pr
	}

	for _, processed_pr := range cachedPrList {
		prs[*processed_pr.Index] = processed_pr
	}

	result.PullRequests = prs
	currentTime := time.Now()
	result.Updated = &currentTime

	return result, nil

}
