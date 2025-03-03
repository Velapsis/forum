package logic

var web Website
var regex Regex

type Website struct {
	Home  string
	Login string
}

type Regex struct {
	Username string
	Email    string
	Password string
}
