package logic

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"time"
	"fmt"
	database "forum/web/database"
)

var db *sql.DB

var oauthStateStrings = make(map[string]time.Time)

func InitSessionDB() {
    db = database.GetDB()
    if db == nil {
        fmt.Println("Erreur: Base de données non initialisée")
    } else {
        fmt.Println("Base de données des sessions initialisée avec succès")
    }
}

// Génération d'un état aléatoire pour OAuth
func GenerateOAuthState() string {
	state := GenerateSessionUUID()
	oauthStateStrings[state] = time.Now().Add(10 * time.Minute) // Expire après 10 minutes
	return state
}

// Génération d'un UUID
func GenerateSessionUUID() string {
	uuid := make([]byte, 16)
	rand.Read(uuid)
	return hex.EncodeToString(uuid)
}

// Gestion des sessions
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

func GetSessionFromCookie(r *http.Request) *Session {
	if db == nil {
        InitSessionDB()
        if db == nil {
            return nil
        }
    }

    cookie, err := r.Cookie("session_id")
    if err != nil {
        return nil
    }

    var session Session
    var expiresAtStr, createdAtStr string // Variables intermédiaires
    
    // Scanner dans des chaînes de caractères
    err = db.QueryRow("SELECT id, user_id, expires_at, created_at FROM sessions WHERE id = ?", cookie.Value).
        Scan(&session.ID, &session.UserID, &expiresAtStr, &createdAtStr)
    if err != nil {
        fmt.Println("Erreur récupération session:", err)
        return nil
    }
    
    // Convertir les chaînes en time.Time
	session.ExpiresAt, err = time.Parse("2006-01-02 15:04:05", expiresAtStr)
    if err != nil {
        fmt.Println("Erreur de conversion expires_at:", err)
        return nil
    }
    
    session.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
    if err != nil {
        fmt.Println("Erreur de conversion created_at:", err)
        return nil
    }
    
    return &session
}

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
webpage.IsConnected = false
webpage.UserID = 0
webpage.Username = ""

return nil
}


func DeleteSession(w http.ResponseWriter, sessionID string) error {
	// Supprimer la session de la base de données
	_, err := db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	if err != nil {
		return err
	}

	// Supprimer le cookie
	DeleteCookie(w, "session_id")

	return nil
}

func DeleteCookie(w http.ResponseWriter, cookieName string) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().AddDate(0, 0, -1), // Définir la date d'expiration dans le passé
	}
	http.SetCookie(w, cookie)
}

func DeleteExpiredSessions() error {
	_, err := db.Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now())
	return err
}

//  A FAIRE : Fonction de vérification de session
//  A FAIRE : Fonction de renouvellement de session
