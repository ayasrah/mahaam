package handler

import (
	"mahaam-api/app/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskHandler interface {
	Create(c *gin.Context)
	Delete(c *gin.Context)
	UpdateDone(c *gin.Context)
	UpdateTitle(c *gin.Context)
	ReOrder(c *gin.Context)
	GetMany(c *gin.Context)
}

type taskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) TaskHandler {
	return &taskHandler{taskService: taskService}
}

func RegisterTaskHandler(router *gin.RouterGroup, h TaskHandler) {
	taskRouter := router.Group("/plans/:planId/tasks", ValidatePlanId)

	taskRouter.POST("/", h.Create)
	taskRouter.DELETE("/:taskId", h.Delete)
	taskRouter.PATCH("/:taskId/done", h.UpdateDone)
	taskRouter.PATCH("/:taskId/title", h.UpdateTitle)
	taskRouter.PATCH("/reorder", h.ReOrder)
	taskRouter.GET("", h.GetMany)
}

func ValidatePlanId(c *gin.Context) {
	parsePathUuid(c, "planId")
}

func (h *taskHandler) Create(c *gin.Context) {
	planID := parsePathUuid(c, "planId")
	title := parseFormParam(c, "title")
	id := h.taskService.Create(planID, title)
	c.JSON(http.StatusCreated, id)
}

func (h *taskHandler) Delete(c *gin.Context) {
	planID := parsePathUuid(c, "planId")
	id := parsePathUuid(c, "taskId")
	h.taskService.Delete(planID, id)
	c.Status(http.StatusNoContent)
}

func (h *taskHandler) UpdateDone(c *gin.Context) {
	planID := parsePathUuid(c, "planId")
	id := parsePathUuid(c, "taskId")
	done := parseFormBool(c, "done")
	h.taskService.UpdateDone(planID, id, done)
	c.Status(http.StatusOK)
}

func (h *taskHandler) UpdateTitle(c *gin.Context) {
	id := parsePathUuid(c, "taskId")
	title := parseFormParam(c, "title")
	h.taskService.UpdateTitle(id, title)
	c.Status(http.StatusOK)
}

func (h *taskHandler) ReOrder(c *gin.Context) {
	planID := parsePathUuid(c, "planId")
	oldOrder := parseFormInt(c, "oldOrder")
	newOrder := parseFormInt(c, "newOrder")
	h.taskService.ReOrder(planID, oldOrder, newOrder)
	c.Status(http.StatusOK)
}

func (h *taskHandler) GetMany(c *gin.Context) {
	planID := parsePathUuid(c, "planId")
	tasks := h.taskService.GetList(planID)
	c.JSON(http.StatusOK, tasks)
}
