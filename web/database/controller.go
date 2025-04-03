package database

import (
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
	GetCreatedAt string
	UpdateUsername string
	UpdateEmail string
	UpdatePassword string
}

var query Query

func DefineRequests() {
	query.InsertUser = `INSERT INTO users (username, email, password, id) VALUES (?, ?, ?, ?)`
	query.GetUserID = `SELECT id FROM users WHERE username = ?`
	query.GetUsername = `SELECT username FROM users WHERE id = ?`
	query.GetEmail = `SELECT email FROM users WHERE id = ?`
	query.GetCreatedAt = `SELECT created_at FROM users WHERE id = ?`
	query.GetPassword = `SELECT password FROM users WHERE id = ?`
	query.UpdateUsername = `UPDATE users SET username = ? WHERE id = ?`
    query.UpdateEmail = `UPDATE users SET email = ? WHERE id = ?`
    query.UpdatePassword = `UPDATE users SET password = ? WHERE id = ?`
    

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
	if id == 0 {
		return ""
	}
	
	var username string
	err := database.QueryRow(query.GetUsername, id).Scan(&username)
	if err != nil {
		println("DB: Error while scanning users: ", err.Error())
	}
	return username
}

func GetEmail(id int) string {
	if id == 0 {
		return ""
	}
	
	var email string
	err := database.QueryRow(query.GetEmail, id).Scan(&email)
	if err != nil {
		println("DB: Error while scanning users: ", err.Error())
	}
	return email
}

func GetCreatedAt(id int) string{
	var createdAtStr string
    err := database.QueryRow(query.GetCreatedAt, id).Scan(&createdAtStr)
    if err != nil {
        println("DB: Error while getting created_at:", err.Error())
        return "Unknown"
    }
    
    // Formater la date pour l'affichage
    if t, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
        return t.Format("January 2, 2006")
    }
    
    return createdAtStr
}
