package user

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
)

func UpdateName(userID uuid.UUID, name string) *model.Err {
	if err := updateUserName(userID, name); err != nil {
		return model.ServerError("failed to update name: " + err.Error())
	}
	return nil
}

func updateUserName(userID uuid.UUID, name string) error {
	query := `UPDATE users SET name = :name, updated_at = current_timestamp WHERE id = :id`
	params := model.Param{"id": userID, "name": name}
	_, err := dbs.Exec(query, params)
	return err
}
