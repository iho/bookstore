package main

import (
	"context"
	"log"
	"net/http"

	"connectrpc.com/connect"
	v1 "github.com/iho/bookstore/protos/gen/books/v1"
	"github.com/iho/bookstore/protos/gen/books/v1/booksv1connect"
)

const (
	serverAddr = "http://localhost:9090/"
)

func list() {
	client := booksv1connect.NewBooksServiceClient(
		http.DefaultClient,
		serverAddr,
	)
	ids := []string{"1", "2", "3"}
	req := connect.NewRequest(&v1.ListBooksRequest{
		Ids: ids,
	})
	res, err := client.ListBooks(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Msg)
}

func get(id string) {
	client := booksv1connect.NewBooksServiceClient(
		http.DefaultClient,
		serverAddr,
	)
	req := connect.NewRequest(&v1.GetBookRequest{
		Id: id,
	})
	res, err := client.GetBook(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Msg)
}

func create() string {
	client := booksv1connect.NewBooksServiceClient(
		http.DefaultClient,
		serverAddr,
	)
	req := connect.NewRequest(&v1.CreateBookRequest{
		Title:         "New Book",
		AuthorId:      "1",
		PublishedDate: "2024-06-12T18:37:04.189Z",
	})
	res, err := client.CreateBook(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Msg)
	return res.Msg.GetBook().Id
}

func delete(id string) {
	client := booksv1connect.NewBooksServiceClient(
		http.DefaultClient,
		serverAddr,
	)
	req := connect.NewRequest(&v1.DeleteBookRequest{
		Id: id,
	})
	res, err := client.DeleteBook(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Msg)
}

func update(id string) {
	client := booksv1connect.NewBooksServiceClient(
		http.DefaultClient,
		serverAddr,
	)
	req := connect.NewRequest(&v1.UpdateBookRequest{
		Id:            id,
		Title:         "Updated Book",
		AuthorId:      "1",
		PublishedDate: "2024-06-12T18:37:04.189Z",
	})
	res, err := client.UpdateBook(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Msg)
}

func main() {
	id := create()
	list()
	get(id)
	update(id)
	list()
	get(id)
	delete(id)
}
