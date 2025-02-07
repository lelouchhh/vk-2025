package delivery

import (
	"context"
	"log"
	"time"

	"pinger/internal/repository"
	"pinger/pinger"
)

func RunCLI(ctx context.Context, dockerRepo repository.DockerRepository, rabbitRepo repository.RabbitMQRepository, interval time.Duration) {
	// Инициализация сервиса
	pingerService := pinger.NewPingerService(dockerRepo, rabbitRepo)

	// Запуск сервиса
	pingerService.Start(ctx, interval)

	log.Printf("Pinger service started with interval %v", interval)
}
