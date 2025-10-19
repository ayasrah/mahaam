package plan

import (
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"
	"mahaam-api/internal/user"

	"github.com/google/uuid"
)

const (
	shareLimit = 20
)

func Share(userID uuid.UUID, planID uuid.UUID, email string) *model.Err {
	// Validate user owns the plan
	if err := validateUserOwnsThePlan(userID, planID); err != nil {
		return model.ForbiddenError("user does not own this plan")
	}

	// Get user by email
	userInfo, err := user.GetUserByEmail(email)
	if err != nil || userInfo == nil {
		return model.ServerError("failed to get user by email: " + err.Error())
	}

	// Check if trying to share with creator
	if userInfo.ID == userID {
		return model.LogicError("not allowed to share with creator", "not_allowed_to_share_with_creator")
	}

	// Check limits
	plan, err := getOne(planID)
	if err != nil {
		return model.ServerError("failed to get plan: " + err.Error())
	}

	if plan.IsShared {
		count, err := getUsersCount(planID)
		if err != nil {
			return model.ServerError("failed to get users count: " + err.Error())
		}
		if count >= shareLimit {
			return model.LogicError("maximum of 20 shares reached", "max_is_20")
		}
	} else {
		count, err := getSharedPlansCount(userID)
		if err != nil {
			return model.ServerError("failed to get shared plans count: " + err.Error())
		}
		if count >= shareLimit {
			return model.LogicError("maximum of 20 shares reached", "max_is_20")
		}
	}

	// Create plan member
	if err := createPlanMember(planID, userInfo.ID); err != nil {
		return model.ServerError("failed to create plan member: " + err.Error())
	}

	// Add suggested emails (no transaction needed)
	createSuggestedEmail(userID, email)

	creator, err := user.GetUser(userID)
	if creator != nil && creator.Email != nil {
		createSuggestedEmail(userInfo.ID, *creator.Email)
	}

	return nil
}

// func getUserByEmail(email string) (*model.User, error) {
// 	query := `SELECT id, name, email FROM users WHERE email = :email`
// 	params := model.Param{"email": email}
// 	user, err := dbs.SelectOne[model.User](query, params)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if user.ID == uuid.Nil {
// 		return nil, nil
// 	}
// 	return &user, nil
// }

// func getUserByID(id uuid.UUID) (*model.User, error) {
// 	query := `SELECT id, name, email FROM users WHERE id = :id`
// 	params := model.Param{"id": id}
// 	user, err := dbs.SelectOne[model.User](query, params)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if user.ID == uuid.Nil {
// 		return nil, nil
// 	}
// 	return &user, nil
// }

func getUsersCount(planID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(1) FROM plan_members WHERE plan_id = :plan_id`
	params := model.Param{"plan_id": planID}
	return dbs.SelectOne[int64](query, params)
}

func getSharedPlansCount(userID uuid.UUID) (int64, error) {
	query := `
		SELECT COUNT(1)
		FROM plan_members cm
		LEFT JOIN plans c ON cm.plan_id = c.id
		WHERE c.user_id = :user_id`
	params := model.Param{"user_id": userID}
	return dbs.SelectOne[int64](query, params)
}

func createPlanMember(planID, userID uuid.UUID) error {
	query := `
		INSERT INTO plan_members (plan_id, user_id, created_at)
        VALUES (:plan_id, :user_id, current_timestamp)`
	params := model.Param{"plan_id": planID, "user_id": userID}
	_, err := dbs.Exec(query, params)
	return err
}

func createSuggestedEmail(userID uuid.UUID, email string) {
	query := `
		INSERT INTO suggested_emails (id, user_id, email, created_at)
		VALUES (:id, :user_id, :email, current_timestamp)
		ON CONFLICT (user_id, email) DO NOTHING`
	id := uuid.New()
	params := model.Param{"id": id, "user_id": userID, "email": email}
	dbs.Exec(query, params)
}
