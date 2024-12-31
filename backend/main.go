package main

import (
	"context"
	"log"
	"net/http"

	"github-pull-request-dashboard/db_pkg"
	"github-pull-request-dashboard/github_pkg"
	"github-pull-request-dashboard/web_pkg"
)

func main() {
	ctx := context.Background()

	db, err := db_pkg.InitDatabase(ctx)
	if err != nil {
		log.Fatalln("Could not start the database: ", err.Error())
	}

	defer db.Close()

	client, owner, defaultRepository, err := github_pkg.InitGithubConnection(ctx)
	if err != nil {
		log.Fatalln("Could not start up github connection: ", err.Error())
	}

	// GETS
	http.HandleFunc("/config/hello_go", web_pkg.HelloGo)
	http.HandleFunc("/config/get_repos", web_pkg.GetRepositories(ctx, db, client, owner))
	http.HandleFunc("/config/get_default_repository", web_pkg.GetDefaultRepository(ctx, defaultRepository))
	http.HandleFunc("/config/get_teams", web_pkg.GetTeams(ctx, db, client, owner, defaultRepository))
	//http.HandleFunc("/config/get_members", web_pkg.GetMembers(ctx, client, owner))
	http.HandleFunc("/dashboard/get_pr_list", web_pkg.GetPullRequests(ctx, db, client, owner, defaultRepository))

	// POSTS
	http.HandleFunc("/config/set_repos", web_pkg.SetRepositories(ctx, db))
	http.HandleFunc("/config/set_teams", web_pkg.SetTeams(ctx, db))

	cors_handler := web_pkg.EnableCors(http.DefaultServeMux)

	// Start the server on port 8080
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", cors_handler); err != nil {
		log.Fatalln("Could not start server: ", err.Error())
	}
}
