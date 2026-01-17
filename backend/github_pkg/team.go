package github_pkg

import (
	"context"
	"database/sql"
	"github-pull-request-dashboard/db_pkg"

	"github.com/google/go-github/v81/github"
)

/*
get list of all teams for a given organisation
*/
func GetTeams(ctx context.Context, db *sql.DB, c *github.Client, owner string) ([]*db_pkg.Team, error) {

	opt := &github.ListOptions{
		PerPage: 100,
	}

	var allTeams []*github.Team
	var err error

	for {
		respTeam, resp, err := c.Teams.ListTeams(ctx, owner, opt)
		if err != nil {
			return nil, err
		}

		allTeams = append(allTeams, respTeam...)
		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage

	}

	teams := make([]*db_pkg.Team, 0)

	for _, team := range allTeams {
		customTeam := new(db_pkg.Team)
		customTeam.Team = team

		teams = append(teams, customTeam)
	}

	err = db_pkg.CreateTeams(ctx, db, teams)
	if err != nil {
		return nil, err
	}

	return teams, nil
}
