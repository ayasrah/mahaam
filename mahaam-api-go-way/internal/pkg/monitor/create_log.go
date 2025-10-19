package monitor

import (
	"log"
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
)

func CreateLog(trafficID uuid.UUID, logType, message string, nodeIP string) {
	query := `
		INSERT INTO x_log (traffic_id, type, message, node_ip, created_at)
		VALUES (:traffic_id, :type, :message, :node_ip, current_timestamp)`

	params := model.Param{
		"traffic_id": trafficID,
		"type":       logType,
		"message":    message,
		"node_ip":    nodeIP,
	}

	// Run asynchronously to avoid blocking
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Error creating log:", r)
			}
		}()
		dbs.Exec(query, params)
	}()
}
