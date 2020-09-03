package session

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

// GetUser - gets session_id and returns the user
func (s AccountStore) GetUser(r *http.Request) (string, error) {
	session, _ := s.store.Get(r, "presence")

	if session.Values["session_id"] == nil {
		err := errors.New("no session id found")
		return "", err
	} else {
		id := session.Values["session_id"].(string)
		return id, nil
	}
}

// AlreadyLoggedIn - tells if the user is logged in or not
func (s AccountStore) AlreadyLoggedIn(r *http.Request) bool {
	session, _ := s.store.Get(r, "presence")

	if session.Values["session_id"] == nil {
		return false
	} else {
		return true
	}
}

// CleanSession - cleans session
func (s AccountStore) CleanSession(w http.ResponseWriter, r *http.Request) error {
	session, _ := s.store.Get(r, "presence")

	session.Options.MaxAge = -1

	err := session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

// SaveSession - saves session
func (s AccountStore) SaveSession(w http.ResponseWriter, r *http.Request, id string) error {
	// Create New Session
	session, _ := s.store.Get(r, "presence")

	// adding session to value
	session.Values["session_id"] = id

	// adding session options
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	err := session.Save(r, w)

	return err
}
