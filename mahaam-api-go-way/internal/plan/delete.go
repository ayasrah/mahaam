package plan

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

func Delete(userID uuid.UUID, planID uuid.UUID) *model.Err {
	if err := validateUserOwnsThePlan(userID, planID); err != nil {
		return model.ForbiddenError("user does not own this plan")
	}
	err := dbs.WithTx(func(tx *sqlx.Tx) error {
		if _, err := RemoveFromOrder(tx, userID, planID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return model.ServerError("failed to delete plan: " + err.Error())
	}
	return nil
}
