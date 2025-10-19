package task

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func UpdateDone(planID, taskID uuid.UUID, done bool) *model.Err {
	err := dbs.WithTx(func(tx *sqlx.Tx) error {
		if _, err := updateTaskDone(tx, taskID, done); err != nil {
			return err
		}
		if _, err := updatePlanDonePercent(tx, planID); err != nil {
			return err
		}

		// Get all tasks to find current position and reorder
		tasks, err := getAllTasks(planID)
		if err != nil {
			return err
		}

		taskIndex := -1
		for i, t := range tasks {
			if t.ID == taskID {
				taskIndex = i
				break
			}
		}

		if taskIndex == -1 {
			return nil // Task not found, skip reorder
		}

		newOrder := 0
		if done {
			newOrder = len(tasks) - 1
		}

		if taskIndex != newOrder {
			if err := reOrderWithTx(tx, planID, taskIndex, newOrder); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return model.ServerError("failed to update task done status: " + err.Error())
	}

	return nil
}

func updateTaskDone(tx *sqlx.Tx, taskID uuid.UUID, done bool) (int64, error) {
	query := `UPDATE tasks SET done = :done, updated_at = current_timestamp WHERE id = :id`
	params := model.Param{"id": taskID, "done": done}
	return dbs.ExecTx(tx, query, params)
}

func getAllTasks(planID uuid.UUID) ([]model.Task, error) {
	query := `SELECT id, plan_id, title, done, sort_order, created_at, updated_at
		FROM tasks WHERE plan_id = :plan_id ORDER BY sort_order DESC`
	params := model.Param{"plan_id": planID}
	return dbs.SelectMany[model.Task](query, params)
}
