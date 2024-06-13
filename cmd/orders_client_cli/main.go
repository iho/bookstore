package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	v1 "github.com/iho/bookstore/protos/gen/orders/v1"
	"github.com/iho/bookstore/protos/gen/orders/v1/ordersv1connect"
)

const (
	serverAddr = "http://localhost:9999/"
)

func list() {
	client := ordersv1connect.NewOrdersServiceClient(
		http.DefaultClient,
		serverAddr,
	)

	req := connect.NewRequest(&v1.ListOrdersRequest{
		Limit:  1000,
		Offset: 0,
	})
	res, err := client.ListOrders(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Msg)
}

func get(id string) {
	client := ordersv1connect.NewOrdersServiceClient(
		http.DefaultClient,
		serverAddr,
	)
	req := connect.NewRequest(&v1.GetOrderRequest{
		Id: id,
	})
	res, err := client.GetOrder(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Msg)
}

func create() string {
	client := ordersv1connect.NewOrdersServiceClient(
		http.DefaultClient,
		serverAddr,
	)
	req := connect.NewRequest(&v1.CreateOrderRequest{
		OrderLines: []*v1.OrderLine{
			{
				BookId:   "1",
				Quantity: 1,
			},
		},
		TotalPrice: 100,
		OrderDate:  "2024-06-12T18:37:04.189Z",
	})
	res, err := client.CreateOrder(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Msg)
	return res.Msg.GetOrder().Id
}

func delete(id string) {
	client := ordersv1connect.NewOrdersServiceClient(
		http.DefaultClient,
		serverAddr,
	)
	req := connect.NewRequest(&v1.DeleteOrderRequest{
		Id: id,
	})
	res, err := client.DeleteOrder(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Msg)
}

func update(id string) {
	client := ordersv1connect.NewOrdersServiceClient(
		http.DefaultClient,
		serverAddr,
	)
	req := connect.NewRequest(&v1.UpdateOrderRequest{
		Id:         id,
		OrderLines: []*v1.OrderLine{{BookId: "1", Quantity: 2}},
		TotalPrice: 200,
		OrderDate:  "2024-06-12T18:37:04.189Z",
	})
	res, err := client.UpdateOrder(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.Msg)
}

func main() {
	fmt.Println("Creating order")
	id := create()
	fmt.Println("Listing orders")
	list()
	fmt.Println("Getting order")
	get(id)
	fmt.Println("Updating order", id)
	update(id)
	fmt.Println("Listing orders")
	list()
	fmt.Println("Getting order")
	get(id)
	fmt.Println("Deleting order")
	delete(id)
}
