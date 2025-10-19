package monitor

import (
	"context"
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"
	"time"

	"github.com/google/uuid"
)

func StartSendingPulses(ctx context.Context, healthID uuid.UUID) {
	go startSendingPulses(ctx, healthID)
}

func startSendingPulses(ctx context.Context, healthID uuid.UUID) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			updatePulse(healthID)
		}
	}
}

func updatePulse(healthID uuid.UUID) {
	query := `UPDATE x_health SET pulsed_at = current_timestamp WHERE id = :id`
	dbs.Exec(query, model.Param{"id": healthID})
}
