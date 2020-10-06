package database

import (
	"context"
	"errors"
	"time"

	"github.com/Hamaiz/go-rest-eg/helper"
	"github.com/Hamaiz/go-rest-eg/model"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// OauthDatabase - struct that hold database conn
type OauthDatabase struct {
	conn *pgxpool.Pool
}

// NewOauthDatabase - return OauthDatabase struct
func NewOauthDatabase(conn *pgxpool.Pool) *OauthDatabase {
	return &OauthDatabase{conn}
}

// InsertGoogle - adds user to database
// it is a big function that take cares of all google users
func (o *OauthDatabase) InsertGoogle(u model.GoogleData) (string, error) {
	ctx := context.Background()

	// check for email
	c, id := o.CheckEmail(u.Email)

	// check if google email already exists
	g := o.CheckGoogle(u.Email)

	if !c {
		// if the user already exists
		if !g {
			return id, nil
		}

		// make confirmed true
		_, err := o.conn.Exec(ctx, "UPDATE addition SET confirmed=$1 WHERE account_id=$2", true, id)

		if err != nil {
			err = errors.New("error occured, try again")
			return "", err
		}

		// insert google data to database
		_, err = o.conn.Exec(ctx, "INSERT INTO google (google_id, google_token, google_email, google_name, account_id) VALUES ($1, $2, $3, $4, $5)", u.ID, u.Token, u.Email, u.Name, id)

		if err != nil {
			err = errors.New("error occured, try again")
			return "", err
		}

		return id, nil
	}

	// get unique name
	un := helper.UniqueName(u.Name)

	// add user to account database
	_, err := o.conn.Exec(ctx, "INSERT INTO account (id, username, email, password, unique_name) VALUES ($1, $2, $3, $4, $5)", u.ID, u.Name, u.Email, " ", un)

	if err != nil {
		err = errors.New("error occured, try again")
		return "", err
	}

	// add user to google
	_, err = o.conn.Exec(ctx, "INSERT INTO google (google_id, google_token, google_email, google_name, account_id) VALUES ($1, $2, $3, $4, $5)", u.ID, u.Token, u.Email, u.Name, u.ID)

	if err != nil {
		err = errors.New("error occured, try again")
		return "", err
	}

	t := time.Now().Local()
	_, err = o.conn.Exec(ctx, "INSERT INTO addition (expires, token, account_id) VALUES ($1, $2, $3)", t, "", u.ID)

	if err != nil {
		err = errors.New("error occured, try again")
		return "", err
	}

	return u.ID, nil
}

// CheckEmail - checks if email exists
func (o *OauthDatabase) CheckEmail(e string) (bool, string) {
	ctx := context.Background()

	// checking if user exists
	var id string
	err := o.conn.QueryRow(ctx, "SELECT id FROM account WHERE email=$1", e).Scan(&id)

	if err == pgx.ErrNoRows {
		return true, ""
	} else {
		return false, id
	}
}

//CheckGoogle - check if google account exists - if there is no item it returns true
func (o *OauthDatabase) CheckGoogle(e string) bool {
	ctx := context.Background()

	// checking if user exists
	var id string
	err := o.conn.QueryRow(ctx, "SELECT google_id FROM google WHERE google_email=$1", e).Scan(&id)

	if err == pgx.ErrNoRows {
		return true
	} else {
		return false
	}
}
