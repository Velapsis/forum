package logic

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

var db sql.DB

var oauthStateStrings = make(map[string]time.Time)

// Génération d'un état aléatoire pour OAuth
func GenerateOAuthState() string {
	state := GenerateSessionUUID()
	oauthStateStrings[state] = time.Now().Add(10 * time.Minute) // Expire après 10 minutes
	return state
}

// Vérification de l'état OAuth
func ValidateOAuthState(state string) bool {
	expiresAt, exists := oauthStateStrings[state]
	if !exists {
		return false
	}

	// Vérifier si l'état a expiré
	if time.Now().After(expiresAt) {
		delete(oauthStateStrings, state)
		return false
	}

	// Supprimer l'état après utilisation (protège contre la réutilisation)
	delete(oauthStateStrings, state)
	return true
}

// Génération d'un UUID plus sécurisé
func GenerateSessionUUID() string {
	uuid := make([]byte, 32) // Augmenter à 32 bytes pour plus de sécurité
	_, err := rand.Read(uuid)
	if err != nil {
		// En cas d'erreur, logger et générer quelque chose d'aléatoire comme fallback
		fmt.Println("Error generating random bytes:", err)
		return hex.EncodeToString([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	}
	return hex.EncodeToString(uuid)
}

// Gestion des sessions
func CreateSession(w http.ResponseWriter, userID int) (*Session, error) {
	// Nettoyer les anciennes sessions
	_, err := db.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	// Créer une nouvelle session
	sessionID := GenerateSessionUUID()
	expiresAt := time.Now().Add(2 * time.Hour) // 2 heures
	createdAt := time.Now()

	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
	}

	_, err = db.Exec("INSERT INTO sessions (id, user_id, expires_at, created_at) VALUES (?, ?, ?, ?)",
		sessionID, userID, expiresAt, createdAt)
	if err != nil {
		return nil, err
	}

	// Définir le cookie avec plus de sécurité
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,                    // Seulement transmis via HTTPS
		SameSite: http.SameSiteStrictMode, // Protection contre CSRF
		MaxAge:   int(2 * time.Hour.Seconds()),
		Expires:  expiresAt,
	}
	http.SetCookie(w, cookie)

	return session, nil
}

// Récupérer une session à partir d'un cookie et vérifier sa validité
func GetSessionFromCookie(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, fmt.Errorf("cookie not found: %w", err)
	}

	// Récupérer la session
	var session Session
	err = db.QueryRow(`
		SELECT id, user_id, expires_at, created_at 
		FROM sessions 
		WHERE id = ? AND expires_at > ?`,
		cookie.Value, time.Now()).Scan(&session.ID, &session.UserID, &session.ExpiresAt, &session.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found or expired")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &session, nil
}

// Vérifier si la session est valide et appartient à l'utilisateur
func IsValidSession(r *http.Request, userID int) bool {
	session, err := GetSessionFromCookie(r)
	if err != nil {
		return false
	}

	// Vérifier si la session appartient à l'utilisateur
	return session.UserID == userID && time.Now().Before(session.ExpiresAt)
}

// Suppression de session avec vérification d'appartenance
func DeleteSession(w http.ResponseWriter, r *http.Request) error {
	// Récupérer l'ID de session
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return fmt.Errorf("cookie not found: %w", err)
	}

	sessionID := cookie.Value

	// Vérifier que la session existe avant de la supprimer
	var exists bool
	err = db.QueryRow("SELECT 1 FROM sessions WHERE id = ?", sessionID).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("database error: %w", err)
	}

	// Supprimer la session de la base de données
	if err != sql.ErrNoRows {
		_, err = db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
		if err != nil {
			return fmt.Errorf("failed to delete session: %w", err)
		}
	}

	// Supprimer le cookie
	DeleteCookie(w, "session_id")

	return nil
}

// Suppression de cookie plus sécurisée
func DeleteCookie(w http.ResponseWriter, cookieName string) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
		Expires:  time.Now().AddDate(0, 0, -1), // Définir la date d'expiration dans le passé
	}
	http.SetCookie(w, cookie)
}

// Nettoyage périodique des sessions expirées
func DeleteExpiredSessions() error {
	_, err := db.Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now())
	return err
}

// Nouvelle fonction pour renouveler une session
func RenewSession(w http.ResponseWriter, r *http.Request) (*Session, error) {
	// Récupérer la session actuelle
	session, err := GetSessionFromCookie(r)
	if err != nil {
		return nil, err
	}

	// Définir une nouvelle date d'expiration
	newExpiresAt := time.Now().Add(2 * time.Hour)

	// Mettre à jour la base de données
	_, err = db.Exec("UPDATE sessions SET expires_at = ? WHERE id = ?",
		newExpiresAt, session.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	// Mettre à jour le cookie
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(2 * time.Hour.Seconds()),
		Expires:  newExpiresAt,
	}
	http.SetCookie(w, cookie)

	// Mettre à jour l'objet session
	session.ExpiresAt = newExpiresAt
	return session, nil
}

// Middleware de vérification de session
func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Chemins qui ne nécessitent pas d'authentification
		publicPaths := map[string]bool{
			"/":           true,
			"/login":      true,
			"/register":   true,
			"/static/":    true,
			"/api/public": true,
		}

		// Vérifier si le chemin est public
		if publicPaths[r.URL.Path] || (len(r.URL.Path) >= 8 && publicPaths[r.URL.Path[:8]]) {
			next.ServeHTTP(w, r)
			return
		}

		// Vérifier la session
		session, err := GetSessionFromCookie(r)
		if err != nil {
			// Rediriger vers la page de connexion
			http.Redirect(w, r, "/login?redirect="+r.URL.Path, http.StatusSeeOther)
			return
		}

		// Vérifier si la session n'est pas sur le point d'expirer et la renouveler si nécessaire
		if time.Until(session.ExpiresAt) < 30*time.Minute {
			RenewSession(w, r)
		}

		// Ajouter l'ID utilisateur au contexte de la requête pour une utilisation ultérieure
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", session.UserID)

		// Continuer avec la requête modifiée
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
