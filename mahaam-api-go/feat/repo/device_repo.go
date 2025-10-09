package repo

import (
	"mahaam-api/infra/dbs"

	"github.com/google/uuid"
)

type DeviceRepo interface {
	GetOne(id UUID) *Device
	GetMany(userID UUID) []Device
	Create(tx Tx, device Device) UUID
	Delete(tx Tx, id UUID) int64
	DeleteByUser(userID, exceptDeviceID UUID) int64
	DeleteByFingerprint(tx Tx, fingerprint string) int64
	UpdateUserID(tx Tx, deviceID, userID UUID) int64
}

type deviceRepo struct {
}

func NewDeviceRepo() DeviceRepo {
	return &deviceRepo{}
}

func (r *deviceRepo) GetOne(id UUID) *Device {
	query := `SELECT id, user_id, platform, fingerprint, info, created_at FROM devices WHERE id = :id`
	dev := dbs.SelectOne[Device](query, Param{"id": id})
	return &dev
}

func (r *deviceRepo) GetMany(userID UUID) []Device {
	query := `SELECT id, user_id, platform, fingerprint, info, created_at 
			FROM devices WHERE user_id = :user_id ORDER BY created_at DESC`
	return dbs.SelectMany[Device](query, Param{"user_id": userID})
}

func (r *deviceRepo) Create(tx Tx, device Device) UUID {
	query := `
		INSERT INTO devices (id, user_id, platform, fingerprint, info, created_at) 
		VALUES (:id, :user_id, :platform, :fingerprint, :info, current_timestamp)`
	device.ID = uuid.New()
	rows := dbs.ExecTx(tx, query, device)
	if rows != 1 {
		panic("failed to create device")
	}
	return device.ID
}

func (r *deviceRepo) Delete(tx Tx, id UUID) int64 {
	query := `DELETE FROM devices WHERE id = :id`
	param := Param{"id": id}
	return dbs.ExecTx(tx, query, param)
}

func (r *deviceRepo) DeleteByUser(userId, exceptDeviceId UUID) int64 {
	query := `DELETE FROM devices WHERE user_id = :user_id AND id != :except_device_id`
	params := Param{"user_id": userId, "except_device_id": exceptDeviceId}
	return dbs.Exec(query, params)
}

func (r *deviceRepo) DeleteByFingerprint(tx Tx, fingerprint string) int64 {
	query := `DELETE FROM devices WHERE fingerprint = :fingerprint`
	params := Param{"fingerprint": fingerprint}
	return dbs.ExecTx(tx, query, params)
}

func (r *deviceRepo) UpdateUserID(tx Tx, deviceId, userId UUID) int64 {
	query := "UPDATE devices SET user_id = :user_id, updated_at = current_timestamp WHERE id = :device_id"
	params := Param{"user_id": userId, "device_id": deviceId}
	return dbs.ExecTx(tx, query, params)
}
