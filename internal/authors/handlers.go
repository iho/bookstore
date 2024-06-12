package authors

import (
	"context"
	"fmt"
	"strconv"

	"connectrpc.com/connect"
	"github.com/iho/bookstore/internal/authors/db"
	v1 "github.com/iho/bookstore/protos/gen/authors/v1"
)

type AuthorsService struct {
	pgDB *db.Queries
}

func NewAuthorsService(pgDB db.DBTX) *AuthorsService {
	return &AuthorsService{
		pgDB: db.New(pgDB),
	}
}

func (as *AuthorsService) ListAuthors(ctx context.Context, req *connect.Request[v1.ListAuthorsRequest]) (*connect.Response[v1.ListAuthorsResponse], error) {
	dbAuthors, err := as.pgDB.ListAuthors(ctx, db.ListAuthorsParams{
		Limit:  req.Msg.Limit,
		Offset: req.Msg.Offset,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list authors: %w", err)
	}

	authors := make([]*v1.Author, 0, len(dbAuthors))
	for _, dbAuthor := range dbAuthors {
		authors = append(authors, &v1.Author{
			Id:   strconv.FormatInt(dbAuthor.ID, 10),
			Name: dbAuthor.Name,
		})
	}

	return &connect.Response[v1.ListAuthorsResponse]{
		Msg: &v1.ListAuthorsResponse{
			Authors: authors,
		},
	}, nil
}

func (as AuthorsService) GetAuthor(ctx context.Context, req *connect.Request[v1.GetAuthorRequest]) (*connect.Response[v1.GetAuthorResponse], error) {
	id, err := strconv.ParseInt(req.Msg.Id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ID: %w", err)
	}

	dbAuthor, err := as.pgDB.GetAuthor(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get author: %w", err)
	}

	return &connect.Response[v1.GetAuthorResponse]{
		Msg: &v1.GetAuthorResponse{
			Author: &v1.Author{
				Id:   strconv.FormatInt(dbAuthor.ID, 10),
				Name: dbAuthor.Name,
			},
		},
	}, nil
}

func (as *AuthorsService) CreateAuthor(ctx context.Context, req *connect.Request[v1.CreateAuthorRequest]) (*connect.Response[v1.CreateAuthorResponse], error) {
	dbAuthor, err := as.pgDB.CreateAuthor(ctx, req.Msg.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create author: %w", err)
	}

	return &connect.Response[v1.CreateAuthorResponse]{
		Msg: &v1.CreateAuthorResponse{
			Author: &v1.Author{
				Id:   strconv.FormatInt(dbAuthor.ID, 10),
				Name: dbAuthor.Name,
			},
		},
	}, nil
}

func (as *AuthorsService) UpdateAuthor(ctx context.Context, req *connect.Request[v1.UpdateAuthorRequest]) (*connect.Response[v1.UpdateAuthorResponse], error) {
	id, err := strconv.ParseInt(req.Msg.Id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ID: %w", err)
	}

	err = as.pgDB.UpdateAuthor(ctx, db.UpdateAuthorParams{
		ID:   id,
		Name: req.Msg.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update author: %w", err)
	}

	// maybe transaction is needed here
	dbAuthor, err := as.pgDB.GetAuthor(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get author: %w", err)
	}

	return &connect.Response[v1.UpdateAuthorResponse]{
		Msg: &v1.UpdateAuthorResponse{
			Author: &v1.Author{
				Id:   strconv.FormatInt(dbAuthor.ID, 10),
				Name: dbAuthor.Name,
			},
		},
	}, nil
}

func (as *AuthorsService) DeleteAuthor(ctx context.Context, req *connect.Request[v1.DeleteAuthorRequest]) (*connect.Response[v1.DeleteAuthorResponse], error) {
	id, err := strconv.ParseInt(req.Msg.Id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ID: %w", err)
	}

	err = as.pgDB.DeleteAuthor(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete author: %w", err)
	}
	return &connect.Response[v1.DeleteAuthorResponse]{
		Msg: &v1.DeleteAuthorResponse{
			Status: true,
		},
	}, nil
}
