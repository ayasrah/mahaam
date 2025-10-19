package plan

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
)

func Update(userID uuid.UUID, plan model.PlanIn) *model.Err {
	// Validate user owns the plan
	if err := validateUserOwnsThePlan(userID, plan.ID); err != nil {
		return model.ForbiddenError("user does not own this plan")
	}

	// Update the plan
	if err := updatePlan(plan); err != nil {
		return model.ServerError("failed to update plan: " + err.Error())
	}

	return nil
}

func updatePlan(plan model.PlanIn) error {
	query := `UPDATE plans SET title = :title, starts = :starts, ends = :ends, updated_at = current_timestamp WHERE id = :id`
	_, err := dbs.Exec(query, plan)
	return err
}
