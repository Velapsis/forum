package database

import (
	"database/sql"
	"fmt"
	"io"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var database *sql.DB
var source string

func Exec(request string) {

	// Defines the Data Source Name (DSN)
	source = "@tcp(db:3306)/forum_db"

	// Opens the database
	database, err := sql.Open("mysql", request+source)
	fmt.Println("SQL Request: ", request+source)
	if err != nil {
		fmt.Println("Error while parsing database: ", err)
	}

	// Read SQL file
	sql, err := os.Open("web/database/users.sql")
	if err != nil {
		fmt.Println("Error while reading SQL file: ", err)
	}

	// Convert SQL instructions to bytes
	sqlBytes, err := io.ReadAll(sql)
	if err != nil {
		fmt.Println("Error while converting SQL file to bytes: ", err)
	}

	// Execute the database
	output, err := database.Exec(string(sqlBytes))
	fmt.Println("Exec database output: ", output)
	if err != nil {
		fmt.Println("Error while executing database: ", err)
	}

	// Closes database when main() stops running
	defer database.Close()

}
