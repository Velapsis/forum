package logic

var regex Regex
var user User

type Regex struct {
	Username string
	Email    string
}

// DEBUG ONLY
type User struct {
	Username string
	Email    string
	Password string
	UUID     int
}
