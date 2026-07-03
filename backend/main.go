package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github-pull-request-dashboard/db_pkg"
	"github-pull-request-dashboard/github_pkg"
	"github-pull-request-dashboard/web_pkg"
)

//go:embed all:static
var staticFiles embed.FS

func main() {
	ctx := context.Background()

	// start the database and create tables if they dont exist
	db, err := db_pkg.InitDatabase(ctx)
	if err != nil {
		log.Fatalln("Could not start the database: ", err.Error())
	}

	defer db.Close()

	// connect to github using the env
	client, owner, defaultRepository, err := github_pkg.InitGithubConnection(ctx)
	if err != nil {
		log.Fatalln("Could not start up github connection: ", err.Error())
	}

	// GETS
	http.HandleFunc("/api/config/hello_go", web_pkg.HelloGo)
	http.HandleFunc("/api/config/rate_limit", web_pkg.GetRateLimit(ctx, client))
	http.HandleFunc("/api/config/get_repos", web_pkg.GetRepositories(ctx, db, client, owner))
	http.HandleFunc("/api/config/get_default_repository", web_pkg.GetDefaultRepository(ctx, defaultRepository))
	http.HandleFunc("/api/config/get_teams", web_pkg.GetTeams(ctx, db, client, owner, defaultRepository))
	http.HandleFunc("/api/config/get_users", web_pkg.GetUsers(ctx, db, client, owner, defaultRepository))
	http.HandleFunc("/api/config/get_title_regex_list", web_pkg.GetTitleRegexList(ctx, db))
	http.HandleFunc("/api/dashboard/get_pr_list", web_pkg.GetPullRequests(ctx, db, client, owner, defaultRepository))

	// POSTS
	http.HandleFunc("/api/config/set_repos", web_pkg.SetRepositories(ctx, db))
	http.HandleFunc("/api/config/set_teams", web_pkg.SetTeams(ctx, db))
	http.HandleFunc("/api/config/set_regex", web_pkg.SetTitleRegex(ctx, db))
	http.HandleFunc("/api/config/delete_regex", web_pkg.DeleteTitleRegex(ctx, db))

	// Serve embedded frontend static files with SPA fallback
	staticSubFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalln("Could not create static file system: ", err.Error())
	}
	fileServer := http.FileServer(http.FS(staticSubFS))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "index.html"
		} else {
			path = path[1:] // strip leading /
		}
		_, err := staticSubFS.Open(path)
		if err != nil {
			// SPA fallback: serve index.html for unknown routes
			http.ServeFileFS(w, r, staticSubFS, "index.html")
			return
		}
		fileServer.ServeHTTP(w, r)
	})

	cors_handler := web_pkg.EnableCors(http.DefaultServeMux)

	// Start the server on port 8080
	log.Println("Starting server on :8080...")
	if err := http.ListenAndServe(":8080", cors_handler); err != nil {
		log.Fatalln("Could not start server: ", err.Error())
	}
}
