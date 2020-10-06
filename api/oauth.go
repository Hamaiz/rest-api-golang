package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Hamaiz/go-rest-eg/helper"
	"github.com/Hamaiz/go-rest-eg/model"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// OauthDatabase - hold all the functions of oauth database
type OauthDatabase interface {
	InsertGoogle(u model.GoogleData) (string, error)
	CheckEmail(e string) (bool, string)
	CheckGoogle(e string) bool
}

// Oauth - struct holds all the functions
type Oauth struct {
	store AccountStore
	conn  OauthDatabase
}

// googleOauthConfig - return google config
func googleOauthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email"},
		RedirectURL:  os.Getenv("URL") + "account/google/callback",
		Endpoint:     google.Endpoint,
	}

}

// NewOauthApi - creates new oauthapi
func NewOauthApi(s AccountStore, c OauthDatabase) *Oauth {
	return &Oauth{s, c}
}

// GoogleLoginHandler - google login handler - /account/google/login
func (o *Oauth) GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Method check
	if r.Method != "GET" {
		helper.ASM(w, 405, "")
		return
	}

	// check if logged in
	if o.store.AlreadyLoggedIn(r) {
		helper.ASM(w, 401, "")
		return
	}

	// make random value for CSRF attacks
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	// make cookie
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}

	// set it as a cookie
	http.SetCookie(w, &cookie)

	//conf := googleOauthConfig.AuthCodeURL(state)
	u := googleOauthConfig().AuthCodeURL(state)

	// redirect to url
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

// GoogleCallbackHandler - google callback - /account/google/callback
func (o *Oauth) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		helper.ASM(w, 405, "")
	}

	// check if logged in
	if o.store.AlreadyLoggedIn(r) {
		helper.ASM(w, 401, "")
		return
	}

	// get oauthstate cookie
	oauthstate, _ := r.Cookie("oauthstate")

	// check if it matches
	if r.FormValue("state") != oauthstate.Value {
		helper.ASM(w, 403, "")
		return
	}

	// exchange for the access token
	token, err := googleOauthConfig().Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		helper.ASM(w, 403, "")
		return
	}

	// get response data with token
	oauthGoogleURLAPI := "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	response, err := http.Get(oauthGoogleURLAPI + token.AccessToken)
	if err != nil {
		helper.ASM(w, 403, "")
		return
	}

	defer response.Body.Close()

	// read all the response body
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		helper.ASM(w, 403, "")
		return
	}

	// convert byte to a struct
	var data = model.GoogleData{}
	data.Token = token.AccessToken

	err = json.Unmarshal(contents, &data)
	if err != nil {
		helper.ASM(w, 403, "")
		return
	}

	// remove oauthstate cookie
	c := &http.Cookie{
		Name:   "oauthstate",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	// add to the database
	var id string
	id, err = o.conn.InsertGoogle(data)
	if err != nil {
		helper.ASM(w, 403, err.Error())
		return
	}

	// create new session
	err = o.store.SaveSession(w, r, id)
	if err != nil {
		helper.ASM(w, 500, "")
		return
	}

	helper.ASM(w, 200, "logged into google")
}
