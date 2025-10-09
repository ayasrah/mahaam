package cache

import (
	"mahaam-api/feat/models"

	"github.com/google/uuid"
)

var (
	NodeIP     string
	NodeName   string
	ApiName    string
	ApiVersion string
	EnvName    string
	HealthID   uuid.UUID
)

func Init(h *models.Health) {
	NodeIP = h.NodeIP
	NodeName = h.NodeName
	ApiName = h.ApiName
	ApiVersion = h.ApiVersion
	EnvName = h.EnvName
	HealthID = h.ID
}
