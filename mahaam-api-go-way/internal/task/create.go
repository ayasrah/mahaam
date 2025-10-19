package task

import (
	"fmt"
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	tasksLimit = 100
)

func Create(planID uuid.UUID, title string) (uuid.UUID, *model.Err) {
	count, err := getTasksCount(planID)
	if err != nil {
		return uuid.Nil, model.ServerError("failed to get tasks count")
	}
	if count >= tasksLimit {
		return uuid.Nil, model.LogicError("maximum number of tasks is 100", "max_is_100")
	}

	var id uuid.UUID
	err = dbs.WithTx(func(tx *sqlx.Tx) error {
		var txErr error
		id, txErr = createTask(tx, planID, title)
		if txErr != nil {
			return txErr
		}
		if _, txErr = updatePlanDonePercent(tx, planID); txErr != nil {
			return txErr
		}
		return nil
	})

	if err != nil {
		return uuid.Nil, model.ServerError("failed to create task: " + err.Error())
	}

	return id, nil
}

func getTasksCount(planID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(1) FROM tasks WHERE plan_id = :plan_id`
	params := model.Param{"plan_id": planID}
	return dbs.SelectOne[int64](query, params)
}

func createTask(tx *sqlx.Tx, planID uuid.UUID, title string) (uuid.UUID, error) {
	id := uuid.New()
	query := `INSERT INTO tasks (id, plan_id, title, done, sort_order, created_at)
		VALUES (:id, :plan_id, :title, :done, (SELECT COUNT(1) FROM tasks WHERE plan_id = :plan_id), current_timestamp)`
	params := model.Param{"id": id, "plan_id": planID, "title": title, "done": false}
	_, err := dbs.ExecTx(tx, query, params)
	return id, err
}

func updatePlanDonePercent(tx *sqlx.Tx, planID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(1) as total, COUNT(CASE WHEN done THEN 1 END) as done FROM tasks WHERE plan_id = :id`
	params := model.Param{"id": planID}
	type Result struct {
		Total int `db:"total"`
		Done  int `db:"done"`
	}
	result, err := dbs.SelectOne[Result](query, params)
	if err != nil {
		return 0, err
	}
	updateQuery := `UPDATE plans SET done_percent = :done_percent WHERE id = :id`
	params["done_percent"] = fmt.Sprintf("%d/%d", result.Done, result.Total)
	return dbs.ExecTx(tx, updateQuery, params)
}
