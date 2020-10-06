package database

import (
	"context"
	"errors"
	"time"

	"github.com/Hamaiz/go-rest-eg/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// AccountDatabase - struct that holds functions for putting users in database
type AccountDatabase struct {
	conn *pgxpool.Pool
}

// NewAccountDatabase - returns AccountDatabase
func NewAccountDatabase(conn *pgxpool.Pool) *AccountDatabase {
	return &AccountDatabase{conn}
}

// CheckingExists - checks if email provided already exists
func (a *AccountDatabase) CheckingExists(e string) bool {
	var id string
	err := a.conn.QueryRow(context.Background(), "SELECT id FROM account WHERE email=$1", e).Scan(&id)

	if err != pgx.ErrNoRows {
		return true
	} else {
		return false
	}
}

// InsertUser - inserts user into database and return error
func (a *AccountDatabase) InsertUser(u model.User, t string) error {
	id := uuid.New()
	expires := time.Now().Local().Add(time.Hour * time.Duration(6))

	_, err := a.conn.Exec(context.Background(), "INSERT INTO account (id, username, email, password, unique_name) VALUES ($1, $2, $3, $4, $5)", id, u.Name, u.Email, u.Password, u.UniqueName)

	_, err = a.conn.Exec(context.Background(), "INSERT INTO addition (expires, token, account_id) VALUES ($1, $2, $3)", expires, t, id)

	return err
}

// GetUser - gets user from databasae
func (a *AccountDatabase) GetUser(id string) (model.UserSend, error) {
	u := model.UserSend{}
	row := a.conn.QueryRow(context.Background(), "SELECT username, email, unique_name FROM account WHERE id=$1", id)
	err := row.Scan(&u.Name, &u.Email, &u.UnqiueName)
	return u, err
}

// GetUserInLogin - gets user when logging in with email
func (a *AccountDatabase) GetUserInLogin(email string) (model.UserGet, error) {
	u := model.UserGet{}
	row := a.conn.QueryRow(context.Background(), "SELECT id, password FROM account WHERE email=$1", email)
	err := row.Scan(&u.ID, &u.Password)
	return u, err
}

// LoginConfirm - checks if login confirms
func (a *AccountDatabase) LoginConfirm(e string) bool {
	var id string

	// with email get id
	row := a.conn.QueryRow(context.Background(), "SELECT id FROM account WHERE email=$1", e)
	err := row.Scan(&id)

	// if email not found return false
	if err != nil {
		return false
	}

	// with id get additional data from database
	var confirm bool
	token := a.conn.QueryRow(context.Background(), "SELECT confirmed FROM addition WHERE account_id=$1", id)
	err = token.Scan(&confirm)

	if err != nil {
		return false
	}

	return confirm
}

// ConfirmEmail - confirms the token from email
func (a *AccountDatabase) ConfirmEmail(t string) (bool, error) {
	et := model.EmailToken{}

	// check in the database if token exists
	row := a.conn.QueryRow(context.Background(), "SELECT * FROM addition WHERE token=$1", t)
	err := row.Scan(&et.Confirmed, &et.Expires, &et.Token, &et.Account_id)

	// if no rows return false and with error
	if err == pgx.ErrNoRows {
		err = errors.New("no user found with token")
		return false, err
	}

	if et.Confirmed {
		err = errors.New("email already confirmed")
		return true, err
	}

	// if no error - check the expiry time (which is 6hrs)
	e, _ := time.Parse("2006-01-02 15:04:05", et.Expires)
	c := time.Now().Local().After(e)

	if !c {
		err = errors.New("token expired")
		return false, err
	}

	// if token not expired, change token to nil and confirmed to false
	_, err = a.conn.Exec(context.Background(), "UPDATE addition SET token=$1, confirmed=$2 WHERE account_id=$3", "", true, et.Account_id)

	if err != nil {
		err = errors.New("problem occured")
		return false, err
	}

	return true, nil
}

// GetToken - gets the token again for email
func (a *AccountDatabase) GetToken(e string) (string, error) {
	var id string
	row := a.conn.QueryRow(context.Background(), "SELECT id FROM account WHERE email=$1", e)
	err := row.Scan(&id)

	if err != nil {
		err = errors.New("not found")
		return "", err
	}

	var token string
	var expires time.Time
	var confirmed bool
	row = a.conn.QueryRow(context.Background(), "SELECT token, expires, confirmed FROM addition WHERE account_id=$1", id)
	err = row.Scan(&token, &expires, &confirmed)

	if err != nil {
		err = errors.New("token not found")
		return "", err
	}

	if confirmed {
		err = errors.New("email already confirmed")
		return "", err
	}

	exp, _ := time.Parse("2006-01-02 15:04:05", expires.String())
	c := time.Now().Local().After(exp)

	if !c {
		err = errors.New("token expired")
		return "", err
	}

	return token, nil
}

// == pass reset ==//

// ForgotToken - add forgot token
func (a *AccountDatabase) ForgotToken(e string, token string) error {
	// get id from email
	var id string
	row := a.conn.QueryRow(context.Background(), "SELECT id FROM account WHERE email=$1", e)
	err := row.Scan(&id)

	if err != nil {
		return err
	}

	t := time.Now().UTC().Add(time.Hour * time.Duration(6))
	_, err = a.conn.Exec(context.Background(), "UPDATE addition SET token=$1, expires=$2 WHERE account_id=$3", token, t, id)

	if err != nil {
		return err
	}

	return nil
}

// ConfirmToken - confirm passwrod token
func (a *AccountDatabase) ConfirmToken(token string) (bool, error) {
	var expires time.Time
	row := a.conn.QueryRow(context.Background(), "SELECT expires FROM addition WHERE token=$1", token)
	err := row.Scan(&expires)

	switch {
	case err == pgx.ErrNoRows:
		err = errors.New("token not found")
		return false, err
	case err != nil:
		err = errors.New("an error occured")
		return false, err
	}

	// check if token expired
	e, _ := time.Parse("2006-01-02 15:04:05", expires.String())
	c := time.Now().UTC().After(e)

	if !c {
		err = errors.New("token expired")
		return false, err
	}

	return true, nil
}

// ResetPass - reset password
func (a *AccountDatabase) ResetPass(p string, t string) error {
	var id string
	row := a.conn.QueryRow(context.Background(), "SELECT account_id FROM addition WHERE token=$1", t)
	err := row.Scan(&id)

	if err != nil {
		err = errors.New("an error occured")
		return err
	}

	_, err = a.conn.Exec(context.Background(), "UPDATE account SET password=$1 WHERE id=$1", p, id)
	if err != nil {
		err = errors.New("error occured changing password")
		return err
	}

	_, err = a.conn.Exec(context.Background(), "UPDATE addition SET token=$1, expires=$2 WHERE account_id=$3", "", nil, id)
	if err != nil {
		err = errors.New("an error occured")
		return err
	}

	return nil
}
