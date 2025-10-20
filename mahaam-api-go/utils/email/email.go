package emails

import (
	"mahaam-api/utils/conf"
	logs "mahaam-api/utils/log"

	"github.com/google/uuid"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

type EmailService interface {
	SendOtp(email string) (string, error)
	VerifyOtp(otp, sid, email string) (string, error)
}

type emailService struct {
	cfg    *conf.Conf
	client *twilio.RestClient
	logger logs.Logger
}

func NewEmailService(cfg *conf.Conf, logger logs.Logger) EmailService {
	var client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.EmailAccountSID,
		Password: cfg.EmailAuthToken,
	})
	return &emailService{cfg: cfg, client: client, logger: logger}
}

func (s *emailService) SendOtp(email string) (string, error) {
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(email)
	params.SetChannel("email")

	verification, err := s.client.VerifyV2.CreateVerification(s.cfg.EmailVerificationServiceSID, params)
	if err != nil {
		s.logger.Error(uuid.Nil, "Error sending OTP to %s: %v", email, err)
		return "", err
	}
	if verification.Sid == nil {
		return "", nil
	}
	return *verification.Sid, nil
}

func (s *emailService) VerifyOtp(otp, sid, email string) (string, error) {
	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(email)
	params.SetCode(otp)
	params.SetVerificationSid(sid)

	check, err := s.client.VerifyV2.CreateVerificationCheck(s.cfg.EmailVerificationServiceSID, params)
	if err != nil {
		s.logger.Info(uuid.Nil, "Error verifying OTP for %s: %v", email, err)
		return "", err
	}
	if check.Status == nil {
		return "", nil
	}
	return *check.Status, nil
}
