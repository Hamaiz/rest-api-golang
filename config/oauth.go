package serve

import (
	"github.com/Hamaiz/go-rest-eg/api"
	"github.com/Hamaiz/go-rest-eg/database"
	"github.com/Hamaiz/go-rest-eg/helper"
	"github.com/Hamaiz/go-rest-eg/session"
	"github.com/globalsign/mgo"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

// NewOauthSubRouter - oauth accounts subrouter
func NewOauthSubRouter(s *mux.Router, dbsess *mgo.Session, conn *pgxpool.Pool) {
	// getting store
	store := session.StoreConn(dbsess)
	newoauth := database.NewOauthDatabase(conn)

	// newaccountstore sending store
	o := api.NewOauthApi(store, newoauth)

	// routes - /accounts
	s.HandleFunc("/google/login", helper.JH(o.GoogleLoginHandler))
	s.HandleFunc("/google/callback", helper.JH(o.GoogleCallbackHandler))
}
