package main

import (
	"log"
	"net/http"
	"os"

	"VisulModerator/handlers"
)

func main() {
	http.HandleFunc("/analyze", handlers.AnalyzeHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server started at :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
