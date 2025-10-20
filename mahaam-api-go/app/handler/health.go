package handler

import (
	"mahaam-api/app/models"
	"mahaam-api/utils/conf"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler interface {
	GetInfo(c *gin.Context)
}

type healthHandler struct {
	cfg *conf.Conf
}

func NewHealthHandler(cfg *conf.Conf) HealthHandler {
	return &healthHandler{cfg: cfg}
}

func RegisterHealthHandler(router *gin.RouterGroup, h HealthHandler) {
	healthRouter := router.Group("/health")
	healthRouter.GET("", h.GetInfo)
}

func (h *healthHandler) GetInfo(c *gin.Context) {
	c.JSON(http.StatusOK, models.Health{
		ID:         conf.Env().HealthID,
		NodeIP:     conf.Env().NodeIP,
		NodeName:   conf.Env().NodeName,
		ApiName:    h.cfg.ApiName,
		ApiVersion: h.cfg.ApiVersion,
		EnvName:    h.cfg.EnvName,
	})
}
