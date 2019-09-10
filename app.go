package main

import (
	"desafio_bemobi/handler"
	"desafio_bemobi/repository"
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	logStderr = log.New(os.Stderr, "bemobi_error: ", log.Lshortfile|log.LstdFlags|log.Ltime)
	logStdout = log.New(os.Stdout, "bemobi: ", log.Lshortfile|log.LstdFlags|log.Ltime)
)

func main() {

	// Create DNS's
	appDNS := fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
	mongoDNS := fmt.Sprintf("%s:%s", os.Getenv("MONGODB_HOST"), os.Getenv("MONGODB_PORT"))

	// Connect on MongoDB
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    []string{mongoDNS},
		Database: os.Getenv("MONGODB_DATABASE"),
		Username: os.Getenv("MONGODB_USER"),
		Password: os.Getenv("MONGODB_PASS"),
		Timeout:  2 * time.Minute,
		AppName:  "Desafio Bemobi",
	})

	if err != nil {
		log.Fatalf("couldn't connect on mongo server: %v", err)
	}

	// Get a DB instance
	db := session.DB(os.Getenv("MONGO_DB"))

	// Create index field for `alias`, when running application again, don't create again,
	// if already exists.
	createIndexes(db)

	// Initialize repository
	repo := repository.New(db)

	// Initialize handlers
	handlers := handler.New(repo, logStdout)

	// Setup routers with handlers
	routers := setupRouters(handlers)

	srv := http.Server{
		Addr:         appDNS,
		Handler:      cors.AllowAll().Handler(routers),
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 1 * time.Minute,
		ErrorLog:     logStderr,
	}

	log.Fatalf("couldn't listen: %v", srv.ListenAndServe())
}

// createIndexes create a index if not exists
func createIndexes(db *mgo.Database) {
	indexes := mgo.Index{
		Key:        []string{"alias"},
		Background: true,
		Unique:     true,
		Name:       "alias_idx",
	}
	_ = db.C(repository.Collection).EnsureIndex(indexes)
}

// setupRouters create all routing for application
func setupRouters(h *handler.Handler) *mux.Router {

	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(handler.NotFoundHandler)
	router.MethodNotAllowedHandler = http.HandlerFunc(handler.MethodNotAllowedHandler)

	router.HandleFunc("/create", h.ShortURLHandler).Methods("POST")
	router.HandleFunc("/moreVisited", h.MoreVisitedURLHandler).Methods("GET")
	router.HandleFunc("/{code}", h.GetCodeAndRedirectHandler).Methods("GET")

	return router
}
