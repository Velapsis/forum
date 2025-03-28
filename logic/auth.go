package logic

import (
	// "database/sql" // Add this import for sql.ErrNoRows

	"fmt"
	database "forum/web/database"
	"hash/fnv"
	"math/rand/v2"
	"regexp"
)

func Login(username string, passwd string) {
	if database.IsUserCorrect(username, passwd) {
		// TODO: LOGIQUE LOGIN
	}
}

// Registers a new user with provided data
func Register(username string, email string, passwd string) {
	println("Attempting to register to the database")
	println("Creds: ", username, email, passwd)
	if IsLegit(username, email, passwd) && database.IsUserAvailable(username, email) {
		println("User is legit, attempting to add to database..")
		database.AddUser(username, email, passwd, GenerateUUID(username))
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
	println("PASS: Null check")

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
	println("PASS: Regex check")

	// Password check
	if len(passwd) < 5 {
		fmt.Println("Password is not valid: Too short")
		return false
	}
	hasCapsLetter, _ := regexp.MatchString(`[A-Z]`, passwd)
	if !hasCapsLetter {
		fmt.Println("Password is not valid: No capital letter")
		return false
	}
	hasOneNumber, _ := regexp.MatchString(`[0-9]`, passwd)
	if !hasOneNumber {
		fmt.Println("Password is not valid: No number")
		return false
	}
	println("PASS: Password check")

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
