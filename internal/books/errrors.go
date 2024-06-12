package books

import "errors"

var (
	ErrInvalidID            = errors.New("books: invalid id")
	ErrInvalidTitle         = errors.New("books: invalid title")
	ErrInvalidAuthorID      = errors.New("books: invalid author id")
	ErrInvalidPublishedDate = errors.New("books: invalid published date")
)
