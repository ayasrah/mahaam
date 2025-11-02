package repo

import (
	"mahaam-api/app/models"
)

type TrafficRepo interface {
	Create(traffic models.Traffic) int64
}

type trafficRepo struct {
	db *AppDB
}

func NewTrafficRepo(db *AppDB) TrafficRepo {
	return &trafficRepo{db: db}
}

func (r *trafficRepo) Create(traffic models.Traffic) int64 {
	query := `
		INSERT INTO monitor.traffic (id, health_id, method, path, code, elapsed, headers, request, response, created_at)
		VALUES (:id, :health_id, :method, :path, :code, :elapsed, :headers, :request, :response, current_timestamp)`
	return execute(r.db, query, traffic)
}
