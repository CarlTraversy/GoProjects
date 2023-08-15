package services

import (
	"errors"
	"io"
	"net/http"
	"regexp"

	"Collectivei.GoProjects/src/domain"
)

type ProjectsService interface {
	GetAll() ([]domain.Project, error)
	FindAll(criteria string) ([]domain.Project, error)
}

type GithubProjectsService struct {
	sourceUrl string
}

func NewGithubProjectsService(sourceUrl string) GithubProjectsService {
	return GithubProjectsService{sourceUrl: sourceUrl}
}

func (projectService GithubProjectsService) GetAll() ([]domain.Project, error) {
	readme, err := getReadmeContent(projectService.sourceUrl)
	if err != nil {
		return nil, err
	}

	var projects []domain.Project
	regex := regexp.MustCompile(`\[.*?\]\((https://[www\.]*github\.com/[^/]*/[^/]*/*)\)`)
	matches := regex.FindAllStringSubmatch(readme, -1)
	for _, match := range matches {
		projects = append(projects, domain.Project{Url: match[1]})
	}

	return projects, nil
}

func (projectService GithubProjectsService) FindAll(criteria string) ([]domain.Project, error) {
	readme, err := getReadmeContent(projectService.sourceUrl)
	if err != nil {
		return nil, err
	}

	var projects []domain.Project
	regex := regexp.MustCompile(`\[.*?\]\((https://[www\.]*github\.com/[^/]*/[^/]*` + regexp.QuoteMeta(criteria) + `[^/]*/*)\)`)
	matches := regex.FindAllStringSubmatch(readme, -1)
	for _, match := range matches {
		projects = append(projects, domain.Project{Url: match[1]})
	}

	return projects, nil
}

func getReadmeContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Github unreachable")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
