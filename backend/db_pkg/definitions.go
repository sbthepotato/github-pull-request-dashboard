package db_pkg

import (
	"time"

	"github.com/google/go-github/v68/github"
)

/*
repository with enabled field
*/
type Repository struct {
	*github.Repository
	Enabled *bool `json:"enabled,omitempty"`
}

/*
Team with review info
*/
type Team struct {
	*github.Team
	ReviewOrder *int `json:"review_order,omitempty"`
}

/*
User with a Custom Team attached
*/
type User struct {
	*github.User
	Team *Team `json:"team,omitempty"`
}

type Review struct {
	User  *User   `json:"user,omitempty"`
	Team  *Team   `json:"team,omitempty"`
	State *string `json:"state,omitempty"`
}

/*
Pull Request with extra fields for custom objects
*/
type PullRequest struct {
	*github.PullRequest
	CreatedBy      *User     `json:"created_by,omitempty"`
	ReviewOverview []*Review `json:"review_overview,omitempty"`
	Awaiting       *string   `json:"awaiting,omitempty"`
	Unassigned     *bool     `json:"unassigned,omitempty"`
	ErrorMessage   *string   `json:"error_message,omitempty"`
	ErrorText      *string   `json:"error_text,omitempty"`
	Index          *int      `json:"-"`
}

/*
Pull Request list with accompanying information for list
*/
type PullRequestInfo struct {
	PullRequests []*PullRequest `json:"pull_requests,omitempty"`
	ReviewTeams  []*Team        `json:"review_teams,omitempty"`
	Users        []*User        `json:"users,omitempty"`
	Updated      *time.Time     `json:"updated,omitempty"`
}
