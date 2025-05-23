package github_pkg

import (
	"context"
	"os"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v72/github"
	"github.com/joho/godotenv"
)

func InitGithubConnection(ctx context.Context) (*github.Client, string, string, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, "", "", err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("token")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return client, os.Getenv("owner"), os.Getenv("repo"), nil

}
