package github_pkg

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"golang.org/x/oauth2"

	"github.com/google/go-github/v68/github"
)

func InitGithubConnection(ctx context.Context) (*Config, *github.Client) {
	config := loadEnv()

	authToken := config.Token

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: authToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	config.Token = ""

	return config, client

}

func loadEnv() *Config {
	content, err := os.ReadFile("./db/config.json")
	if err != nil {
		log.Fatal("Error when opening config: ", err)
	}

	var payload Config
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal of config: ", err)
	}

	return &payload
}
