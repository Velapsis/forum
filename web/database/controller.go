package database

import (
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Query struct {
	InsertUser string
	GetUserID  string
}

var query Query

func DefineRequests() {
	query.InsertUser = `INSERT INTO users (username, email, password, id) VALUES (?, ?, ?, ?)`
	query.GetUserID = `SELECT id FROM users WHERE username = ?`
	// Sql.UpdateUsernameRequest = `UPDATE user SET username = ? WHERE id = ?`
	// Sql.UpdateEmailRequest = `UPDATE user SET email = ? WHERE id = ?`
	// Sql.UpdatePasswordRequest = `UPDATE user SET password = ? WHERE id = ?`
}

// User représente un utilisateur du forum
func AddUser(username string, email string, password string, id int) {
	args := []string{username, email, password}
	cred := strings.Join(args, " ")
	println("DB: Attempting to add a new user : [", cred, " ]")

	database.Exec(query.InsertUser, username, email, password, id)
}

func IsUserAvailable(username string, email string) bool {
	rows, err := database.Query("SELECT username, email FROM users")
	if err != nil {
		println("DB: Error querying users table:", err.Error())
		return false
	}

	for rows.Next() {
		var dbUsername, dbEmail string
		if err := rows.Scan(&dbUsername, &dbEmail); err != nil {
			println("DB: Error scanning users: ", err.Error())
			return false
		}
		if dbUsername == username {
			println("Username ", dbUsername, " is already taken")
			return false
		} else if dbEmail == email {
			println("Email ", dbEmail, " is already taken")
		}
	}

	println("PASS: Availaibility check")
	return true
}

func IsUserCorrect(username string, password string) bool {
	rows, err := database.Query("SELECT username, password FROM users")
	if err != nil {
		println("DB: Error querying users table:", err.Error())
		return false
	}

	var isCorrect bool
	for rows.Next() {
		var dbUsername, dbPassword string
		if err := rows.Scan(&dbUsername); err != nil {
			println("DB: Error scanning users: ", err.Error())
			return false
		}
		if dbUsername == username && dbPassword == password {
			isCorrect = true
		} else {
			isCorrect = false
		}

	}
	return isCorrect
}
