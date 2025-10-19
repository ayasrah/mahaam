package task

import (
	"fmt"
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func ReOrder(planID uuid.UUID, oldOrder, newOrder int) *model.Err {
	if oldOrder == newOrder {
		return nil
	}

	err := dbs.WithTx(func(tx *sqlx.Tx) error {
		return reOrderWithTx(tx, planID, oldOrder, newOrder)
	})

	if err != nil {
		return model.ServerError("failed to reorder tasks: " + err.Error())
	}

	return nil
}

func reOrderWithTx(tx *sqlx.Tx, planID uuid.UUID, oldOrder, newOrder int) error {
	count, err := getTasksCount(planID)
	if err != nil {
		return err
	}

	if int64(oldOrder) > count || int64(newOrder) > count {
		return fmt.Errorf("oldOrder and newOrder should be less than %d", count)
	}

	query := `
		UPDATE tasks
		SET sort_order = CASE
			WHEN sort_order = :old_index THEN :new_index
			WHEN sort_order > :old_index AND sort_order <= :new_index THEN sort_order - 1
			WHEN sort_order >= :new_index AND sort_order < :old_index THEN sort_order + 1
			ELSE sort_order
		END
		WHERE plan_id = :plan_id`
	params := model.Param{
		"old_index": oldOrder,
		"new_index": newOrder,
		"plan_id":   planID,
	}
	_, err = dbs.ExecTx(tx, query, params)
	return err
}
