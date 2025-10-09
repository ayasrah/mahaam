package repo

import (
	"fmt"

	"mahaam-api/feat/models"
	"mahaam-api/infra/dbs"

	"github.com/google/uuid"
)

type PlanRepo interface {
	GetOne(id UUID) *Plan
	GetMany(userID UUID, planType string) []Plan
	Create(tx Tx, userID UUID, plan PlanIn) UUID
	Update(plan *PlanIn) int64
	Delete(tx Tx, id UUID) int64
	UpdateDonePercent(tx Tx, id UUID) int64
	RemoveFromOrder(tx Tx, userID, id UUID) int64
	UpdateOrder(userID UUID, planType string, oldOrder, newOrder int) int64
	UpdateType(tx Tx, userID, id UUID, planType string) error
	GetCount(userID UUID, planType string) int64
	UpdateUserID(tx Tx, oldUserID, newUserID UUID) int64
}

type planRepo struct {
}

func NewPlanRepo() PlanRepo {
	return &planRepo{}
}

func (r *planRepo) Create(tx Tx, userID UUID, plan PlanIn) UUID {
	id := uuid.New()
	query := `
		INSERT INTO plans (id, user_id, title, starts, ends, type, status, done_percent, sort_order, created_at)
		VALUES (:id, :user_id, :title, :starts, :ends, :type, :status, '0/0', 
		(SELECT COUNT(1) FROM plans WHERE user_id = :user_id AND type = :type), current_timestamp)`
	params := Param{
		"id":      id,
		"user_id": userID,
		"title":   plan.Title,
		"starts":  plan.Starts,
		"ends":    plan.Ends,
		"type":    models.PlanTypeMain,
		"status":  "Open",
	}
	dbs.ExecTx(tx, query, params)
	return id
}

func (r *planRepo) Update(plan *PlanIn) int64 {
	query := `UPDATE plans SET title = :title, starts = :starts, ends = :ends, updated_at = current_timestamp WHERE id = :id`
	return dbs.Exec(query, plan)
}

func (r *planRepo) GetOne(id UUID) *Plan {
	query := `
		SELECT c.id, c.title, c.starts, c.ends, c.type, c.done_percent, c.sort_order,
			EXISTS(SELECT 1 FROM plan_members cm WHERE cm.plan_id = c.id) AS is_shared,
			u.id "user.id", u.email "user.email", u.name "user.name"
		FROM plans c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.id = :id`
	param := Param{"id": id}
	c := dbs.SelectOne[Plan](query, param)
	return &c
}

func (r *planRepo) GetMany(userID UUID, planType string) []Plan {
	query := `
		SELECT c.id, c.title, c.starts, c.ends, c.type, c.done_percent, c.sort_order,
			EXISTS(SELECT 1 FROM plan_members cm WHERE cm.plan_id = c.id) AS is_shared,
			u.id "user.id", u.email "user.email", u.name "user.name"
		FROM plans c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.user_id = :user_id AND c.type = :type
		ORDER BY c.sort_order DESC`

	params := Param{"user_id": userID, "type": planType}
	return dbs.SelectMany[Plan](query, params)
}

func (r *planRepo) Delete(tx Tx, id UUID) int64 {
	query := `DELETE FROM plans WHERE id = :id`
	return dbs.ExecTx(tx, query, Param{"id": id})
}

// UpdateDonePercent updates the done percentage for a plan based on tasks
func (r *planRepo) UpdateDonePercent(tx Tx, id UUID) int64 {
	query := `SELECT COUNT(1) as total, COUNT(CASE WHEN done THEN 1 END) as done FROM tasks WHERE plan_id = :id`
	params := Param{"id": id}
	type Result struct {
		Total int
		Done  int
	}
	result := dbs.SelectOne[Result](query, params)
	donePercent := fmt.Sprintf("%d/%d", result.Done, result.Total)
	updateQuery := `UPDATE plans SET done_percent = :done_percent WHERE id = :id`
	params["done_percent"] = donePercent
	return dbs.ExecTx(tx, updateQuery, params)
}

// RemoveFromOrder decrements sort_order for plans after deletion
func (r *planRepo) RemoveFromOrder(tx Tx, userID, id UUID) int64 {
	query := `
		UPDATE plans SET sort_order = sort_order - 1
		WHERE user_id = :user_id
		AND type = (SELECT type FROM plans WHERE id = :id)
		AND sort_order > (SELECT sort_order FROM plans WHERE id = :id)`
	params := Param{"user_id": userID, "id": id}
	return dbs.ExecTx(tx, query, params)
}

func (r *planRepo) UpdateOrder(userID UUID, planType string, oldOrder, newOrder int) int64 {
	query := `
		UPDATE plans SET sort_order = 
			CASE 
				WHEN sort_order = :oldOrder THEN :newOrder
				WHEN sort_order > :oldOrder AND sort_order <= :newOrder THEN sort_order - 1
				WHEN sort_order >= :newOrder AND sort_order < :oldOrder THEN sort_order + 1
				ELSE sort_order
			END
		WHERE user_id = :user_id AND type = :type`
	params := Param{"user_id": userID, "type": planType, "oldOrder": oldOrder, "newOrder": newOrder}
	return dbs.Exec(query, params)
}

func (r *planRepo) UpdateType(tx Tx, userID, id UUID, planType string) error {
	query := `UPDATE plans SET type = :type, 
		sort_order = (SELECT COUNT(1) FROM plans WHERE user_id = :user_id AND type = :type), 
		updated_at = current_timestamp WHERE id = :id`
	params := Param{"id": id, "type": planType, "user_id": userID}
	_, err := tx.NamedExec(query, params)
	return err
}

func (r *planRepo) GetCount(userID UUID, planType string) int64 {
	query := `SELECT COUNT(1) FROM plans WHERE user_id = :user_id AND type = :type`
	params := Param{"user_id": userID, "type": planType}
	return dbs.SelectOne[int64](query, params)
}

func (r *planRepo) UpdateUserID(tx Tx, oldUserID, newUserID UUID) int64 {
	query := `
		UPDATE plans
		SET user_id = :newUserID,
			sort_order = sort_order + (SELECT COUNT(1) FROM plans WHERE user_id = :newUserID),
			updated_at = current_timestamp
		WHERE user_id = :oldUserID`
	params := Param{"newUserID": newUserID, "oldUserID": oldUserID}
	return dbs.ExecTx(tx, query, params)
}
