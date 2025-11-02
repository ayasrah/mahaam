package repo

import (
	"mahaam-api/app/models"

	"github.com/google/uuid"
)

type HealthRepo interface {
	Create(health *models.Health) int64
	UpdatePulse(id uuid.UUID) int64
	UpdateStopped(id uuid.UUID) int64
}

type healthRepo struct {
	db *AppDB
}

func NewHealthRepo(db *AppDB) HealthRepo {
	return &healthRepo{db: db}
}

func (r *healthRepo) Create(health *models.Health) int64 {
	query := `
		INSERT INTO monitor.health (id, api_name, api_version, env_name, node_ip, node_name, started_at)
		VALUES (:id, :api_name, :api_version, :env_name, :node_ip, :node_name, current_timestamp)`
	rows := execute(r.db, query, health)
	if rows != 1 {
		panic(models.ServerError("Health record not created"))
	}
	return rows
}

func (r *healthRepo) UpdatePulse(id uuid.UUID) int64 {
	query := `UPDATE monitor.health SET pulsed_at = current_timestamp WHERE id = :id`
	return execute(r.db, query, Param{"id": id})
}

func (r *healthRepo) UpdateStopped(id uuid.UUID) int64 {
	query := `UPDATE monitor.health SET stopped_at = current_timestamp WHERE id = :id`
	return execute(r.db, query, Param{"id": id})
}
