package task

import (
	"mahaam-api/internal/model"

	"github.com/google/uuid"
)

func GetMany(planID uuid.UUID) ([]model.Task, *model.Err) {
	tasks, err := getAllTasks(planID)
	if err != nil {
		return nil, model.ServerError("failed to get tasks: " + err.Error())
	}
	return tasks, nil
}
