package main

import (
	"context"
	"errors"
	"fmt"
	"grpc-demo/grpc-demo/proto/orderpb"
	"grpc-demo/internal/repository"
	"grpc-demo/internal/service"
	transportgrpc "grpc-demo/internal/transport/gRPC"
	"grpc-demo/internal/transport/rest"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
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
	grpcHandler := transportgrpc.NewOrderHandler(orderService)
	restHandler := rest.NewOrderHandler(orderService)

	grpcSrv := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(grpcSrv, grpcHandler)

	restSrv := &http.Server{
		Addr:    ":8080",
		Handler: restHandler.Routes(),
	}

	go func() {
		fmt.Println("gRPC server listening on :50051")

		if err := grpcSrv.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			fmt.Println("Failed to serve:", err)
		}
	}()

	go func() {
		fmt.Println("REST server listening on :8080")

		if err := restSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Failed to serve REST API:", err)
		}
	}()

	waitForShutdown(grpcSrv, restSrv)
}

func waitForShutdown(grpcSrv *grpc.Server, restSrv *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	fmt.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serversStopped := make(chan struct{})

	go func() {
		fmt.Println("starting graceful shutdown")

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			grpcSrv.GracefulStop()
		}()

		go func() {
			defer wg.Done()
			_ = restSrv.Shutdown(shutdownCtx)
		}()

		wg.Wait()
		close(serversStopped)
	}()

	select {
	case <-serversStopped:
		fmt.Println("gRPC and REST servers gracefully stopped")

	case <-shutdownCtx.Done():
		fmt.Println("graceful shutdown timeout, forcing stop")
		grpcSrv.Stop()
		_ = restSrv.Close()
	}
}
