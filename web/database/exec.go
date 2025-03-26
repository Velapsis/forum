package database

import (
	"database/sql"
	"io"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var database *sql.DB
var source string

func Init() {

	DefineRequests()
	var err error
	// user:passwd@tcp(addr:port)/db_name
	source := "forum_user:forum_password@tcp(localhost:3306)/forum_db"

	// Connects to the database
	database, err = sql.Open("mysql", source)
	if err != nil {
		println("Error while connecting to the database: ", err)
	}

	// Read SQL file
	sql, err := os.Open("web/database/init-db.sql")
	if err != nil {
		println("Error while reading SQL file: ", err)
	}

	// Convert SQL instructions to bytes
	sqlBytes, err := io.ReadAll(sql)
	if err != nil {
		println("Error while converting SQL file to bytes: ", err)
	}

	// Execute the database
	output, err := database.Exec(string(sqlBytes))
	println("DB: Exec database output: ", output)
	if err != nil {
		println("Error while executing database: ", err)
	}

	// Closes database when main() stops running
	defer database.Close()

}

func Exec(query string, username string, email string, password string) {
	output, err := database.Exec(query, username, email, password)
	if err != nil {
		println("DB: Error while executing database: ")
		println("  User: ", username, email, password)
		println("  Query: ", query)
		println("  Error: ", err)
	}
	println("DB: Exec database output: ", output)
}

func PingTest() {
	err := database.Ping()
	if err != nil {
		println("Error while pinging database: ", err)
	}
}
