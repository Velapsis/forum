package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// User représente un utilisateur du forum
func AddUser(request string, username string, email string, password string) {

	// Execute database with provided request
	output, err := database.Exec(request, username, email, password)
	fmt.Println("Exec database output: ", output)
	if err != nil {
		fmt.Println("Error while attempting to add user ", username, " to database")
		fmt.Println("Error: ", err)
	}
}
