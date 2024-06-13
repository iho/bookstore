package graph

import (
	"net/http"

	"github.com/iho/bookstore/internal/cfg"
	"github.com/iho/bookstore/protos/gen/authors/v1/authorsv1connect"
	"github.com/iho/bookstore/protos/gen/books/v1/booksv1connect"
)

type Resolver struct {
	cfg              *cfg.Config
	booksv1connect   booksv1connect.BooksServiceClient
	authorsv1connect authorsv1connect.AuthorsServiceClient
}

func NewResolver(cfg *cfg.Config) *Resolver {
	return &Resolver{
		cfg:              cfg,
		booksv1connect:   booksv1connect.NewBooksServiceClient(http.DefaultClient, cfg.BookServiceUrl),
		authorsv1connect: authorsv1connect.NewAuthorsServiceClient(http.DefaultClient, cfg.AuthorSericeUrl),
	}
}
