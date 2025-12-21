package routes

import (
	"net/http"

	"go-api/internal/handler"
	"go-api/internal/middleware"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	router := mux.NewRouter()

	// Serve web interface
	router.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir("./web"))))

	// Serve the main index at root to make the web UI available at '/'
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	}).Methods("GET")

	// Public routes (no authentication needed)
	router.HandleFunc("/login", handler.Login).Methods("POST")
	router.HandleFunc("/create", handler.CreateUser).Methods("POST")

	// Protected routes (require JWT token)
	router.HandleFunc("/users", middleware.JWTAuth(handler.GetUsers)).Methods("GET")
	router.HandleFunc("/update", middleware.JWTAuth(handler.UpdateUser)).Methods("PUT")
	router.HandleFunc("/delete", middleware.JWTAuth(handler.DeleteUser)).Methods("DELETE")

	return router
}
