-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = $1 LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name limit $1 offset $2;

-- name: CreateAuthor :one
INSERT INTO authors (
  name
) VALUES (
  $1
)
RETURNING *;

-- name: UpdateAuthor :exec
UPDATE authors
  set name = $2
WHERE id = $1
RETURNING *;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = $1
RETURNING *;