package handler

import (
	"mahaam-api/internal/model"
	tasks "mahaam-api/internal/task"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterTasks(router *gin.RouterGroup) {
	taskRouter := router.Group("/plans/:planId/tasks")

	taskRouter.POST("", CreateTask)
	taskRouter.DELETE("/:taskId", DeleteTask)
	taskRouter.PATCH("/:taskId/done", UpdateTaskDone)
	taskRouter.PATCH("/:taskId/title", UpdateTaskTitle)
	taskRouter.PATCH("/reorder", ReOrderTasks)
	taskRouter.GET("", GetManyTasks)
}

func CreateTask(c *gin.Context) {
	var err *model.Err

	planID, err := PathUuid(c, "planId")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	title, err := Form(c, "title")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	id, err := tasks.Create(planID, title)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusCreated, id)
}

func DeleteTask(c *gin.Context) {
	var err *model.Err

	planID, err := PathUuid(c, "planId")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	taskID, err := PathUuid(c, "taskId")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	if err := tasks.Delete(planID, taskID); err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func UpdateTaskDone(c *gin.Context) {
	var err *model.Err

	planID, err := PathUuid(c, "planId")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	taskID, err := PathUuid(c, "taskId")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	done, err := FormBool(c, "done")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	if err := tasks.UpdateDone(planID, taskID, done); err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.Status(http.StatusOK)
}

func UpdateTaskTitle(c *gin.Context) {
	var err *model.Err

	taskID, err := PathUuid(c, "taskId")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	title, err := Form(c, "title")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	if err := tasks.UpdateTitle(taskID, title); err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.Status(http.StatusOK)
}

func ReOrderTasks(c *gin.Context) {
	var input struct {
		OldOrder int `form:"oldOrder" binding:"gte=0"`
		NewOrder int `form:"newOrder" binding:"gte=0"`
	}
	var err *model.Err

	planID, err := PathUuid(c, "planId")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	if err := Parse(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := tasks.ReOrder(planID, input.OldOrder, input.NewOrder); err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.Status(http.StatusOK)
}

func GetManyTasks(c *gin.Context) {
	var err *model.Err

	planID, err := PathUuid(c, "planId")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	tasksList, err := tasks.GetMany(planID)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusOK, tasksList)
}
