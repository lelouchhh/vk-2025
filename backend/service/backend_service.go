package service

import (
	"context"
	"log"
	"sync"

	"backend/domain"
	"backend/internal/repository"
)

type BackendService struct {
	rabbitRepo repository.RabbitMQRepository
	dbRepo     repository.PostgresRepository
}

func NewBackendService(rabbitRepo repository.RabbitMQRepository, dbRepo repository.PostgresRepository) *BackendService {
	return &BackendService{
		rabbitRepo: rabbitRepo,
		dbRepo:     dbRepo,
	}
}

func (s *BackendService) StartConsuming(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	resultsChan, err := s.rabbitRepo.ConsumePingResults(ctx)
	if err != nil {
		log.Printf("Failed to start consuming ping results: %v", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping result consumption...")
			return
		case result, ok := <-resultsChan:
			if !ok {
				log.Println("Results channel closed")
				return
			}
			result.Status = true // Assuming success if we received the result
			if err := s.dbRepo.SavePingResult(ctx, result); err != nil {
				log.Printf("Failed to save ping result: %v", err)
			}
		}
	}
}

func (s *BackendService) GetAllContainers(ctx context.Context) ([]domain.Container, error) {
	return s.dbRepo.GetAllContainers(ctx)
}
