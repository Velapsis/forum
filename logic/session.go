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
	expiresAt := time.Now().Add(2 * time.Hour ) // 2 heures
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

func deleteSession(w http.ResponseWriter, sessionID string) error {
	// Supprimer la session de la base de données
	_, err := db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	if err != nil {
	return err
	}
   
  
	// Supprimer le cookie
	deleteCookie(w, "session_id")
   
  
	return nil
   }
   

func deleteCookie(w http.ResponseWriter, cookieName string) {
	cookie := &http.Cookie{
	Name: cookieName,
	Value: "",
	Path: "/",
	HttpOnly: true,
	Expires: time.Now().AddDate(0, 0, -1), // Définir la date d'expiration dans le passé
	}
	http.SetCookie(w, cookie)
   }


   func deleteExpiredSessions() error {
	_, err := db.Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now())
	return err
   }

//  A FAIRE : Fonction de vérification de session
//  A FAIRE : Fonction de renouvellement de session
