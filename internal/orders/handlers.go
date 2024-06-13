package orders

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	v1 "github.com/iho/bookstore/protos/gen/orders/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// booksV1 "github.com/iho/bookstore/protos/gen/books/v1"
	// authorsV1 "github.com/iho/bookstore/protos/gen/authors/v1"
)

const (
	bookStoreKey       = "bookstore"
	orderCollectionKey = "orders"
)

type OrdersService struct {
	client *mongo.Client
}

func NewOrdersService(client *mongo.Client) *OrdersService {
	return &OrdersService{client: client}
}

func (os *OrdersService) ListOrders(ctx context.Context, req *connect.Request[v1.ListOrdersRequest]) (*connect.Response[v1.ListOrdersResponse], error) {
	var limit int64 = 10
	var offset int64 = 0
	if req.Msg.Limit < 1 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("limit must be greater than 0"))
	}

	if req.Msg.Offset < 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("offset must be greater than or equal to 0"))
	}
	filter := bson.D{}
	if req.Msg.BookId != "" {
		filter = bson.D{{"order_lines.book_id", req.Msg.BookId}}
	}

	limit = int64(req.Msg.Limit)
	offset = int64(req.Msg.Offset)

	orders := make([]*v1.Order, 0, req.Msg.Limit)

	cursor, err := os.client.Database(bookStoreKey).Collection(orderCollectionKey).Find(ctx, filter, &options.FindOptions{
		Limit: &limit,
		Skip:  &offset,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		order := new(Order)
		if err := cursor.Decode(order); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}

		orderLines := make([]*v1.OrderLine, 0, len(order.OrderLines))
		for _, line := range order.OrderLines {
			orderLines = append(orderLines, &v1.OrderLine{
				BookId:   line.BookId,
				Quantity: line.Quantity,
			})
		}

		orders = append(orders, &v1.Order{
			Id:         order.ID.Hex(),
			OrderLines: orderLines,
			TotalPrice: order.TotalPrice,
			OrderDate:  order.OrderDate,
		})
	}

	return &connect.Response[v1.ListOrdersResponse]{
		Msg: &v1.ListOrdersResponse{
			Orders: orders,
		},
	}, nil
}

func (os *OrdersService) GetOrder(ctx context.Context, req *connect.Request[v1.GetOrderRequest]) (*connect.Response[v1.GetOrderResponse], error) {
	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("order ID must be provided"))
	}

	order := new(Order)
	id, err := primitive.ObjectIDFromHex(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	filter := bson.D{{"_id", id}}
	if err := os.client.Database(bookStoreKey).Collection(orderCollectionKey).FindOne(ctx, filter).Decode(order); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	orderLines := make([]*v1.OrderLine, 0, len(order.OrderLines))
	for _, line := range order.OrderLines {
		orderLines = append(orderLines, &v1.OrderLine{
			BookId:   line.BookId,
			Quantity: line.Quantity,
		})
	}

	return &connect.Response[v1.GetOrderResponse]{
		Msg: &v1.GetOrderResponse{
			Order: &v1.Order{
				Id:         order.ID.Hex(),
				OrderLines: orderLines,
				TotalPrice: order.TotalPrice,
				OrderDate:  order.OrderDate,
			},
		},
	}, nil
}

func (os *OrdersService) CreateOrder(ctx context.Context, req *connect.Request[v1.CreateOrderRequest]) (*connect.Response[v1.CreateOrderResponse], error) {
	if req.Msg.TotalPrice < 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("total price must be greater than 0"))
	}

	if req.Msg.OrderDate == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("order date must be provided"))
	}

	if len(req.Msg.GetOrderLines()) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("order must have at least one line"))
	}

	orderLines := make([]*OrderLine, 0, len(req.Msg.GetOrderLines()))
	for _, line := range req.Msg.GetOrderLines() {
		if line.Quantity < 1 {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("quantity must be greater than 0"))
		}

		if line.BookId == "" {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("book ID must be provided"))
		}

		orderLines = append(orderLines, NewOrderLine(line.GetBookId(), line.GetQuantity()))
	}

	order := NewOrder(orderLines, req.Msg.GetTotalPrice(), req.Msg.GetOrderDate())
	res, err := os.client.Database(bookStoreKey).Collection(orderCollectionKey).InsertOne(ctx, order)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &connect.Response[v1.CreateOrderResponse]{
		Msg: &v1.CreateOrderResponse{
			Order: &v1.Order{
				Id:         res.InsertedID.(primitive.ObjectID).Hex(),
				OrderLines: req.Msg.GetOrderLines(),
				TotalPrice: req.Msg.GetTotalPrice(),
				OrderDate:  req.Msg.GetOrderDate(),
			},
		},
	}, nil
}

func (os *OrdersService) UpdateOrder(ctx context.Context, req *connect.Request[v1.UpdateOrderRequest]) (*connect.Response[v1.UpdateOrderResponse], error) {
	if req.Msg.TotalPrice < 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("total price must be greater than 0"))
	}

	if req.Msg.OrderDate == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("order date must be provided"))
	}

	if len(req.Msg.GetOrderLines()) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("order must have at least one line"))
	}
	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("order ID must be provided"))
	}
	id, err := primitive.ObjectIDFromHex(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	orderLines := make([]*OrderLine, 0, len(req.Msg.GetOrderLines()))
	for _, line := range req.Msg.GetOrderLines() {
		if line.Quantity < 1 {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("quantity must be greater than 0"))
		}

		if line.BookId == "" {
			return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("book ID must be provided"))
		}

		orderLines = append(orderLines, NewOrderLine(line.GetBookId(), line.GetQuantity()))
	}

	order := NewOrder(orderLines, req.Msg.GetTotalPrice(), req.Msg.GetOrderDate())
	order.ID = primitive.NilObjectID
	update := bson.M{
		"$set": order,
	}
	res, err := os.client.Database(bookStoreKey).Collection(orderCollectionKey).UpdateByID(ctx, id, update)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if res.ModifiedCount == 0 {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("order not found"))
	}

	return &connect.Response[v1.UpdateOrderResponse]{
		Msg: &v1.UpdateOrderResponse{
			Order: &v1.Order{
				Id:         id.Hex(),
				OrderLines: req.Msg.GetOrderLines(),
				TotalPrice: req.Msg.GetTotalPrice(),
				OrderDate:  req.Msg.GetOrderDate(),
			},
		},
	}, nil
}

func (os *OrdersService) DeleteOrder(ctx context.Context, req *connect.Request[v1.DeleteOrderRequest]) (*connect.Response[v1.DeleteOrderResponse], error) {
	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("order ID must be provided"))
	}
	id, err := primitive.ObjectIDFromHex(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	filter := bson.D{{"_id", id}}
	res, err := os.client.Database(bookStoreKey).Collection(orderCollectionKey).DeleteOne(ctx, filter)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &connect.Response[v1.DeleteOrderResponse]{
		Msg: &v1.DeleteOrderResponse{
			Status: res.DeletedCount > 0,
		},
	}, nil
}
