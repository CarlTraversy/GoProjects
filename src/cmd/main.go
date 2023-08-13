package main

import (
	"net/http"

	"Collectivei.GoProjects/src/github"
	"Collectivei.GoProjects/src/handlers"
	"Collectivei.GoProjects/src/services"
)

func main() {
	githubProjectsService := services.NewGithubProjectsService("https://raw.githubusercontent.com/avelino/awesome-go/master/README.md", github.GithubReadmeProjectsUrlParser{})
	projectsHandler := handlers.NewProjectsHandler(githubProjectsService)

	http.Handle("/projects", projectsHandler)

	http.ListenAndServe(":8080", nil)
}