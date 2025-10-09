package repo

import (
	"mahaam-api/feat/models"
	"mahaam-api/infra/dbs"
)

type TrafficRepo interface {
	Create(traffic models.Traffic) int64
}

type trafficRepo struct {
}

func NewTrafficRepo() TrafficRepo {
	return &trafficRepo{}
}

func (r *trafficRepo) Create(traffic models.Traffic) int64 {
	query := `
		INSERT INTO x_traffic (id, health_id, method, path, code, elapsed, headers, request, response, created_at)
		VALUES (:id, :health_id, :method, :path, :code, :elapsed, :headers, :request, :response, current_timestamp)`
	return dbs.Exec(query, traffic)
}
