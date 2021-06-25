package entity

//
//const usernameRegexp = "^[a-zA-Z][a-zA-Z0-9_]{1,41}$"
//const firstNameRegexp = "^[a-zA-Z ]{0,42}$"

// User is, well, a struct depicting a user
type User struct {
	ID       int    `json:"-"`
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname,omitempty"`
	Email    string `json:"email,omitempty"`
	About    string `json:"about,omitempty"`
}
