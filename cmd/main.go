package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"ttages/pkg/database"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"ttages/internal/config"
	appGRPC "ttages/internal/file/handlers/gRPC"
	"ttages/internal/file/repository/postgres"
	appRedis "ttages/internal/file/repository/redis"
	"ttages/internal/file/usecase"
	"ttages/proto/pb"
)

func main() {
	cfg := config.MustLoad()

	// Postgres
	pg := database.NewDatabase(cfg.Postgres)
	defer pg.Close()

	// Redis
	log.Printf("Connecting to Redis at: %s", cfg.Redis.Addr)

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Fatal("Redis connection failed:", err)
	}
	defer rdb.Close()

	// Dependencies
	fileRepo := postgres.NewFileRepository(pg)
	fileCache := appRedis.NewFileCache(rdb, 120*time.Minute)
	fileUC := usecase.NewFileUsecase(fileRepo, fileCache, cfg)

	// gRPC Server
	grpcServer := grpc.NewServer()
	pb.RegisterFileServiceServer(grpcServer, appGRPC.NewFileHandler(fileUC))

	// Запуск сервера в отдельной горутине
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatal("Failed to listen:", err)
	}

	go func() {
		log.Printf("Starting gRPC server on port %s", cfg.GRPCPort)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatal("Failed to serve:", err)
		}
	}()

	// Ожидание сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Остановка gRPC сервера
	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		log.Println("Server stopped gracefully")
	case <-ctx.Done():
		log.Println("Forcing shutdown after timeout")
		grpcServer.Stop()
	}

	// Закрытие соединений
	if err := pg.Close(); err != nil {
		log.Println("Error closing PostgreSQL connection:", err)
	}
	if err := rdb.Close(); err != nil {
		log.Println("Error closing Redis connection:", err)
	}

	log.Println("Server exited properly")
}
