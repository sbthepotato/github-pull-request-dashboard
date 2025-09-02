package github_pkg

import (
	"context"
	"database/sql"
	"github-pull-request-dashboard/db_pkg"

	"github.com/google/go-github/v74/github"
)

/*
get all repositories for the currently set org
*/
func GetRepositories(ctx context.Context, db *sql.DB, c *github.Client, owner string) ([]*db_pkg.Repository, error) {

	opt := &github.RepositoryListByOrgOptions{
		Sort:        "full_name",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var allRepos []*github.Repository

	for {
		respRepo, resp, err := c.Repositories.ListByOrg(ctx, owner, opt)
		if err != nil {
			return nil, err
		}

		allRepos = append(allRepos, respRepo...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage

	}

	repositories := make([]*db_pkg.Repository, 0)
	enabled := false

	for _, repo := range allRepos {
		repository := new(db_pkg.Repository)
		repository.Repository = repo
		repository.Enabled = &enabled
		repositories = append(repositories, repository)
	}

	db_pkg.CreateRepositories(ctx, db, repositories)

	return repositories, nil
}
