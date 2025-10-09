package repo

import (
	"errors"
	"mahaam-api/feat/models"
	"mahaam-api/infra/dbs"
)

type HealthRepo interface {
	Create(health *models.Health) int64
	UpdatePulse(id UUID) int64
	UpdateStopped(id UUID) int64
}

type healthRepo struct {
}

func NewHealthRepo() HealthRepo {
	return &healthRepo{}
}

func (r *healthRepo) Create(health *models.Health) int64 {
	query := `
		INSERT INTO x_health (id, api_name, api_version, env_name, node_ip, node_name, started_at)
		VALUES (:id, :api_name, :api_version, :env_name, :node_ip, :node_name, current_timestamp)`
	rows := dbs.Exec(query, health)
	if rows != 1 {
		panic(DBErr{Err: errors.New("Health record not created")})
	}
	return rows
}

func (r *healthRepo) UpdatePulse(id UUID) int64 {
	query := `UPDATE x_health SET pulsed_at = current_timestamp WHERE id = :id`
	return dbs.Exec(query, Param{"id": id})
}

func (r *healthRepo) UpdateStopped(id UUID) int64 {
	query := `UPDATE x_health SET stopped_at = current_timestamp WHERE id = :id`
	return dbs.Exec(query, Param{"id": id})
}
