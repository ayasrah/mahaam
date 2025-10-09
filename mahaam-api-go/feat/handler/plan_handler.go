package handler

import (
	"mahaam-api/feat/models"
	"mahaam-api/feat/service"
	logs "mahaam-api/infra/log"
	"mahaam-api/infra/req"
	"strings"

	"github.com/google/uuid"
)

type PlanHandler interface {
	Create(ctx Ctx)
	Update(ctx Ctx)
	Delete(ctx Ctx)
	Share(ctx Ctx)
	Unshare(ctx Ctx)
	Leave(ctx Ctx)
	UpdateType(ctx Ctx)
	ReOrder(ctx Ctx)
	GetOne(ctx Ctx)
	GetMany(ctx Ctx)
}

type planHandler struct {
	planService service.PlanService
}

func NewPlanHandler(service service.PlanService) PlanHandler {
	return &planHandler{planService: service}
}

func RegisterPlanHandler(router Router, h PlanHandler) {
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

func (h *planHandler) Create(ctx Ctx) {
	var plan models.PlanIn
	req.Parse(ctx, &plan)
	if plan.Title == nil && plan.Starts == nil && plan.Ends == nil {
		panic(models.InputErr("At least one of title, starts, or ends is required"))
	}
	meta := ExtractMeta(ctx)
	id := h.planService.Create(meta.UserID, plan)
	ctx.JSON(Created, id)
}

func (h *planHandler) Update(ctx Ctx) {
	var plan models.PlanIn
	req.Parse(ctx, &plan)
	req.Required(plan.ID, "id")

	if plan.Title == nil && plan.Starts == nil && plan.Ends == nil {
		panic(models.InputErr("At least one of title, starts, or ends is required"))
	}
	meta := ExtractMeta(ctx)
	h.planService.Update(meta.UserID, &plan)
	ctx.Status(OK)
}

func (h *planHandler) Delete(ctx Ctx) {
	id := req.PathUuid(ctx, "planId")
	meta := ExtractMeta(ctx)
	h.planService.Delete(meta.UserID, id)
	ctx.Status(NoContent)
}

func (h *planHandler) Share(ctx Ctx) {
	id := req.PathUuid(ctx, "planId")
	email := req.Form(ctx, "email")
	meta := ExtractMeta(ctx)
	h.planService.Share(meta.UserID, id, email)
	ctx.Status(OK)
}

func (h *planHandler) Unshare(ctx Ctx) {
	id := req.PathUuid(ctx, "planId")
	email := req.Form(ctx, "email")
	meta := ExtractMeta(ctx)
	h.planService.Unshare(meta.UserID, id, email)
	ctx.Status(OK)
}

func (h *planHandler) Leave(ctx Ctx) {
	id := req.PathUuid(ctx, "planId")
	meta := ExtractMeta(ctx)
	h.planService.Leave(meta.UserID, id)
	logs.Info(ExtractTrafficID(ctx), "user %s left plan %s", meta.UserID, id)
	ctx.Status(OK)
}

func (h *planHandler) UpdateType(ctx Ctx) {
	id := req.PathUuid(ctx, "planId")
	planType := req.Form(ctx, "type")
	ValidatePlanType(planType)
	meta := ExtractMeta(ctx)
	h.planService.UpdateType(meta.UserID, id, planType)
	ctx.Status(OK)
}

func (h *planHandler) ReOrder(ctx Ctx) {
	var input struct {
		Type     string `form:"type" binding:"gte=0"`
		OldIndex int    `form:"oldOrder" binding:"gte=0"`
		NewIndex int    `form:"newOrder" binding:"required"`
	}
	req.Parse(ctx, &input)

	// planType := req.Form(c, "type")
	// oldOrder := req.FormInt(c, "oldOrder")
	// newOrder := req.FormInt(c, "newOrder")

	ValidatePlanType(input.Type)
	meta := ExtractMeta(ctx)
	h.planService.ReOrder(meta.UserID, input.Type, input.OldIndex, input.NewIndex)
	ctx.Status(OK)
}

func (h *planHandler) GetOne(ctx Ctx) {
	id := req.PathUuid(ctx, "planId")
	plan := h.planService.GetOne(id)
	if plan.ID == uuid.Nil {
		ctx.JSON(NotFound, res{"error": "Plan not found"})
		return
	}
	ctx.JSON(OK, plan)
}

func (h *planHandler) GetMany(ctx Ctx) {
	planType := req.Query(ctx, "type")
	ValidatePlanType(planType)
	meta := ExtractMeta(ctx)
	plans := h.planService.GetMany(meta.UserID, planType)
	ctx.JSON(OK, plans)
}

func ValidatePlanType(t string) {
	validTypes := []string{string(models.PlanTypeMain), string(models.PlanTypeArchived)} // Add other types if defined
	for _, vt := range validTypes {
		if strings.EqualFold(t, vt) {
			return
		}
	}
	panic(models.InputErr("Invalid plan type"))
}
