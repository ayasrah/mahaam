package repo

import (
	"mahaam-api/infra/dbs"

	"github.com/google/uuid"
)

type TaskRepo interface {
	GetAll(planID UUID) []Task
	GetOne(id UUID) Task
	Create(tx Tx, planID UUID, title string) UUID
	DeleteOne(tx Tx, id UUID) int64
	UpdateDone(tx Tx, id UUID, done bool) int64
	UpdateTitle(id UUID, title string) int64
	UpdateOrder(tx Tx, planID UUID, oldOrder, newOrder int) int64
	UpdateOrderBeforeDelete(tx Tx, planID UUID, id UUID) int64
	GetCount(planID UUID) int64
}

type taskRepo struct {
}

func NewTaskRepo() TaskRepo {
	return &taskRepo{}
}

func (r *taskRepo) GetAll(planID UUID) []Task {
	query := `SELECT id, plan_id, title, done, sort_order, created_at, updated_at
		FROM tasks WHERE plan_id = :plan_id ORDER BY sort_order DESC`
	param := Param{"plan_id": planID}
	return dbs.SelectMany[Task](query, param)
}

func (r *taskRepo) GetOne(id UUID) Task {
	query := `SELECT id, plan_id, title, done, sort_order, created_at, updated_at FROM tasks WHERE id = :id`
	param := Param{"id": id}
	return dbs.SelectOne[Task](query, param)
}

func (r *taskRepo) Create(tx Tx, planID UUID, title string) UUID {
	id := uuid.New()
	query := `INSERT INTO tasks (id, plan_id, title, done, sort_order, created_at)
		VALUES (:id, :plan_id, :title, :done, (SELECT COUNT(1) FROM tasks WHERE plan_id = :plan_id), current_timestamp)`
	params := Param{"id": id, "plan_id": planID, "title": title, "done": false}
	dbs.ExecTx(tx, query, params)
	return id
}

func (r *taskRepo) DeleteOne(tx Tx, id UUID) int64 {
	query := `DELETE FROM tasks WHERE id = :id`
	param := Param{"id": id}
	return dbs.ExecTx(tx, query, param)
}

func (r *taskRepo) UpdateDone(tx Tx, id UUID, done bool) int64 {
	query := `UPDATE tasks SET done = :done, updated_at = current_timestamp WHERE id = :id`
	params := Param{"id": id, "done": done}
	return dbs.ExecTx(tx, query, params)
}

func (r *taskRepo) UpdateTitle(id UUID, title string) int64 {
	query := `UPDATE tasks SET title = :title, updated_at = current_timestamp WHERE id = :id`
	params := Param{"id": id, "title": title}
	return dbs.Exec(query, params)
}

func (r *taskRepo) UpdateOrderBeforeDelete(tx Tx, planID, id UUID) int64 {
	query := `UPDATE tasks SET sort_order = sort_order - 1
		WHERE plan_id = :plan_id 
		AND sort_order > (SELECT sort_order FROM tasks WHERE id = :id)`
	params := Param{"id": id, "plan_id": planID}
	return dbs.ExecTx(tx, query, params)
}

func (r *taskRepo) UpdateOrder(tx Tx, planID UUID, oldOrder, newOrder int) int64 {
	query := `
		UPDATE tasks
		SET sort_order = CASE
			WHEN sort_order = :old_index THEN :new_index
			WHEN sort_order > :old_index AND sort_order <= :new_index THEN sort_order - 1
			WHEN sort_order >= :new_index AND sort_order < :old_index THEN sort_order + 1
			ELSE sort_order
		END
		WHERE plan_id = :plan_id`
	params := Param{
		"old_index": oldOrder,
		"new_index": newOrder,
		"plan_id":   planID,
	}
	return dbs.ExecTx(tx, query, params)
}

func (r *taskRepo) GetCount(planID UUID) int64 {
	query := `SELECT COUNT(1) FROM tasks WHERE plan_id = :plan_id`
	param := Param{"plan_id": planID}
	return dbs.SelectOne[int64](query, param)
}
