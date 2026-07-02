package repository

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Order struct {
	ID     string
	UserID string
	Symbol string
	Price  float64
	Amount float64
	Status string
}

type OrderRepository struct {
	mu     sync.Mutex
	orders map[string]Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders: make(map[string]Order),
	}
}

func (r *OrderRepository) Create(ctx context.Context, order Order) (Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order.ID = fmt.Sprintf("order-%d", time.Now().UnixNano())
	order.Status = "created"

	r.orders[order.ID] = order

	return order, nil
}

func (r *OrderRepository) GetByID(ctx context.Context, orderId string) (Order, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, exists := r.orders[orderId]

	return order, exists
}
