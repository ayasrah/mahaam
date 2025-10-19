package user

import (
	"mahaam-api/internal/pkg/configs"
	logs "mahaam-api/internal/pkg/log"

	"github.com/google/uuid"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

var client *twilio.RestClient

func InitEmail() {
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: configs.EmailAccountSID,
		Password: configs.EmailAuthToken,
	})
}

func sendEmailOtp(email string) (string, error) {
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(email)
	params.SetChannel("email")

	verification, err := client.VerifyV2.CreateVerification(configs.EmailVerificationServiceSID, params)
	if err != nil {
		logs.Error(uuid.Nil, "Error sending OTP to %s: %v", email, err)
		return "", err
	}
	if verification.Sid == nil {
		return "", nil
	}
	return *verification.Sid, nil
}

func verifyEmailOtp(otp, sid, email string) (string, error) {
	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(email)
	params.SetCode(otp)
	params.SetVerificationSid(sid)

	check, err := client.VerifyV2.CreateVerificationCheck(configs.EmailVerificationServiceSID, params)
	if err != nil {
		logs.Info(uuid.Nil, "Error verifying OTP for %s: %v", email, err)
		return "", err
	}
	if check.Status == nil {
		return "", nil
	}
	return *check.Status, nil
}
