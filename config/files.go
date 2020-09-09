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
func NewFilesSubRouter(s *mux.Router, dbsess *mgo.Session, conn *pgxpool.Pool) {
	// getting store
	store := session.StoreConn(dbsess)
	newFiles := database.NewFilesDatabase(conn)

	// newaccountstore sending store
	f := api.NewFilesApi(store, newFiles)

	// Routes - /accounts
	s.HandleFunc("/question", helper.JH(f.GetQuestionsHandler))
	s.HandleFunc("/question/{slug}", helper.JH(f.SendQuestionHandler))
	s.HandleFunc("/add-question", helper.JH(f.CreatePostHandler))
	s.HandleFunc("/edit-question/{q}", helper.JH(f.EditQuestionHandler))
	s.HandleFunc("/answer/{slug}", helper.JH(f.SendAnswersHandler))
	s.HandleFunc("/add-answer/{ans}", helper.JH(f.CreateAnswerHandler))
	s.HandleFunc("/edit-answer/{ans}", helper.JH(f.EditAnswerHandler))
}
