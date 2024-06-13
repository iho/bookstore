package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/iho/bookstore/internal/cfg"
	"github.com/iho/bookstore/internal/gateway/graph"
	"github.com/iho/bookstore/internal/gateway/loaders"
	"github.com/rs/cors"
)

func main() {
	var cfg = &cfg.Config{}
	cfg.AuthorSericeUrl = "http://authors:9090"
	cfg.BookServiceUrl = "http://books:8080"

	// create the query handler
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver(cfg)}))

	router := chi.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	})

	srv.Use(extension.Introspection{})

	router.Handle("/", playground.Handler("My GraphQL App", "/app"))
	router.Handle("/app", loaders.Middleware(cfg, c.Handler(srv)))

	// register the wrapped handler
	fmt.Println("Server is running on port 10000")
	if err := http.ListenAndServe(":10000", router); err != nil {
		log.Fatal(err)
	}
}
