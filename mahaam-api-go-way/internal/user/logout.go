package user

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func Logout(userID, deviceID uuid.UUID) *model.Err {
	device, err := getDevice(deviceID)
	if err != nil {
		return model.ServerError("failed to get device: " + err.Error())
	}

	if device.UserID != userID {
		return model.ForbiddenError("invalid deviceId")
	}

	err = dbs.WithTx(func(tx *sqlx.Tx) error {
		return deleteDevice(tx, deviceID)
	})

	if err != nil {
		return model.ServerError("failed to logout: " + err.Error())
	}

	return nil
}

func getDevice(deviceID uuid.UUID) (*model.Device, error) {
	query := `SELECT id, user_id, platform, fingerprint, info, created_at FROM devices WHERE id = :id`
	params := model.Param{"id": deviceID}
	device, err := dbs.SelectOne[model.Device](query, params)
	if err != nil {
		return nil, err
	}
	return &device, nil
}
