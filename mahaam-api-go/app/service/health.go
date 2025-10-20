package service

import (
	"context"
	"fmt"
	"time"

	"mahaam-api/app/models"
	"mahaam-api/app/repo"
	"mahaam-api/utils/conf"
	logs "mahaam-api/utils/log"

	"github.com/google/uuid"
)

type HealthService interface {
	ServerStarted(health *models.Health)
	StartSendingPulses(ctx context.Context)
	ServerStopped()
}

type healthService struct {
	healthRepo repo.HealthRepo
	cfg        *conf.Conf
	logger     logs.Logger
}

func NewHealthService(healthRepo repo.HealthRepo, cfg *conf.Conf, logger logs.Logger) HealthService {
	return &healthService{healthRepo: healthRepo, cfg: cfg, logger: logger}
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
			s.healthRepo.UpdatePulse(conf.Env().HealthID)
		}
	}
}

func (s *healthService) ServerStopped() {
	s.healthRepo.UpdateStopped(conf.Env().HealthID)
	stopMsg := fmt.Sprintf("✓ %s-v%s/%s-%s stopped with healthID=%s", s.cfg.ApiName, s.cfg.ApiVersion, conf.Env().NodeIP, conf.Env().NodeName, conf.Env().HealthID.String())
	s.logger.Info(uuid.Nil, stopMsg)
}
