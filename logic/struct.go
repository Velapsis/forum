package logic

import "database/sql"

var website Website
var user User

var regex Regex
var Sql SQL
var db *sql.DB

type Website struct {
	Port     string
	Database string

	Home  string
	Login string
}

type User struct {
	Username string
	Email    string
	Password string
	UUID     int
}

type Regex struct {
	Username string
	Email    string
}

type SQL struct {
	InsertRequest string
}
