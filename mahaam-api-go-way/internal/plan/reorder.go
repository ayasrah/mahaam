package plan

import (
	"fmt"
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
)

func ReOrder(userID uuid.UUID, planType string, oldOrder, newOrder int) *model.Err {
	count, err := getPlansCount(userID, model.PlanType(planType))
	if err != nil {
		return model.ServerError("failed to get plans count")
	}

	if oldOrder > int(count) || newOrder > int(count) {
		errMessage := fmt.Sprintf("oldOrder and newOrder should be less than %d", count)
		return model.InputError(errMessage)
	}

	if err := updateOrder(userID, planType, oldOrder, newOrder); err != nil {
		return model.ServerError("failed to reorder plans: " + err.Error())
	}

	return nil
}

func updateOrder(userID uuid.UUID, planType string, oldOrder, newOrder int) error {
	query := `
		UPDATE plans SET sort_order = 
			CASE 
				WHEN sort_order = :oldOrder THEN :newOrder
				WHEN sort_order > :oldOrder AND sort_order <= :newOrder THEN sort_order - 1
				WHEN sort_order >= :newOrder AND sort_order < :oldOrder THEN sort_order + 1
				ELSE sort_order
			END
		WHERE user_id = :user_id AND type = :type`
	params := model.Param{
		"user_id":  userID,
		"type":     planType,
		"oldOrder": oldOrder,
		"newOrder": newOrder,
	}
	_, err := dbs.Exec(query, params)
	return err
}
