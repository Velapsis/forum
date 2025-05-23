package logic

import (
	"crypto/tls"
	"fmt"
	"forum/web/database"
	"html/template"
	"net/http"

	"golang.org/x/crypto/acme/autocert"
)

func InitWebsite() {
	InitSessionDB()
	// Define website data
	website.Database = "web/database/forum.db"
	website.Port = ":8080"

	// Define website routes
	website.Home = ""
	website.Login = ""

	CreateWebsite()
}

func CreateWebsite() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/login", LoginPage)
	http.HandleFunc("/register", RegisterPage)
	http.HandleFunc("/profile", ProfilePage)
	http.HandleFunc("/post", NewPostPage)
	http.HandleFunc("/topic", NewTopicPage)
	http.HandleFunc("/logout", LogoutHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	http.ListenAndServe(website.Port, nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	session, err := GetSessionFromCookie(r)

	if session != nil {
		webpage = WebPage{
			IsConnected: true,
			UserID:      session.UserID,
			Username:    database.GetUsernameByID(session.UserID),
		}
		fmt.Println(err)
	} else {
		webpage = WebPage{
			IsConnected: false,
			UserID:      0,
			Username:    "",
		}
	}

	ParseTemplate(w, "web/index.html")
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	println("From HTML: ", r.FormValue("username"), r.FormValue("passwd"))
	Login(r.FormValue("username"), r.FormValue("passwd"), r)
	if webpage.UserID != 0 {
		// Créer une session
		_, err := CreateSession(w, webpage.UserID)
		if err != nil {
			fmt.Println("Erreur lors de la création de la session:", err)
		} else {
			// Rediriger vers la page d'accueil
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
	ParseTemplate(w, "web/login.html")
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {

	println("From HTML: ", r.FormValue("username"), r.FormValue("email"), r.FormValue("passwd"))
	Register(r.FormValue("username"), r.FormValue("email"), r.FormValue("passwd"), r)
	if webpage.UserID != 0 {
		// Créer une session
		_, err := CreateSession(w, webpage.UserID)
		if err != nil {
			fmt.Println("Erreur lors de la création de la session:", err)
		} else {
			// Rediriger vers la page d'accueil
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
	ParseTemplate(w, "web/register.html")
}

func ProfilePage(w http.ResponseWriter, r *http.Request) {
	session, err := GetSessionFromCookie(r)
	if session == nil {
		// Rediriger vers la page de connexion si non connecté
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		fmt.Println(err)
		return
	}

	// Récupérer les informations de l'utilisateur
	userID := session.UserID
	username := database.GetUsernameByID(userID)
	email := database.GetEmail(userID)
	createdAt := database.GetCreatedAt(userID)

	// Mettre à jour les données de la page
	webpage = WebPage{
		IsConnected: true,
		UserID:      userID,
		Username:    username,
		Email:       email,
		CreatedAt:   createdAt,
	}

	// Afficher la page de profil
	ParseTemplate(w, "web/profile.html")
}

func NewPostPage(w http.ResponseWriter, r *http.Request) {
	ParseTemplate(w, "web/post.html")
	CreatePost(webpage.UserID, r.FormValue("topic"), r.FormValue("title"), r.FormValue("postcontent"))
}

func NewTopicPage(w http.ResponseWriter, r *http.Request) {
	ParseTemplate(w, "web/topic.html")
	CreateTopic(webpage.UserID, r.FormValue("category"), r.FormValue("title"), r.FormValue("desc"))
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	err := Logout(w, r)
	if err != nil {
		fmt.Println("Erreur lors de la déconnexion:", err)
	}

	// Important: réinitialiser explicitement les données de l'utilisateur
	webpage = WebPage{
		IsConnected: false,
		UserID:      0,
		Username:    "",
	}

	// Rediriger vers la page d'accueil après la déconnexion
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func ParseTemplate(w http.ResponseWriter, tempPath string) {
	tmpl, err := template.ParseFiles(tempPath)
	println("HTTP: Parsing template: ", tempPath)

	// Error management
	if tmpl == nil {
		fmt.Println("Error parsing template: ", tempPath, " does not exist")
	} else if err != nil {
		fmt.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, webpage)
	if err != nil {
		fmt.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
func setupHTTPS() {
	// Gestionnaire de certificats
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		Cache:      autocert.DirCache("certs"),                 // Dossier pour stocker les certificats
		HostPolicy: autocert.HostWhitelist("votredomaine.com"), // Remplacer par votre domaine
	}

	// Configuration du serveur
	server := &http.Server{
		Addr:    ":443",
		Handler: nil, // Votre gestionnaire ici
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
			MinVersion:     tls.VersionTLS12, // Exiger TLS 1.2 au minimum
			CipherSuites: []uint16{ // Liste de suites de chiffrement sécurisées
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				// Ajouter d'autres suites selon vos besoins
			},
		},
	}

	// Démarrer le serveur
	go http.ListenAndServe(":80", certManager.HTTPHandler(nil))
	server.ListenAndServeTLS("", "")
}
