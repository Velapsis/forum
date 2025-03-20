package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "gdhs")
		log.Print("bjr")
	})

	fmt.Println("Serveur démarré sur le port 8080")
	http.ListenAndServe(":8080", nil)
}
