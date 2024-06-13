package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/iho/bookstore/internal/orders"
	"github.com/iho/bookstore/protos/gen/orders/v1/ordersv1connect"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	var uri string
	if uri = os.Getenv("MONGODB_URI"); uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable. See\n\t https://www.mongodb.com/docs/drivers/go/current/usage-examples/#environment-variable")
	}
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	ordersService := orders.NewOrdersService(client)
	mux := http.NewServeMux()
	mux.Handle(ordersv1connect.NewOrdersServiceHandler(ordersService))
	fmt.Println("Starting server on :9999")

	if err := http.ListenAndServe(":9999", h2c.NewHandler(mux, &http2.Server{})); err != nil {
		log.Fatal(err)
	}
}
