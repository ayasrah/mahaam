package middleware

import (
	"bytes"
	"net/http"
	"strings"

	logs "mahaam-api/utils/log"
	token "mahaam-api/utils/token"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func AuthMiddleware(tokenService *token.TokenService, logger logs.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		trafficId := c.Value("trafficID").(uuid.UUID)
		// Validate path base
		pathBase := c.Request.URL.Path
		if !strings.HasPrefix(pathBase, "/mahaam-api") {
			c.AbortWithStatusJSON(http.StatusNotFound, "invalid path base")
			return
		}

		// Validate headers
		appStore := c.GetHeader("x-app-store")
		appVersion := c.GetHeader("x-app-version")
		if appStore == "" || appVersion == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Required headers not exists")
			logger.Error(trafficId, "Required headers not exists")
			return
		}

		// Check bypass paths
		path := c.Request.URL.Path
		path = strings.ReplaceAll(path, "/mahaam-api", "")
		bypassAuthPaths := []string{"/swagger", "/health", "/users/create", "/audit/info", "/audit/error"}
		requiresAuth := true
		for _, bypassPath := range bypassAuthPaths {
			if strings.HasPrefix(path, bypassPath) {
				requiresAuth = false
				break
			}
		}

		// Authenticate if required
		var userId, deviceId uuid.UUID
		if requiresAuth {
			var err error
			userId, deviceId, err = (*tokenService).Parse(c)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, "invalid jwt")
				logger.Error(trafficId, err.Error())
				return
			}
			c.Set("userId", userId)
			c.Set("deviceId", deviceId)
		}

		c.Next()

	}
}
