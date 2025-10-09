package service

import (
	"fmt"
	"mahaam-api/feat/models"
	"mahaam-api/feat/repo"
	"mahaam-api/infra/dbs"
	"slices"

	"github.com/jmoiron/sqlx"
)

type TaskService interface {
	Create(planID UUID, title string) UUID
	GetList(planID UUID) []Task
	Delete(planID, id UUID)
	UpdateDone(planID, id UUID, done bool)
	UpdateTitle(id UUID, title string)
	ReOrder(planID UUID, oldOrder, newOrder int)
}

type taskService struct {
	taskRepo repo.TaskRepo
	planRepo repo.PlanRepo
}

func NewTaskService(taskRepo repo.TaskRepo, planRepo repo.PlanRepo) TaskService {
	return &taskService{
		taskRepo: taskRepo,
		planRepo: planRepo,
	}
}

func (s *taskService) Create(planID UUID, title string) UUID {
	count := s.taskRepo.GetCount(planID)
	if count >= 100 {
		panic(models.LogicErr("maximum number of tasks is 100", "max_is_100"))
	}

	var id UUID
	err := dbs.WithTx(func(tx *sqlx.Tx) error {
		id = s.taskRepo.Create(tx, planID, title)
		s.planRepo.UpdateDonePercent(tx, planID)
		return nil
	})
	if err != nil {
		panic(models.LogicErr(err.Error(), "error_creating_task"))
	}
	return id
}

func (s *taskService) GetList(planID UUID) []Task {
	return s.taskRepo.GetAll(planID)
}

func (s *taskService) Delete(planID, id UUID) {
	txFunc := func(tx *sqlx.Tx) error {
		s.taskRepo.UpdateOrderBeforeDelete(tx, planID, id)
		s.taskRepo.DeleteOne(tx, id)
		s.planRepo.UpdateDonePercent(tx, planID)
		return nil
	}

	if err := dbs.WithTx(txFunc); err != nil {
		panic(models.LogicErr(err.Error(), "error_deleting_task"))
	}
}

func (s *taskService) UpdateDone(planID, id UUID, done bool) {
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
	dbs.WithTx(txFunc)
}

func (s *taskService) UpdateTitle(id UUID, title string) {
	s.taskRepo.UpdateTitle(id, title)
}

func (s *taskService) ReOrder(planID UUID, oldOrder, newOrder int) {
	txFunc := func(tx *sqlx.Tx) error {
		s.reOrderWithTx(planID, oldOrder, newOrder, tx)
		return nil
	}
	dbs.WithTx(txFunc)
}

func (s *taskService) reOrderWithTx(planID UUID, oldOrder, newOrder int, tx *sqlx.Tx) {
	if oldOrder == newOrder {
		return
	}
	count := s.taskRepo.GetCount(planID)
	if int64(oldOrder) > count || int64(newOrder) > count {
		panic(models.InputErr(fmt.Sprintf("oldOrder and newOrder should be less than %d", count)))
	}

	s.taskRepo.UpdateOrder(tx, planID, oldOrder, newOrder)
}
