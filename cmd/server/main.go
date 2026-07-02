package main

import (
	"context"
	"fmt"
	"grpc-demo/grpc-demo/proto/orderpb"
	"grpc-demo/internal/repository"
	"grpc-demo/internal/service"
	transportgrpc "grpc-demo/internal/transport/gRPC"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println("Failed to listen:", err)
		return
	}

	orderRepo := repository.NewOrderRepository()
	orderService := service.NewOrderService(orderRepo)
	orderHandler := transportgrpc.NewOrderHandler(orderService)

	gprcSrv := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(gprcSrv, orderHandler)

	go func() {
		fmt.Println("gRPC server listening on :50051")

		if err := gprcSrv.Serve(listener); err != nil {
			fmt.Println("Failed to serve:", err)
		}
	}()

	waitForShutdown(gprcSrv)
}

func waitForShutdown(grpcSrv *grpc.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	fmt.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	grpcStopped := make(chan struct{})

	go func() {
		fmt.Println("starting graceful shutdown")
		grpcSrv.GracefulStop()
		close(grpcStopped)
	}()

	select {
	case <-grpcStopped:
		fmt.Println("gRPC server gracefully stopped")

	case <-shutdownCtx.Done():
		fmt.Println("graceful shutdown timeout, forcing stop")
		grpcSrv.Stop()
	}
}
