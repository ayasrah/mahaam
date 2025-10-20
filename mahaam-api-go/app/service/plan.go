package service

import (
	"fmt"
	"mahaam-api/app/models"
	"mahaam-api/app/repo"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PlanService interface {
	GetOne(planID uuid.UUID) *Plan
	GetMany(userID uuid.UUID, planType string) []Plan
	Create(userID uuid.UUID, plan PlanIn) uuid.UUID
	Update(userID uuid.UUID, plan *PlanIn)
	Delete(userID uuid.UUID, id uuid.UUID)
	Share(userID uuid.UUID, id uuid.UUID, email string)
	Unshare(userID uuid.UUID, id uuid.UUID, email string)
	Leave(userID uuid.UUID, id uuid.UUID)
	UpdateType(userID uuid.UUID, id uuid.UUID, planType string)
	ReOrder(userID uuid.UUID, planType string, oldOrder, newOrder int)
	ValidateUserOwnsThePlan(userID uuid.UUID, planID uuid.UUID)
}

type planService struct {
	planRepo            repo.PlanRepo
	planMembersRepo     repo.PlanMembersRepo
	userRepo            repo.UserRepo
	suggestedEmailsRepo repo.SuggestedEmailRepo
	db                  *repo.AppDB
}

func NewPlanService(
	db *repo.AppDB,
	planRepo repo.PlanRepo,
	planMembersRepo repo.PlanMembersRepo,
	userRepo repo.UserRepo,
	suggestedEmailsRepo repo.SuggestedEmailRepo) PlanService {

	return &planService{
		planRepo:            planRepo,
		planMembersRepo:     planMembersRepo,
		userRepo:            userRepo,
		suggestedEmailsRepo: suggestedEmailsRepo,
		db:                  db,
	}
}

func (s *planService) GetOne(planID uuid.UUID) *Plan {
	plan := s.planRepo.GetOne(planID)
	if plan.IsShared {
		users := s.planMembersRepo.GetUsers(planID)
		plan.Members = users
	}
	return plan
}

func (s *planService) GetMany(userID uuid.UUID, planType string) []Plan {
	plans := s.planRepo.GetMany(userID, planType)
	sharedPlans := s.planMembersRepo.GetOtherPlans(userID)
	return append(plans, sharedPlans...)
}

const plansLimit = 100

func (s *planService) Create(userID uuid.UUID, plan PlanIn) uuid.UUID {
	plansCount := s.planRepo.GetCount(userID, string(models.PlanTypeMain))
	if plansCount >= plansLimit {
		panic(models.LogicError("maximum plans limit reached", "max_plans_limit_reached"))
	}

	var planID uuid.UUID
	err := repo.WithTransaction(s.db, func(tx *sqlx.Tx) error {
		planID = s.planRepo.Create(tx, userID, plan)
		return nil
	})

	if err != nil {
		return uuid.Nil
	}
	return planID
}

func (s *planService) Update(userID uuid.UUID, plan *PlanIn) {
	s.ValidateUserOwnsThePlan(userID, plan.ID)
	s.planRepo.Update(plan)
}

func (s *planService) Delete(userID uuid.UUID, id uuid.UUID) {
	s.ValidateUserOwnsThePlan(userID, id)
	repo.WithTransaction(s.db, func(tx *sqlx.Tx) error {
		s.planRepo.RemoveFromOrder(tx, userID, id)
		s.planRepo.Delete(tx, id)
		return nil
	})
}

func (s *planService) Share(userID uuid.UUID, id uuid.UUID, email string) {
	s.ValidateUserOwnsThePlan(userID, id)
	user := s.userRepo.GetOneByEmail(email)
	if user == nil {
		panic(models.NotFoundError("email not found"))
	}
	if user.ID == userID {
		panic(models.LogicError("not allowed to share with creator", "not_allowed_to_share_with_creator"))
	}

	const sharedPlanUsersLimit = 20
	plan := s.planRepo.GetOne(id)
	if plan.IsShared {
		sharedUsersCount := s.planMembersRepo.GetUsersCount(id)
		if sharedUsersCount >= sharedPlanUsersLimit {
			panic(models.LogicError("maximum of 20 shares reached", "max_is_20"))
		}
	} else {
		count := s.planMembersRepo.GetPlansCount(userID)
		if count >= sharedPlanUsersLimit {
			panic(models.LogicError("maximum of 20 shares reached", "max_is_20"))
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

func (s *planService) Unshare(userID uuid.UUID, id uuid.UUID, email string) {
	s.ValidateUserOwnsThePlan(userID, id)
	user := s.userRepo.GetOneByEmail(email)
	if user == nil {
		panic(models.NotFoundError("email not found"))
	}
	s.planMembersRepo.Delete(id, user.ID)
}

// Leave allows a user to leave a shared plan
func (s *planService) Leave(userID uuid.UUID, id uuid.UUID) {
	rows := s.planMembersRepo.Delete(id, userID)
	if rows != 1 {
		panic(models.LogicError(fmt.Sprintf("user cannot leave plan: userId=%s, planId=%s", userID, id), "user_cannot_leave_plan"))
	}
}

func (s *planService) UpdateType(userID uuid.UUID, id uuid.UUID, planType string) {
	s.ValidateUserOwnsThePlan(userID, id)
	count := s.planRepo.GetCount(userID, planType)
	if count >= 100 {
		panic(models.LogicError("maximum of 100 plans reached", "max_is_100"))
	}

	repo.WithTransaction(s.db, func(tx *sqlx.Tx) error {
		s.planRepo.RemoveFromOrder(tx, userID, id)
		s.planRepo.UpdateType(tx, userID, id, planType)
		return nil
	})
}

func (s *planService) ReOrder(userID uuid.UUID, planType string, oldOrder, newOrder int) {
	count := s.planRepo.GetCount(userID, planType)
	if oldOrder > int(count) || newOrder > int(count) {
		panic(models.InputError(fmt.Sprintf("oldOrder and newOrder should be less than %d", count)))
	}
	s.planRepo.UpdateOrder(userID, planType, oldOrder, newOrder)
}

func (s *planService) ValidateUserOwnsThePlan(userID uuid.UUID, planID uuid.UUID) {
	plan := s.planRepo.GetOne(planID)
	if plan == nil {
		panic(models.NotFoundError("plan not found"))
	}

	if plan.User.ID != userID {
		panic(models.ForbiddenError("user does not own this plan"))
	}
}
