package main

import (
	"context"
	"log"
	"net/http"

	"github-pull-request-dashboard/github_pkg"
	"github-pull-request-dashboard/web_pkg"
)

func main() {
	ctx := context.Background()

	client, owner, repo, err := github_pkg.InitGithubConnection(ctx)
	if err != nil {
		log.Fatalln("Could not start up github connection: ", err.Error())
	}

	// GETS
	http.HandleFunc("/config/hello_go", web_pkg.HelloGo)
	http.HandleFunc("/config/get_repos", web_pkg.GetRepos(ctx, client, owner))
	http.HandleFunc("/config/get_teams", web_pkg.GetTeams(ctx, client, owner))
	http.HandleFunc("/config/get_members", web_pkg.GetMembers(ctx, client, owner))
	http.HandleFunc("/dashboard/get_pr_list", web_pkg.GetPrList(ctx, client, owner, repo))

	// POSTS
	http.HandleFunc("/config/set_teams", web_pkg.SetTeams)
	http.HandleFunc("/config/set_repos", web_pkg.SetRepos)

	cors_handler := web_pkg.EnableCors(http.DefaultServeMux)

	// Start the server on port 8080
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", cors_handler); err != nil {
		log.Fatalln("Could not start server: ", err.Error())
	}
}
