package handler_test

import (
	"desafio_bemobi/handler"
	"desafio_bemobi/handler/mock"
	"desafio_bemobi/model"
	"github.com/edermanoel94/rest-go"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	logger = log.New(os.Stdout, "bemobi_test: ", log.Lshortfile|log.LstdFlags|log.Ltime)
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal(err)
	}
}

func TestHandler_ShortURLHandler(t *testing.T) {

	repositoryMock := new(mock.RepositoryMock)

	h := handler.New(repositoryMock, logger)

	url := model.Url{}

	originalUrl := "http://google.com"

	request, _ := http.NewRequest("POST", "/create?url="+originalUrl, nil)

	recorder := httptest.NewRecorder()

	h.ShortURLHandler(recorder, request)

	result := recorder.Result()

	defer result.Body.Close()

	err := rest.GetBody(result.Body, &url)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, originalUrl, url.Original)
}

func TestHandler_GetCodeAndRedirectHandler(t *testing.T) {

	repositoryMock := new(mock.RepositoryMock)

	h := handler.New(repositoryMock, logger)

	customAlias := "eder"

	request, _ := http.NewRequest("GET", "/", nil)

	request = mux.SetURLVars(request, map[string]string{
		"code": customAlias,
	})

	recorder := httptest.NewRecorder()

	h.GetCodeAndRedirectHandler(recorder, request)

	result := recorder.Result()

	assert.Equal(t, http.StatusTemporaryRedirect, result.StatusCode)
}

func TestHandler_MoreVisitedURLHandler(t *testing.T) {

	repositoryMock := new(mock.RepositoryMock)

	h := handler.New(repositoryMock, logger)

	request, _ := http.NewRequest("GET", "/moreVisited", nil)

	recorder := httptest.NewRecorder()

	h.MoreVisitedURLHandler(recorder, request)

	result := recorder.Result()

	defer result.Body.Close()

	var urls []model.Url

	err := rest.GetBody(result.Body, &urls)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, result.StatusCode)
}
