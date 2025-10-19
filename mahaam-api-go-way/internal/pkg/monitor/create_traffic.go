package monitor

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"
)

func CreateTraffic(traffic model.Traffic) *model.Err {
	if err := createTrafficRecord(traffic); err != nil {
		return model.ServerError("failed to create traffic record: " + err.Error())
	}
	return nil
}

func createTrafficRecord(traffic model.Traffic) error {
	query := `
		INSERT INTO x_traffic (id, health_id, method, path, code, elapsed, headers, request, response, created_at)
		VALUES (:id, :health_id, :method, :path, :code, :elapsed, :headers, :request, :response, current_timestamp)`
	_, err := dbs.Exec(query, traffic)
	return err
}
