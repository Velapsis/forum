package logic
import ("database/sql"
		"time"
		)

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
	// Provider   string
	// ProviderID string
	// CreatedAt  time.Time
}
// type GoogleUserInfo struct {
// 	ID            string `json:"id"`
// 	Email         string `json:"email"`
// 	VerifiedEmail bool   `json:"verified_email"`
// 	Name          string `json:"name"`
// 	Picture       string `json:"picture"`
// }

// type GithubUserInfo struct {
// 	ID        int    `json:"id"`
// 	Login     string `json:"login"`
// 	Name      string `json:"name"`
// 	Email     string `json:"email"`
// 	AvatarURL string `json:"avatar_url"`
// }
type Session struct {
	ID        string
	UserID    int
	ExpiresAt time.Time
	CreatedAt time.Time
}

type SQL struct {
	InsertRequest string
}
