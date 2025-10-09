package handler

import (
	"io"
	"mahaam-api/feat/models"
	"mahaam-api/feat/service"
	logs "mahaam-api/infra/log"
	"mahaam-api/infra/middleware"
	"mahaam-api/infra/req"
	"net/url"
)

type UserHandler interface {
	Create(c Ctx)
	SendMeOtp(c Ctx)
	VerifyOtp(c Ctx)
	RefreshToken(c Ctx)
	UpdateName(c Ctx)
	Logout(c Ctx)
	Delete(c Ctx)
	GetDevices(c Ctx)
	GetSuggestedEmails(c Ctx)
	DeleteSuggestedEmail(c Ctx)
}

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &userHandler{userService: userService}
}

func RegisterUserHandler(router Router, h UserHandler) {
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

func (r *userHandler) Create(c Ctx) {
	platform := req.Form(c, "platform")
	isPhysicalDevice := req.FormBool(c, "isPhysicalDevice")
	deviceFingerprint := req.Form(c, "deviceFingerprint")
	deviceInfo := req.Form(c, "deviceInfo")

	if !isPhysicalDevice {
		panic(models.InputErr("Device should be real, not a simulator"))
	}

	device := models.Device{
		Platform:    platform,
		Fingerprint: deviceFingerprint,
		Info:        deviceInfo,
	}
	createdUser := r.userService.Create(device)
	logs.Info(ExtractTrafficID(c), "User Created with id:%s, deviceId:%s.", createdUser.ID, createdUser.DeviceID)

	c.JSON(OK, createdUser)
}

func (r *userHandler) SendMeOtp(c Ctx) {
	email := req.Form(c, "email")
	verificationSid := r.userService.SendMeOtp(email)
	logs.Info(ExtractTrafficID(c), "OTP sent to %s", email)
	c.JSON(OK, verificationSid)
}

func (r *userHandler) VerifyOtp(c Ctx) {
	email := req.Form(c, "email")
	sid := req.Form(c, "sid")
	otp := req.Form(c, "otp")
	meta := ExtractMeta(c)
	verifiedUser := r.userService.VerifyOtp(meta, email, sid, otp)
	logs.Info(ExtractTrafficID(c), "OTP verified for %s", email)
	c.JSON(OK, verifiedUser)
}

func (r *userHandler) RefreshToken(c Ctx) {
	meta := ExtractMeta(c)
	verifiedUser := r.userService.RefreshToken(meta)
	c.JSON(OK, verifiedUser)
}

func (r *userHandler) UpdateName(c Ctx) {
	name := req.Form(c, "name")
	meta := ExtractMeta(c)
	r.userService.UpdateName(meta.UserID, name)
	c.Status(OK)
}

func (r *userHandler) Logout(c Ctx) {
	deviceId := req.FormUuid(c, "deviceId")
	meta := ExtractMeta(c)
	r.userService.Logout(meta.UserID, deviceId)
	c.Status(OK)
}

func (r *userHandler) Delete(c Ctx) {
	// Parsed this way as the body is not parsed by the framework for DELETE requests
	body, _ := io.ReadAll(c.Request.Body)
	values, _ := url.ParseQuery(string(body))
	sid := values.Get("sid")
	otp := values.Get("otp")
	if sid == "" || otp == "" {
		panic(models.InputErr("sid and otp are required"))
	}

	meta := ExtractMeta(c)
	r.userService.Delete(meta.UserID, sid, otp)
	c.Status(NoContent)
}

func (r *userHandler) GetDevices(c Ctx) {
	meta := ExtractMeta(c)
	devices := r.userService.GetDevices(meta.UserID)
	c.JSON(OK, devices)
}

func (r *userHandler) GetSuggestedEmails(c Ctx) {
	meta := ExtractMeta(c)
	suggestedEmails := r.userService.GetSuggestedEmails(meta.UserID)
	c.JSON(OK, suggestedEmails)
}

func (r *userHandler) DeleteSuggestedEmail(c Ctx) {
	suggestedEmailId := req.FormUuid(c, "suggestedEmailId")
	meta := ExtractMeta(c)
	r.userService.DeleteSuggestedEmail(meta.UserID, suggestedEmailId)
	c.Status(NoContent)
}
