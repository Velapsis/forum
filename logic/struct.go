package logic

var web Website
var regex Regex

var user User

type Website struct {
	Home  string
	Login string
}

type Regex struct {
	Username string
	Email    string
	Password string
	UUID     int
}

// DEBUG ONLY
type User struct {
	Username string
	Email    string
	Password string
	UUID     int
}
