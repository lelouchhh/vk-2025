package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"backend/internal/delivery"
	"backend/internal/repository"
	"backend/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Чтение переменных окружения из docker-compose
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatal("One or more required environment variables are not set")
	}

	// Формирование строки подключения к PostgreSQL
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbName)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	// Инициализация репозиториев
	dbRepo := repository.NewPostgresRepository(db)

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		log.Fatal("RABBITMQ_URL environment variable is not set")
	}

	rabbitRepo, err := repository.NewRabbitMQRepository(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ repository: %v", err)
	}
	defer rabbitRepo.Close()

	accountRepo := repository.NewAccountRepository(db)

	// Инициализация сервисов
	authService := service.NewAuthService(accountRepo, os.Getenv("mysecretkey"))
	backendService := service.NewBackendService(rabbitRepo, dbRepo)

	// Запуск потребителя RabbitMQ
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go backendService.StartConsuming(ctx, &wg)

	// Инициализация HTTP-сервера
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3001"}, // Разрешаем запросы только с фронтенда
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}))

	handler := delivery.NewHTTPHandler(authService, backendService)

	// Регистрация маршрутов
	e.POST("/register", handler.Register)
	e.POST("/login", handler.Login)

	protected := e.Group("/protected")
	protected.Use(handler.AuthMiddleware)
	protected.GET("/containers", handler.GetContainers)

	// Запуск HTTP-сервера в отдельной горутине
	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Обработка сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("Received signal: %v. Shutting down gracefully...", sig)

	// Отмена контекста для остановки всех фоновых процессов
	cancel()
	wg.Wait()

	// Остановка HTTP-сервера
	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()
	if err := e.Shutdown(ctxShutdown); err != nil {
		log.Printf("HTTP server shutdown failed: %v", err)
	}

	log.Println("Shutdown complete.")
}
