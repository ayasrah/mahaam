package task

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func Delete(planID, taskID uuid.UUID) *model.Err {
	err := dbs.WithTx(func(tx *sqlx.Tx) error {
		if _, err := updateOrderBeforeDelete(tx, planID, taskID); err != nil {
			return err
		}
		if _, err := deleteTask(tx, taskID); err != nil {
			return err
		}
		if _, err := updatePlanDonePercent(tx, planID); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return model.ServerError("failed to delete task: " + err.Error())
	}

	return nil
}

func updateOrderBeforeDelete(tx *sqlx.Tx, planID, taskID uuid.UUID) (int64, error) {
	query := `UPDATE tasks SET sort_order = sort_order - 1
		WHERE plan_id = :plan_id 
		AND sort_order > (SELECT sort_order FROM tasks WHERE id = :id)`
	params := model.Param{"id": taskID, "plan_id": planID}
	return dbs.ExecTx(tx, query, params)
}

func deleteTask(tx *sqlx.Tx, taskID uuid.UUID) (int64, error) {
	query := `DELETE FROM tasks WHERE id = :id`
	params := model.Param{"id": taskID}
	return dbs.ExecTx(tx, query, params)
}
