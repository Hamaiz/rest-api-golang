package model

// User - login user
type User struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	UniqueName string `json:"uniquename"`
	Password   string `json:"password:`
}

// UserGet retruns for login handler
type UserGet struct {
	ID       string `json:"id"`
	Password string `json:"password:`
}

// UserSend - the struct sending to user
type UserSend struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	UnqiueName string `json:"uniquename"`
}

// EmailToken - email token struct
type EmailToken struct {
	Confirmed  bool
	Expires    string
	Token      string
	Account_id string
}
