package logic

import (
	"fmt"
	database "forum/web/database"
	"hash/fnv"
	"math/rand/v2"
	"net/http"
	"regexp"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt" // Importation de bcrypt pour le hachage des mots de passe
)

// Structure pour le rate limiting
type RateLimiter struct {
	ips    map[string][]time.Time
	mu     sync.Mutex
	rate   int           // Nombre maximal de requêtes
	window time.Duration // Fenêtre de temps pour le rate limit
}

// Singleton du rate limiter
var limiter = NewRateLimiter(5, time.Minute) // 5 tentatives par minute

// Initialiser un nouveau rate limiter
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		ips:    make(map[string][]time.Time),
		rate:   rate,
		window: window,
	}
}

// Vérifier si une IP est rate limited
func (rl *RateLimiter) IsLimited(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Nettoyer les anciennes requêtes
	now := time.Now()
	if _, exists := rl.ips[ip]; exists {
		var validReqs []time.Time
		for _, t := range rl.ips[ip] {
			if now.Sub(t) <= rl.window {
				validReqs = append(validReqs, t)
			}
		}
		rl.ips[ip] = validReqs
	} else {
		rl.ips[ip] = []time.Time{}
	}

	// Vérifier si on dépasse le rate limit
	if len(rl.ips[ip]) >= rl.rate {
		return true
	}

	// Ajouter la requête actuelle
	rl.ips[ip] = append(rl.ips[ip], now)
	return false
}

// Middleware pour le rate limiting des requêtes d'authentification
func RateLimitMiddleware(next http.HandlerFunc, action string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		
		// Ajouter un identifiant d'action pour séparer les limits par type d'action
		actionKey := ip + ":" + action
		
		if limiter.IsLimited(actionKey) {
			http.Error(w, "Too many attempts, please try again later", http.StatusTooManyRequests)
			return
		}
		
		next.ServeHTTP(w, r)
	}
}

// Hachage du mot de passe avec bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Vérification du mot de passe haché
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Login modifié pour utiliser bcrypt
func Login(username string, passwd string, r *http.Request) (bool, string) {
	// Vérifier le rate limiting
	if limiter.IsLimited(r.RemoteAddr + ":login") {
		return false, "Too many login attempts, please try again later"
	}
	
	// Récupérer le hash du mot de passe de la base de données
	hashedPassword := database.GetUserPasswordHash(username)
	if hashedPassword == "" {
		return false, "Username or password incorrect"
	}
	
	// Vérifier le mot de passe avec bcrypt
	if CheckPasswordHash(passwd, hashedPassword) {
		webpage.UserID = database.GetUserID(username)
	} else {
		println("Username or password incorrect")
	}
}

// Register modifié pour utiliser bcrypt
func Register(username string, email string, passwd string, r *http.Request) (bool, string) {
	// Vérifier le rate limiting
	if limiter.IsLimited(r.RemoteAddr + ":register") {
		return false, "Too many registration attempts, please try again later"
	}
	
	println("Attempting to register to the database")
	println("Creds: ", username, email, "********") // Ne pas logger le mot de passe en clair
	
	if IsLegit(username, email, passwd) && database.IsUserAvailable(username, email) {
		println("User is legit, attempting to add to database..")
		
		// Hasher le mot de passe avant stockage
		hashedPassword, err := HashPassword(passwd)
		if err != nil {
			return false, "Error processing password"
		}
		
		// Stocker l'utilisateur avec le mot de passe haché
		database.AddUser(username, email, hashedPassword, GenerateUUID(username))
		
		// Connecter l'utilisateur
		success, msg := Login(username, passwd, r)
		return success, msg
	}
	
	return false, "Registration failed"
}

// IsLegit reste inchangé
func IsLegit(username string, email string, passwd string) bool {
	// Null check
	if username == "" {
		fmt.Println("Username is null")
		return false
	} else if email == "" {
		fmt.Println("Email is null")
		return false
	}
	println("PASS: Null check")

	// Regex check
	isUsernameValid, _ := regexp.MatchString(regex.Username, username)
	isEmailValid, _ := regexp.MatchString(regex.Email, email)
	if !isUsernameValid { 
		fmt.Println("Username is not valid")
		return false
	} else if !isEmailValid { 
		fmt.Println("Email is not valid")
		return false
	}
	println("PASS: Regex check")

	// Password check
	if len(passwd) < 5 {
		fmt.Println("Password is not valid: Too short")
		return false
	}
	hasCapsLetter, _ := regexp.MatchString(`[A-Z]`, passwd)
	if !hasCapsLetter {
		fmt.Println("Password is not valid: No capital letter")
		return false
	}
	hasOneNumber, _ := regexp.MatchString(`[0-9]`, passwd)
	if !hasOneNumber {
		fmt.Println("Password is not valid: No number")
		return false
	}
	println("PASS: Password check")

	return true
}

// GenerateUUID reste inchangé
func GenerateUUID(username string) int {
	if username == "" {
		fmt.Println("Error: Cannot generate UUID if username is null")
		return 0
	}
	uuid := fnv.New32a()
	uuid.Write([]byte(username))
	random := rand.IntN(9999999)
	return int(uuid.Sum32()) + random
}