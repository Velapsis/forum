package database

import (
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// User représente un utilisateur du forum
func AddUser(username string, email string, password string) {
	args := []string{username, email, password}
	token := strings.Join(args, ":")
	Exec(token)
}
