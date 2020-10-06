package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Hamaiz/go-rest-eg/helper"
	"github.com/Hamaiz/go-rest-eg/model"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

// FilesDatabase - holds all the function - interface
type FilesDatabase interface {
	GetSearchedQuestions(l string) ([]model.GetQuestions, error)
	GetQuestions() ([]model.GetQuestions, error)
	PostQuestion(p model.FilesQuestion) error
	GetQuest(s string) (model.FilesSend, error)
	GetQuestion(s string) (model.FilesQuestion, error)
	EditQuestion(s string, nq string, slug string) error
	AddAnswer(a model.FilesComment) error
	GetAnswer(s string, c string) (model.FilesComment, error)
	GetOneAnswer(s string) (string, error)
	EditAnswer(s string, na string) error
	GetAnswers(s string) ([]model.GetAnswers, error)
	Like(id string, u string) error
	Dislike(id string, u string) error
	GetLikes(id string) (int, error)
}

// Account - account store struct
type Media struct {
	store AccountStore
	conn  FilesDatabase
}

// NewAccountStore - creates new store
func NewFilesApi(s AccountStore, c FilesDatabase) *Media {
	return &Media{s, c}
}

// SearchQuestionHandler - search questions - @POST - /api/search
func (m *Media) SearchQuestionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		helper.ASM(w, 405, "")
		return
	}

	// FormValue
	s := r.FormValue("search")

	// extract all the important words
	l := helper.KeyExtract(s)

	// get question
	fqs, err := m.conn.GetSearchedQuestions(l)
	if err != nil {
		helper.ASM(w, 403, err.Error())
		return
	} else {
		json.NewEncoder(w).Encode(fqs)
		return
	}
}

// GetQuestionsHandler - get all questions - @GET | @OPTIONS - /api/question
func (m *Media) GetQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// check for header
		fgq := r.Header.Get("files-get-questions")
		if fgq == "" {
			helper.ASM(w, 401, "")
			return
		}

		fqs, err := m.conn.GetQuestions()
		if err != nil {
			helper.ASM(w, 403, err.Error())
			return
		}

		json.NewEncoder(w).Encode(fqs)
		return
	case "OPTIONS":
		helper.ASM(w, 204, "")
		return

	default:
		helper.ASM(w, 405, "")
		return
	}
}

// SendQuestionHandler - send all question - @GET | @OPTIONS- /api/question/:slug
func (m *Media) SendQuestionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// check for header
		fgq := r.Header.Get("files-get-question")
		if fgq == "" {
			helper.ASM(w, 401, "")
			return
		}

		// get param from request
		param := mux.Vars(r)
		slug := param["slug"]

		// get question with slug
		q, err := m.conn.GetQuest(slug)
		if err != nil {
			helper.ASM(w, 404, err.Error())
			return
		}

		json.NewEncoder(w).Encode(q)

	case "OPTIONS":
		helper.ASM(w, 204, "")
		return

	default:
		helper.ASM(w, 405, "")
		return
	}

}

// CreatePostHandler - create posts - @POST - /api/add-question
func (m *Media) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		helper.ASM(w, 405, "")
		return
	}

	if !m.store.AlreadyLoggedIn(r) {
		helper.ASM(w, 401, "")
		return
	}

	// get user id
	id, err := m.store.GetUser(r)
	if err != nil {
		helper.ASM(w, 403, err.Error())
		return
	}

	// Form Value
	q := r.FormValue("question")
	t := time.Now().UTC().Format(time.RFC3339)
	qs := helper.UniqueQuestion(q)
	qi := uuid.New().String()

	// model hold items
	fq := model.FilesQuestion{qi, q, id, qs, t, t}

	// add item to database
	err = m.conn.PostQuestion(fq)
	if err != nil {
		helper.ASM(w, 500, "")
		return
	}

	helper.ASM(w, 201, "post made")
}

// EditQuestionHandler - edits the already edited - @PUT - /api/edit-question/:q
func (m *Media) EditQuestionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		helper.ASM(w, 405, "")
		return
	}

	if !m.store.AlreadyLoggedIn(r) {
		helper.ASM(w, 401, "")
		return
	}

	// form value
	nq := r.FormValue("question")
	qs := helper.UniqueQuestion(nq)

	if nq == "" {
		helper.ASM(w, 403, "no question found")
		return
	}

	// mux vars
	param := mux.Vars(r)
	q := param["q"]

	if q == "" {
		helper.ASM(w, 403, "no slug")
		return
	}

	// get question from database
	fq, err := m.conn.GetQuestion(q)
	if err != nil {
		helper.ASM(w, 404, err.Error())
		return
	}

	// get user id from the session cookie
	var id string
	id, err = m.store.GetUser(r)
	if err != nil {
		helper.ASM(w, 500, "")
		return
	}

	// if id and question poster return unauthorized
	if id != fq.Poster {
		helper.ASM(w, 401, "")
		return
	}

	// edit question database
	err = m.conn.EditQuestion(q, nq, qs)
	if err != nil {
		helper.ASM(w, 500, "")
		return
	}

	helper.ASM(w, 201, "question edited")
}

