package services

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"GoProjects/src/domain"

	"github.com/stretchr/testify/assert"
)

func TestGetAllShouldReturnProjectsWhenReadmeContainsValidProjectUrls(t *testing.T) {
	response := `## Projects
	Some text
	- [Contents](#contents)
	[anotherLink](https://github.com/avelino/awesome-go/raw/main/tmpl/assets/logo.png)
	- [Project1](https://github.com/user/project1) - Desc1
	- [Project2](https://www.github.com/user2/project2) - Desc2`
	serverStub := httptest.NewServer(makeHTTPHandlerStub(http.StatusOK, response))
	defer serverStub.Close()
	projectsService := makeProjectsService(serverStub.URL)
	expectedResult := []domain.Project{
		{Url: "https://github.com/user/project1"},
		{Url: "https://www.github.com/user2/project2"},
	}

	actual, err := projectsService.GetAll()

	assert.Nil(t, err)
	assert.ElementsMatch(t, expectedResult, actual)
}

func TestGetAllShouldReturnProjectsWhenTwoProjectsOnSameLine(t *testing.T) {
	response := `- [Project1](https://github.com/user/project1) - Desc1 and - [Project2](https://www.github.com/user2/project2)- Desc2`
	serverStub := httptest.NewServer(makeHTTPHandlerStub(http.StatusOK, response))
	defer serverStub.Close()
	projectsService := makeProjectsService(serverStub.URL)
	expectedResult := []domain.Project{
		{Url: "https://github.com/user/project1"},
		{Url: "https://www.github.com/user2/project2"},
	}

	actual, err := projectsService.GetAll()

	assert.Nil(t, err)
	assert.ElementsMatch(t, expectedResult, actual)
}

func TestGetAllShouldReturnProjectsWhenEndingWithSlash(t *testing.T) {
	response := `- [Project1](https://github.com/user/project1/) - Desc1`
	serverStub := httptest.NewServer(makeHTTPHandlerStub(http.StatusOK, response))
	defer serverStub.Close()
	projectsService := makeProjectsService(serverStub.URL)
	expectedResult := []domain.Project{
		{Url: "https://github.com/user/project1/"},
	}

	actual, err := projectsService.GetAll()

	assert.Nil(t, err)
	assert.ElementsMatch(t, expectedResult, actual)
}

func TestGetAllShouldReturnNilWhenReadmeDoesNotContainValidProjectUrls(t *testing.T) {
	response := `## No projects here
	Some text
	https://github.com/avelino/awesome-go/raw/main/tmpl/assets/logo.png
	- [NotAProject](https://somewbsite.com) - Desc1
	- [NotAProjectAgain](https://www.someOtherWebsite)- Desc2`
	serverStub := httptest.NewServer(makeHTTPHandlerStub(http.StatusOK, response))
	defer serverStub.Close()
	projectsService := makeProjectsService(serverStub.URL)

	actual, err := projectsService.GetAll()

	assert.Nil(t, err)
	assert.Nil(t, actual)
}

func TestGetAllShouldReturnErrorWhenGetNotStatusCodeOk(t *testing.T) {
	serverStub := httptest.NewServer(makeHTTPHandlerStub(http.StatusRequestTimeout, ""))
	defer serverStub.Close()
	projectsService := makeProjectsService(serverStub.URL)

	_, err := projectsService.GetAll()

	assert.NotNil(t, err)
}

func TestGetAllShouldReturnErrorWhenGetReturnsError(t *testing.T) {
	serverStub := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		http.Error(response, "", http.StatusRequestTimeout)
	}))
	projectsService := makeProjectsService(serverStub.URL)

	_, err := projectsService.GetAll()

	assert.NotNil(t, err)
}

func TestFindAllShouldReturnSomeProjectsWhenPassedCriteriaMatches(t *testing.T) {
	response := `## Projects
	Some text
	- [Contents](#contents)
	https://github.com/avelino/awesome-go/raw/main/tmpl/assets/logo.png
	- [super Project1](https://github.com/super/superproject1) - super Desc1
	- [Project2 super](https://www.github.com/user2/project2super)- super Desc2`
	serverStub := httptest.NewServer(makeHTTPHandlerStub(http.StatusOK, response))
	defer serverStub.Close()
	projectsService := makeProjectsService(serverStub.URL)
	expectedResult := []domain.Project{
		{Url: "https://github.com/super/superproject1"},
		{Url: "https://www.github.com/user2/project2super"},
	}

	actual, err := projectsService.FindAll("super")

	assert.Nil(t, err)
	assert.ElementsMatch(t, expectedResult, actual)
}

func TestFindAllShouldReturnNilWhenNoProjectMatchCriteria(t *testing.T) {
	response := `## No projects here
	Some text
	https://github.com/avelino/awesome-go/raw/main/tmpl/assets/logo.png
	- [NotAProject](https://somewbsite.com) - Desc1
	- [NotAProjectAgain](https://www.someOtherWebsite)- Desc2`
	serverStub := httptest.NewServer(makeHTTPHandlerStub(http.StatusOK, response))
	defer serverStub.Close()
	projectsService := makeProjectsService(serverStub.URL)

	actual, err := projectsService.FindAll("thisDoesntMatchAnything")

	assert.Nil(t, err)
	assert.Nil(t, actual)
}

func TestFindAllShouldReturnErrorWhenGetNotStatusCodeOk(t *testing.T) {
	serverStub := httptest.NewServer(makeHTTPHandlerStub(http.StatusRequestTimeout, ""))
	defer serverStub.Close()
	projectsService := makeProjectsService(serverStub.URL)

	_, err := projectsService.FindAll("")

	assert.NotNil(t, err)
}

func TestFindAllShouldReturnErrorWhenGetReturnsError(t *testing.T) {
	serverStub := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		http.Error(response, "", http.StatusRequestTimeout)
	}))
	projectsService := makeProjectsService(serverStub.URL)

	_, err := projectsService.FindAll("")

	assert.NotNil(t, err)
}

func makeProjectsService(url string) GithubProjectsService {
	return NewGithubProjectsService(url)
}

func makeHTTPHandlerStub(statusCode int, responseBody string) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "text/plain")
		response.WriteHeader(statusCode)
		_, _ = response.Write([]byte(responseBody))
	}
}
