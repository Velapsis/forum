package logic

import (
	"net/http"
	"time"
	"crypto/rand"
	"encoding/hex"
)

var oauthStateStrings = make(map[string]time.Time)

// Génération d'un état aléatoire pour OAuth
func generateOAuthState() string {
	state := generateUUID()
	oauthStateStrings[state] = time.Now().Add(10 * time.Minute) // Expire après 10 minutes
	return state
}
// Génération d'un UUID
func generateUUID() string {
	uuid := make([]byte, 16)
	rand.Read(uuid)
	return hex.EncodeToString(uuid)
}
// Gestion des sessions
func createSession(w http.ResponseWriter, userID int) (*Session, error) {
	// Nettoyer les anciennes sessions
	_, err := db.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	// Créer une nouvelle session
	sessionID := generateUUID()
	expiresAt := time.Now().Add(2 * time.Hour ) // 7 jours
	createdAt := time.Now()

	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt : createdAt,
	}

	_, err = db.Exec("INSERT INTO sessions (id, user_id, expires_at, created_at) VALUES (?, ?, ?, ?)",
		sessionID, userID, expiresAt, createdAt)
	if err != nil {
		return nil, err
	}

	// Définir le cookie
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Expires:  expiresAt,
	}
	http.SetCookie(w, cookie)

	return session, nil
}

func getSessionFromCookie(r *http.Request) *Session {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil
	}

	var session Session
	err = db.QueryRow("SELECT id, user_id, expires_at, created_at FROM sessions WHERE id = ?", cookie.Value).
		Scan(session.ID, session.UserID, session.ExpiresAt, session.CreatedAt)
	if err != nil {
		return nil
	}
	return &session
}
