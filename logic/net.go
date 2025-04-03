package logic

import (
	"fmt"
	"forum/web/database"
	"html/template"
	"net/http"
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
	session := GetSessionFromCookie(r)
    
    if session != nil {
        webpage = WebPage{
            IsConnected: true,
            UserID:      session.UserID,
            Username:    database.GetUsername(session.UserID),
        }
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
	if r.Method == "POST" {
        username := r.FormValue("username")
        passwd := r.FormValue("passwd")
        
        println("Username: ", username, " Password: ", passwd)
        Login(username, passwd)
        
        // Si la connexion est réussie
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
    }
    
    // Afficher le formulaire de connexion
    ParseTemplate(w, "web/login.html")
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
        username := r.FormValue("username")
        email := r.FormValue("email")
        passwd := r.FormValue("passwd")
        
        println("From HTML: ", username, email, passwd)
        Register(username, email, passwd)
        
        // Si l'inscription est réussie
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
    }
    
    // Afficher le formulaire d'inscription
    ParseTemplate(w, "web/register.html")
}

func ProfilePage(w http.ResponseWriter, r *http.Request) {
    session := GetSessionFromCookie(r)
    if session == nil {
        // Rediriger vers la page de connexion si non connecté
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    
    // Récupérer les informations de l'utilisateur
    userID := session.UserID
    username := database.GetUsername(userID)
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
