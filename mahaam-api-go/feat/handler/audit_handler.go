package handler

import (
	"mahaam-api/feat/repo"
	logs "mahaam-api/infra/log"
	"mahaam-api/infra/req"
)

type AuditHandler interface {
	Error(c Ctx)
	Trace(c Ctx)
}

type auditHandler struct {
	logRepo repo.LogRepo
}

func NewAuditHandler(logRepo repo.LogRepo) AuditHandler {
	return &auditHandler{logRepo: logRepo}
}

func RegisterAuditHandler(router Router, h AuditHandler) {
	rg := router.Group("/audit")
	rg.POST("/error", h.Error)
	rg.POST("/info", h.Trace)
}

func (h *auditHandler) Error(c Ctx) {
	errorVal := req.Form(c, "error")
	logs.Error(ExtractTrafficID(c), errorVal)
	c.Status(Created)
}

func (r *auditHandler) Trace(c Ctx) {
	info := req.Form(c, "info")
	logs.Info(ExtractTrafficID(c), "mahaam-mb: "+info)
	c.Status(Created)
}
