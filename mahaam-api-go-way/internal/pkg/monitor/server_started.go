package monitor

import (
	"errors"
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"
)

func ServerStarted(health *model.Health) *model.Err {
	if err := createHealth(health); err != nil {
		return model.ServerError("failed to create health record: " + err.Error())
	}
	return nil
}

func createHealth(health *model.Health) error {
	query := `
		INSERT INTO x_health (id, api_name, api_version, env_name, node_ip, node_name, started_at)
		VALUES (:id, :api_name, :api_version, :env_name, :node_ip, :node_name, current_timestamp)`
	rows, err := dbs.Exec(query, health)
	if err != nil {
		return err
	}
	if rows != 1 {
		return errors.New("health record not created")
	}
	return nil
}
