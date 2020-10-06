package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

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
	ForgotToken(e string, token string) error
	ConfirmToken(token string) (bool, error)
	ResetPass(p string, t string) error
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

// SignUpHandler - signing in route - @POST - /account/signup
func (s *Account) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		helper.ASM(w, 405, "")
		return
	}

	// if already logged in
	if s.store.AlreadyLoggedIn(r) {
		helper.ASM(w, 401, "")
		return
	}

	// getting form value
	n := r.FormValue("name")
	e := r.FormValue("email")
	p := r.FormValue("password")

	// if empty return error
	if n == "" || e == "" || p == "" {
		helper.ASM(w, 403, "missing credentials")
		return
	}

	// checking if already exists
	check := s.conn.CheckingExists(e)
	if check {
		helper.ASM(w, 403, "email already exists")
		return
	}

	// bcrypt password - hashing
	bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		helper.ASM(w, 500, "")
		return
	}

	// unique name
	unique := helper.UniqueName(n)

	// email token
	token := uniuri.NewLen(50)

	// email url
	host := os.Getenv("URL")
	url := host + "account/confirm/" + token

	// user model
	u := model.User{n, e, unique, string(bs)}

	// inser user to database
	err = s.conn.InsertUser(u, token)

	if err != nil {
		helper.ASM(w, 500, "")
		return
	}

	// sending email
	err = email.SignUpEmail(e, n, url)
	if err != nil {
		helper.ASM(w, 500, "")
		return
	}

	helper.ASM(w, 200, "verify email to continue")
}

// LogInHandler - logging in route - @POST - /account/login
func (s *Account) LogInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		helper.ASM(w, 405, "")
		return
	}

	// If already logged in send unauthorizedError
	if s.store.AlreadyLoggedIn(r) {
		helper.ASM(w, 401, "")
		return
	}

	// Get From Values
	e := r.FormValue("email")
	p := r.FormValue("password")

	// Check if empty
	if e == "" || p == "" {
		helper.ASM(w, 403, "missing credentials")
		return
	}

	u, err := s.conn.GetUserInLogin(e)
	switch {
	case err == pgx.ErrNoRows:
		helper.ASM(w, 404, "email not found")
		return
	case err != nil:
		helper.ASM(w, 500, "")
		return
	}

	confirm := s.conn.LoginConfirm(e)
	if !confirm {
		helper.ASM(w, 403, "confirm your email to continue")
		return
	}

	// Bcrypt Pasword - Decrypting
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	if err != nil {
		helper.ASM(w, 404, "Email and/or password do not match")
		return
	}

	// Create New Session
	err = s.store.SaveSession(w, r, u.ID)
	if err != nil {
		helper.ASM(w, 500, "")
		return
	}

	helper.ASM(w, 200, "logged in successfully")
}

// LogoutHandler - logs user out & removes cookie - @DELETE - /account/logout
func (s *Account) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		helper.ASM(w, 405, "")
		return
	}

	// if not logged in
	if !s.store.AlreadyLoggedIn(r) {
		helper.ASM(w, 401, "")
		return
	}

	// clean session
	err := s.store.CleanSession(w, r)
	if err != nil {
		helper.ASM(w, 404, "")
		return
	}

	helper.ASM(w, 200, "logged out")
}

// GetUserHandler - sends user - @GET | @OPTIONS - /account/getUser
func (s *Account) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// get and check for header
		fgu := r.Header.Get("files-get-user")
		if fgu == "" {
			helper.ASM(w, 401, "")
			return
		}

		if !s.store.AlreadyLoggedIn(r) {
			helper.ASM(w, 401, "you are not logged in")
			return
		}

		// get user id
		id, err := s.store.GetUser(r)
		if err != nil {
			helper.ASM(w, 404, "")
			return
		}

		// get user from database
		var u model.UserSend
		u, err = s.conn.GetUser(id)
		if err != nil {
			log.Println(err)
			helper.ASM(w, 404, "")
			return
		}

		json.NewEncoder(w).Encode(u)
		return
	case "OPTIONS":
		helper.ASM(w, 204, "")
		return
	default:
		helper.ASM(w, 405, "")
		return
	}

}

