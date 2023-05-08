package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Setup() {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	SetupRoutes(r)

	srv := &http.Server{
		Addr:         "localhost:8000",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	defer srv.Close()

	fmt.Println("Starting server on http://localhost:8000")
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		fmt.Fprintln(os.Stderr, "Failed to start server:", err)
		os.Exit(1)
	}
}
