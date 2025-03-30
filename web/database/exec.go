package database

import (
	"database/sql"
	"io"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var database *sql.DB

func Exec(query string, args ...interface{}) {
	output, err := database.Exec(query, args...)
	if err != nil {
		println("DB: Error while executing database: ")
		println("  Parameters: ", args)
		println("  Query: ", query)
		println("  Error: ", err.Error())
	}
	println("DB: Exec database output: ", output)
}

func Connect() {

	var err error

	source := "root:rootpassword@tcp(db:3306)/forum_db?multiStatements=true"

	// Connects to the database
	println("Open database using dsn: ", source)
	database, err = sql.Open("mysql", source)
	if err != nil {
		println("Error while connecting to the database: ", err.Error())
	}

	println("Connecting to database..")
	if err = database.Ping(); err != nil {
		println("Error while pinging the database: ", err.Error())
	} else {
		println("Successfully connected to the database!")
	}

	// Read SQL file
	println("Reading SQL file..")
	sql, err := os.Open("web/database/init-db.sql")
	if err != nil {
		println("Error while reading SQL file: ", err)
	}

	// Convert SQL instructions to bytes
	println("Converting SQL file to bytes..")
	sqlBytes, err := io.ReadAll(sql)
	if err != nil {
		println("Error while converting SQL file to bytes: ", err)
	}

	// Execute the database
	output, err := database.Exec(string(sqlBytes))
	println("Exec database output: ", output)
	if err != nil {
		println("Error while executing database: ", err.Error())
	}

}
