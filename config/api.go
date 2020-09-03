package serve

import (
	//	"log"
	"net/http"

	"github.com/Hamaiz/go-rest-eg/database"
	"github.com/Hamaiz/go-rest-eg/helper"
	"github.com/Hamaiz/go-rest-eg/middleware"
	"github.com/Hamaiz/go-rest-eg/session"
	"github.com/gorilla/mux"
)

// New - starts the api
func New() (*mux.Router, error) {
	// mongodb connection session
	dbsess, err := session.DBConn()
	if err != nil {
		return nil, err
	}

	// postgresql connection
	conn, err := database.DBConn()
	if err != nil {
		return nil, err
	}

	// initializing mux router
	r := mux.NewRouter()

	// subrouter - apiAccounts
	apiAccounts := r.PathPrefix("/account").Subrouter()

	// middlewares
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.UsefulHeaders)

	// account router - /account
	NewAccountSubRouter(apiAccounts, dbsess, conn)

	// static files
	helper.AllStaticFiles(r)

	// custom not found handler
	r.NotFoundHandler = http.HandlerFunc(helper.NotFound)

	return r, nil
}
