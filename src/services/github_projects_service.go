package services

import (
	"io"
	"net/http"

	"Collectivei.GoProjects/src/domain"
	"Collectivei.GoProjects/src/github"
)

type ProjectsService interface {
	GetAll() ([]domain.Project, error)
}

type GithubProjectsService struct {
	sourceUrl          string
	githubReadmeParser github.GithubReadmeProjectsUrlParser
}

func NewGithubProjectsService(sourceUrl string, readmeParser github.GithubReadmeProjectsUrlParser) GithubProjectsService {
	return GithubProjectsService{sourceUrl: sourceUrl, githubReadmeParser: readmeParser}
}

func (projectService GithubProjectsService) GetAll() ([]domain.Project, error) {
	resp, err := http.Get(projectService.sourceUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	urls := projectService.githubReadmeParser.ParseReadme(string(data))
	var projects []domain.Project
	for _, url := range urls {
		projects = append(projects, domain.Project{Url: url})
	}

	return projects, nil
}
