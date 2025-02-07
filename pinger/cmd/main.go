package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"pinger/internal/delivery"
	"pinger/internal/repository"
)

func main() {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		log.Fatal("RABBITMQ_URL environment variable is not set")
	}
	seconds, err := strconv.Atoi(os.Getenv("PING_TIME"))
	if err != nil {
		log.Fatal("PING_TIME environment variable is not set")
	}
	interval := time.Duration(seconds) * time.Second

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализация репозиториев
	dockerRepo, err := repository.NewDockerRepository()
	if err != nil {
		log.Fatal(err)
	}
	rabbitRepo, err := repository.NewRabbitMQRepository(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ repository: %v", err)
	}
	defer func() {
		if err := rabbitRepo.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ connection: %v", err)
		}
	}()

	// Запуск CLI
	go func() {
		delivery.RunCLI(ctx, dockerRepo, rabbitRepo, interval)
	}()

	// Обработка сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Received signal: %v. Shutting down gracefully...", sig)

	cancel()
	log.Println("Waiting for cleanup...")
	time.Sleep(2 * time.Second)
	log.Println("Shutdown complete.")
}
