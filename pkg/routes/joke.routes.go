package routes

import (
	"net/http"
	"strings"

	controller "github.com/imbradyboy/go-firestore-crud/pkg/controllers"
)

// manual route handler
func InitJokeRoutes() {
	http.HandleFunc("/joke", jokeHandler)
	http.HandleFunc("/joke/", jokeHandler)
}

func jokeHandler(w http.ResponseWriter, r *http.Request) {
	// check for each allowed route and send to controller
	switch {
	// routes with no path params
	case r.URL.Path == "/joke" && r.Method == "GET":
		controller.GetAllJokes(w, r)
	case r.URL.Path == "/joke" && r.Method == "POST":
		controller.AddJoke(w, r)

	// routes with path param
	case strings.HasPrefix(r.URL.Path, "/joke/") && r.Method == "GET":
		controller.GetJokeById(w, r)
	case strings.HasPrefix(r.URL.Path, "/joke/") && r.Method == "PUT":
		controller.UpdateJoke(w, r)
	case strings.HasPrefix(r.URL.Path, "/joke/") && r.Method == "DELETE":
		controller.DeleteJoke(w, r)
	default:
		http.Error(w, "404 not found", http.StatusNotFound)
	}
}
