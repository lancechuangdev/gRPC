package grpc

import (
	"context"
	"grpc-demo/grpc-demo/proto/orderpb"
	"grpc-demo/internal/service"
)

type OrderHandler struct {
	orderpb.UnimplementedOrderServiceServer

	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	input := service.CreateOrderInput{
		UserID: req.UserId,
		Symbol: req.Symbol,
		Price:  req.Price,
		Amount: req.Amount,
	}

	order, err := h.orderService.CreateOrder(ctx, input)
	if err != nil {
		return nil, err
	}

	return &orderpb.CreateOrderResponse{
		OrderId: order.ID,
		Status:  order.Status,
	}, nil
}
