package database

import (
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Query struct {
	InsertUser string
}

var query Query

func DefineRequests() {
	query.InsertUser = `INSERT INTO user (username, email, password) VALUES (?, ?, ?)`
	// Sql.UpdateUsernameRequest = `UPDATE user SET username = ? WHERE id = ?`
	// Sql.UpdateEmailRequest = `UPDATE user SET email = ? WHERE id = ?`
	// Sql.UpdatePasswordRequest = `UPDATE user SET password = ? WHERE id = ?`
}

// User représente un utilisateur du forum
func AddUser(username string, email string, password string) {
	args := []string{username, email, password}
	cred := strings.Join(args, " ")
	println("DB: Attempting to add a new user : [", cred, " ]")

	Exec(query.InsertUser, username, email, password)
}
