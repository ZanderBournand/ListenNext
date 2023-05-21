package api

import (
	"fmt"
	"log"
	"main/middlewares"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

const defaultPort = "8000"

func Setup() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middlewares.AuthMiddleware)

	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Replace with your client's origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value for preflight OPTIONS request cache (in seconds)
	})
	r.Use(corsConfig.Handler)

	SetupRoutes(r)

	fmt.Printf("Starting server on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
