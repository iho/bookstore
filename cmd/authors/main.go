package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/iho/bookstore/internal/authors"
	"github.com/iho/bookstore/protos/gen/authors/v1/authorsv1connect"
	"github.com/jackc/pgx/v5"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func run() error {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://bookstore_user:bookstore_password@postgres:5432/bookstore?sslmode=disable")
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	authorsService := authors.NewAuthorsService(conn)

	mux := http.NewServeMux()
	mux.Handle(
		authorsv1connect.NewAuthorsServiceHandler(authorsService),
	)
	fmt.Println("Starting server on :8080")

	return http.ListenAndServe(":8080", h2c.NewHandler(mux, &http2.Server{}))
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
