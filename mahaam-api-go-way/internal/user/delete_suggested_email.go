package user

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
)

func DeleteSuggestedEmail(userID, suggestedEmailID uuid.UUID) *model.Err {
	suggestedEmail, err := getSuggestedEmail(suggestedEmailID)
	if err != nil {
		return model.ServerError("failed to get suggested email: " + err.Error())
	}

	if suggestedEmail == nil || suggestedEmail.UserID != userID {
		return model.ForbiddenError("invalid suggestedEmailId")
	}

	if err := deleteSuggestedEmail(suggestedEmailID); err != nil {
		return model.ServerError("failed to delete suggested email: " + err.Error())
	}

	return nil
}

func getSuggestedEmail(suggestedEmailID uuid.UUID) (*model.SuggestedEmail, error) {
	query := `
		SELECT id, user_id, email, created_at
		FROM suggested_emails
		WHERE id = :id`
	params := model.Param{"id": suggestedEmailID}
	email, err := dbs.SelectOne[model.SuggestedEmail](query, params)
	if err != nil {
		return nil, err
	}
	if email.ID == uuid.Nil {
		return nil, nil
	}
	return &email, nil
}

func deleteSuggestedEmail(suggestedEmailID uuid.UUID) error {
	query := `DELETE FROM suggested_emails WHERE id = :id`
	params := model.Param{"id": suggestedEmailID}
	_, err := dbs.Exec(query, params)
	return err
}
