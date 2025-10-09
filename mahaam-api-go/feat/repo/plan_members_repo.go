package repo

import (
	"mahaam-api/infra/dbs"
)

type PlanMembersRepo interface {
	Create(planID, userID UUID) int64
	Delete(planID, userID UUID) int64
	GetOtherPlans(userID UUID) []Plan
	GetUsers(planID UUID) []User
	GetPlansCount(userID UUID) int64
	GetUsersCount(planID UUID) int64
}

type planMembersRepo struct {
}

func NewPlanMembersRepo() PlanMembersRepo {
	return &planMembersRepo{}
}

func (r *planMembersRepo) Create(planID, userID UUID) int64 {
	query := `
		INSERT INTO plan_members (plan_id, user_id, created_at)
        VALUES (:plan_id, :user_id, current_timestamp)`
	params := Param{"plan_id": planID, "user_id": userID}
	return dbs.Exec(query, params)
}

func (r *planMembersRepo) Delete(planID, userID UUID) int64 {
	query := `
		DELETE FROM plan_members
		WHERE plan_id = :plan_id AND user_id = :user_id`
	params := Param{"plan_id": planID, "user_id": userID}
	return dbs.Exec(query, params)
}

func (r *planMembersRepo) GetOtherPlans(userID UUID) []Plan {
	query := `
		SELECT c.id, c.title, c.starts, c.ends, c.type, c.done_percent, c.sort_order, 
			true AS is_shared, u.id as "user.id",u.email as "user.email",u.name as "user.name"
		FROM plan_members cm
		LEFT JOIN plans c ON cm.plan_id = c.id
		LEFT JOIN users u ON c.user_id = u.id
		WHERE cm.user_id = :user_id
		ORDER BY c.sort_order ASC`
	params := Param{"user_id": userID}
	return dbs.SelectMany[Plan](query, params)
}

func (r *planMembersRepo) GetUsers(planID UUID) []User {
	query := `
		SELECT u.id, u.email, u.name
		FROM plan_members cm
		LEFT JOIN users u ON cm.user_id = u.id
		WHERE cm.plan_id = :plan_id
		ORDER BY u.created_at DESC`
	params := Param{"plan_id": planID}
	return dbs.SelectMany[User](query, params)
}

func (r *planMembersRepo) GetPlansCount(userID UUID) int64 {
	query := `
		SELECT COUNT(1)
		FROM plan_members cm
		LEFT JOIN plans c ON cm.plan_id = c.id
		WHERE c.user_id = :user_id`
	param := Param{"user_id": userID}
	return dbs.SelectOne[int64](query, param)
}

func (r *planMembersRepo) GetUsersCount(planID UUID) int64 {
	query := `
		SELECT COUNT(1)
		FROM plan_members
		WHERE plan_id = :plan_id`
	param := Param{"plan_id": planID}
	return dbs.SelectOne[int64](query, param)

}
