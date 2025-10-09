package handler

import (
	"mahaam-api/feat/models"
	"mahaam-api/infra/cache"
)

type HealthHandler interface {
	GetInfo(c Ctx)
}

type healthHandler struct {
}

func NewHealthHandler() HealthHandler {
	return &healthHandler{}
}

func RegisterHealthHandler(router Router, h HealthHandler) {
	healthRouter := router.Group("/health")
	healthRouter.GET("", h.GetInfo)
}

func (h *healthHandler) GetInfo(c Ctx) {
	c.JSON(OK, models.Health{
		ID:         cache.HealthID,
		ApiName:    cache.ApiName,
		ApiVersion: cache.ApiVersion,
		EnvName:    cache.EnvName,
		NodeIP:     cache.NodeIP,
		NodeName:   cache.NodeName,
	})
}
