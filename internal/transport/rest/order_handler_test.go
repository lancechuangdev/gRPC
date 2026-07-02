package rest

import (
	"encoding/json"
	"grpc-demo/internal/repository"
	"grpc-demo/internal/service"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateOrder(t *testing.T) {
	handler := newTestHandler()
	req := httptest.NewRequest(
		http.MethodPost,
		"/orders",
		strings.NewReader(`{"user_id":"user-1","symbol":"BTCUSD","price":50000,"amount":0.1}`),
	)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, recorder.Code)
	}

	var response createOrderResponse
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if response.OrderID == "" || response.Status != "created" {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestCreateOrderRejectsInvalidInput(t *testing.T) {
	handler := newTestHandler()
	req := httptest.NewRequest(
		http.MethodPost,
		"/orders",
		strings.NewReader(`{"user_id":"user-1","symbol":"BTCUSD","price":0,"amount":0.1}`),
	)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}
}

func newTestHandler() http.Handler {
	repo := repository.NewOrderRepository()
	orderService := service.NewOrderService(repo)

	return NewOrderHandler(orderService).Routes()
}
