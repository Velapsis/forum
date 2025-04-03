package logic

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	database "forum/web/database"
)

var db *sql.DB
var oauthStateStrings = make(map[string]time.Time)

// InitSessionDB initialise la connexion à la base de données pour les sessions
func InitSessionDB() {
	db = database.GetDB()
	if db == nil {
		fmt.Println("Erreur: Base de données non initialisée")
	} else {
		fmt.Println("Base de données des sessions initialisée avec succès")
	}
}

// GenerateOAuthState génère un état aléatoire pour OAuth
func GenerateOAuthState() string {
	state := GenerateSessionUUID()
	oauthStateStrings[state] = time.Now().Add(10 * time.Minute) // Expire après 10 minutes
	return state
}

// ValidateOAuthState vérifie la validité d'un état OAuth
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

// GenerateSessionUUID génère un UUID sécurisé pour les sessions
func GenerateSessionUUID() string {
	uuid := make([]byte, 32) // 32 bytes pour plus de sécurité
	_, err := rand.Read(uuid)
	if err != nil {
		// En cas d'erreur, logger et générer quelque chose d'aléatoire comme fallback
		fmt.Println("Error generating random bytes:", err)
		return hex.EncodeToString([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	}
	return hex.EncodeToString(uuid)
}

// CreateSession crée une nouvelle session pour un utilisateur
func CreateSession(w http.ResponseWriter, userID int) (*Session, error) {
	if db == nil {
		InitSessionDB()
		if db == nil {
			return nil, fmt.Errorf("erreur: base de données non initialisée")
		}
	}

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

	expiresAtStr := expiresAt.Format("2006-01-02 15:04:05")
	createdAtStr := createdAt.Format("2006-01-02 15:04:05")

	_, err = db.Exec("INSERT INTO sessions (id, user_id, expires_at, created_at) VALUES (?, ?, ?, ?)",
		sessionID, userID, expiresAtStr, createdAtStr)
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

// GetSessionFromCookie récupère une session à partir d'un cookie
func GetSessionFromCookie(r *http.Request) (*Session, error) {
	if db == nil {
		InitSessionDB()
		if db == nil {
			return nil, fmt.Errorf("erreur: base de données non initialisée")
		}
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, fmt.Errorf("cookie not found: %w", err)
	}

	var session Session
	var expiresAtStr, createdAtStr string // Variables intermédiaires

	// Scanner dans des chaînes de caractères
	err = db.QueryRow("SELECT id, user_id, expires_at, created_at FROM sessions WHERE id = ?", cookie.Value).
		Scan(&session.ID, &session.UserID, &expiresAtStr, &createdAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Convertir les chaînes en time.Time
	session.ExpiresAt, err = time.Parse("2006-01-02 15:04:05", expiresAtStr)
	if err != nil {
		return nil, fmt.Errorf("erreur de conversion expires_at: %w", err)
	}

	session.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("erreur de conversion created_at: %w", err)
	}

	// Vérifier si la session a expiré
	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("session expired")
	}

	return &session, nil
}

// IsValidSession vérifie si une session est valide et appartient à l'utilisateur
func IsValidSession(r *http.Request, userID int) bool {
	session, err := GetSessionFromCookie(r)
	if err != nil {
		return false
	}

	// Vérifier si la session appartient à l'utilisateur
	return session.UserID == userID && time.Now().Before(session.ExpiresAt)
}

// Logout déconnecte un utilisateur en supprimant sa session
func Logout(w http.ResponseWriter, r *http.Request) error {
	// S'assurer que db est initialisé
	if db == nil {
		InitSessionDB()
		if db == nil {
			return fmt.Errorf("erreur: base de données non initialisée")
		}
	}

	// Récupérer le cookie de session
	cookie, err := r.Cookie("session_id")
	if err != nil {
		// Si pas de cookie, pas besoin de déconnexion
		return nil
	}

	// Supprimer la session de la base de données
	_, err = db.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)
	if err != nil {
		return err
	}

	// Supprimer le cookie
	DeleteCookie(w, "session_id")

	// Réinitialiser l'état de l'utilisateur
	// Assurez-vous que webpage est bien défini dans votre code
	/*
		webpage.IsConnected = false
		webpage.UserID = 0
		webpage.Username = ""
	*/

	return nil
}

// DeleteSession supprime une session spécifique
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

// DeleteCookie supprime un cookie du navigateur client
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

// DeleteExpiredSessions nettoie les sessions expirées dans la base de données
func DeleteExpiredSessions() error {
	if db == nil {
		InitSessionDB()
		if db == nil {
			return fmt.Errorf("erreur: base de données non initialisée")
		}
	}

	_, err := db.Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now().Format("2006-01-02 15:04:05"))
	return err
}

// RenewSession renouvelle une session existante
func RenewSession(w http.ResponseWriter, r *http.Request) (*Session, error) {
	// Récupérer la session actuelle
	session, err := GetSessionFromCookie(r)
	if err != nil {
		return nil, err
	}

	// Définir une nouvelle date d'expiration
	newExpiresAt := time.Now().Add(2 * time.Hour)
	newExpiresAtStr := newExpiresAt.Format("2006-01-02 15:04:05")

	// Mettre à jour la base de données
	_, err = db.Exec("UPDATE sessions SET expires_at = ? WHERE id = ?",
		newExpiresAtStr, session.ID)
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

// SessionMiddleware middleware pour vérifier l'authentification sur les routes protégées
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
		if publicPaths[r.URL.Path] || (len(r.URL.Path) >= 8 && r.URL.Path[:7] == "/static") {
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
