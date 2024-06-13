package loaders

import (
	"context"

	"connectrpc.com/connect"
	"github.com/iho/bookstore/internal/gateway/graph/model"
	authorsV1 "github.com/iho/bookstore/protos/gen/authors/v1"
	"github.com/iho/bookstore/protos/gen/authors/v1/authorsv1connect"
)

type authorLoader struct {
	authorsv1connect authorsv1connect.AuthorsServiceClient
}

func (l *authorLoader) getAuthors(ctx context.Context, keys []string) ([]*model.Author, []error) {
	authors := make([]*model.Author, len(keys))
	errors := make([]error, len(keys))
	for i, key := range keys {
		req := connect.NewRequest(&authorsV1.GetAuthorRequest{Id: key})
		res, err := l.authorsv1connect.GetAuthor(ctx, req)
		if err != nil {
			errors[i] = err
			continue
		}
		authors[i] = &model.Author{
			ID:   res.Msg.Author.Id,
			Name: res.Msg.Author.Name,
		}
	}
	return authors, errors
}

// GetAuthor returns single author by id efficiently
func GetAuthor(ctx context.Context, authorID string) (*model.Author, error) {
	loaders := For(ctx)
	return loaders.AuthourLoader.Load(ctx, authorID)
}

// GetAuthours returns many authors by ids efficiently
func GetAuthours(ctx context.Context, authorIDs []string) ([]*model.Author, error) {
	loaders := For(ctx)
	return loaders.AuthourLoader.LoadAll(ctx, authorIDs)
}
