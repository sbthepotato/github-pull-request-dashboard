package github_pkg

import (
	"context"
	"github-pull-request-dashboard/db_pkg"

	"github.com/google/go-github/v68/github"
)

/*
get all repositories for the currently set org
*/
func GetRepositories(ctx context.Context, c *github.Client, owner string) ([]*db_pkg.Repository, error) {

	opt := &github.RepositoryListByOrgOptions{
		Sort:        "full_name",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var all_repos []*github.Repository

	for {
		repos, resp, err := c.Repositories.ListByOrg(ctx, owner, opt)
		if err != nil {
			return nil, err
		}

		all_repos = append(all_repos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage

	}

	var custom_repos []*db_pkg.Repository

	for _, repo := range all_repos {
		custom_repo := new(db_pkg.Repository)
		custom_repo.Repository = repo
		db_pkg.CreateRepository(ctx, custom_repo)
		custom_repos = append(custom_repos, custom_repo)
	}

	return custom_repos, nil
}
