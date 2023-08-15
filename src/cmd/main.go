package main

import (
	"net/http"

	"GoProjects/src/handlers"
	"GoProjects/src/services"
)

func main() {
	githubProjectsService := services.NewGithubProjectsService("https://raw.githubusercontent.com/avelino/awesome-go/master/README.md")
	projectsHandler := handlers.NewProjectsHandler(githubProjectsService)

	http.Handle(handlers.ProjectsHandlerRoute, projectsHandler)

	http.ListenAndServe(":8080", nil)
}
