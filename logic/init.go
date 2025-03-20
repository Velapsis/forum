package logic

import (
	"forum/web"
)

func Init() {
	// database.Init()
	web.Init()

	// Define regex
	regex.Username = `[a-zA-Z0-9_]$`
	regex.Email = `[a-zA-Z0-9]@[a-z].[a-z]$`

	// Define SQL requests
	Sql.InsertRequest = `INSERT INTO user (username, email, password) VALUES (?, ?, ?)`
	// sql.UpdateUsernameRequest = `UPDATE user SET username = ? WHERE id = ?`
	// sql.UpdateEmailRequest = `UPDATE user SET email = ? WHERE id = ?`
	// sql.UpdatePasswordRequest = `UPDATE user SET password = ? WHERE id = ?`
}
