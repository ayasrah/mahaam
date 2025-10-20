package handler

import (
	"io"
	"mahaam-api/app/models"
	"mahaam-api/app/service"
	logs "mahaam-api/utils/log"
	"mahaam-api/utils/middleware"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	Create(c *gin.Context)
	SendMeOtp(c *gin.Context)
	VerifyOtp(c *gin.Context)
	RefreshToken(c *gin.Context)
	UpdateName(c *gin.Context)
	Logout(c *gin.Context)
	Delete(c *gin.Context)
	GetDevices(c *gin.Context)
	GetSuggestedEmails(c *gin.Context)
	DeleteSuggestedEmail(c *gin.Context)
}

type userHandler struct {
	userService service.UserService
	logger      logs.Logger
}

func NewUserHandler(userService service.UserService, logger logs.Logger) UserHandler {
	return &userHandler{userService: userService, logger: logger}
}

func RegisterUserHandler(router *gin.RouterGroup, h UserHandler) {
	rg := router.Group("/users")

	// 5 requests per minute per IP
	createLimiter := middleware.RateLimiterMW()

	rg.POST("/create", createLimiter, h.Create)
	rg.POST("/send-me-otp", createLimiter, h.SendMeOtp)
	rg.POST("/verify-otp", h.VerifyOtp)
	rg.POST("/refresh-token", h.RefreshToken)
	rg.PATCH("/name", h.UpdateName)
	rg.POST("/logout", h.Logout)
	rg.DELETE("", h.Delete)
	rg.GET("/devices", h.GetDevices)
	rg.GET("/suggested-emails", h.GetSuggestedEmails)
	rg.DELETE("/suggested-emails", h.DeleteSuggestedEmail)
}

func (r *userHandler) Create(c *gin.Context) {
	platform := parseFormParam(c, "platform")
	isPhysicalDevice := parseFormBool(c, "isPhysicalDevice")
	deviceFingerprint := parseFormParam(c, "deviceFingerprint")
	deviceInfo := parseFormParam(c, "deviceInfo")

	if !isPhysicalDevice {
		panic(models.InputError("Device should be real, not a simulator"))
	}

	device := models.Device{
		Platform:    platform,
		Fingerprint: deviceFingerprint,
		Info:        deviceInfo,
	}
	createdUser := r.userService.Create(device)
	r.logger.Info(parseTrafficID(c), "User Created with id:%s, deviceId:%s.", createdUser.ID, createdUser.DeviceID)

	c.JSON(http.StatusOK, createdUser)
}

func (r *userHandler) SendMeOtp(c *gin.Context) {
	email := parseFormParam(c, "email")
	verificationSid := r.userService.SendMeOtp(email)
	r.logger.Info(parseTrafficID(c), "OTP sent to %s", email)
	c.JSON(http.StatusOK, verificationSid)
}

func (r *userHandler) VerifyOtp(c *gin.Context) {
	email := parseFormParam(c, "email")
	sid := parseFormParam(c, "sid")
	otp := parseFormParam(c, "otp")
	meta := parseRequestMeta(c)
	verifiedUser := r.userService.VerifyOtp(meta, email, sid, otp)
	r.logger.Info(parseTrafficID(c), "OTP verified for %s", email)
	c.JSON(http.StatusOK, verifiedUser)
}

func (r *userHandler) RefreshToken(c *gin.Context) {
	meta := parseRequestMeta(c)
	verifiedUser := r.userService.RefreshToken(meta)
	c.JSON(http.StatusOK, verifiedUser)
}

func (r *userHandler) UpdateName(c *gin.Context) {
	name := parseFormParam(c, "name")
	meta := parseRequestMeta(c)
	r.userService.UpdateName(meta.UserID, name)
	c.Status(http.StatusOK)
}

func (r *userHandler) Logout(c *gin.Context) {
	deviceId := parseFormUuid(c, "deviceId")
	meta := parseRequestMeta(c)
	r.userService.Logout(meta.UserID, deviceId)
	c.Status(http.StatusOK)
}

func (r *userHandler) Delete(c *gin.Context) {
	// Parsed this way as the body is not parsed by the framework for DELETE requests
	body, _ := io.ReadAll(c.Request.Body)
	values, _ := url.ParseQuery(string(body))
	sid := values.Get("sid")
	otp := values.Get("otp")
	if sid == "" || otp == "" {
		panic(models.InputError("sid and otp are required"))
	}

	meta := parseRequestMeta(c)
	r.userService.Delete(meta.UserID, sid, otp)
	c.Status(http.StatusNoContent)
}

func (r *userHandler) GetDevices(c *gin.Context) {
	meta := parseRequestMeta(c)
	devices := r.userService.GetDevices(meta.UserID)
	c.JSON(http.StatusOK, devices)
}

func (r *userHandler) GetSuggestedEmails(c *gin.Context) {
	meta := parseRequestMeta(c)
	suggestedEmails := r.userService.GetSuggestedEmails(meta.UserID)
	c.JSON(http.StatusOK, suggestedEmails)
}

func (r *userHandler) DeleteSuggestedEmail(c *gin.Context) {
	suggestedEmailId := parseFormUuid(c, "suggestedEmailId")
	meta := parseRequestMeta(c)
	r.userService.DeleteSuggestedEmail(meta.UserID, suggestedEmailId)
	c.Status(http.StatusNoContent)
}
