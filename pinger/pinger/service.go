package pinger

import (
	"context"
	"log"
	"os/exec"
	"pinger/domain"
	"pinger/internal/repository"
	"time"
)

type PingerService struct {
	dockerRepo repository.DockerRepository
	rabbitRepo repository.RabbitMQRepository
}

func NewPingerService(dockerRepo repository.DockerRepository, rabbitRepo repository.RabbitMQRepository) *PingerService {
	return &PingerService{
		dockerRepo: dockerRepo,
		rabbitRepo: rabbitRepo,
	}
}

func (s *PingerService) GetContainers(ctx context.Context) ([]domain.Container, error) {
	return s.dockerRepo.GetContainers()
}

func (s *PingerService) PingContainer(ctx context.Context, ip string) (domain.PingResult, error) {
	start := time.Now()
	cmd := exec.CommandContext(ctx, "ping", "-c", "1", "-W", "1", ip)
	err := cmd.Run()
	pingTime := time.Since(start).Seconds()

	result := domain.PingResult{
		IP:       ip,
		PingTime: pingTime,
	}

	if err == nil {
		result.LastSuccess = time.Now()
	}

	return result, nil
}

func (s *PingerService) StorePingResult(ctx context.Context, result domain.PingResult) error {
	return s.rabbitRepo.PublishPingResult(result)
}

func (s *PingerService) Run(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping pinger service...")
			return
		case <-ticker.C:
			containers, err := s.GetContainers(ctx)
			if err != nil {
				log.Printf("Error fetching containers: %v", err)
				continue
			}

			for _, container := range containers {
				result, err := s.PingContainer(ctx, container.IP)
				if err != nil {
					log.Printf("Error pinging container %s: %v", container.IP, err)
					continue
				}

				err = s.StorePingResult(ctx, result)
				if err != nil {
					log.Printf("Error storing ping result for container %s: %v", container.IP, err)
				}
			}
		}
	}
}

func (s *PingerService) Start(ctx context.Context, interval time.Duration) {
	go s.Run(ctx, interval)
}
