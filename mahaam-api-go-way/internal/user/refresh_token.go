package user

import (
	"mahaam-api/internal/model"
)

func RefreshToken(meta model.Meta) (*model.VerifiedUser, *model.Err) {
	user, err := GetUser(meta.UserID)
	if err != nil {
		return nil, model.ServerError("failed to get user: " + err.Error())
	}

	jwt, err := CreateToken(meta.UserID, meta.DeviceID)
	if err != nil {
		return nil, model.ServerError("failed to create token: " + err.Error())
	}

	return &model.VerifiedUser{
		UserID:       meta.UserID,
		DeviceID:     meta.DeviceID,
		Jwt:          jwt,
		UserFullName: user.Name,
		Email:        user.Email,
	}, nil
}
