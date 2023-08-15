package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"Collectivei.GoProjects/src/domain"
	"Collectivei.GoProjects/src/services"
	"github.com/stretchr/testify/assert"
)

func TestServeHttpShouldReturnsInternalServerErrorWhenGetAllReturnsError(t *testing.T) {
	handler := makeProjectsHandler(makeProjectsServiceStub(nil, nil, errors.New("someError")))

	assert.HTTPStatusCode(t, handler.ServeHTTP, http.MethodGet, ProjectsHandlerRoute, nil, http.StatusInternalServerError)
}

func TestServeHttpShouldReturnsInternalServerErrorWhenFindAllReturnsError(t *testing.T) {
	handler := makeProjectsHandler(makeProjectsServiceStub(nil, nil, errors.New("someError")))
	values := url.Values{}
	values.Add("name", "criteria")

	assert.HTTPStatusCode(t, handler.ServeHTTP, http.MethodGet, ProjectsHandlerRoute, values, http.StatusInternalServerError)
}

func TestServeHttpShouldReturnsAllProjectsWhenPassedInvalidParam(t *testing.T) {
	expectedUrl := "https://github.com/user/expectedProject"
	handler := makeProjectsHandler(makeProjectsServiceStub([]domain.Project{{Url: expectedUrl}}, nil, nil))
	request := httptest.NewRequest(http.MethodGet, ProjectsHandlerRoute+"?invalidParam", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	var actual domain.ProjectsResponse
	err := json.NewDecoder(response.Body).Decode(&actual)

	assert.Nil(t, err)
	assert.Equal(t, expectedUrl, actual.Projects[0].Url)
}

func TestServeHttpShouldCallFindAllWhenNameParam(t *testing.T) {
	expectedUrl := "https://github.com/user/expectedProject"
	handler := makeProjectsHandler(makeProjectsServiceStub(nil, []domain.Project{{Url: expectedUrl}}, nil))
	request := httptest.NewRequest(http.MethodGet, ProjectsHandlerRoute+"?name=someValue", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	var actual domain.ProjectsResponse
	err := json.NewDecoder(response.Body).Decode(&actual)

	assert.Nil(t, err)
	assert.Equal(t, expectedUrl, actual.Projects[0].Url)
}

func TestServeHttpShouldCallGetAllWhenNoQueryParam(t *testing.T) {
	expectedUrl := "https://github.com/user/expectedProject"
	handler := makeProjectsHandler(makeProjectsServiceStub([]domain.Project{{Url: expectedUrl}}, nil, nil))
	request := httptest.NewRequest(http.MethodGet, ProjectsHandlerRoute, nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	var actual domain.ProjectsResponse
	err := json.NewDecoder(response.Body).Decode(&actual)

	assert.Nil(t, err)
	assert.Equal(t, expectedUrl, actual.Projects[0].Url)
}

func TestServeHttpShouldReturnEmptySliceWhenNoProjects(t *testing.T) {
	handler := makeProjectsHandler(makeProjectsServiceStub(nil, nil, nil))
	request := httptest.NewRequest(http.MethodGet, ProjectsHandlerRoute, nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	var actual domain.ProjectsResponse
	err := json.NewDecoder(response.Body).Decode(&actual)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(actual.Projects))
}

func TestServeHttpShouldReturnEmptySliceWhenNoProjectsMatchCriteria(t *testing.T) {
	handler := makeProjectsHandler(makeProjectsServiceStub(nil, nil, nil))
	request := httptest.NewRequest(http.MethodGet, ProjectsHandlerRoute+"?name=someCriteria", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	var actual domain.ProjectsResponse
	err := json.NewDecoder(response.Body).Decode(&actual)

	assert.Nil(t, err)
	assert.Equal(t, 0, len(actual.Projects))
}

func makeProjectsHandler(projectsService services.ProjectsService) ProjectsHandler {
	return NewProjectsHandler(projectsService)
}

func makeProjectsServiceStub(getAllResult []domain.Project, findAllResult []domain.Project, err error) services.ProjectsService {
	return projectsServiceStub{getAllResult: getAllResult, findAllResult: findAllResult, err: err}
}

type projectsServiceStub struct {
	getAllResult  []domain.Project
	findAllResult []domain.Project
	err           error
}

func (stub projectsServiceStub) GetAll() ([]domain.Project, error) {
	return stub.getAllResult, stub.err
}

func (stub projectsServiceStub) FindAll(criteria string) ([]domain.Project, error) {
	return stub.findAllResult, stub.err
}
