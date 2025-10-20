package service

import (
	"fmt"
	"mahaam-api/app/models"
	"mahaam-api/app/repo"
	"slices"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TaskService interface {
	Create(planID uuid.UUID, title string) uuid.UUID
	GetList(planID uuid.UUID) []Task
	Delete(planID, id uuid.UUID)
	UpdateDone(planID, id uuid.UUID, done bool)
	UpdateTitle(id uuid.UUID, title string)
	ReOrder(planID uuid.UUID, oldOrder, newOrder int)
}

type taskService struct {
	taskRepo repo.TaskRepo
	planRepo repo.PlanRepo
	db       *repo.AppDB
}

func NewTaskService(db *repo.AppDB, taskRepo repo.TaskRepo, planRepo repo.PlanRepo) TaskService {
	return &taskService{
		db:       db,
		taskRepo: taskRepo,
		planRepo: planRepo,
	}
}

const maxTasksLimit = 100

func (s *taskService) Create(planID uuid.UUID, title string) uuid.UUID {
	tasksCount := s.taskRepo.GetCount(planID)
	if tasksCount >= maxTasksLimit {
		panic(models.LogicError("maximum tasks limit reached", "max_tasks_limit_reached"))

	}

	var id uuid.UUID
	err := repo.WithTransaction(s.db, func(tx *sqlx.Tx) error {
		id = s.taskRepo.Create(tx, planID, title)
		s.planRepo.UpdateDonePercent(tx, planID)
		return nil
	})
	if err != nil {
		panic(models.LogicError(err.Error(), "error_creating_task"))
	}
	return id
}

func (s *taskService) GetList(planID uuid.UUID) []Task {
	return s.taskRepo.GetAll(planID)
}

func (s *taskService) Delete(planID, id uuid.UUID) {
	txFunc := func(tx *sqlx.Tx) error {
		s.taskRepo.UpdateOrderBeforeDelete(tx, planID, id)
		s.taskRepo.DeleteOne(tx, id)
		s.planRepo.UpdateDonePercent(tx, planID)
		return nil
	}

	if err := repo.WithTransaction(s.db, txFunc); err != nil {
		panic(models.LogicError(err.Error(), "error_deleting_task"))
	}
}

func (s *taskService) UpdateDone(planID, id uuid.UUID, done bool) {
	txFunc := func(tx *sqlx.Tx) error {
		s.taskRepo.UpdateDone(tx, id, done)
		s.planRepo.UpdateDonePercent(tx, planID)
		tasks := s.taskRepo.GetAll(planID)
		taskIndex := slices.IndexFunc(tasks, func(t Task) bool {
			return t.ID == id
		})
		newOrder := 0
		if done {
			newOrder = len(tasks) - 1
		}
		s.reOrderWithTx(planID, taskIndex, newOrder, tx)
		return nil
	}
	repo.WithTransaction(s.db, txFunc)
}

func (s *taskService) UpdateTitle(id uuid.UUID, title string) {
	s.taskRepo.UpdateTitle(id, title)
}

func (s *taskService) ReOrder(planID uuid.UUID, oldOrder, newOrder int) {
	txFunc := func(tx *sqlx.Tx) error {
		s.reOrderWithTx(planID, oldOrder, newOrder, tx)
		return nil
	}
	repo.WithTransaction(s.db, txFunc)
}

func (s *taskService) reOrderWithTx(planID uuid.UUID, oldOrder, newOrder int, tx *sqlx.Tx) {
	if oldOrder == newOrder {
		return
	}
	count := s.taskRepo.GetCount(planID)
	if int64(oldOrder) > count || int64(newOrder) > count {
		panic(models.InputError(fmt.Sprintf("oldOrder and newOrder should be less than %d", count)))
	}

	s.taskRepo.UpdateOrder(tx, planID, oldOrder, newOrder)
}
