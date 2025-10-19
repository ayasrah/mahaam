package plan

import (
	"fmt"
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func UpdateType(userID uuid.UUID, planID uuid.UUID, planType string) *model.Err {
	// Validate user owns the plan
	if err := validateUserOwnsThePlan(userID, planID); err != nil {
		return model.ForbiddenError("user does not own this plan")
	}

	// Check plan count for the target type
	count, err := getPlansCount(userID, model.PlanType(planType))
	if err != nil {
		return model.ServerError("failed to get plans count")
	}
	if count >= plansLimit {
		errMessage := fmt.Sprintf("maximum of %d plans reached", plansLimit)
		return model.LogicError(errMessage, "max_is_100")
	}

	// Update plan type with transaction
	err = dbs.WithTx(func(tx *sqlx.Tx) error {
		if _, err := RemoveFromOrder(tx, userID, planID); err != nil {
			return err
		}
		if err := updatePlanType(tx, userID, planID, planType); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return model.ServerError("failed to update plan type: " + err.Error())
	}

	return nil
}

func updatePlanType(tx *sqlx.Tx, userID, planID uuid.UUID, planType string) error {
	query := `UPDATE plans SET type = :type, 
		sort_order = (SELECT COUNT(1) FROM plans WHERE user_id = :user_id AND type = :type), 
		updated_at = current_timestamp WHERE id = :id`
	params := model.Param{"id": planID, "type": planType, "user_id": userID}
	_, err := dbs.ExecTx(tx, query, params)
	return err
}
