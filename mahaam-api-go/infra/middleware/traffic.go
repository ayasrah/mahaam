package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"mahaam-api/feat/models"
	"mahaam-api/feat/repo"
	"mahaam-api/infra/cache"
	"mahaam-api/infra/configs"
	logs "mahaam-api/infra/log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TrafficMiddleware(trafficRepo repo.TrafficRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		trafficId := uuid.New()
		c.Set("trafficID", trafficId)

		var payload []byte

		if c.Request.Body != nil && configs.LogReqEnabled {
			payload, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewReader(payload))
		}

		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = writer

		// Process request
		defer func() {
			// Log traffic (skip for swagger, health, audit paths)
			path := c.Request.URL.Path
			path = strings.ReplaceAll(path, "/mahaam-api", "")
			if !strings.HasPrefix(path, "/swagger") && !strings.EqualFold(path, "/health") && !strings.HasPrefix(path, "/audit") {

				code := c.Writer.Status()
				var reqBody, resBody string

				if code >= 400 {
					if configs.LogReqEnabled {
						reqBody = string(payload)
					}
					resBody = writer.body.String()
				}

				if strings.HasPrefix(path, "/user") && code < 400 {
					resBody = ""
				}

				elapsed := time.Since(start).Milliseconds()

				userId, _ := c.Value("userId").(uuid.UUID)
				deviceId, _ := c.Value("deviceId").(uuid.UUID)
				// Create traffic headers JSON
				headers := models.TrafficHeaders{
					UserID:     userId,
					DeviceID:   deviceId,
					AppVersion: c.GetHeader("x-app-version"),
					AppStore:   c.GetHeader("x-app-store"),
				}
				headersJSON, _ := json.Marshal(headers)
				headersStr := string(headersJSON)

				traffic := models.Traffic{
					ID:       trafficId,
					HealthID: cache.HealthID,
					Method:   c.Request.Method,
					Path:     path + c.Request.URL.RawQuery,
					Code:     code,
					Elapsed:  elapsed,
					Headers:  headersStr,
					Request:  reqBody,
					Response: resBody,
				}
				go func() {
					// Note: recovery is needed here as defer in recovery middleware will not be called if a panic occurs here.
					// This is because the goroutine where the panic occurred is different from the one where the defer statement was executed.
					// Panic recovery only works within the same goroutine where the defer statement was executed.

					defer func() {
						if r := recover(); r != nil {
							logs.Error(trafficId, "Traffic creation failed: %v", r)
						}
					}()
					trafficRepo.Create(traffic)
				}()
			}
		}()

		// Call next handler
		c.Next()

	}
}
