package github_pkg

import (
	"time"

	"github.com/google/go-github/v68/github"
)

/*
repository with enabled field
*/
type CustomRepo struct {
	*github.Repository
	Enabled *bool `json:"enabled,omitempty"`
}

/*
Team with review info
*/
type CustomTeam struct {
	*github.Team
	ReviewEnabled *bool `json:"review_enabled,omitempty"`
	ReviewOrder   *int  `json:"review_order,omitempty"`
}

/*
User with a Custom Team attached
*/
type CustomUser struct {
	*github.User
	Team *CustomTeam `json:"team,omitempty"`
}

type Review struct {
	User  *CustomUser `json:"user,omitempty"`
	Team  *CustomTeam `json:"team,omitempty"`
	State *string     `json:"state,omitempty"`
}

/*
Pull Request with extra fields for custom objects
*/
type CustomPullRequest struct {
	*github.PullRequest
	CreatedBy      *CustomUser `json:"created_by,omitempty"`
	ReviewOverview []*Review   `json:"review_overview,omitempty"`
	Awaiting       *string     `json:"awaiting,omitempty"`
	Unassigned     *bool       `json:"unassigned,omitempty"`
	ErrorMessage   *string     `json:"error_message,omitempty"`
	ErrorText      *string     `json:"error_text,omitempty"`
	Index          *int        `json:"-"`
}

/*
pull request list with accompanying information for list
*/
type PullRequestInfo struct {
	PullRequests []*CustomPullRequest `json:"pull_requests,omitempty"`
	ReviewTeams  []*CustomTeam        `json:"review_teams,omitempty"`
	Users        []*CustomUser        `json:"users,omitempty"`
	Updated      *time.Time           `json:"updated,omitempty"`
}
