package logic

import (
	"fmt"
	"forum/web/database"
	"html/template"
	"net/http"
)

func InitWebsite() {

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
	http.HandleFunc("/post", PostPage)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	http.ListenAndServe(website.Port, nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	ParseTemplate(w, "web/index.html")
	webpage = WebPage{
		IsConnected: true,
		Username:    database.GetUsername(webpage.UserID),
	}
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	ParseTemplate(w, "web/login.html")
	println("Username: ", r.FormValue("username"), " Password: ", r.FormValue("passwd"))
	Login(r.FormValue("username"), r.FormValue("passwd"))

}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	ParseTemplate(w, "web/register.html")
	println("From HTML: ", r.FormValue("username"), r.FormValue("email"), r.FormValue("passwd"))
	Register(r.FormValue("username"), r.FormValue("email"), r.FormValue("passwd"))
}

func ProfilePage(w http.ResponseWriter, r *http.Request) {
	ParseTemplate(w, "web/profile.html")
}

func PostPage(w http.ResponseWriter, r *http.Request) {
	ParseTemplate(w, "web/post.html")
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
