package logic

import (
	"database/sql" // Add this import for sql.ErrNoRows
	"fmt"
	"hash/fnv"
	"math/rand/v2"
	"regexp"
)

func Login(username string, passwd string) {
}

// Registers a new user with provided data
func Register(username string, email string, passwd string) {
	if IsLegit(username, email, passwd) {
		// web.AddUser(Sql.InsertRequest, username, email, passwd)
	}
}

// Checks if the user inputs already exist in the database
func IsLegit(username string, email string, passwd string) bool {
	// Null check
	if username == "" {
		fmt.Println("Username is null")
		return false
	} else if email == "" {
		fmt.Println("Email is null")
		return false
	}

	// Regex check
	isUsernameValid, _ := regexp.MatchString(regex.Username, username)
	isEmailValid, _ := regexp.MatchString(regex.Email, email)
	if !isUsernameValid { // Fixed inverted logic
		fmt.Println("Username is not valid")
		return false
	} else if !isEmailValid { // Fixed inverted logic
		fmt.Println("Email is not valid")
		return false
	}

	// Password check
	if len(passwd) < 5 {
		fmt.Println("Password is not valid: Too short")
		return false
	}
	hasCapsLetter, _ := regexp.MatchString(`[A-Z]{1}$`, passwd)
	if !hasCapsLetter {
		fmt.Println("Password is not valid: No capital letter")
		return false
	}
	hasOneNumber, _ := regexp.MatchString(`[0-9]{1}$`, passwd)
	if !hasOneNumber {
		fmt.Println("Password is not valid: No number")
		return false
	}

	// Database check for username
	var existingUsername string
	err := db.QueryRow("SELECT username FROM user WHERE username = ?", username).Scan(&existingUsername)
	if err == nil {
		fmt.Println("Username already exists")
		return false
	} else if err != sql.ErrNoRows {
		fmt.Println("Database error:", err)
		return false
	}

	// Database check for email
	var existingEmail string
	err = db.QueryRow("SELECT email FROM user WHERE email = ?", email).Scan(&existingEmail)
	if err == nil {
		fmt.Println("Email already exists")
		return false
	} else if err != sql.ErrNoRows {
		fmt.Println("Database error:", err)
		return false
	}

	return true
}

func GenerateUUID(username string) int {
	if username == "" {
		fmt.Println("Error: Cannot generate UUID if username is null")
		return 0
	}
	uuid := fnv.New32a()
	uuid.Write([]byte(username))
	random := rand.IntN(9999999)
	return int(uuid.Sum32()) + random
}
