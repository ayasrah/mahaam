package task

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
)

func UpdateTitle(taskID uuid.UUID, title string) *model.Err {
	if err := updateTaskTitle(taskID, title); err != nil {
		return model.ServerError("failed to update task title: " + err.Error())
	}
	return nil
}

func updateTaskTitle(taskID uuid.UUID, title string) error {
	query := `UPDATE tasks SET title = :title, updated_at = current_timestamp WHERE id = :id`
	params := model.Param{"id": taskID, "title": title}
	_, err := dbs.Exec(query, params)
	return err
}
