package plan

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
)

func GetOne(id uuid.UUID) (model.Plan, *model.Err) {
	plan, err := getPlan(id)
	if err != nil {
		return model.Plan{}, model.ServerError("failed to get plan: " + err.Error())
	}
	if plan.IsShared {
		users, err := getPlanUsers(id)
		if err != nil {
			return model.Plan{}, model.ServerError("failed to get plan users: " + err.Error())
		}
		plan.Members = users
	}
	return plan, nil
}

func getPlan(id uuid.UUID) (model.Plan, error) {
	query := `
		SELECT c.id, c.title, c.starts, c.ends, c.type, c.done_percent, c.sort_order,
			EXISTS(SELECT 1 FROM plan_members cm WHERE cm.plan_id = c.id) AS is_shared,
			u.id "user.id", u.email "user.email", u.name "user.name"
		FROM plans c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.id = :id`
	param := model.Param{"id": id}
	return dbs.SelectOne[model.Plan](query, param)
}

func getPlanUsers(planID uuid.UUID) ([]model.User, error) {
	query := `
		SELECT u.id, u.email, u.name
		FROM plan_members cm
		LEFT JOIN users u ON cm.user_id = u.id
		WHERE cm.plan_id = :plan_id
		ORDER BY u.created_at DESC`
	params := model.Param{"plan_id": planID}
	return dbs.SelectMany[model.User](query, params)
}
