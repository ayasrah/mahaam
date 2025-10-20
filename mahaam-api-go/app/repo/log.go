package repo

import (
	"log"
	"mahaam-api/utils/conf"

	"github.com/google/uuid"
)

type LogRepo interface {
	Create(trafficId uuid.UUID, logType, message string)
}

type logRepo struct {
	db *AppDB
}

func NewLogRepo(db *AppDB) LogRepo {
	return &logRepo{db: db}
}

func (r *logRepo) Create(trafficId uuid.UUID, logType, message string) {
	query := `
		INSERT INTO x_log (traffic_id, type, message, node_ip, created_at)
		VALUES (:traffic_id, :type, :message, :node_ip, current_timestamp)`

	params := Param{
		"traffic_id": trafficId,
		"type":       logType,
		"message":    message,
		"node_ip":    conf.Env().NodeIP,
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		execute(r.db, query, params)
	}()
}
