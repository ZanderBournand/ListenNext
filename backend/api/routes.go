package api

import (
	"fmt"
	"main/config"
	"main/directives"
	"main/graph"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"golang.org/x/oauth2"
)

func SetupRoutes(r *chi.Mux) {
	c := graph.Config{Resolvers: &graph.Resolver{}}
	c.Directives.Auth = directives.Auth
	c.Directives.Spotify = directives.SpotifyAuth
	graphHandler := handler.NewDefaultServer(graph.NewExecutableSchema(c))

	r.Get("/callback", callbackHandler)
	r.Get("/login", loginHandler)
	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", graphHandler)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	fmt.Println("SPOTIFY CALLBACK CODE:", code)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	url := config.SpotifyOAuth.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
