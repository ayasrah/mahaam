package service

import (
	"mahaam-api/app/repo"
	"mahaam-api/utils/conf"
	emails "mahaam-api/utils/email"
	logs "mahaam-api/utils/log"
	token "mahaam-api/utils/token"
	"slices"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserService interface {
	Create(device Device) *CreatedUser
	SendMeOtp(email string) string
	VerifyOtp(meta Meta, email, sid, otp string) *VerifiedUser
	RefreshToken(meta Meta) *VerifiedUser
	UpdateName(userID uuid.UUID, name string) int64
	Logout(userID uuid.UUID, deviceId uuid.UUID) int64
	Delete(userID uuid.UUID, sid, otp string)
	GetDevices(userID uuid.UUID) []Device
	GetSuggestedEmails(userID uuid.UUID) []SuggestedEmail
	DeleteSuggestedEmail(userID uuid.UUID, suggestedEmailId uuid.UUID)
}

type userService struct {
	userRepo            repo.UserRepo
	deviceRepo          repo.DeviceRepo
	planRepo            repo.PlanRepo
	suggestedEmailsRepo repo.SuggestedEmailRepo
	tokenService        token.TokenService
	emailService        emails.EmailService
	db                  *repo.AppDB
	cfg                 *conf.Conf
	logger              logs.Logger
}

func NewUserService(db *repo.AppDB,
	userRepo repo.UserRepo,
	deviceRepo repo.DeviceRepo,
	planRepo repo.PlanRepo,
	suggestedEmailsRepo repo.SuggestedEmailRepo,
	tokenService token.TokenService,
	emailService emails.EmailService,
	cfg *conf.Conf,
	logger logs.Logger,
) UserService {
	return &userService{
		userRepo:            userRepo,
		deviceRepo:          deviceRepo,
		planRepo:            planRepo,
		suggestedEmailsRepo: suggestedEmailsRepo,
		tokenService:        tokenService,
		emailService:        emailService,
		db:                  db,
		cfg:                 cfg,
		logger:              logger,
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
		if jwt, err = s.tokenService.Create(userId, deviceId); err != nil {
			return err
		}
		return nil
	}

	if err = repo.WithTransaction(s.db, txFn); err != nil {
		panic(err)
	}

	return &CreatedUser{ID: userId, DeviceID: deviceId, Jwt: jwt}
}

func (s *userService) SendMeOtp(email string) string {
	var verifySid string
	if slices.Contains(s.cfg.TestEmails, email) {
		verifySid = s.cfg.TestSID
	} else {
		var err error
		verifySid, err = s.emailService.SendOtp(email)
		if err != nil {
			panic(err)
		}
	}
	return verifySid
}

func (s *userService) VerifyOtp(meta Meta, email, sid, otp string) *VerifiedUser {
	var otpStatus string
	var err error

	if slices.Contains(s.cfg.TestEmails, email) && sid == s.cfg.TestSID && otp == s.cfg.TestOTP {
		otpStatus = "approved"
	} else {
		otpStatus, err = s.emailService.VerifyOtp(otp, sid, email)
	}

	if otpStatus != "approved" {
		panic("OTP not verified for " + email + ", status: " + otpStatus)
	}

	var jwt string
	var newUserId uuid.UUID
	var user *User

	txFn := func(tx *sqlx.Tx) error {
		user = s.userRepo.GetOneByEmail(email)
		if user == nil {
			s.userRepo.UpdateEmail(tx, meta.UserID, email)
			newUserId = meta.UserID
			s.logger.Info(uuid.Nil, "User loggedIn for %s", email)
		} else {
			s.planRepo.UpdateUserID(tx, meta.UserID, user.ID)
			devices := s.deviceRepo.GetMany(user.ID)

			if len(devices) >= 5 {
				s.deviceRepo.Delete(tx, devices[len(devices)-1].ID)
			}
			s.deviceRepo.UpdateUserID(tx, meta.DeviceID, user.ID)
			s.userRepo.Delete(tx, meta.UserID)
			newUserId = user.ID
			s.logger.Info(uuid.Nil, "Merging userId:%s to %s", meta.UserID, user.ID)
		}

		jwt, err = s.tokenService.Create(newUserId, meta.DeviceID)
		return err
	}

	if err = repo.WithTransaction(s.db, txFn); err != nil {
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
	jwt, err := s.tokenService.Create(meta.UserID, meta.DeviceID)
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

func (s *userService) UpdateName(userID uuid.UUID, name string) int64 {
	return s.userRepo.UpdateName(userID, name)
}

func (s *userService) Logout(userID uuid.UUID, deviceId uuid.UUID) int64 {

	device := s.deviceRepo.GetOne(deviceId)
	if device.UserID != userID {
		panic("invalid deviceId")
	}

	var rows int64 = 0
	repo.WithTransaction(s.db, func(tx *sqlx.Tx) error {
		rows = s.deviceRepo.Delete(tx, deviceId)
		return nil
	})
	return int64(rows)
}

func (s *userService) Delete(userID uuid.UUID, sid, otp string) {
	user := s.userRepo.GetOne(userID)
	var otpStatus string
	var err error
	if slices.Contains(s.cfg.TestEmails, *user.Email) && sid == s.cfg.TestSID && otp == s.cfg.TestOTP {
		otpStatus = "approved"
	} else {
		otpStatus, err = s.emailService.VerifyOtp(otp, sid, *user.Email)
		if err != nil {
			panic("OTP not verified for " + *user.Email + ", status: " + otpStatus)
		}
	}

	if otpStatus != "approved" {
		panic("OTP not verified for " + *user.Email + ", status: " + otpStatus)
	}

	repo.WithTransaction(s.db, func(tx *sqlx.Tx) error {
		s.suggestedEmailsRepo.DeleteManyByEmail(*user.Email)
		s.userRepo.Delete(tx, userID)
		return nil
	})

}

func (s *userService) GetDevices(userID uuid.UUID) []Device {
	return s.deviceRepo.GetMany(userID)
}

func (s *userService) GetSuggestedEmails(userID uuid.UUID) []SuggestedEmail {
	return s.suggestedEmailsRepo.GetMany(userID)
}

func (s *userService) DeleteSuggestedEmail(userID uuid.UUID, suggestedEmailId uuid.UUID) {
	suggestedEmail := s.suggestedEmailsRepo.GetOne(suggestedEmailId)
	if suggestedEmail == nil || suggestedEmail.UserID != userID {
		panic("invalid suggestedEmailId")
	}
	s.suggestedEmailsRepo.Delete(suggestedEmailId)
}
