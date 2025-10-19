package monitor

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
)

func ServerStopped(healthID uuid.UUID) *model.Err {
	if err := updateStopped(healthID); err != nil {
		return model.ServerError("failed to update stopped timestamp: " + err.Error())
	}
	return nil
}

func updateStopped(healthID uuid.UUID) error {
	query := `UPDATE x_health SET stopped_at = current_timestamp WHERE id = :id`
	_, err := dbs.Exec(query, model.Param{"id": healthID})
	return err
}
