package plan

import (
	"fmt"
	"mahaam-api/internal/model"
	"mahaam-api/internal/pkg/dbs"

	"github.com/google/uuid"
)

const (
	plansLimit = 100
)

func Create(userID uuid.UUID, plan model.PlanIn) (uuid.UUID, *model.Err) {
	plansCount, err := getPlansCount(userID, model.PlanTypeMain)
	if err != nil {
		return uuid.Nil, model.ServerError("failed to get plans count")
	}
	if plansCount >= plansLimit {
		errMessage := fmt.Sprintf("maximum of %d plans reached", plansLimit)
		return uuid.Nil, model.LogicError(errMessage, "plans_limit_reached")
	}

	id, err := create(userID, plan)
	if err != nil {
		return uuid.Nil, model.ServerError("failed to create plan: " + err.Error())
	}
	return id, nil
}

func getPlansCount(userID uuid.UUID, planType model.PlanType) (int64, error) {
	query := `SELECT COUNT(1) FROM plans WHERE user_id = :user_id AND type = :type`
	params := model.Param{"user_id": userID, "type": string(planType)}
	return dbs.SelectOne[int64](query, params)
}

func create(userID uuid.UUID, plan model.PlanIn) (uuid.UUID, error) {
	id := uuid.New()
	query := `
		INSERT INTO plans (id, user_id, title, starts, ends, type, status, done_percent, sort_order, created_at)
		VALUES (:id, :user_id, :title, :starts, :ends, :type, :status, '0/0', 
		(SELECT COUNT(1) FROM plans WHERE user_id = :user_id AND type = :type), current_timestamp)`
	params := model.Param{
		"id":      id,
		"user_id": userID,
		"title":   plan.Title,
		"starts":  plan.Starts,
		"ends":    plan.Ends,
		"type":    model.PlanTypeMain,
		"status":  "Open",
	}
	_, err := dbs.Exec(query, params)
	return id, err
}
