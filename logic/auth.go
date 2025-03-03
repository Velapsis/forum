package logic

import (
	"fmt"
	"hash/fnv"
)

func Login(username string, passwd string) {
}

// Registers a new user with provided data
func Register(username string, email string, passwd string) {
	if IsLegit(username, email, passwd) {
		// DEBUG ONLY
		user.Username = username
		user.Email = email
		user.Password = passwd
		user.UUID = GenerateUUID(username)
		// TODO: INSERT USER IN DATABASE
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
	if username != regex.Username {
		fmt.Println("Username is not valid")
		return false
	} else if email != regex.Email {
		fmt.Println("Email is not valid")
		return false
	} else if passwd != regex.Password {
		fmt.Println("Password is not valid")
		return false
	}

	// Database check
	// TODO: CHECK FOR INPUT USERNAME IN THE DATABASE
	// TODO: CHECK FOR INPUT EMAIL IN THE DATABASE

	return true
}

func GenerateUUID(username string) int {
	uuid := fnv.New32a()
	uuid.Write([]byte(username))
	return int(uuid.Sum32())
}
