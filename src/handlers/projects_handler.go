package handlers

import (
	"encoding/json"
	"net/http"

	"Collectivei.GoProjects/src/domain"
	"Collectivei.GoProjects/src/services"
)

type ProjectsHandler struct {
	ProjectService services.ProjectsService
}

func NewProjectsHandler(projectsService services.ProjectsService) http.Handler {
	return &ProjectsHandler{ProjectService: projectsService}
}

func (handler *ProjectsHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	filter := request.URL.Query().Get("name")
	var (
		projects []domain.Project
		err      error
	)

	if len(filter) != 0 {
		projects, err = handler.ProjectService.FindAll(filter)
	} else {
		projects, err = handler.ProjectService.GetAll()
	}

	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(map[string]interface{}{"projects": projects})
}
