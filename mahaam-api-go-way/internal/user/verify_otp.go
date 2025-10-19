package user

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/configs"
	"mahaam-api/internal/pkg/dbs"
	"slices"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func VerifyOtp(meta model.Meta, email, sid, otp string) (*model.VerifiedUser, *model.Err) {
	var otpStatus string
	var err error

	// Verify OTP
	if slices.Contains(configs.TestEmails, email) && sid == configs.TestSID && otp == configs.TestOTP {
		otpStatus = "approved"
	} else {
		otpStatus, err = verifyEmailOtp(otp, sid, email)
		if err != nil || otpStatus != "approved" {
			return nil, model.UnauthorizedError("OTP not verified for " + email)
		}
	}

	var jwt string
	var newUserID uuid.UUID
	var user *model.User

	txFn := func(tx *sqlx.Tx) error {
		// Check if user with email exists
		user, err = getUserByEmail(email)
		if err != nil {
			return err
		}

		if user == nil {
			// First time login - update the anonymous user's email
			if err = updateUserEmail(tx, meta.UserID, email); err != nil {
				return err
			}
			newUserID = meta.UserID
		} else {
			// Existing user - merge accounts
			// Transfer plans from anonymous user to existing user
			if err = updatePlanUserID(tx, meta.UserID, user.ID); err != nil {
				return err
			}

			// Get devices to check limit
			devices, err := getDevices(user.ID)
			if err != nil {
				return err
			}

			// Delete oldest device if limit reached
			if len(devices) >= 5 {
				if err = deleteDevice(tx, devices[len(devices)-1].ID); err != nil {
					return err
				}
			}

			// Transfer device to existing user
			if err = updateDeviceUserID(tx, meta.DeviceID, user.ID); err != nil {
				return err
			}

			// Delete anonymous user
			if err = deleteUser(tx, meta.UserID); err != nil {
				return err
			}

			newUserID = user.ID
		}

		return nil
	}

	if err = dbs.WithTx(txFn); err != nil {
		return nil, model.ServerError("failed to verify OTP: " + err.Error())
	}

	jwt, err = CreateToken(newUserID, meta.DeviceID)
	if err != nil {
		return nil, model.ServerError("failed to create token: " + err.Error())
	}

	userFullName := ""
	if user != nil && user.Name != nil {
		userFullName = *user.Name
	}

	return &model.VerifiedUser{
		UserID:       newUserID,
		DeviceID:     meta.DeviceID,
		Jwt:          jwt,
		UserFullName: &userFullName,
		Email:        &email,
	}, nil
}

func getUserByEmail(email string) (*model.User, error) {
	query := `SELECT id, name, email FROM users WHERE email = :email`
	params := model.Param{"email": email}
	user, err := dbs.SelectOne[model.User](query, params)
	if err != nil {
		return nil, err
	}
	if user.ID == uuid.Nil {
		return nil, nil
	}
	return &user, nil
}

func updateUserEmail(tx *sqlx.Tx, userID uuid.UUID, email string) error {
	query := `UPDATE users SET email = :email, updated_at = current_timestamp WHERE id = :id`
	params := model.Param{"id": userID, "email": email}
	_, err := dbs.ExecTx(tx, query, params)
	return err
}

func updatePlanUserID(tx *sqlx.Tx, oldUserID, newUserID uuid.UUID) error {
	query := `
		UPDATE plans
		SET user_id = :newUserID,
			sort_order = sort_order + (SELECT COUNT(1) FROM plans WHERE user_id = :newUserID),
			updated_at = current_timestamp
		WHERE user_id = :oldUserID`
	params := model.Param{"newUserID": newUserID, "oldUserID": oldUserID}
	_, err := dbs.ExecTx(tx, query, params)
	return err
}

func getDevices(userID uuid.UUID) ([]model.Device, error) {
	query := `SELECT id, user_id, platform, fingerprint, info, created_at 
			FROM devices WHERE user_id = :user_id ORDER BY created_at DESC`
	params := model.Param{"user_id": userID}
	return dbs.SelectMany[model.Device](query, params)
}

func deleteDevice(tx *sqlx.Tx, deviceID uuid.UUID) error {
	query := `DELETE FROM devices WHERE id = :id`
	params := model.Param{"id": deviceID}
	_, err := dbs.ExecTx(tx, query, params)
	return err
}

func updateDeviceUserID(tx *sqlx.Tx, deviceID, userID uuid.UUID) error {
	query := "UPDATE devices SET user_id = :user_id, updated_at = current_timestamp WHERE id = :device_id"
	params := model.Param{"user_id": userID, "device_id": deviceID}
	_, err := dbs.ExecTx(tx, query, params)
	return err
}

func deleteUser(tx *sqlx.Tx, userID uuid.UUID) error {
	query := `DELETE FROM users WHERE id = :id`
	params := model.Param{"id": userID}
	_, err := dbs.ExecTx(tx, query, params)
	return err
}
