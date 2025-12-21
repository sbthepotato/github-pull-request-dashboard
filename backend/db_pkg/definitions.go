package db_pkg

import (
	"time"

	"github.com/google/go-github/v80/github"
)

/*
	The following definitions don't fit in any of the normal packages
*/

type Review struct {
	User  *User   `json:"user,omitempty"`
	Team  *Team   `json:"team,omitempty"`
	State *string `json:"state,omitempty"`
}

type PullRequest struct {
	*github.PullRequest
	HtmlTitle      *string   `json:"html_title,omitempty"`
	CreatedBy      *User     `json:"created_by,omitempty"`
	ReviewOverview []*Review `json:"review_overview,omitempty"`
	Awaiting       *string   `json:"awaiting,omitempty"`
	Unassigned     *bool     `json:"unassigned,omitempty"`
	Error          *string   `json:"error,omitempty"`
	Index          *int      `json:"-"`
}

type PullRequestInfo struct {
	PullRequests []*PullRequest `json:"pull_requests,omitempty"`
	ReviewTeams  []*Team        `json:"review_teams,omitempty"`
	Users        []*User        `json:"users,omitempty"`
	Updated      *time.Time     `json:"updated,omitempty"`
}
