package handler

import (
	"mahaam-api/internal/model"
	plans "mahaam-api/internal/plan"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterPlans(router *gin.RouterGroup) {
	planRouter := router.Group("/plans")

	planRouter.POST("", Create)
	planRouter.PUT("", Update)
	planRouter.DELETE("/:planId", Delete)
	planRouter.PATCH("/:planId/share", Share)
	planRouter.PATCH("/:planId/unshare", Unshare)
	planRouter.PATCH("/:planId/leave", Leave)
	planRouter.PATCH("/:planId/type", UpdateType)
	planRouter.PATCH("/reorder", ReOrder)
	planRouter.GET("/:planId", GetOne)
	planRouter.GET("", GetMany)
}

func Create(c *gin.Context) {
	var plan model.PlanIn
	var userId uuid.UUID
	var err *model.Err
	var id uuid.UUID

	if err := Parse(c, &plan); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if userId, err = ExtractUserID(c); err != nil {
		c.JSON(err.Code, err)
		return
	}
	if id, err = plans.Create(userId, plan); err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusCreated, id)
}

func Update(c *gin.Context) {
	var plan model.PlanIn
	var userId uuid.UUID
	var err *model.Err

	if err := Parse(c, &plan); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	if userId, err = ExtractUserID(c); err != nil {
		c.JSON(err.Code, err)
		return
	}

	if err := Required(plan.ID, "id"); err != nil {
		c.JSON(err.Code, err)
		return
	}

	if err := plans.Update(userId, plan); err != nil {
		c.JSON(err.Code, err)
		return
	}
	c.Status(http.StatusOK)
}

func Delete(c *gin.Context) {
	var userId uuid.UUID
	var err *model.Err

	id, err := PathUuid(c, "planId")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	if userId, err = ExtractUserID(c); err != nil {
		c.JSON(err.Code, err)
		return
	}

	if err := plans.Delete(userId, id); err != nil {
		c.JSON(err.Code, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func Share(c *gin.Context) {
	var userId uuid.UUID
	var err *model.Err

	id, err := PathUuid(c, "planId")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	email, err := Form(c, "email")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	if userId, err = ExtractUserID(c); err != nil {
		c.JSON(err.Code, err)
		return
	}

	if err := plans.Share(userId, id, email); err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.Status(http.StatusOK)
}

func Unshare(c *gin.Context) {
	var userId uuid.UUID
	var err *model.Err

	id, err := PathUuid(c, "planId")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	email, err := Form(c, "email")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	if userId, err = ExtractUserID(c); err != nil {
		c.JSON(err.Code, err)
		return
	}

	if err := plans.Unshare(userId, id, email); err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.Status(http.StatusOK)
}

func Leave(c *gin.Context) {
	var userId uuid.UUID
	var err *model.Err

	id, err := PathUuid(c, "planId")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	if userId, err = ExtractUserID(c); err != nil {
		c.JSON(err.Code, err)
		return
	}

	if err := plans.Leave(userId, id); err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.Status(http.StatusOK)
}

func UpdateType(c *gin.Context) {
	var userId uuid.UUID
	var err *model.Err

	id, err := PathUuid(c, "planId")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	planType, err := Form(c, "type")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	if err := validatePlanType(planType); err != nil {
		c.JSON(err.Code, err)
		return
	}

	if userId, err = ExtractUserID(c); err != nil {
		c.JSON(err.Code, err)
		return
	}

	if err := plans.UpdateType(userId, id, planType); err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.Status(http.StatusOK)
}

func ReOrder(c *gin.Context) {
	var input struct {
		Type     string `form:"type" binding:"required"`
		OldOrder int    `form:"oldOrder" binding:"gte=0"`
		NewOrder int    `form:"newOrder" binding:"gte=0"`
	}
	var userId uuid.UUID
	var err *model.Err

	if err := Parse(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if err := validatePlanType(input.Type); err != nil {
		c.JSON(err.Code, err)
		return
	}

	if userId, err = ExtractUserID(c); err != nil {
		c.JSON(err.Code, err)
		return
	}

	if err := plans.ReOrder(userId, input.Type, input.OldOrder, input.NewOrder); err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.Status(http.StatusOK)
}

func GetOne(c *gin.Context) {
	id, err := PathUuid(c, "planId")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}
	var plan model.Plan
	if plan, err = plans.GetOne(id); err != nil {
		c.JSON(err.Code, err)
		return
	}
	c.JSON(http.StatusOK, plan)
}

func GetMany(c *gin.Context) {
	planType, err := Query(c, "type")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	if err := validatePlanType(planType); err != nil {
		c.JSON(err.Code, err)
		return
	}

	var userId uuid.UUID
	if userId, err = ExtractUserID(c); err != nil {
		c.JSON(err.Code, err)
		return
	}

	plans, err := plans.GetMany(userId, planType)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}
	c.JSON(http.StatusOK, plans)
}

func validatePlanType(t string) *model.Err {
	validTypes := []model.PlanType{model.PlanTypeMain, model.PlanTypeArchived}
	for _, vt := range validTypes {
		if strings.EqualFold(t, string(vt)) {
			return nil
		}
	}
	return model.InputError("Invalid plan type")
}
