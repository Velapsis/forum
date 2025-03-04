package web

import (
	"database/sql"
	"fmt"
	"io"
	"os"
)

func Init() {

	// Opens the database
	database, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		fmt.Println("Error while parsing database: ", err)
	}

	// Read SQL file
	sql, err := os.Open("users.sql")
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
