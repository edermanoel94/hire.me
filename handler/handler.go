package handler

import (
	"desafio_bemobi/model"
	"desafio_bemobi/repository"
	"desafio_bemobi/util"
	"fmt"
	"github.com/edermanoel94/rest-go"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"github.com/teris-io/shortid"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Handler struct {
	Repository repository.Repository
	Log        *log.Logger
}

func New(repository repository.Repository, logger *log.Logger) *Handler {
	return &Handler{
		Repository: repository,
		Log:        logger,
	}
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	err := fmt.Errorf("not found the following path: %s", r.URL.Path)
	_, _ = rest.Error(w, err, http.StatusNotFound)
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	err := fmt.Errorf("this method %s is not allowed on path: %s", r.Method, r.URL.Path)
	_, _ = rest.Error(w, err, http.StatusMethodNotAllowed)
}

func (h *Handler) ShortURLHandler(w http.ResponseWriter, r *http.Request) {

	longUrl := r.FormValue("url")

	if !util.IsUrl(longUrl) {
		rest.Error(w, util.NewError("", "003", util.InvalidUrl), http.StatusBadRequest)
		return
	}

	startTime := time.Now()

	customAlias := r.FormValue("CUSTOM_ALIAS")

	customAlias = strings.TrimSpace(customAlias)

	var shortUrl string

	baseUrl := fmt.Sprintf("http://%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))

	if customAlias != "" {
		if h.Repository.ExistByAlias(customAlias) {
			rest.Error(w, util.NewError(customAlias, "001", util.AlreadyExistsUrl), http.StatusConflict)
			return
		}
		shortUrl = baseUrl + "/" + customAlias
	} else {
		customAlias, _ = shortid.Generate()
		shortUrl = baseUrl + "/" + customAlias
	}

	totalTime := time.Since(startTime)

	u := &model.Url{
		ID:        bson.NewObjectId(),
		Original:  longUrl,
		Short:     shortUrl,
		Alias:     customAlias,
		TimeTaken: totalTime,
	}

	err := h.Repository.Create(u)

	if err != nil {
		h.Log.Printf("couln't create url: %v", err)
		rest.Error(w, err, http.StatusInternalServerError)
		return
	}

	h.Log.Printf("create url successful with CUSTOM_ALIAS=%s and short url: %s", u.Alias, u.Short)
	rest.Marshalled(w, u, http.StatusOK)
}

func (h *Handler) GetCodeAndRedirectHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	err := rest.CheckPathVariables(params, "code")

	if err != nil {
		rest.Error(w, err, http.StatusBadRequest)
		return
	}

	code := params["code"]

	if !h.Repository.ExistByAlias(code) {
		rest.Error(w, util.NewError(code, "002", util.NotFoundUrl), http.StatusNotFound)
		return
	}

	url := model.Url{}

	err = h.Repository.FindByAlias(code, &url)

	if err != nil {
		h.Log.Printf("couln't find url: %v", err)
		rest.Error(w, err, http.StatusInternalServerError)
		return
	}

	// Visited
	url.Visited++
	err = h.Repository.Update(url.ID.Hex(), &url)

	if err != nil {
		h.Log.Printf("couln't update url: %v", err)
		rest.Error(w, err, http.StatusInternalServerError)
		return
	}

	h.Log.Printf("redirect url successful to %s", url.Original)
	http.Redirect(w, r, url.Original, http.StatusTemporaryRedirect)
}

func (h *Handler) MoreVisitedURLHandler(w http.ResponseWriter, r *http.Request) {

	urls := make([]model.Url, 0)

	err := h.Repository.MoreVisited(&urls)

	if err != nil {
		h.Log.Printf("couln't get more visited url's: %v", err)
		rest.Error(w, err, http.StatusInternalServerError)
		return
	}

	rest.Marshalled(w, &urls, http.StatusOK)
}
