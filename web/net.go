package web

import (
	"net/http"
	"html/template"
)

func CreateWebsite() {
	http.HandleFunc("/", Index)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	http.ListenAndServe(":8080", nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	ParseTemplate("web\templates\index.html")
}

func ParseTemplate(w http.ResponseWriter, template string) {
	tmpl, err := template.ParseFiles(template)

	if err != nil {
		fmt.Println("Error parsing template:", err)
		fmt.Println("Template path: ", web.Template)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, artists)
	if err != nil {
		fmt.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}