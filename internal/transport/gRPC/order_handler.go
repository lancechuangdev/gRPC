package grpc

import (
	"context"
	"errors"
	"grpc-demo/grpc-demo/proto/orderpb"
	"grpc-demo/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		if errors.Is(err, service.ErrInvalidInput) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}

		return nil, err
	}

	return &orderpb.CreateOrderResponse{
		OrderId: order.ID,
		Status:  order.Status,
	}, nil
}
