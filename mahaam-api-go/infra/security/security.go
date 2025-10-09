package security

import (
	"errors"
	"mahaam-api/feat/repo"
	"mahaam-api/infra/configs"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService struct {
	deviceRepo repo.DeviceRepo
	userRepo   repo.UserRepo
}

func NewAuthService(deviceRepo repo.DeviceRepo, userRepo repo.UserRepo) AuthService {
	return AuthService{
		deviceRepo: deviceRepo,
		userRepo:   userRepo,
	}
}

func (s *AuthService) ValidateAndExtractJwt(r *gin.Context) (uuid.UUID, uuid.UUID, error) {
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

func (s *AuthService) CreateToken(userId, deviceId uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"userId":   userId,
		"deviceId": deviceId,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iss":      "mahaam-api",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(configs.TokenSecretKey))
	if err != nil {
		return "", errors.New("failed to sign token: " + err.Error())
	}
	return signedToken, nil
}

func (s *AuthService) validateJwt(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(configs.TokenSecretKey), nil
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
