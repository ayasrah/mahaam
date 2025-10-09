package service

import (
	"context"
	"fmt"
	"time"

	"mahaam-api/feat/models"
	"mahaam-api/feat/repo"
	"mahaam-api/infra/cache"
	logs "mahaam-api/infra/log"

	"github.com/google/uuid"
)

type HealthService interface {
	ServerStarted(health *models.Health)
	StartSendingPulses(ctx context.Context)
	ServerStopped()
}

type healthService struct {
	healthRepo repo.HealthRepo
}

func NewHealthService(healthRepo repo.HealthRepo) HealthService {
	return &healthService{healthRepo: healthRepo}
}

func (s *healthService) ServerStarted(health *models.Health) {
	s.healthRepo.Create(health)
}

func (s *healthService) StartSendingPulses(ctx context.Context) {
	go s.startSendingPulses(ctx)
}

func (s *healthService) startSendingPulses(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.healthRepo.UpdatePulse(cache.HealthID)
		}
	}
}

func (s *healthService) ServerStopped() {
	s.healthRepo.UpdateStopped(cache.HealthID)
	stopMsg := fmt.Sprintf("✓ %s-v%s/%s-%s stopped with healthID=%s", cache.ApiName, cache.ApiVersion, cache.NodeIP, cache.NodeName, cache.HealthID.String())
	logs.Info(uuid.Nil, stopMsg)
}
