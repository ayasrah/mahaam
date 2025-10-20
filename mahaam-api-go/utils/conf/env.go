package conf

import (
	"mahaam-api/app/models"

	"github.com/google/uuid"
)

var env *Environment

func Env() *Environment {
	return env
}

func NewEnvironment(h *models.Health) {
	env = &Environment{
		NodeIP:   h.NodeIP,
		NodeName: h.NodeName,
		HealthID: h.ID,
	}
}

type Environment struct {
	NodeIP   string
	NodeName string
	HealthID uuid.UUID
}
