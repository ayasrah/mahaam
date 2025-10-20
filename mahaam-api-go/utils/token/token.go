package security

import (
	"errors"
	"mahaam-api/app/repo"
	"mahaam-api/utils/conf"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService interface {
	Parse(r *gin.Context) (uuid.UUID, uuid.UUID, error)
	Create(userId, deviceId uuid.UUID) (string, error)
}

type tokenService struct {
	deviceRepo repo.DeviceRepo
	userRepo   repo.UserRepo
	cfg        *conf.Conf
}

func NewTokenService(deviceRepo repo.DeviceRepo, userRepo repo.UserRepo, cfg *conf.Conf) TokenService {
	return &tokenService{
		deviceRepo: deviceRepo,
		userRepo:   userRepo,
		cfg:        cfg,
	}
}

func (s *tokenService) Parse(r *gin.Context) (uuid.UUID, uuid.UUID, error) {
	authorization := r.GetHeader("Authorization")
	if authorization == "" {
		return uuid.Nil, uuid.Nil, errors.New("authorization header not exists")
	}

	if !strings.HasPrefix(authorization, "Bearer ") {
		return uuid.Nil, uuid.Nil, errors.New("invalid authorization header format")
	}
	tokenString := strings.TrimPrefix(authorization, "Bearer ")

	claims, err := s.validateJwt(tokenString)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	userIdStr, ok := claims["userId"].(string)
	if !ok || userIdStr == "" {
		return uuid.Nil, uuid.Nil, errors.New("userId is required")
	}
	userId, err := uuid.Parse(userIdStr)
	if err != nil || userId == uuid.Nil {
		return uuid.Nil, uuid.Nil, errors.New("userId is empty")
	}

	deviceIdStr, ok := claims["deviceId"].(string)
	if !ok || deviceIdStr == "" {
		return uuid.Nil, uuid.Nil, errors.New("deviceId is required")
	}
	deviceId, err := uuid.Parse(deviceIdStr)
	if err != nil || deviceId == uuid.Nil {
		return uuid.Nil, uuid.Nil, errors.New("deviceId is empty")
	}

	if r.Request.URL.Path != "/user/logout" {
		device := s.deviceRepo.GetOne(deviceId)
		if device.UserID != userId {
			return uuid.Nil, uuid.Nil, errors.New("invalid user info")
		}
	}

	user := s.userRepo.GetOne(userId)
	if user == nil || user.ID != userId {
		panic("user not found")
	}

	return userId, deviceId, nil
}

func (s *tokenService) Create(userId, deviceId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"userId":   userId,
		"deviceId": deviceId,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iss":      "mahaam-api",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.cfg.TokenSecretKey))
	if err != nil {
		return "", errors.New("failed to sign token: " + err.Error())
	}
	return signedToken, nil
}

func (s *tokenService) validateJwt(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.cfg.TokenSecretKey), nil
	})
	if err != nil {
		return nil, errors.New("invalid token: " + err.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
