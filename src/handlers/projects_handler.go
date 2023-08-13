package handlers

import (
	"encoding/json"
	"net/http"

	"Collectivei.GoProjects/src/services"
)

type ProjectsHandler struct {
	ProjectService services.ProjectsService
}

func NewProjectsHandler(projectsService services.ProjectsService) http.Handler {
	return &ProjectsHandler{ProjectService: projectsService}
}

func (handler *ProjectsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	projects, err := handler.ProjectService.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"projects": projects})
}
