package rest

import (
	"encoding/json"
	"errors"
	"grpc-demo/internal/service"
	"net/http"
)

type OrderHandler struct {
	orderService *service.OrderService
}

type createOrderRequest struct {
	UserID string  `json:"user_id"`
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
	Amount float64 `json:"amount"`
}

type createOrderResponse struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /orders", h.CreateOrder)

	return mux
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req createOrderRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}

	order, err := h.orderService.CreateOrder(r.Context(), service.CreateOrderInput{
		UserID: req.UserID,
		Symbol: req.Symbol,
		Price:  req.Price,
		Amount: req.Amount,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusCreated, createOrderResponse{
		OrderID: order.ID,
		Status:  order.Status,
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(value)
}
