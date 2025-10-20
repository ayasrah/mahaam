package handler

import (
	"mahaam-api/app/models"
	"mahaam-api/app/service"
	logs "mahaam-api/utils/log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PlanHandler interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	Share(c *gin.Context)
	Unshare(c *gin.Context)
	Leave(c *gin.Context)
	UpdateType(c *gin.Context)
	ReOrder(c *gin.Context)
	GetOne(c *gin.Context)
	GetMany(c *gin.Context)
}

type planHandler struct {
	planService service.PlanService
	logger      logs.Logger
}

func NewPlanHandler(service service.PlanService, logger logs.Logger) PlanHandler {
	return &planHandler{planService: service, logger: logger}
}

func RegisterPlanHandler(router *gin.RouterGroup, h PlanHandler) {
	planRouter := router.Group("/plans")

	planRouter.POST("", h.Create)
	planRouter.PUT("", h.Update)
	planRouter.DELETE("/:planId", h.Delete)
	planRouter.PATCH("/:planId/share", h.Share)
	planRouter.PATCH("/:planId/unshare", h.Unshare)
	planRouter.PATCH("/:planId/leave", h.Leave)
	planRouter.PATCH("/:planId/type", h.UpdateType)
	planRouter.PATCH("/reorder", h.ReOrder)
	planRouter.GET("/:planId", h.GetOne)
	planRouter.GET("", h.GetMany)

}

func (h *planHandler) Create(c *gin.Context) {
	var plan PlanIn
	parse(c, &plan)
	if plan.Title == nil && plan.Starts == nil && plan.Ends == nil {
		panic(models.InputError("At least one of title, starts, or ends is required"))
	}
	meta := parseRequestMeta(c)
	id := h.planService.Create(meta.UserID, plan)
	c.JSON(http.StatusCreated, id)
}

func (h *planHandler) Update(c *gin.Context) {
	var plan PlanIn
	parse(c, &plan)
	requiredParam(plan.ID, "id")

	if plan.Title == nil && plan.Starts == nil && plan.Ends == nil {
		panic(models.InputError("At least one of title, starts, or ends is required"))
	}
	meta := parseRequestMeta(c)
	h.planService.Update(meta.UserID, &plan)
	c.Status(http.StatusOK)
}

func (h *planHandler) Delete(c *gin.Context) {
	id := parsePathUuid(c, "planId")
	meta := parseRequestMeta(c)
	h.planService.Delete(meta.UserID, id)
	c.Status(http.StatusNoContent)
}

func (h *planHandler) Share(c *gin.Context) {
	id := parsePathUuid(c, "planId")
	email := parseFormParam(c, "email")
	meta := parseRequestMeta(c)
	h.planService.Share(meta.UserID, id, email)
	c.Status(http.StatusOK)
}

func (h *planHandler) Unshare(c *gin.Context) {
	id := parsePathUuid(c, "planId")
	email := parseFormParam(c, "email")
	meta := parseRequestMeta(c)
	h.planService.Unshare(meta.UserID, id, email)
	c.Status(http.StatusOK)
}

func (h *planHandler) Leave(c *gin.Context) {
	id := parsePathUuid(c, "planId")
	meta := parseRequestMeta(c)
	h.planService.Leave(meta.UserID, id)
	h.logger.Info(parseTrafficID(c), "user %s left plan %s", meta.UserID, id)
	c.Status(http.StatusOK)
}

func (h *planHandler) UpdateType(c *gin.Context) {
	id := parsePathUuid(c, "planId")
	planType := parseFormParam(c, "type")
	validatePlanType(planType)
	meta := parseRequestMeta(c)
	h.planService.UpdateType(meta.UserID, id, planType)
	c.Status(http.StatusOK)
}

func (h *planHandler) ReOrder(c *gin.Context) {
	var input struct {
		Type     string `form:"type" binding:"gte=0"`
		OldIndex int    `form:"oldOrder" binding:"gte=0"`
		NewIndex int    `form:"newOrder" binding:"required"`
	}
	parse(c, &input)

	validatePlanType(input.Type)
	meta := parseRequestMeta(c)
	h.planService.ReOrder(meta.UserID, input.Type, input.OldIndex, input.NewIndex)
	c.Status(http.StatusOK)
}

func (h *planHandler) GetOne(c *gin.Context) {
	id := parsePathUuid(c, "planId")
	plan := h.planService.GetOne(id)
	if plan.ID == uuid.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		return
	}
	c.JSON(http.StatusOK, plan)
}

func (h *planHandler) GetMany(c *gin.Context) {
	planType := parseQueryParam(c, "type")
	validatePlanType(planType)
	meta := parseRequestMeta(c)
	plans := h.planService.GetMany(meta.UserID, planType)
	c.JSON(http.StatusOK, plans)
}

func validatePlanType(t string) {
	validTypes := []string{string(PlanTypeMain), string(PlanTypeArchived)} // Add other types if defined
	for _, vt := range validTypes {
		if strings.EqualFold(t, vt) {
			return
		}
	}
	panic(models.InputError("Invalid plan type"))
}
