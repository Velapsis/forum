package logic

import "database/sql"

var website Website
var webpage WebPage

var regex Regex
var Sql SQL
var db *sql.DB

type Website struct {
	Port     string
	Database string

	Home  string
	Login string
}

type WebPage struct {
	IsConnected bool
}

type Regex struct {
	Username string
	Email    string
}

type SQL struct {
	InsertRequest string
}
