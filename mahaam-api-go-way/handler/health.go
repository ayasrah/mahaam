package handler

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/configs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHealth(router *gin.RouterGroup) {
	healthRouter := router.Group("/health")
	healthRouter.GET("", GetHealthInfo)
}

func GetHealthInfo(c *gin.Context) {
	c.JSON(http.StatusOK, model.Health{
		ID:         configs.HealthID,
		ApiName:    configs.ApiName,
		ApiVersion: configs.ApiVersion,
		EnvName:    configs.EnvName,
		NodeIP:     configs.NodeIP,
		NodeName:   configs.NodeName,
	})
}
