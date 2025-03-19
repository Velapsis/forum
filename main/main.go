package main

import (
	"fmt"
	logic "main/logic/"
	"net/http"
)

func main() {
	logic.Init()
	web.CreateWebsite()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Bienvenue sur le Forum - Configuration en cours")
	})

	fmt.Println("Serveur démarré sur le port 8080")
	http.ListenAndServe(":8080", nil)
}
