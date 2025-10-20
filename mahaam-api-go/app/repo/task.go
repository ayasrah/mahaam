package repo

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TaskRepo interface {
	GetAll(planID uuid.UUID) []Task
	GetOne(id uuid.UUID) Task
	Create(tx *sqlx.Tx, planID uuid.UUID, title string) uuid.UUID
	DeleteOne(tx *sqlx.Tx, id uuid.UUID) int64
	UpdateDone(tx *sqlx.Tx, id uuid.UUID, done bool) int64
	UpdateTitle(id uuid.UUID, title string) int64
	UpdateOrder(tx *sqlx.Tx, planID uuid.UUID, oldOrder, newOrder int) int64
	UpdateOrderBeforeDelete(tx *sqlx.Tx, planID uuid.UUID, id uuid.UUID) int64
	GetCount(planID uuid.UUID) int64
}

type taskRepo struct {
	db *AppDB
}

func NewTaskRepo(db *AppDB) TaskRepo {
	return &taskRepo{db: db}
}

func (r *taskRepo) GetAll(planID uuid.UUID) []Task {
	query := `SELECT id, plan_id, title, done, sort_order, created_at, updated_at
		FROM tasks WHERE plan_id = :plan_id ORDER BY sort_order DESC`
	param := Param{"plan_id": planID}
	return selectMany[Task](r.db, query, param)
}

func (r *taskRepo) GetOne(id uuid.UUID) Task {
	query := `SELECT id, plan_id, title, done, sort_order, created_at, updated_at FROM tasks WHERE id = :id`
	param := Param{"id": id}
	return selectOne[Task](r.db, query, param)
}

func (r *taskRepo) Create(tx *sqlx.Tx, planID uuid.UUID, title string) uuid.UUID {
	id := uuid.New()
	query := `INSERT INTO tasks (id, plan_id, title, done, sort_order, created_at)
		VALUES (:id, :plan_id, :title, :done, (SELECT COUNT(1) FROM tasks WHERE plan_id = :plan_id), current_timestamp)`
	params := Param{"id": id, "plan_id": planID, "title": title, "done": false}
	executeTransaction(tx, query, params)
	return id
}

func (r *taskRepo) DeleteOne(tx *sqlx.Tx, id uuid.UUID) int64 {
	query := `DELETE FROM tasks WHERE id = :id`
	param := Param{"id": id}
	return executeTransaction(tx, query, param)
}

func (r *taskRepo) UpdateDone(tx *sqlx.Tx, id uuid.UUID, done bool) int64 {
	query := `UPDATE tasks SET done = :done, updated_at = current_timestamp WHERE id = :id`
	params := Param{"id": id, "done": done}
	return executeTransaction(tx, query, params)
}

func (r *taskRepo) UpdateTitle(id uuid.UUID, title string) int64 {
	query := `UPDATE tasks SET title = :title, updated_at = current_timestamp WHERE id = :id`
	params := Param{"id": id, "title": title}
	return execute(r.db, query, params)
}

func (r *taskRepo) UpdateOrderBeforeDelete(tx *sqlx.Tx, planID, id uuid.UUID) int64 {
	query := `UPDATE tasks SET sort_order = sort_order - 1
		WHERE plan_id = :plan_id 
		AND sort_order > (SELECT sort_order FROM tasks WHERE id = :id)`
	params := Param{"id": id, "plan_id": planID}
	return executeTransaction(tx, query, params)
}

func (r *taskRepo) UpdateOrder(tx *sqlx.Tx, planID uuid.UUID, oldOrder, newOrder int) int64 {
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
	return executeTransaction(tx, query, params)
}

func (r *taskRepo) GetCount(planID uuid.UUID) int64 {
	query := `SELECT COUNT(1) FROM tasks WHERE plan_id = :plan_id`
	param := Param{"plan_id": planID}
	return selectOne[int64](r.db, query, param)
}
