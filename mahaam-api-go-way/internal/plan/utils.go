package plan

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func validateUserOwnsThePlan(userID uuid.UUID, planID uuid.UUID) error {
	plan, err := getOne(planID)
	if err != nil {
		return err
	}
	if plan.ID == uuid.Nil {
		return model.NotFoundError("plan not found")
	}

	if plan.User.ID != userID {
		return model.ForbiddenError("user does not own this plan")
	}
	return nil
}

func getOne(id uuid.UUID) (model.Plan, error) {
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

func RemoveFromOrder(tx *sqlx.Tx, userID, id uuid.UUID) (int64, error) {
	query := `
		UPDATE plans SET sort_order = sort_order - 1
		WHERE user_id = :user_id
		AND type = (SELECT type FROM plans WHERE id = :id)
		AND sort_order > (SELECT sort_order FROM plans WHERE id = :id)`
	params := model.Param{"user_id": userID, "id": id}
	return dbs.ExecTx(tx, query, params)
}
