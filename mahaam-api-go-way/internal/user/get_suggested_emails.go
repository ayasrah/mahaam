package user

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
)

func GetSuggestedEmails(userID uuid.UUID) ([]model.SuggestedEmail, *model.Err) {
	suggestedEmails, err := getSuggestedEmailsByUserID(userID)
	if err != nil {
		return nil, model.ServerError("failed to get suggested emails: " + err.Error())
	}
	return suggestedEmails, nil
}

func getSuggestedEmailsByUserID(userID uuid.UUID) ([]model.SuggestedEmail, error) {
	query := `
		SELECT id, user_id, email, created_at
		FROM suggested_emails
		WHERE user_id = :user_id
		ORDER BY created_at DESC`
	params := model.Param{"user_id": userID}
	return dbs.SelectMany[model.SuggestedEmail](query, params)
}
