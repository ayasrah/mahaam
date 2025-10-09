package repo

import (
	"log"
	"mahaam-api/infra/cache"
	"mahaam-api/infra/dbs"
)

type LogRepo interface {
	Create(trafficId UUID, logType, message string)
}

type logRepo struct {
}

func NewLogRepo() LogRepo {
	return &logRepo{}
}

func (r *logRepo) Create(trafficId UUID, logType, message string) {
	query := `
		INSERT INTO x_log (traffic_id, type, message, node_ip, created_at)
		VALUES (:traffic_id, :type, :message, :node_ip, current_timestamp)`

	params := Param{
		"traffic_id": trafficId,
		"type":       logType,
		"message":    message,
		"node_ip":    cache.NodeIP,
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		dbs.Exec(query, params)
	}()
}
