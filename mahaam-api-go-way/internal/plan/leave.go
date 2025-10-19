package plan

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
)

func Leave(userID uuid.UUID, planID uuid.UUID) *model.Err {
	rowsDeleted, err := deletePlanMember(planID, userID)
	if err != nil || rowsDeleted != 1 {
		return model.ServerError("failed to leave plan: " + err.Error())
	}
	return nil
}

func deletePlanMember(planID, userID uuid.UUID) (int64, error) {
	query := `
		DELETE FROM plan_members
		WHERE plan_id = :plan_id AND user_id = :user_id`
	params := model.Param{"plan_id": planID, "user_id": userID}
	return dbs.Exec(query, params)
}
