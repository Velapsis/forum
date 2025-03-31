package logic

import (
	"time"
)

var website Website
var webpage WebPage
var user User

var regex Regex

type Website struct {
	Port     string
	Database string

	Home  string
	Login string
}

type User struct {
	Username string
	Email    string
	Password string
}

type WebPage struct {
	IsConnected bool
	Username    string
	UserID      int
}

type Regex struct {
	Username string
	Email    string
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
