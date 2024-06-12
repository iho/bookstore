package books

import "time"

type Book struct {
	ID            int64
	Title         string
	AuthorID      int64
	PublishedDate time.Time
}

func NewBook(id int64, title string, authorID int64, publishedDate time.Time) (*Book, error) {
	if id == 0 {
		return nil, ErrInvalidID
	}
	if title == "" {
		return nil, ErrInvalidTitle
	}
	if authorID == 0 {
		return nil, ErrInvalidAuthorID
	}
	if publishedDate.IsZero() {
		return nil, ErrInvalidPublishedDate
	}

	return &Book{
		ID:            id,
		Title:         title,
		AuthorID:      authorID,
		PublishedDate: publishedDate,
	}, nil
}
