package handler

import (
	"mahaam-api/feat/service"
	"mahaam-api/infra/req"
)

type TaskHandler interface {
	Create(c Ctx)
	Delete(c Ctx)
	UpdateDone(c Ctx)
	UpdateTitle(c Ctx)
	ReOrder(c Ctx)
	GetMany(c Ctx)
}

type taskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) TaskHandler {
	return &taskHandler{taskService: taskService}
}

func RegisterTaskHandler(router Router, h TaskHandler) {
	taskRouter := router.Group("/plans/:planId/tasks", ValidatePlanId)

	taskRouter.POST("/", h.Create)
	taskRouter.DELETE("/:taskId", h.Delete)
	taskRouter.PATCH("/:taskId/done", h.UpdateDone)
	taskRouter.PATCH("/:taskId/title", h.UpdateTitle)
	taskRouter.PATCH("/reorder", h.ReOrder)
	taskRouter.GET("", h.GetMany)
}

func ValidatePlanId(c Ctx) {
	req.PathUuid(c, "planId")
}

func (h *taskHandler) Create(c Ctx) {
	planID := req.PathUuid(c, "planId")
	title := req.Form(c, "title")
	id := h.taskService.Create(planID, title)
	c.JSON(Created, id)
}

func (h *taskHandler) Delete(c Ctx) {
	planID := req.PathUuid(c, "planId")
	id := req.PathUuid(c, "taskId")
	h.taskService.Delete(planID, id)
	c.Status(NoContent)
}

func (h *taskHandler) UpdateDone(c Ctx) {
	planID := req.PathUuid(c, "planId")
	id := req.PathUuid(c, "taskId")
	done := req.FormBool(c, "done")
	h.taskService.UpdateDone(planID, id, done)
	c.Status(OK)
}

func (h *taskHandler) UpdateTitle(c Ctx) {
	id := req.PathUuid(c, "taskId")
	title := req.Form(c, "title")
	h.taskService.UpdateTitle(id, title)
	c.Status(OK)
}

func (h *taskHandler) ReOrder(c Ctx) {
	planID := req.PathUuid(c, "planId")
	oldOrder := req.FormInt(c, "oldOrder")
	newOrder := req.FormInt(c, "newOrder")
	h.taskService.ReOrder(planID, oldOrder, newOrder)
	c.Status(OK)
}

func (h *taskHandler) GetMany(c Ctx) {
	planID := req.PathUuid(c, "planId")
	tasks := h.taskService.GetList(planID)
	c.JSON(OK, tasks)
}
