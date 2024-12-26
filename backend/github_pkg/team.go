package github_pkg

import (
	"context"
	"encoding/json"
	"github-pull-request-dashboard/db_pkg"

	"github.com/google/go-github/v68/github"
)

/*
get list of all teams for a given organisation
*/
func GetTeams(ctx context.Context, c *github.Client, owner string) ([]*CustomTeam, error) {

	opt := &github.ListOptions{
		PerPage: 100,
	}

	teams, _, err := c.Teams.ListTeams(ctx, owner, opt)
	if err != nil {
		return nil, err
	}

	customTeams := make([]*CustomTeam, 0)
	defaultReview := false
	defaultOrder := 0

	for _, team := range teams {
		customTeam := new(CustomTeam)
		customTeam.Team = team
		customTeam.ReviewEnabled = &defaultReview
		customTeam.ReviewOrder = &defaultOrder

		customTeams = append(customTeams, customTeam)
	}

	err = writeTeams(customTeams)
	if err != nil {
		return nil, err
	}

	return customTeams, nil
}

func writeTeams(teams []*CustomTeam) error {
	jsonData, err := json.Marshal(teams)
	if err != nil {
		return err
	}

	db_pkg.Write("teams", jsonData)

	return nil

}

func readTeams() ([]*CustomTeam, error) {
	jsonData, err := db_pkg.Read("teams")
	if err != nil {
		return nil, err
	}

	teams := make([]*CustomTeam, 0)

	err = json.Unmarshal(jsonData, teams)
	if err != nil {
		return nil, err
	}

	return teams, nil
}
