package loaders

import (
	"context"

	"connectrpc.com/connect"
	"github.com/iho/bookstore/internal/gateway/graph/model"
	ordersV1 "github.com/iho/bookstore/protos/gen/orders/v1"
	"github.com/iho/bookstore/protos/gen/orders/v1/ordersv1connect"
)

type orderLoader struct {
	ordersv1connect ordersv1connect.OrdersServiceClient
}

func (l *orderLoader) getOrders(ctx context.Context, keys []string) ([]*model.Order, []error) {
	orders := make([]*model.Order, len(keys))
	errors := make([]error, len(keys))
	for i, key := range keys {
		req := connect.NewRequest(&ordersV1.GetOrderRequest{Id: key})
		res, err := l.ordersv1connect.GetOrder(ctx, req)
		if err != nil {
			errors[i] = err
			continue
		}

		orderLines := make([]*model.OrderLine, len(res.Msg.Order.OrderLines))

		for i, line := range res.Msg.Order.OrderLines {
			orderLines[i] = &model.OrderLine{
				BookID:   line.BookId,
				Quantity: int(line.Quantity),
			}
		}

		orders[i] = &model.Order{
			ID:         res.Msg.Order.Id,
			Quantity:   len(res.Msg.Order.OrderLines),
			OrderLines: orderLines,
			TotalPrice: int(res.Msg.Order.TotalPrice),
			OrderDate:  res.Msg.Order.OrderDate,
		}
	}
	return orders, errors
}

// GetOrder returns single author by id efficiently
func GetOrder(ctx context.Context, authorID string) (*model.Order, error) {
	loaders := For(ctx)
	return loaders.OrderLoader.Load(ctx, authorID)
}

// GetOrders returns many authors by ids efficiently
func GetOrders(ctx context.Context, authorIDs []string) ([]*model.Order, error) {
	loaders := For(ctx)
	return loaders.OrderLoader.LoadAll(ctx, authorIDs)
}
