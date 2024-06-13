package loaders

import (
	"context"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"github.com/iho/bookstore/internal/cfg"
	"github.com/iho/bookstore/internal/gateway/graph/model"
	"github.com/iho/bookstore/protos/gen/authors/v1/authorsv1connect"
	"github.com/iho/bookstore/protos/gen/books/v1/booksv1connect"
	"github.com/iho/bookstore/protos/gen/orders/v1/ordersv1connect"
	"github.com/vikstrous/dataloadgen"

	booksV1 "github.com/iho/bookstore/protos/gen/books/v1"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

// Middleware injects data loaders into the context
func Middleware(conn *cfg.Config, next http.Handler) http.Handler {
	// return a middleware that injects the loader to the request context
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader := NewLoaders(conn)
		r = r.WithContext(context.WithValue(r.Context(), loadersKey, loader))
		next.ServeHTTP(w, r)
	})
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

// Get returns single book by id efficiently
func GetBook(ctx context.Context, bookID string) (*model.Book, error) {
	loaders := For(ctx)
	return loaders.BookLoader.Load(ctx, bookID)
}

// GetBooks returns many books by ids efficiently
func GetBooks(ctx context.Context, bookIDs []string) ([]*model.Book, error) {
	loaders := For(ctx)
	return loaders.BookLoader.LoadAll(ctx, bookIDs)
}

// Loaders wrap your data loaders to inject via middleware
type Loaders struct {
	BookLoader    *dataloadgen.Loader[string, *model.Book]
	AuthourLoader *dataloadgen.Loader[string, *model.Author]
	OrderLoader   *dataloadgen.Loader[string, *model.Order]
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders(cfg *cfg.Config) *Loaders {
	bl := &bookLoader{
		booksv1connect: booksv1connect.NewBooksServiceClient(
			http.DefaultClient,
			cfg.BookServiceUrl,
		),
	}

	al := &authorLoader{
		authorsv1connect: authorsv1connect.NewAuthorsServiceClient(
			http.DefaultClient,
			cfg.AuthorSericeUrl,
		),
	}

	ol := &orderLoader{
		ordersv1connect: ordersv1connect.NewOrdersServiceClient(
			http.DefaultClient,
			cfg.OrderServiceUrl,
		),
	}

	return &Loaders{
		BookLoader:    dataloadgen.NewLoader(bl.getBooks, dataloadgen.WithWait(time.Millisecond)),
		AuthourLoader: dataloadgen.NewLoader(al.getAuthors, dataloadgen.WithWait(time.Millisecond)),
		OrderLoader:   dataloadgen.NewLoader(ol.getOrders, dataloadgen.WithWait(time.Millisecond)),
	}
}

type bookLoader struct {
	booksv1connect booksv1connect.BooksServiceClient
}

func (l *bookLoader) getBooks(ctx context.Context, keys []string) ([]*model.Book, []error) {
	books := make([]*model.Book, len(keys))
	errors := make([]error, len(keys))
	req := connect.NewRequest(&booksV1.ListBooksRequest{Ids: keys})
	bookResp, err := l.booksv1connect.ListBooks(ctx, req)
	if err != nil {
		for i := range errors {
			errors[i] = err
		}
		return books, errors
	}

	for i, book := range bookResp.Msg.GetBooks() {
		books[i] = &model.Book{
			ID:            book.Id,
			Title:         book.Title,
			PublishedDate: book.PublishedDate,
		}
	}

	return books, errors
}
