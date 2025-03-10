package logic
import "database/sql"

var regex Regex
var user User
 var Sql SQL
 var db *sql.DB

type Regex struct {
	Username string
	Email    string
}

// DEBUG ONLY
type User struct {
	Username string
	Email    string
	Password string
	UUID     int
}

type SQL struct {
	InsertRequest string
}
