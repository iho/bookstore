package books

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"strconv"
	"time"

	"connectrpc.com/connect"
	v1 "github.com/iho/bookstore/protos/gen/books/v1"
	redis "github.com/redis/go-redis/v9"
)

const (
	booksKey          = "books"
	booksIDCounterKey = "books_id"
	JSONDateFormat    = "2006-01-02T15:04:05.000Z"
)

type BooksService struct {
	rdb *redis.Client
}

func NewBooksService(rdb *redis.Client) *BooksService {
	return &BooksService{
		rdb: rdb,
	}
}

func (bs *BooksService) ListBooks(ctx context.Context, req *connect.Request[v1.ListBooksRequest]) (*connect.Response[v1.ListBooksResponse], error) {
	ids := make([]string, 0, len(req.Msg.GetIds()))
	for _, redisId := range req.Msg.GetIds() {
		id, err := strconv.ParseInt(redisId, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ID: [id=%s] %w", redisId, err)
		}
		ids = append(ids, newBookID(id))
	}
	redisBooks, err := bs.rdb.MGet(ctx, ids...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to list books: %w", err)
	}

	books := make([]*v1.Book, 0, len(redisBooks))
	for _, book := range redisBooks {
		var bookObj Book

		bookBytes, ok := book.(string)
		if !ok {
			return nil, fmt.Errorf("failed to convert book to []byte")
		}

		if err := gob.NewDecoder(bytes.NewReader([]byte(bookBytes))).Decode(&bookObj); err != nil {
			return nil, fmt.Errorf("failed to decode book: %w", err)
		}

		books = append(books, &v1.Book{
			Id:            strconv.FormatInt(bookObj.ID, 10),
			Title:         bookObj.Title,
			AuthorId:      strconv.FormatInt(bookObj.AuthorID, 10),
			PublishedDate: bookObj.PublishedDate.Format(time.RFC3339),
		})
	}

	return &connect.Response[v1.ListBooksResponse]{
		Msg: &v1.ListBooksResponse{
			Books: books,
		},
	}, nil
}

func (bs *BooksService) GetBook(ctx context.Context, req *connect.Request[v1.GetBookRequest]) (*connect.Response[v1.GetBookResponse], error) {
	id, err := strconv.ParseInt(req.Msg.Id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ID: %w", err)
	}
	book, err := bs.rdb.Get(ctx, newBookID(id)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get book: [id=%d] %w", id, err)
	}

	var bookObj Book
	if err := gob.NewDecoder(bytes.NewReader([]byte(book))).Decode(&bookObj); err != nil {
		return nil, fmt.Errorf("failed to decode book: %w", err)
	}

	return &connect.Response[v1.GetBookResponse]{
		Msg: &v1.GetBookResponse{
			Book: &v1.Book{
				Id:            strconv.FormatInt(bookObj.ID, 10),
				Title:         bookObj.Title,
				AuthorId:      strconv.FormatInt(bookObj.AuthorID, 10),
				PublishedDate: bookObj.PublishedDate.Format(time.RFC3339),
			},
		},
	}, nil
}

func (bs *BooksService) CreateBook(ctx context.Context, req *connect.Request[v1.CreateBookRequest]) (*connect.Response[v1.CreateBookResponse], error) {
	authorID, err := strconv.ParseInt(req.Msg.AuthorId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse author ID: [author_id=%d] %w", authorID, err)
	}

	publishedDate, err := time.Parse(JSONDateFormat, req.Msg.PublishedDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse published date: [published_date=%s] %w", req.Msg.PublishedDate, err)
	}

	bookID := bs.GetNewRedisID()

	book, err := NewBook(bookID, req.Msg.Title, authorID, publishedDate)
	if err != nil {
		return nil, fmt.Errorf("failed to create book: %w", err)
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(book); err != nil {
		return nil, fmt.Errorf("failed to encode book: %w", err)
	}

	if err := bs.rdb.Set(ctx, newBookID(book.ID), buf.Bytes(), 0).Err(); err != nil {
		return nil, fmt.Errorf("failed to save book: [id=%d] %w", book.ID, err)
	}

	return &connect.Response[v1.CreateBookResponse]{
		Msg: &v1.CreateBookResponse{
			Book: &v1.Book{
				Id:            strconv.FormatInt(book.ID, 10),
				Title:         book.Title,
				AuthorId:      strconv.FormatInt(book.AuthorID, 10),
				PublishedDate: book.PublishedDate.Format(time.RFC3339),
			},
		},
	}, nil
}

func (bs *BooksService) UpdateBook(ctx context.Context, req *connect.Request[v1.UpdateBookRequest]) (*connect.Response[v1.UpdateBookResponse], error) {
	authorID, err := strconv.ParseInt(req.Msg.AuthorId, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse author ID: [author_id=%d] %w", authorID, err)
	}

	publishedDate, err := time.Parse(JSONDateFormat, req.Msg.PublishedDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse published date: [published_date=%s] %w", req.Msg.PublishedDate, err)
	}

	bookID := bs.GetNewRedisID()

	book, err := NewBook(bookID, req.Msg.Title, authorID, publishedDate)
	if err != nil {
		return nil, fmt.Errorf("failed to create book: %w", err)
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(book); err != nil {
		return nil, fmt.Errorf("failed to encode book: %w", err)
	}

	if err := bs.rdb.Set(ctx, newBookID(book.ID), buf.Bytes(), 0).Err(); err != nil {
		return nil, fmt.Errorf("failed to set book: [id=%d] %w", book.ID, err)
	}

	return &connect.Response[v1.UpdateBookResponse]{
		Msg: &v1.UpdateBookResponse{
			Book: &v1.Book{
				Id:            strconv.FormatInt(book.ID, 10),
				Title:         book.Title,
				AuthorId:      strconv.FormatInt(book.AuthorID, 10),
				PublishedDate: book.PublishedDate.Format(time.RFC3339),
			},
		},
	}, nil
}

func (bs *BooksService) DeleteBook(ctx context.Context, req *connect.Request[v1.DeleteBookRequest]) (*connect.Response[v1.DeleteBookResponse], error) {
	id, err := strconv.ParseInt(req.Msg.Id, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ID: %w", err)
	}

	if err := bs.rdb.Del(ctx, newBookID(id)).Err(); err != nil {
		return nil, fmt.Errorf("failed to delete book: [id=%s] %w", req.Msg.Id, err)
	}

	return &connect.Response[v1.DeleteBookResponse]{
		Msg: &v1.DeleteBookResponse{
			Status: true,
		},
	}, nil
}

func (bs *BooksService) GetNewRedisID() int64 {
	return bs.rdb.Incr(context.Background(), booksIDCounterKey).Val()
}

func newBookID(id int64) string {
	return booksKey + ":" + strconv.FormatInt(id, 10)
}
