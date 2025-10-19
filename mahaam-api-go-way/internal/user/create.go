package user

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func Create(device model.Device) (*model.CreatedUser, *model.Err) {
	var jwt string
	var userID uuid.UUID
	var deviceID uuid.UUID
	var err error

	txFn := func(tx *sqlx.Tx) error {
		userID, err = createUser(tx)
		if err != nil {
			return err
		}

		if err = deleteDeviceByFingerprint(tx, device.Fingerprint); err != nil {
			return err
		}

		device.UserID = userID
		deviceID, err = createDevice(tx, device)
		if err != nil {
			return err
		}

		return nil
	}

	if err = dbs.WithTx(txFn); err != nil {
		return nil, model.ServerError("failed to create user: " + err.Error())
	}

	jwt, err = CreateToken(userID, deviceID)
	if err != nil {
		return nil, model.ServerError("failed to create token: " + err.Error())
	}

	return &model.CreatedUser{ID: userID, DeviceID: deviceID, Jwt: jwt}, nil
}

func createUser(tx *sqlx.Tx) (uuid.UUID, error) {
	id := uuid.New()
	query := `INSERT INTO users (id, created_at) VALUES (:id, current_timestamp)`
	params := model.Param{"id": id}
	rows, err := dbs.ExecTx(tx, query, params)
	if err != nil {
		return uuid.Nil, err
	}
	if rows != 1 {
		return uuid.Nil, model.ServerError("error creating user")
	}
	return id, nil
}

func createDevice(tx *sqlx.Tx, device model.Device) (uuid.UUID, error) {
	query := `
		INSERT INTO devices (id, user_id, platform, fingerprint, info, created_at) 
		VALUES (:id, :user_id, :platform, :fingerprint, :info, current_timestamp)`
	device.ID = uuid.New()
	rows, err := dbs.ExecTx(tx, query, device)
	if err != nil {
		return uuid.Nil, err
	}
	if rows != 1 {
		return uuid.Nil, model.ServerError("failed to create device")
	}
	return device.ID, nil
}

func deleteDeviceByFingerprint(tx *sqlx.Tx, fingerprint string) error {
	query := `DELETE FROM devices WHERE fingerprint = :fingerprint`
	params := model.Param{"fingerprint": fingerprint}
	_, err := dbs.ExecTx(tx, query, params)
	return err
}
