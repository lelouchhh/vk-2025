package repository

import (
	"context"
	"errors"
	"github.com/docker/docker/api/types/container"
	"log"

	"github.com/docker/docker/client"
	"pinger/domain"
)

type DockerRepository interface {
	GetContainers() ([]domain.Container, error)
}

type dockerRepository struct {
	dockerClient *client.Client
}

func NewDockerRepository() (DockerRepository, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &dockerRepository{dockerClient: cli}, nil
}

func (r *dockerRepository) GetContainers() ([]domain.Container, error) {
	if r.dockerClient == nil {
		return nil, errors.New("Docker client is not initialized")
	}

	ctx := context.Background()

	// Получаем список контейнеров
	containers, err := r.dockerClient.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		log.Printf("Error listing containers: %v", err)
		return nil, err
	}

	var result []domain.Container
	for _, container := range containers {
		inspect, err := r.dockerClient.ContainerInspect(ctx, container.ID)
		if err != nil {
			log.Printf("Error inspecting container %s: %v", container.ID, err)
			continue
		}

		// Ищем первый доступный IP-адрес из всех сетей
		var ip string
		for networkName, network := range inspect.NetworkSettings.Networks {
			if network.IPAddress != "" {
				ip = network.IPAddress
				log.Printf("Container %s is connected to network %s with IP %s", container.ID, networkName, ip)
				break
			}
		}

		// Если IP-адрес не найден, пропускаем контейнер
		if ip == "" {
			log.Printf("Container %s has no IP address in any network", container.ID)
			continue
		}

		result = append(result, domain.Container{
			ID: container.ID,
			IP: ip,
		})
	}

	return result, nil
}