// ConfirmEmailHandler - confirms the email - @GET | @OPTIONS - /account/confirm/:token
func (s *Account) ConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		// check if already logged in
		if s.store.AlreadyLoggedIn(r) {
			helper.ASM(w, 401, "")
			return
		}

		// get params & check for token
		param := mux.Vars(r)
		token := param["token"]

		if token == "" {
			helper.ASM(w, 404, "")
			return
		}

		confirm, err := s.conn.ConfirmEmail(token)

		if err != nil {
			helper.ASM(w, 403, err.Error())
			return
		}

		if confirm {
			helper.ASM(w, 200, "email is confirmed")
			return
		} else {
			helper.ASM(w, 403, "token expired")
			return
		}
		return
	case "OPTIONS":
		helper.ASM(w, 204, "")
		return
	default:
		helper.ASM(w, 405, "")
		return
	}

}

// ForgotHandler - forget handler - @POST - /account/forgot
func (s *Account) ForgotHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		helper.ASM(w, 405, "")
		return
	}

	if s.store.AlreadyLoggedIn(r) {
		helper.ASM(w, 401, "")
		return
	}

	// form value
	e := r.FormValue("email")
	if e == "" {
		helper.ASM(w, 403, "email is empty")
		return
	}

	// check if user exist
	c := s.conn.CheckingExists(e)
	if !c {
		helper.ASM(w, 404, "user not found with the email")
		return
	}

	// email token
	token := uniuri.NewLen(50)

	// email url
	host := os.Getenv("URL")
	url := host + "account/confirm-pass/" + token

	// inser user to database
	err := s.conn.ForgotToken(e, token)
	if err != nil {
		helper.ASM(w, 422, "an error occured")
		return
	}

	// sending email
	err = email.ForgotEmail(e, url)
	if err != nil {
		helper.ASM(w, 500, "an error occured while sending email")
		return
	}

	helper.ASM(w, 200, "email has been sent to your account")
}

// ConfirmPassHandler - confirm pasword token - @GET | @OPTIONS - /account/confirm-pass/:token
func (s *Account) ConfirmPassHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		param := mux.Vars(r)
		t := param["token"]

		if t == "" {
			helper.ASM(w, 403, "token is empty")
			return
		}

		c, err := s.conn.ConfirmToken(t)
		if err != nil {
			helper.ASM(w, 422, err.Error())
			return
		}

		if !c {
			helper.ASM(w, 403, "token expired")
			return
		}

		helper.ASM(w, 200, "token verified, type your new password")

		return
	case "OPTIONS":
		helper.ASM(w, 204, "")
		return
	default:
		helper.ASM(w, 405, "")
		return
	}
}

// ResetHandler - resets the password - @PUT | @POST - /account/reset
func (s *Account) ResetHandler(w http.ResponseWriter, r *http.Request) {
	if s.store.AlreadyLoggedIn(r) {
		helper.ASM(w, 401, "")
		return
	}

	switch r.Method {
	case "POST":
		// Form value
		t := r.FormValue("token")

		if t == "" {
			helper.ASM(w, 403, "token is empty")
			return
		}

		c, err := s.conn.ConfirmToken(t)
		if err != nil {
			helper.ASM(w, 422, err.Error())
			return
		}

		if !c {
			helper.ASM(w, 403, "token expired")
			return
		}

		helper.ASM(w, 200, "token verified, type your new password")

	case "PUT":
		p := r.FormValue("pass")
		cp := r.FormValue("confirmPass")
		t := r.FormValue("token")

		if cp == "" || p == "" || t == "" {
			helper.ASM(w, 204, "missing credentials")
			return
		}

		if cp != p {
			helper.ASM(w, 403, "password not match")
			return
		}

		c, err := s.conn.ConfirmToken(t)
		if err != nil {
			helper.ASM(w, 422, err.Error())
			return
		}

		if !c {
			helper.ASM(w, 403, "token expired")
			return
		}

		bs, errs := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
		if errs != nil {
			helper.ASM(w, 500, "")
			return
		}

		err = s.conn.ResetPass(string(bs), t)
		if err != nil {
			helper.ASM(w, 403, err.Error())
			return
		}

		helper.ASM(w, 200, "your password changed")
	default:
		helper.ASM(w, 405, "")
	}
}

// EmailAgain - send email again - @POST - /account/againemail
func (s *Account) EmailAgain(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		helper.ASM(w, 405, "")
		return
	}

	if s.store.AlreadyLoggedIn(r) {
		helper.ASM(w, 401, "")
		return
	}

	// Get From Values
	e := r.FormValue("email")
	n := r.FormValue("name")

	// get token from database
	token, err := s.conn.GetToken(e)
	if err != nil {
		helper.ASM(w, 404, err.Error())
		return
	}

	// email url
	url := "http://localhost:9002/account/confirm/" + token

	err = email.SignUpEmail(e, n, url)
	if err != nil {
		helper.ASM(w, 500, "")
		return
	}

	helper.ASM(w, 200, "email sent again")
}
