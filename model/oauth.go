package model

// GoogleData - holds data of google
type GoogleData struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	Verified_Email bool   `json:"verified_email"`
	Name           string `json:"name"`
	Given_Name     string `json:"given_name"`
	Family_Name    string `json:"family_name"`
	Picture        string `json:"picture"`
	Locale         string `json:"locale"`
	Token          string `json:"token"`
}

// InsertGoolge - struct that holds adding user
type InsertGoogle struct {
	Google_Id    string `json:"google_id"`
	Google_Token string `json:"google_token"`
	Google_Email string `json:"google_email"`
	Google_Name  string `json:"google_name"`
	Account_ID   string `json:"account_id"`
}
