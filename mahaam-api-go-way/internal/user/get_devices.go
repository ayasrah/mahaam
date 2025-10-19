package user

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
)

func GetDevices(userID uuid.UUID) ([]model.Device, *model.Err) {
	devices, err := getDevices(userID)
	if err != nil {
		return nil, model.ServerError("failed to get devices: " + err.Error())
	}
	return devices, nil
}

func GetDevice(deviceID uuid.UUID) (*model.Device, error) {
	query := `SELECT id, user_id, platform, fingerprint, info, created_at FROM devices WHERE id = :id`
	params := model.Param{"id": deviceID}
	device, err := dbs.SelectOne[model.Device](query, params)
	if err != nil {
		return nil, err
	}
	return &device, nil
}
