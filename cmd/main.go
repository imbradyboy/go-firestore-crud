package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/imbradyboy/go-firestore-crud/pkg/config"
	"github.com/imbradyboy/go-firestore-crud/pkg/routes"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	// connect to Firebase services
	config.InitializeFirebaseApp()

	// declare joke routes, send to handler to do manual switching on methods and subpaths
	routes.InitJokeRoutes()

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
