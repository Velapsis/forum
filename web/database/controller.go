package database

import (
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Query struct {
	InsertUser  string
	GetUserID   string
	GetUsername string
	GetEmail    string
	GetPassword string

	InsertPost string
}

var query Query

func DefineRequests() {
	query.InsertUser = `INSERT INTO users (username, email, password, id) VALUES (?, ?, ?, ?)`
	query.GetUserID = `SELECT id FROM users WHERE username = ?`
	query.GetUsername = `SELECT username FROM users WHERE id = ?`
	query.InsertPost = `INSERT INTO posts (id, title, content, topic_id, created_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	// Sql.UpdateUsernameRequest = `UPDATE user SET username = ? WHERE id = ?`
	// Sql.UpdateEmailRequest = `UPDATE user SET email = ? WHERE id = ?`
	// Sql.UpdatePasswordRequest = `UPDATE user SET password = ? WHERE id = ?`
}

// USERS
// User représente un utilisateur du forum
func AddUser(username string, email string, password string, id int) {
	args := []string{username, email, password}
	cred := strings.Join(args, " ")
	println("DB: Attempting to add a new user : [", cred, " ]")

	Exec(query.InsertUser, username, email, password, id)
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
		if err := rows.Scan(&dbUsername, &dbPassword); err != nil {
			println("DB: Error scanning users: ", err.Error())
			return false
		}
		if dbUsername == username && dbPassword == password {
			isCorrect = true
			break
		} else {
			isCorrect = false
		}

	}
	return isCorrect
}

func GetUserID(username string) int {
	var id int
	err := database.QueryRow(query.GetUserID, username).Scan(&id)
	if err != nil {
		println("DB: Error while scanning users: ", err.Error())
	}
	return id
}

func GetUsername(id int) string {
	var username string
	err := database.QueryRow(query.GetUsername, id).Scan(&username)
	if err != nil {
		println("DB: Error while scanning users: ", err.Error())
	}
	return username
}

// POSTS
func AddPost(id int, title string, content string, topic_id string, created_by string, created_at time.Time, updated_at time.Time) {
	args := []string{strconv.Itoa(id), content, topic_id, created_by, created_by, created_at.String(), updated_at.String()}
	data := strings.Join(args, " ")
	println("DB: Attempting to add a new post: [", data, " ]")
	if title != "" && content != "" && topic_id != "" {
		Exec(query.InsertPost, id, title, content, topic_id, created_by, created_at, updated_at)
	}
}
