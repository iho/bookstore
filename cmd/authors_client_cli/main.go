package main

import (
	"context"
	"log"
	"net/http"

	"connectrpc.com/connect"
	v1 "github.com/iho/bookstore/protos/gen/authors/v1"
	"github.com/iho/bookstore/protos/gen/authors/v1/authorsv1connect"
)

func list() {
	client := authorsv1connect.NewAuthorsServiceClient(
		http.DefaultClient,
		"http://localhost:8080/",
	)
	req := connect.NewRequest(&v1.ListAuthorsRequest{
		Limit:  100,
		Offset: 0,
	})
	res, err := client.ListAuthors(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Msg)
}

func get(id string) {
	client := authorsv1connect.NewAuthorsServiceClient(
		http.DefaultClient,
		"http://localhost:8080/",
	)
	req := connect.NewRequest(&v1.GetAuthorRequest{
		Id: id,
	})
	res, err := client.GetAuthor(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Msg)
}

func create() string {
	client := authorsv1connect.NewAuthorsServiceClient(
		http.DefaultClient,
		"http://localhost:8080/",
	)
	req := connect.NewRequest(&v1.CreateAuthorRequest{
		Name: "New Author",
	})
	res, err := client.CreateAuthor(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Msg)
	return res.Msg.GetAuthor().Id
}

func delete(id string) {
	client := authorsv1connect.NewAuthorsServiceClient(
		http.DefaultClient,
		"http://localhost:8080/",
	)
	req := connect.NewRequest(&v1.DeleteAuthorRequest{
		Id: id,
	})
	res, err := client.DeleteAuthor(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Msg)
}

func update(id string) {
	client := authorsv1connect.NewAuthorsServiceClient(
		http.DefaultClient,
		"http://localhost:8080/",
	)
	req := connect.NewRequest(&v1.UpdateAuthorRequest{
		Id:   id,
		Name: "Updated Author",
	})
	res, err := client.UpdateAuthor(context.Background(), req)
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
