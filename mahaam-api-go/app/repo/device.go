package repo

import (
	"mahaam-api/app/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DeviceRepo interface {
	GetOne(id uuid.UUID) *Device
	GetMany(userID uuid.UUID) []Device
	Create(tx *sqlx.Tx, device Device) uuid.UUID
	Delete(tx *sqlx.Tx, id uuid.UUID) int64
	DeleteByUser(userID, exceptDeviceID uuid.UUID) int64
	DeleteByFingerprint(tx *sqlx.Tx, fingerprint string) int64
	UpdateUserID(tx *sqlx.Tx, deviceID, userID uuid.UUID) int64
}

type deviceRepo struct {
	db *AppDB
}

func NewDeviceRepo(db *AppDB) DeviceRepo {
	return &deviceRepo{db: db}
}

func (r *deviceRepo) GetDB() *AppDB {
	return r.db
}

func (r *deviceRepo) GetOne(id uuid.UUID) *Device {
	query := `SELECT id, user_id, platform, fingerprint, info, created_at FROM devices WHERE id = :id`
	dev := selectOne[Device](r.db, query, Param{"id": id})
	return &dev
}

func (r *deviceRepo) GetMany(userID uuid.UUID) []Device {
	query := `SELECT id, user_id, platform, fingerprint, info, created_at 
			FROM devices WHERE user_id = :user_id ORDER BY created_at DESC`
	return selectMany[Device](r.db, query, Param{"user_id": userID})
}

func (r *deviceRepo) Create(tx *sqlx.Tx, device Device) uuid.UUID {
	query := `
		INSERT INTO devices (id, user_id, platform, fingerprint, info, created_at) 
		VALUES (:id, :user_id, :platform, :fingerprint, :info, current_timestamp)`
	device.ID = uuid.New()
	rows := executeTransaction(tx, query, device)
	if rows != 1 {
		panic(models.ServerError("failed to create device"))
	}
	return device.ID
}

func (r *deviceRepo) Delete(tx *sqlx.Tx, id uuid.UUID) int64 {
	query := `DELETE FROM devices WHERE id = :id`
	param := Param{"id": id}
	return executeTransaction(tx, query, param)
}

func (r *deviceRepo) DeleteByUser(userId, exceptDeviceId uuid.UUID) int64 {
	query := `DELETE FROM devices WHERE user_id = :user_id AND id != :except_device_id`
	params := Param{"user_id": userId, "except_device_id": exceptDeviceId}
	return execute(r.db, query, params)
}

func (r *deviceRepo) DeleteByFingerprint(tx *sqlx.Tx, fingerprint string) int64 {
	query := `DELETE FROM devices WHERE fingerprint = :fingerprint`
	params := Param{"fingerprint": fingerprint}
	return executeTransaction(tx, query, params)
}

func (r *deviceRepo) UpdateUserID(tx *sqlx.Tx, deviceId, userId uuid.UUID) int64 {
	query := "UPDATE devices SET user_id = :user_id, updated_at = current_timestamp WHERE id = :device_id"
	params := Param{"user_id": userId, "device_id": deviceId}
	return executeTransaction(tx, query, params)
}
