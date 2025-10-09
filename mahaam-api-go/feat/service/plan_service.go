package service

import (
	"fmt"
	"mahaam-api/feat/models"
	"mahaam-api/feat/repo"
	"mahaam-api/infra/dbs"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PlanService interface {
	GetOne(planID UUID) *Plan
	GetMany(userID UUID, planType string) []Plan
	Create(userID UUID, plan PlanIn) UUID
	Update(userID UUID, plan *PlanIn)
	Delete(userID UUID, id UUID)
	Share(userID UUID, id UUID, email string)
	Unshare(userID UUID, id UUID, email string)
	Leave(userID UUID, id UUID)
	UpdateType(userID UUID, id UUID, planType string)
	ReOrder(userID UUID, planType string, oldOrder, newOrder int)
	ValidateUserOwnsThePlan(userID UUID, planID UUID)
}

type planService struct {
	planRepo            repo.PlanRepo
	planMembersRepo     repo.PlanMembersRepo
	userRepo            repo.UserRepo
	suggestedEmailsRepo repo.SuggestedEmailRepo
}

func NewPlanService(db *sqlx.DB, planRepo repo.PlanRepo, planMembersRepo repo.PlanMembersRepo, userRepo repo.UserRepo, suggestedEmailsRepo repo.SuggestedEmailRepo) PlanService {
	return &planService{
		planRepo:            planRepo,
		planMembersRepo:     planMembersRepo,
		userRepo:            userRepo,
		suggestedEmailsRepo: suggestedEmailsRepo,
	}
}

func (s *planService) GetOne(planID UUID) *Plan {
	plan := s.planRepo.GetOne(planID)
	if plan.IsShared {
		users := s.planMembersRepo.GetUsers(planID)
		plan.Members = users
	}
	return plan
}

func (s *planService) GetMany(userID UUID, planType string) []Plan {
	plans := s.planRepo.GetMany(userID, planType)
	sharedPlans := s.planMembersRepo.GetOtherPlans(userID)
	return append(plans, sharedPlans...)
}

func (s *planService) Create(userID UUID, plan PlanIn) UUID {
	count := s.planRepo.GetCount(userID, string(models.PlanTypeMain))
	if count >= 100 {
		panic(models.LogicErr("maximum of 100 plans reached", "max_is_100"))
	}

	var planID UUID
	err := dbs.WithTx(func(tx *sqlx.Tx) error {
		planID = s.planRepo.Create(tx, userID, plan)
		return nil
	})

	if err != nil {
		return uuid.Nil
	}
	return planID
}

func (s *planService) Update(userID UUID, plan *PlanIn) {
	s.ValidateUserOwnsThePlan(userID, plan.ID)
	s.planRepo.Update(plan)
}

func (s *planService) Delete(userID UUID, id UUID) {
	s.ValidateUserOwnsThePlan(userID, id)
	dbs.WithTx(func(tx *sqlx.Tx) error {
		s.planRepo.RemoveFromOrder(tx, userID, id)
		s.planRepo.Delete(tx, id)
		return nil
	})
}

func (s *planService) Share(userID UUID, id UUID, email string) {
	s.ValidateUserOwnsThePlan(userID, id)
	user := s.userRepo.GetOneByEmail(email)
	if user == nil {
		panic(models.NotFoundErr("email not found"))
	}
	if user.ID == userID {
		panic(models.LogicErr("not allowed to share with creator", "not_allowed_to_share_with_creator"))
	}

	const limit = 20
	plan := s.planRepo.GetOne(id)
	if plan.IsShared {
		count := s.planMembersRepo.GetUsersCount(id)
		if count >= limit {
			panic(models.LogicErr("maximum of 20 shares reached", "max_is_20"))
		}
	} else {
		count := s.planMembersRepo.GetPlansCount(userID)
		if count >= limit {
			panic(models.LogicErr("maximum of 20 shares reached", "max_is_20"))
		}
	}

	s.planMembersRepo.Create(id, user.ID)
	// No transaction needed for suggested emails, as it's just a suggestion
	s.suggestedEmailsRepo.Create(userID, email)
	creator := s.userRepo.GetOne(userID)
	if creator.Email != nil {
		s.suggestedEmailsRepo.Create(user.ID, *creator.Email)
	}
}

func (s *planService) Unshare(userID UUID, id UUID, email string) {
	s.ValidateUserOwnsThePlan(userID, id)
	user := s.userRepo.GetOneByEmail(email)
	if user == nil {
		panic(models.NotFoundErr("email not found"))
	}
	s.planMembersRepo.Delete(id, user.ID)
}

// Leave allows a user to leave a shared plan
func (s *planService) Leave(userID UUID, id UUID) {
	rows := s.planMembersRepo.Delete(id, userID)
	if rows != 1 {
		panic(models.LogicErr(fmt.Sprintf("user cannot leave plan: userId=%s, planId=%s", userID, id), "user_cannot_leave_plan"))
	}
}

func (s *planService) UpdateType(userID UUID, id UUID, planType string) {
	s.ValidateUserOwnsThePlan(userID, id)
	count := s.planRepo.GetCount(userID, planType)
	if count >= 100 {
		panic(models.LogicErr("maximum of 100 plans reached", "max_is_100"))
	}

	dbs.WithTx(func(tx *sqlx.Tx) error {
		s.planRepo.RemoveFromOrder(tx, userID, id)
		s.planRepo.UpdateType(tx, userID, id, planType)
		return nil
	})
}

func (s *planService) ReOrder(userID UUID, planType string, oldOrder, newOrder int) {
	count := s.planRepo.GetCount(userID, planType)
	if oldOrder > int(count) || newOrder > int(count) {
		panic(models.InputErr(fmt.Sprintf("oldOrder and newOrder should be less than %d", count)))
	}
	s.planRepo.UpdateOrder(userID, planType, oldOrder, newOrder)
}

func (s *planService) ValidateUserOwnsThePlan(userID UUID, planID UUID) {
	plan := s.planRepo.GetOne(planID)
	if plan == nil {
		panic(models.NotFoundErr("plan not found"))
	}

	if plan.User.ID != userID {
		panic(models.ForbiddenErr("user does not own this plan"))
	}
}