// SendAsnwersHandler - send all the answer to desired question
// @GET | @OPTIONS - /api/answers/:slug
func (m *Media) SendAnswersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// check for header
		fga := r.Header.Get("files-get-answers")
		if fga == "" {
			helper.ASM(w, 401, "")
			return
		}

		// get param from request
		param := mux.Vars(r)
		slug := param["slug"]

		// get answers with all slug
		q, err := m.conn.GetAnswers(slug)
		if err != nil {
			helper.ASM(w, 404, err.Error())
			return
		}

		json.NewEncoder(w).Encode(q)
		return
	case "OPTIONS":
		helper.ASM(w, 204, "")
		return
	default:
		helper.ASM(w, 405, "")
		return
	}

}

// SendAsnwerHandler - send all the answer to desired question
// @GET | @OPTIONS - /api/answer/:slug
func (m *Media) SendAnswerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// check for header
		fga := r.Header.Get("files-get-answer")
		if fga == "" {
			helper.ASM(w, 401, "")
			return
		}

		// get param from request
		param := mux.Vars(r)
		slug := param["slug"]

		// get answers with all slug
		q, err := m.conn.GetOneAnswer(slug)
		if err != nil {
			helper.ASM(w, 404, err.Error())
			return
		}

		json.NewEncoder(w).Encode(q)
		return
	case "OPTIONS":
		helper.ASM(w, 204, "")
		return
	default:
		helper.ASM(w, 405, "")
		return
	}

}

// CreateAnswerHandler - create answer for desired post - @POST - /api/add-answer/:ans
func (m *Media) CreateAnswerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		helper.ASM(w, 405, "")
		return
	}

	if !m.store.AlreadyLoggedIn(r) {
		helper.ASM(w, 403, "")
		return
	}

	// form values
	a := r.FormValue("answer")
	t := time.Now().UTC().Format(time.RFC3339)

	if a == "" {
		helper.ASM(w, 403, "answer is empty")
		return
	}

	// get param ans
	param := mux.Vars(r)
	ans := param["ans"]

	// if no param
	if ans == "" {
		helper.ASM(w, 404, "")
		return
	}

	// get user id
	id, err := m.store.GetUser(r)
	if err != nil {
		helper.ASM(w, 403, err.Error())
		return
	}

	// FileComment
	c := model.FilesComment{ans, a, id, t, t}

	// if already answered
	_, err = m.conn.GetAnswer(ans, id)
	switch {
	case err == pgx.ErrNoRows:
		err = m.conn.AddAnswer(c)
		if err != nil {
			helper.ASM(w, 500, "")
			return
		}
		helper.ASM(w, 201, "answer made")
		return
	case err != nil:
		helper.ASM(w, 404, "try again")
		return
	default:
		helper.ASM(w, 403, "you have already answered the question")
		return
	}
}

// EditAnswerHandler - edit answer - @PUT -  /api/edit-answer/:ans
func (m *Media) EditAnswerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		helper.ASM(w, 405, "")
		return
	}

	if !m.store.AlreadyLoggedIn(r) {
		helper.ASM(w, 403, "")
		return
	}

	// form value
	a := r.FormValue("answer")

	if a == "" {
		helper.ASM(w, 403, "answer is empty")
		return
	}

	// get ans param
	param := mux.Vars(r)
	ans := param["ans"]

	if ans == "" {
		helper.ASM(w, 404, "")
		return
	}

	// get user id from the session cookie
	id, err := m.store.GetUser(r)
	if err != nil {
		helper.ASM(w, 500, "")
		return
	}

	// get question from database
	var fc model.FilesComment
	fc, err = m.conn.GetAnswer(ans, id)
	switch {
	case err == pgx.ErrNoRows:
		helper.ASM(w, 404, "no answer found")
		return
	case err != nil:
		helper.ASM(w, 404, "try again")
		return
	}

	// if id and question poster return unauthorized
	if id != fc.Commenter {
		helper.ASM(w, 401, "")
		return
	}

	// edit question database
	err = m.conn.EditAnswer(ans, a)
	if err != nil {
		helper.ASM(w, 500, "")
		return
	}

	helper.ASM(w, 201, "post edited")
}
