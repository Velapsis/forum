package database

import (
	"database/sql" // Ajout de l'import nécessaire
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Query struct {
	InsertPost        string
	InsertTopic       string
	InsertUser        string
	GetUserID         string
	GetUserByUsername string
	GetUserByEmail    string
	GetPasswordHash   string
	GetCreatedAt      string
	UpdateUsername    string
	UpdateEmail       string
	UpdatePassword    string
}

var query Query

func DefineRequests() {
	query.InsertUser = `INSERT INTO users (username, email, password, id) VALUES (?, ?, ?, ?)`
	query.GetUserID = `SELECT id FROM users WHERE username = ?`
	query.GetUserByUsername = `SELECT username FROM users WHERE username = ?`
	query.GetUserByEmail = `SELECT email FROM users WHERE email = ?`
	query.GetPasswordHash = `SELECT password FROM users WHERE username = ?`
	// Sql.UpdateUsernameRequest = `UPDATE user SET username = ? WHERE id = ?`
	// Sql.UpdateEmailRequest = `UPDATE user SET email = ? WHERE id = ?`
	// Sql.UpdatePasswordRequest = `UPDATE user SET password = ? WHERE id = ?`
	query.InsertPost = `INSERT INTO posts (id, title, content, topic_id, created_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	query.InsertTopic = `INSERT INTO topics (id, title, content, category_id, created_by, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	query.GetCreatedAt = `SELECT created_at FROM users WHERE id = ?`
	query.UpdateUsername = `UPDATE users SET username = ? WHERE id = ?`
	query.UpdatePassword = `UPDATE users SET password = ? WHERE id = ?`

}

// AddUser ajoute un nouvel utilisateur avec mot de passe déjà haché
func AddUser(username string, email string, hashedPassword string, id int) {
	// Ne pas afficher le mot de passe haché dans les logs
	args := []string{username, email, "********"}
	cred := strings.Join(args, " ")
	println("DB: Attempting to add a new user : [", cred, " ]")

	_, err := database.Exec(query.InsertUser, username, email, hashedPassword, id)
	if err != nil {
		println("DB: Error adding user:", err.Error())
	}
}

// IsUserAvailable vérifie si un nom d'utilisateur et un email sont disponibles
func IsUserAvailable(username string, email string) bool {
	// Vérifier si le nom d'utilisateur existe déjà
	var existingUsername string
	err := database.QueryRow(query.GetUserByUsername, username).Scan(&existingUsername)
	if err != sql.ErrNoRows {
		println("Username", username, "is already taken")
		return false
	}

	// Vérifier si l'email existe déjà
	var existingEmail string
	err = database.QueryRow(query.GetUserByEmail, email).Scan(&existingEmail)
	if err != sql.ErrNoRows {
		println("Email", email, "is already taken")
		return false
	}

	println("PASS: Availability check")
	return true
}

// IsUserCorrect - Cette fonction est obsolète avec bcrypt
// À la place, nous allons récupérer le hash du mot de passe et le vérifier ailleurs
func IsUserCorrect(username string, password string) bool {
	// Cette fonction n'est plus utilisée directement
	// Voir GetUserPasswordHash et la vérification dans le package logic
	println("WARNING: Using deprecated IsUserCorrect function")
	return false
}

// GetUserPasswordHash récupère le hash du mot de passe d'un utilisateur
func GetUserPasswordHash(username string) string {
	var passwordHash string
	err := database.QueryRow(query.GetPasswordHash, username).Scan(&passwordHash)
	if err != nil {
		if err != sql.ErrNoRows {
			println("DB: Error retrieving password hash:", err.Error())
		}
		return ""
	}
	return passwordHash
}

// GetUserID récupère l'ID d'un utilisateur à partir de son nom d'utilisateur
func GetUserID(username string) int {
	var id int
	database.QueryRow(query.GetUserID, username).Scan(&id)
	return id
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

func AddTopic(id int, title string, content string, category_id string, created_by string, created_at time.Time, updated_at time.Time) {
	args := []string{strconv.Itoa(id), content, category_id, created_by, created_by, created_at.String(), updated_at.String()}
	data := strings.Join(args, " ")
	println("DB: Attempting to add a new topic: [", data, " ]")
	if title != "" && content != "" && category_id != "" {
		Exec(query.InsertTopic, id, title, content, category_id, created_by, created_at, updated_at)
	}
	err := database.QueryRow(query.GetUserID, username).Scan(&id)
	if err != nil {
		println("DB: Error while scanning users:", err.Error())
		return 0
	}
	return id
}

// UserExists vérifie si un utilisateur existe par son nom d'utilisateur
func UserExists(username string) bool {
	var existingUsername string
	err := database.QueryRow(query.GetUserByUsername, username).Scan(&existingUsername)
	return err != sql.ErrNoRows
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

func GetCreatedAt(id int) string {
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
