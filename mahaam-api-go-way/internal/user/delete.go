package user

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/configs"
	"mahaam-api/internal/pkg/dbs"
	"slices"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func Delete(userID uuid.UUID, sid, otp string) *model.Err {
	user, err := GetUser(userID)
	if err != nil {
		return model.ServerError("failed to get user: " + err.Error())
	}

	if user == nil || user.Email == nil {
		return model.NotFoundError("user not found")
	}

	// Verify OTP
	var otpStatus string
	if slices.Contains(configs.TestEmails, *user.Email) && sid == configs.TestSID && otp == configs.TestOTP {
		otpStatus = "approved"
	} else {
		var verifyErr error
		otpStatus, verifyErr = verifyEmailOtp(otp, sid, *user.Email)
		if verifyErr != nil || otpStatus != "approved" {
			return model.UnauthorizedError("OTP not verified for " + *user.Email)
		}
	}

	// Delete suggested emails and user
	err = dbs.WithTx(func(tx *sqlx.Tx) error {
		if deleteErr := deleteSuggestedEmailsByEmail(*user.Email); deleteErr != nil {
			return deleteErr
		}
		return deleteUser(tx, userID)
	})

	if err != nil {
		return model.ServerError("failed to delete user: " + err.Error())
	}

	return nil
}

func deleteSuggestedEmailsByEmail(email string) error {
	query := `DELETE FROM suggested_emails WHERE email = :email`
	params := model.Param{"email": email}
	_, err := dbs.Exec(query, params)
	return err
}
