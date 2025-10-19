package user

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/configs"
	"slices"
)

func SendOtp(email string) (string, *model.Err) {
	var verifySid string

	if slices.Contains(configs.TestEmails, email) {
		verifySid = configs.TestSID
	} else {
		var err error
		verifySid, err = sendEmailOtp(email)
		if err != nil {
			return "", model.ServerError("failed to send OTP: " + err.Error())
		}
	}

	return verifySid, nil
}
