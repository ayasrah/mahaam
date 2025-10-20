package handler

import (
	logs "mahaam-api/utils/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuditHandler interface {
	Error(c *gin.Context)
	Trace(c *gin.Context)
}

type auditHandler struct {
	logger logs.Logger
}

func NewAuditHandler(logger logs.Logger) AuditHandler {
	return &auditHandler{logger: logger}
}

func RegisterAuditHandler(router *gin.RouterGroup, h AuditHandler) {
	rg := router.Group("/audit")
	rg.POST("/error", h.Error)
	rg.POST("/info", h.Trace)
}

func (h *auditHandler) Error(c *gin.Context) {
	errorVal := parseFormParam(c, "error")
	h.logger.Error(parseTrafficID(c), errorVal)
	c.Status(http.StatusCreated)
}

func (h *auditHandler) Trace(c *gin.Context) {
	info := parseFormParam(c, "info")
	h.logger.Info(parseTrafficID(c), "mahaam-mb: "+info)
	c.Status(http.StatusCreated)
}
