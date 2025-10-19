package user

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
)

func GetUser(userID uuid.UUID) (*model.User, error) {
	query := `SELECT id, name, email FROM users WHERE id = :id`
	params := model.Param{"id": userID}
	user, err := dbs.SelectOne[model.User](query, params)
	if err != nil {
		return nil, model.ServerError("failed to get user: " + err.Error())
	}
	if user.ID == uuid.Nil {
		return nil, model.NotFoundError("user not found")
	}
	return &user, nil
}

func GetUserByEmail(email string) (*model.User, error) {
	query := `SELECT id, name, email FROM users WHERE email = :email`
	params := model.Param{"email": email}
	user, err := dbs.SelectOne[model.User](query, params)
	if err != nil {
		return nil, model.ServerError("failed to get user: " + err.Error())
	}
	if user.ID == uuid.Nil {
		return nil, model.NotFoundError("user not found")
	}
	return &user, nil
}
