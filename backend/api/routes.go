package api

import (
	"context"
	"fmt"
	"main/directives"
	"main/graph"
	"main/services"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"golang.org/x/oauth2"
)

var config = &oauth2.Config{
	ClientID:     os.Getenv("SPOTIFY_API_GENERAL_CLIENT_ID"),
	ClientSecret: os.Getenv("SPOTIFY_API_GENERAL_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:8000/callback",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://accounts.spotify.com/authorize",
		TokenURL: "https://accounts.spotify.com/api/token",
	},
	Scopes: []string{"user-top-read"},
}

func SetupRoutes(r *chi.Mux) {
	c := graph.Config{Resolvers: &graph.Resolver{}}
	c.Directives.Auth = directives.Auth
	// var graphHandler = handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	graphHandler := handler.NewDefaultServer(graph.NewExecutableSchema(c))

	r.Get("/callback", callbackHandler)
	r.Get("/login", loginHandler)
	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", graphHandler)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}
	client := config.Client(context.Background(), token)

	releases := services.GetRecommendations(client, "week")
	fmt.Println("FOUND", len(releases), "NEW RECOMMENDATIONS!!!")

	fmt.Println("--------------------------------")
	for _, release := range releases {
		fmt.Print("Artists:", release.Artists)
		fmt.Println("\nTitle:", release.Title)
		fmt.Println("Date:", release.Date)
		fmt.Println("Type:", release.Type)
		fmt.Println("Trending Score:", release.TrendingScore)
		fmt.Println("--------------------------------")
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	url := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
