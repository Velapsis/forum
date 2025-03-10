package web

import (
	"fmt"
	"html/template"
	"net/http"
)

func Init() {

	// Define website data
	web.Database = "web/database/forum.db"
	web.Port = ":8080"

	// Define website routes
	web.Home = ""
	web.Login = ""

	CreateWebsite()
}

func CreateWebsite() {
	http.HandleFunc("/", Index)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	http.ListenAndServe(web.Port, nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	ParseTemplate(w, "web/login.html")
	println("Executing index on port: ", web.Port)
}

func ParseTemplate(w http.ResponseWriter, tempPath string) {
	tmpl, err := template.ParseFiles(tempPath)
	println("Parsing template: ", tempPath)

	// Error management
	if tmpl == nil {
		fmt.Println("Error parsing template: ", tempPath, " does not exist")
	} else if err != nil {
		fmt.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, "")
	if err != nil {
		fmt.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
