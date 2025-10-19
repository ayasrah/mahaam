package handler

import (
	"mahaam-api/internal/model"
	logs "mahaam-api/internal/pkg/log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterAudit(router *gin.RouterGroup) {
	auditRouter := router.Group("/audit")
	auditRouter.POST("/error", AuditError)
	auditRouter.POST("/info", AuditTrace)
}

func AuditError(c *gin.Context) {
	var err *model.Err

	errorVal, err := Form(c, "error")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	trafficID, err := ExtractTrafficID(c)
	if err != nil {
		trafficID = uuid.Nil
	}

	logs.Error(trafficID, "mahaam-mb: "+errorVal)

	c.Status(http.StatusCreated)
}

func AuditTrace(c *gin.Context) {
	var err *model.Err

	info, err := Form(c, "info")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	trafficID, err := ExtractTrafficID(c)
	if err != nil {
		trafficID = uuid.Nil
	}

	logs.Info(trafficID, "mahaam-mb: "+info)
	c.Status(http.StatusCreated)
}
