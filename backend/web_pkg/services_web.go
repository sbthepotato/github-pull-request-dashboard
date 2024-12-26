package web_pkg

import (
	"sync"
	"time"

	"github-pull-request-dashboard/github_pkg"
)

var mu sync.Mutex

var last_fetched_teams time.Time
var cached_teams []*github_pkg.CustomTeam

var last_fetched_members time.Time
var cached_members []*github_pkg.CustomUser

var last_fetched_prs time.Time
var cached_prs *github_pkg.PullRequestInfo
