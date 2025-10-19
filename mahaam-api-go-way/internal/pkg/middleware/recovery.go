package middleware

import (
	"fmt"
	"mahaam-api/internal/model"
	logs "mahaam-api/internal/pkg/log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				trafficID, ok := c.Value("trafficID").(uuid.UUID)
				if !ok || trafficID == uuid.Nil {
					trafficID = uuid.Nil
				}

				if e, ok := err.(*model.Err); ok {
					logs.Error(trafficID, e.Error())
					if e.Key == "" {
						c.JSON(e.Code, e.Message)
					} else {
						c.JSON(e.Code, gin.H{"error": e.Message, "key": e.Key})
					}
					c.Abort()
					return
				} else {
					logs.Error(trafficID, fmt.Sprintf("%v", err))
				}

				c.JSON(http.StatusInternalServerError, gin.H{"error": err})
				c.Abort()
			}
		}()
		c.Next()
	}
}
