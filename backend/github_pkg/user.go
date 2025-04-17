package github_pkg

import (
	"context"
	"database/sql"
	"github-pull-request-dashboard/db_pkg"
	"log"
	"sync"

	"github.com/google/go-github/v71/github"
)

/**** private ****/

/*
process a member into the member channel
*/
func processUser(userChannel chan<- *db_pkg.User, wg *sync.WaitGroup, ctx context.Context, c *github.Client, login string) {
	defer wg.Done()

	ghUser, _, err := c.Users.Get(ctx, login)
	if err != nil {
		log.Println("error fetching user", login, err)
	}

	user := new(db_pkg.User)
	user.User = ghUser

	userChannel <- user
}

/**** public ****/

/*
get users and link them up to one of the active teams
*/
func GetUsers(ctx context.Context, db *sql.DB, c *github.Client, owner string) ([]*db_pkg.User, error) {

	listMembersOpt := &github.ListMembersOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	allUsers := make([]*github.User, 0)
	var err error

	for {
		respUsers, resp, err := c.Organizations.ListMembers(ctx, owner, listMembersOpt)
		if err != nil {
			return nil, err
		}

		allUsers = append(allUsers, respUsers...)

		if resp.NextPage == 0 {
			break
		}
		listMembersOpt.Page = resp.NextPage
	}

	var wg sync.WaitGroup
	userChannel := make(chan *db_pkg.User)

	for _, user := range allUsers {
		wg.Add(1)
		go processUser(userChannel, &wg, ctx, c, *user.Login)
	}

	go func() {
		wg.Wait()
		close(userChannel)
	}()

	finalUsers := make([]*db_pkg.User, 0)

	for user := range userChannel {
		finalUsers = append(finalUsers, user)
	}

	err = db_pkg.CreateUsers(ctx, db, finalUsers)
	if err != nil {
		return nil, err
	}

	return finalUsers, nil
}

/*
get the users that are in a team for a given repository
return them as map where team slug is the key and list of users is the value
*/
func GetUserTeams(ctx context.Context, db *sql.DB, c *github.Client, owner string, repositoryName string) (map[string][]*github.User, error) {

	teams, err := db_pkg.GetTeams(ctx, db, repositoryName)
	if err != nil {
		return nil, err
	}

	userTeams := make(map[string][]*github.User)

	for _, team := range teams {
		if *team.ReviewOrder > 0 {
			teamMembersOpt := &github.TeamListTeamMembersOptions{
				ListOptions: github.ListOptions{PerPage: 100},
			}

			allMembers := make([]*github.User, 0)

			for {
				respMembers, resp, err := c.Teams.ListTeamMembersBySlug(ctx, owner, *team.Slug, teamMembersOpt)
				if err != nil {
					return nil, err
				}

				allMembers = append(allMembers, respMembers...)

				if resp.NextPage == 0 {
					break
				}
				teamMembersOpt.Page = resp.NextPage
			}

			userTeams[*team.Slug] = allMembers
		}
	}

	err = db_pkg.UpsertUserTeams(ctx, db, userTeams, repositoryName)
	if err != nil {
		return nil, err
	}

	return userTeams, nil
}
