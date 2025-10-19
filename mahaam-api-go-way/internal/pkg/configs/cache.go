package configs

import (
	"mahaam-api/internal/model"

	"github.com/google/uuid"
)

var (
	NodeIP   string
	NodeName string
	HealthID uuid.UUID
)

func InitCache(h *model.Health) {
	NodeIP = h.NodeIP
	NodeName = h.NodeName
	HealthID = h.ID
}
