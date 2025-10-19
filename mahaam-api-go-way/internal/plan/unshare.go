package plan

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/user"

	"github.com/google/uuid"
)

func Unshare(userID uuid.UUID, planID uuid.UUID, email string) *model.Err {
	// Validate user owns the plan
	if err := validateUserOwnsThePlan(userID, planID); err != nil {
		return model.ForbiddenError("user does not own this plan")
	}

	// Get user by email
	user, err := user.GetUserByEmail(email)
	if err != nil {
		return model.ServerError("failed to get user by email: " + err.Error())
	}
	if user == nil {
		return model.NotFoundError("email not found")
	}

	// Delete plan member (using function from leave.go)
	_, error := deletePlanMember(planID, user.ID)
	if error != nil {
		return model.ServerError("failed to delete plan member: " + error.Error())
	}

	return nil
}
