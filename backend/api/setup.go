package api

import (
	"fmt"
	"log"
	"main/middlewares"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

	SetupRoutes(r)

	fmt.Printf("Starting server on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
