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

// NewAccountSubRouter - accounts subrouter
func NewAccountSubRouter(s *mux.Router, dbsess *mgo.Session, conn *pgxpool.Pool) {
	// getting store
	store := session.StoreConn(dbsess)
	newAccount := database.NewAccountDatabase(conn)

	// newaccountstore sending store
	a := api.NewAccountStore(store, newAccount)

	// Routes - /accounts
	s.HandleFunc("/getUser", helper.JH(a.GetUserHandler))
	s.HandleFunc("/login", helper.JH(a.LogInHandler))
	s.HandleFunc("/signup", helper.JH(a.SignUpHandler))
	s.HandleFunc("/logout", helper.JH(a.LogoutHandler))
	s.HandleFunc("/confirm/{token}", helper.JH(a.ConfirmEmailHandler))
	s.HandleFunc("/againemail", helper.JH(a.EmailAgain))
}
