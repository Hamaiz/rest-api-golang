package api

import (
	"encoding/json"
	"net/http"

	"github.com/Hamaiz/go-rest-eg/email"
	"github.com/Hamaiz/go-rest-eg/helper"
	"github.com/Hamaiz/go-rest-eg/model"
	"github.com/dchest/uniuri"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

// AccountStore - account store interface
type AccountStore interface {
	GetUser(r *http.Request) (string, error)
	AlreadyLoggedIn(r *http.Request) bool
	CleanSession(w http.ResponseWriter, r *http.Request) error
	SaveSession(w http.ResponseWriter, r *http.Request, id string) error
}

// AccountDatabase - hold database functions
type AccountDatabase interface {
	CheckingExists(e string) bool
	InsertUser(u model.User, t string) error
	GetUser(id string) (model.UserSend, error)
	GetUserInLogin(email string) (model.UserGet, error)
	LoginConfirm(e string) bool
	ConfirmEmail(t string) (bool, error)
	GetToken(e string) (string, error)
}

// Account - account store struct
type Account struct {
	store AccountStore
	conn  AccountDatabase
}

// NewAccountStore - creates new store
func NewAccountStore(s AccountStore, c AccountDatabase) *Account {
	return &Account{s, c}
}

// SignUpHandler - signing in route - /account/signup
func (s *Account) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	helper.MethodCheck(w, r, "POST")

	// if already logged in
	if s.store.AlreadyLoggedIn(r) {
		helper.UnauthorizedError(w)
		return
	}

	// getting form value
	n := r.FormValue("name")
	e := r.FormValue("email")
	p := r.FormValue("password")

	// if empty return error
	if n == "" || e == "" || p == "" {
		helper.ForbiddenError(w, "missing credentials")
		return
	}

	// checking if already exists
	check := s.conn.CheckingExists(e)
	if check {
		helper.ForbiddenError(w, "Email already exists")
		return
	}

	// bcrypt password - hashing
	bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		helper.ServerError(w)
		return
	}

	// unique name
	unique := helper.UniqueName(n)

	// email token
	token := uniuri.NewLen(50)

	// email url
	url := "http://localhost:9002/account/confirm/" + token

	// user model
	u := model.User{n, e, unique, string(bs)}

	// inser user to database
	err = s.conn.InsertUser(u, token)

	if err != nil {
		helper.ServerError(w)
		return
	}

	// sending email
	err = email.SignUpEmail(e, n, url)
	if err != nil {
		helper.ServerError(w)
		return
	}

	helper.Ok(w, "verify email to continue")
}

// LogInHandler - logging in route - /account/login
func (s *Account) LogInHandler(w http.ResponseWriter, r *http.Request) {
	helper.MethodCheck(w, r, "POST")

	// If already logged in send unauthorizedError
	if s.store.AlreadyLoggedIn(r) {
		helper.UnauthorizedError(w)
		return
	}

	// Get From Values
	e := r.FormValue("email")
	p := r.FormValue("password")

	// Check if empty
	if e == "" || p == "" {
		helper.ForbiddenError(w, "Missing Credentials")
		return
	}

	confirm := s.conn.LoginConfirm(e)
	if !confirm {
		helper.ForbiddenError(w, "confirm your email to continue")
		return
	}

	u, err := s.conn.GetUserInLogin(e)
	switch {
	case err == pgx.ErrNoRows:
		helper.NotFound(w, r)
		return
	case err != nil:
		helper.ServerError(w)
		return
	}

	// Bcrypt Pasword - Decrypting
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	if err != nil {
		helper.ForbiddenError(w, "Email and/or password do not match")
		return
	}

	// Create New Session
	err = s.store.SaveSession(w, r, u.ID)
	if err != nil {
		helper.ServerError(w)
		return
	}

	helper.Ok(w, "logged in successfully")
}

// LogoutHandler - logs user out & removes cookie - /account/logout
func (s *Account) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	helper.MethodCheck(w, r, "DELETE")

	// if not logged in
	if !s.store.AlreadyLoggedIn(r) {
		helper.UnauthorizedError(w)
		return
	}

	// clean session
	err := s.store.CleanSession(w, r)
	if err != nil {
		helper.NotFound(w, r)
		return
	}

	helper.Ok(w, "logged out")
}

// GetUserHandler - sends user - /account/getUser
func (s *Account) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	helper.MethodCheck(w, r, "GET")

	// get and check for header
	fgu := r.Header.Get("files-get-user")
	if fgu == "" {
		helper.UnauthorizedError(w)
		return
	}

	// get user id
	id, err := s.store.GetUser(r)
	if err != nil {
		helper.NotFound(w, r)
		return
	}

	// get user from database
	var u model.UserSend
	u, err = s.conn.GetUser(id)
	if err != nil {
		helper.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(u)
}

// ConfirmEmailHandler - confirms the email - /account/confirm/:token
func (s *Account) ConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	helper.MethodCheck(w, r, "GET")

	// check if already logged in
	if s.store.AlreadyLoggedIn(r) {
		helper.UnauthorizedError(w)
		return
	}

	// get params & check for token
	param := mux.Vars(r)
	token := param["token"]

	if token == "" {
		helper.NotFound(w, r)
		return
	}

	confirm, err := s.conn.ConfirmEmail(token)

	if err != nil {
		helper.ForbiddenError(w, err.Error())
		return
	}

	if confirm {
		helper.Ok(w, "email is confirmed")
		return
	} else {
		helper.ForbiddenError(w, "token expired")
		return
	}
}

// EmailAgain - send email again - /account/againemail
func (s *Account) EmailAgain(w http.ResponseWriter, r *http.Request) {
	helper.MethodCheck(w, r, "POST")

	if s.store.AlreadyLoggedIn(r) {
		helper.UnauthorizedError(w)
		return
	}

	// Get From Values
	e := r.FormValue("email")
	n := r.FormValue("name")

	// get token from database
	token, err := s.conn.GetToken(e)
	if err != nil {
		helper.NotFound(w, r)
		return
	}

	// email url
	url := "http://localhost:9002/account/confirm/" + token

	err = email.SignUpEmail(e, n, url)
	if err != nil {
		helper.ServerError(w)
		return
	}

	helper.Ok(w, "email sent again")
}
