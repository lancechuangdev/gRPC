package main

import (
	"context"
	"grpc-demo/grpc-demo/proto/orderpb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := orderpb.NewOrderServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := client.CreateOrder(ctx, &orderpb.CreateOrderRequest{
		UserId: "user-1",
		Symbol: "BTCUSD",
		Price:  50000,
		Amount: 0.1,
	})
	if err != nil {
		panic(err)
	}

	println("Order created with ID:", resp.OrderId)
}
