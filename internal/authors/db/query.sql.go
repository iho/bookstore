// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package db

import (
	"context"
)

const createAuthor = `-- name: CreateAuthor :one
INSERT INTO authors (
  name
) VALUES (
  $1
)
RETURNING id, name
`

func (q *Queries) CreateAuthor(ctx context.Context, name string) (Author, error) {
	row := q.db.QueryRow(ctx, createAuthor, name)
	var i Author
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const deleteAuthor = `-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1
RETURNING id, name
`

func (q *Queries) DeleteAuthor(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteAuthor, id)
	return err
}

const getAuthor = `-- name: GetAuthor :one
SELECT id, name FROM authors
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAuthor(ctx context.Context, id int64) (Author, error) {
	row := q.db.QueryRow(ctx, getAuthor, id)
	var i Author
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const listAuthors = `-- name: ListAuthors :many
SELECT id, name FROM authors
ORDER BY name limit $1 offset $2
`

type ListAuthorsParams struct {
	Limit  int32
	Offset int32
}

func (q *Queries) ListAuthors(ctx context.Context, arg ListAuthorsParams) ([]Author, error) {
	rows, err := q.db.Query(ctx, listAuthors, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Author
	for rows.Next() {
		var i Author
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAuthor = `-- name: UpdateAuthor :exec
UPDATE authors
  set name = $2
WHERE id = $1
RETURNING id, name
`

type UpdateAuthorParams struct {
	ID   int64
	Name string
}

func (q *Queries) UpdateAuthor(ctx context.Context, arg UpdateAuthorParams) error {
	_, err := q.db.Exec(ctx, updateAuthor, arg.ID, arg.Name)
	return err
}
