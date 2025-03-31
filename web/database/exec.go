package database

import (
	"database/sql"
	"io"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var database *sql.DB

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

func Connect() {

	var err error

	source := "root:rootpassword@tcp(db:3306)/forum_db?multiStatements=true"

	// Connects to the database
	println("Open database using dsn: ", source)
	database, err = sql.Open("mysql", source)
	if err != nil {
		println("Error while connecting to the database: ", err)
	}

	println("Connecting to database..")
	if err = database.Ping(); err != nil {
		println("Error while pinging the database: ", err)
	}
	println("Successfully connected to the database!")

	// Read SQL file
	println("Reading SQL file..")
	sql, err := os.Open("web/database/init-db.sql")
	if err != nil {
		println("Error while reading SQL file: ", err)
	}

	println("SQL file: ", sql)

	// Convert SQL instructions to bytes
	println("Converting SQL file to bytes..")
	sqlBytes, err := io.ReadAll(sql)
	if err != nil {
		println("Error while converting SQL file to bytes: ", err)
	}

	println("Compiled SQL file: ", sqlBytes)
	println("Decompiled SQL file: ", string(sqlBytes))

	// Execute the database
	println("Attempting to execute database..")
	output, err := database.Exec(string(sqlBytes))
	println("DB: Exec database output: ", output)
	if err != nil {
		println("Error while executing database: ", err)
	}

	defer database.Close()

}
