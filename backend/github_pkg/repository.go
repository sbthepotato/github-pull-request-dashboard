package github_pkg

import (
	"context"

	"github.com/google/go-github/v68/github"
)

/*
get all repositories for the currently set org
*/
func GetRepositories(ctx context.Context, c *github.Client, owner string) ([]*CustomRepo, error) {

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

	var custom_repos []*CustomRepo

	for _, repo := range all_repos {
		custom_repo := new(CustomRepo)

		custom_repo.Repository = repo
		custom_repos = append(custom_repos, custom_repo)

	}

	return custom_repos, nil
}
