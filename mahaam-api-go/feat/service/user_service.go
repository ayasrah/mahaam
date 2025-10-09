package service

import (
	"mahaam-api/feat/repo"
	"mahaam-api/infra/configs"
	"mahaam-api/infra/dbs"
	"mahaam-api/infra/emails"
	logs "mahaam-api/infra/log"
	"mahaam-api/infra/security"
	"slices"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserService interface {
	Create(device Device) *CreatedUser
	SendMeOtp(email string) string
	VerifyOtp(meta Meta, email, sid, otp string) *VerifiedUser
	RefreshToken(meta Meta) *VerifiedUser
	UpdateName(userID UUID, name string) int64
	Logout(userID UUID, deviceId UUID) int64
	Delete(userID UUID, sid, otp string)
	GetDevices(userID UUID) []Device
	GetSuggestedEmails(userID UUID) []SuggestedEmail
	DeleteSuggestedEmail(userID UUID, suggestedEmailId UUID)
}

type userService struct {
	userRepo            repo.UserRepo
	deviceRepo          repo.DeviceRepo
	planRepo            repo.PlanRepo
	suggestedEmailsRepo repo.SuggestedEmailRepo
	authService         security.AuthService
}

func NewUserService(db *sqlx.DB,
	userRepo repo.UserRepo,
	deviceRepo repo.DeviceRepo,
	planRepo repo.PlanRepo,
	suggestedEmailsRepo repo.SuggestedEmailRepo,
	authService security.AuthService,
) UserService {
	return &userService{
		userRepo:            userRepo,
		deviceRepo:          deviceRepo,
		planRepo:            planRepo,
		suggestedEmailsRepo: suggestedEmailsRepo,
		authService:         authService,
	}
}

func (s *userService) Create(device Device) *CreatedUser {
	var jwt string
	var userId uuid.UUID
	var deviceId uuid.UUID
	var err error

	txFn := func(tx *sqlx.Tx) error {
		userId = s.userRepo.Create(tx)
		s.deviceRepo.DeleteByFingerprint(tx, device.Fingerprint)
		device.UserID = userId
		deviceId = s.deviceRepo.Create(tx, device)
		if jwt, err = s.authService.CreateToken(userId, deviceId); err != nil {
			return err
		}
		return nil
	}

	if err = dbs.WithTx(txFn); err != nil {
		panic(err)
	}

	return &CreatedUser{ID: userId, DeviceID: deviceId, Jwt: jwt}
}

func (s *userService) SendMeOtp(email string) string {
	var verifySid string
	if slices.Contains(configs.TestEmails, email) {
		verifySid = configs.TestSID
	} else {
		var err error
		verifySid, err = emails.SendOtp(email)
		if err != nil {
			panic(err)
		}
	}
	return verifySid
}

func (s *userService) VerifyOtp(meta Meta, email, sid, otp string) *VerifiedUser {
	var otpStatus string
	var err error

	if slices.Contains(configs.TestEmails, email) && sid == configs.TestSID && otp == configs.TestOTP {
		otpStatus = "approved"
	} else {
		otpStatus, err = emails.VerifyOtp(otp, sid, email)
	}

	if otpStatus != "approved" {
		panic("OTP not verified for " + email + ", status: " + otpStatus)
	}

	var jwt string
	var newUserId UUID
	var user *User

	txFn := func(tx *sqlx.Tx) error {
		user = s.userRepo.GetOneByEmail(email)
		if user == nil {
			s.userRepo.UpdateEmail(tx, meta.UserID, email)
			newUserId = meta.UserID
			logs.Info(uuid.Nil, "User loggedIn for %s", email)
		} else {
			s.planRepo.UpdateUserID(tx, meta.UserID, user.ID)
			devices := s.deviceRepo.GetMany(user.ID)

			if len(devices) >= 5 {
				s.deviceRepo.Delete(tx, devices[len(devices)-1].ID)
			}
			s.deviceRepo.UpdateUserID(tx, meta.DeviceID, user.ID)
			s.userRepo.Delete(tx, meta.UserID)
			newUserId = user.ID
			logs.Info(uuid.Nil, "Merging userId:%s to %s", meta.UserID, user.ID)
		}

		jwt, err = s.authService.CreateToken(newUserId, meta.DeviceID)
		return err
	}

	if err = dbs.WithTx(txFn); err != nil {
		panic(err)
	}

	userFullName := ""
	if user != nil && user.Name != nil {
		userFullName = *user.Name
	}
	return &VerifiedUser{
		UserID:       newUserId,
		DeviceID:     meta.DeviceID,
		Jwt:          jwt,
		UserFullName: &userFullName,
		Email:        &email,
	}
}

func (s *userService) RefreshToken(meta Meta) *VerifiedUser {
	user := s.userRepo.GetOne(meta.UserID)
	jwt, err := s.authService.CreateToken(meta.UserID, meta.DeviceID)
	if err != nil {
		return nil
	}
	return &VerifiedUser{
		UserID:       meta.UserID,
		DeviceID:     meta.DeviceID,
		Jwt:          jwt,
		UserFullName: user.Name,
		Email:        user.Email,
	}
}

func (s *userService) UpdateName(userID UUID, name string) int64 {
	return s.userRepo.UpdateName(userID, name)
}

func (s *userService) Logout(userID UUID, deviceId UUID) int64 {

	device := s.deviceRepo.GetOne(deviceId)
	if device.UserID != userID {
		panic("invalid deviceId")
	}

	var rows int64 = 0
	dbs.WithTx(func(tx *sqlx.Tx) error {
		rows = s.deviceRepo.Delete(tx, deviceId)
		return nil
	})
	return int64(rows)
}

func (s *userService) Delete(userID UUID, sid, otp string) {
	user := s.userRepo.GetOne(userID)
	var otpStatus string
	var err error
	if slices.Contains(configs.TestEmails, *user.Email) && sid == configs.TestSID && otp == configs.TestOTP {
		otpStatus = "approved"
	} else {
		otpStatus, err = emails.VerifyOtp(otp, sid, *user.Email)
		if err != nil {
			panic("OTP not verified for " + *user.Email + ", status: " + otpStatus)
		}
	}

	if otpStatus != "approved" {
		panic("OTP not verified for " + *user.Email + ", status: " + otpStatus)
	}

	dbs.WithTx(func(tx *sqlx.Tx) error {
		s.suggestedEmailsRepo.DeleteManyByEmail(*user.Email)
		s.userRepo.Delete(tx, userID)
		return nil
	})

}

func (s *userService) GetDevices(userID UUID) []Device {
	return s.deviceRepo.GetMany(userID)
}

func (s *userService) GetSuggestedEmails(userID UUID) []SuggestedEmail {
	return s.suggestedEmailsRepo.GetMany(userID)
}

func (s *userService) DeleteSuggestedEmail(userID UUID, suggestedEmailId UUID) {
	suggestedEmail := s.suggestedEmailsRepo.GetOne(suggestedEmailId)
	if suggestedEmail == nil || suggestedEmail.UserID != userID {
		panic("invalid suggestedEmailId")
	}
	s.suggestedEmailsRepo.Delete(suggestedEmailId)
}
