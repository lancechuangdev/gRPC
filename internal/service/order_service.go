package service

import (
	"context"
	"errors"
	"grpc-demo/internal/repository"
)

type CreateOrderInput struct {
	UserID string
	Symbol string
	Price  float64
	Amount float64
}

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, input CreateOrderInput) (repository.Order, error) {
	if input.UserID == "" || input.Symbol == "" || input.Price <= 0 || input.Amount <= 0 {
		return repository.Order{}, errors.New("invalid input")
	}

	order := repository.Order{
		UserID: input.UserID,
		Symbol: input.Symbol,
		Price:  input.Price,
		Amount: input.Amount,
	}

	return s.repo.Create(ctx, order)
}
