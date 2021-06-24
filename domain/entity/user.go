package entity
//
//const usernameRegexp = "^[a-zA-Z][a-zA-Z0-9_]{1,41}$"
//const firstNameRegexp = "^[a-zA-Z ]{0,42}$"

// User is, well, a struct depicting a user
type User struct {
	ID       int     `json:"-"`
	Nickname string  `json:"nickname"`
	Fullname *string `json:"fullname,omitempty"`
	Email    *string `json:"email,omitempty"`
	About    *string `json:"about,omitempty"`
}

//// UserOutput is used to marshal JSON with users' data
//type UserOutput struct {
//	UserID     int    `json:"ID"`
//	Username   string `json:"username,omitempty"`
//	Email      string `json:"email,omitempty"`
//	FirstName  string `json:"firstName,omitempty"`
//	LastName   string `json:"lastName,omitempty"`
//	Avatar     string `json:"avatarLink,omitempty"`
//	Following  int    `json:"following"`
//	FollowedBy int    `json:"followers"`
//	Followed   *bool  `json:"followed,omitempty"` // pointer because we need to not send this sometimes
//}
//
//// UserRegInput is used when parsing JSON in user/signup handler
//type UserRegInput struct {
//	Username  string `json:"username" valid:"username"`
//	Password  string `json:"password" valid:"stringlength(8|30)"`
//	Email     string `json:"email" valid:"email"`
//	FirstName string `json:"firstName" valid:"name,optional"`
//	LastName  string `json:"lastName" valid:"name,optional"`
//}
//
//// UserLoginInput is used when parsing JSON in user/login handler
//type UserLoginInput struct {
//	Username string `json:"username"`
//	Password string `json:"password"`
//}
//
//// UserPassChangeInput is used when parsing JSON in profile/password handler
//type UserPassChangeInput struct {
//	Password string `json:"password" valid:"stringlength(8|30)"`
//}
//
//// UserEditInput is used when parsing JSON in profile/edit handler
//type UserEditInput struct {
//	Username  string `json:"username" valid:"username,optional"`
//	Email     string `json:"email" valid:"email,optional"`
//	FirstName string `json:"firstName" valid:"name,optional"`
//	LastName  string `json:"lastName" valid:"name,optional"`
//	Avatar    string `json:"avatarLink" valid:"filepath,optional"`
//}
//

