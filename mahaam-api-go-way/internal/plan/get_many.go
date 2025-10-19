package plan

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
)

func GetMany(userID uuid.UUID, planType string) ([]model.Plan, *model.Err) {
	plans, err := getUserPlans(userID, planType)
	if err != nil {
		return nil, model.ServerError("failed to get plans: " + err.Error())
	}
	sharedPlans, err := getOthersSharedPlans(userID)
	if err != nil {
		return nil, model.ServerError("failed to get others shared plans: " + err.Error())
	}
	return append(plans, sharedPlans...), nil
}

func getUserPlans(userID uuid.UUID, planType string) ([]model.Plan, error) {
	query := `
		SELECT c.id, c.title, c.starts, c.ends, c.type, c.done_percent, c.sort_order,
			EXISTS(SELECT 1 FROM plan_members cm WHERE cm.plan_id = c.id) AS is_shared,
			u.id "user.id", u.email "user.email", u.name "user.name"
		FROM plans c
		LEFT JOIN users u ON c.user_id = u.id
		WHERE c.user_id = :user_id AND c.type = :type
		ORDER BY c.sort_order DESC`

	params := model.Param{"user_id": userID, "type": planType}
	return dbs.SelectMany[model.Plan](query, params)
}

func getOthersSharedPlans(userID uuid.UUID) ([]model.Plan, error) {
	query := `
		SELECT c.id, c.title, c.starts, c.ends, c.type, c.done_percent, c.sort_order, 
			true AS is_shared, u.id as "user.id",u.email as "user.email",u.name as "user.name"
		FROM plan_members cm
		LEFT JOIN plans c ON cm.plan_id = c.id
		LEFT JOIN users u ON c.user_id = u.id
		WHERE cm.user_id = :user_id
		ORDER BY c.sort_order ASC`
	params := model.Param{"user_id": userID}
	return dbs.SelectMany[model.Plan](query, params)
}
