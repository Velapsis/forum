package logic

var regex Regex
var user User
var sql SQL

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

type SQL struct {
	InsertRequest string
}
