package grpc

import (
	"context"
	"log/slog"

	"github.com/destinyhover/microservices-proto/golang/order"
	"github.com/destinyhover/microservices/order/internal/application/core/domain"
)

func (a Adapter) Create(ctx context.Context, request *order.CreateOrderRequest) (*order.GetOrderResponse, error) {
	slog.InfoContext(ctx, "Creating order..")
	var orderItems []domain.OrderItem
	for _, v := range request.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: v.ProductCode,
			UnitPrice:   v.UnitPrice,
			Quantity:    v.Quantity,
		})
	}
	newOrder := domain.NewOrder(request.UserID, orderItems)
	result, err := a.api.PlaceOrder(ctx, newOrder)
	if err != nil {
		return nil, err
	}
	return &order.GetOrderResponse{UserID: result.ID}, nil
}

func (a Adapter) Get(ctx context.Context, request *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	result, err := a.api.GetOrder(ctx, request.OrderID)
	if err != nil {
		return nil, err
	}
	var orderItems []*order.OrderItem
	for _, v := range result.OrderItems {
		orderItems = append(orderItems, &order.OrderItem{
			ProductCode: v.ProductCode,
			UnitPrice:   v.UnitPrice,
			Quantity:    v.Quantity,
		})
	}
	return &order.GetOrderResponse{UserID: result.CustomerID, OrderItems: orderItems}, nil
}
