package handler

import (
	"io"
	"mahaam-api/internal/model"
	users "mahaam-api/internal/user"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterUsers(router *gin.RouterGroup) {
	userRouter := router.Group("/users")

	userRouter.POST("/create", CreateUser)
	userRouter.POST("/send-me-otp", SendMeOtp)
	userRouter.POST("/verify-otp", VerifyOtp)
	userRouter.POST("/refresh-token", RefreshTokenHandler)
	userRouter.PATCH("/name", UpdateUserName)
	userRouter.POST("/logout", LogoutUser)
	userRouter.DELETE("", DeleteUser)
	userRouter.GET("/devices", GetUserDevices)
	userRouter.GET("/suggested-emails", GetUserSuggestedEmails)
	userRouter.DELETE("/suggested-emails", DeleteUserSuggestedEmail)
}

func CreateUser(c *gin.Context) {
	var err *model.Err

	platform, err := Form(c, "platform")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	isPhysicalDevice, err := FormBool(c, "isPhysicalDevice")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	if !isPhysicalDevice {
		c.JSON(http.StatusBadRequest, model.InputError("Device should be real, not a simulator"))
		return
	}

	deviceFingerprint, err := Form(c, "deviceFingerprint")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	deviceInfo, err := Form(c, "deviceInfo")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	device := model.Device{
		Platform:    platform,
		Fingerprint: deviceFingerprint,
		Info:        deviceInfo,
	}

	createdUser, err := users.Create(device)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusOK, createdUser)
}

func SendMeOtp(c *gin.Context) {
	var err *model.Err

	email, err := Form(c, "email")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	verificationSid, err := users.SendOtp(email)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusOK, verificationSid)
}

func VerifyOtp(c *gin.Context) {
	var err *model.Err

	email, err := Form(c, "email")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	sid, err := Form(c, "sid")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	otp, err := Form(c, "otp")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	meta, extractErr := ExtractMeta(c)
	if extractErr != nil {
		c.JSON(http.StatusUnauthorized, model.UnauthorizedError("failed to extract meta"))
		return
	}

	verifiedUser, err := users.VerifyOtp(*meta, email, sid, otp)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusOK, verifiedUser)
}

func RefreshTokenHandler(c *gin.Context) {
	meta, err := ExtractMeta(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, model.UnauthorizedError("failed to extract meta"))
		return
	}

	verifiedUser, modelErr := users.RefreshToken(*meta)
	if modelErr != nil {
		c.JSON(modelErr.Code, modelErr)
		return
	}

	c.JSON(http.StatusOK, verifiedUser)
}

func UpdateUserName(c *gin.Context) {
	var err *model.Err
	var userID uuid.UUID

	name, err := Form(c, "name")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	if userID, err = ExtractUserID(c); err != nil {
		c.JSON(err.Code, err)
		return
	}

	if err := users.UpdateName(userID, name); err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.Status(http.StatusOK)
}

func LogoutUser(c *gin.Context) {
	var err *model.Err
	var userID uuid.UUID

	deviceID, err := FormUuid(c, "deviceId")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	if userID, err = ExtractUserID(c); err != nil {
		c.JSON(err.Code, err)
		return
	}

	if err := users.Logout(userID, deviceID); err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.Status(http.StatusOK)
}

func DeleteUser(c *gin.Context) {
	var err *model.Err
	var userID uuid.UUID

	// Parse body manually for DELETE request
	body, _ := io.ReadAll(c.Request.Body)
	values, _ := url.ParseQuery(string(body))
	sid := values.Get("sid")
	otp := values.Get("otp")

	if sid == "" || otp == "" {
		c.JSON(http.StatusBadRequest, model.InputError("sid and otp are required"))
		return
	}

	if userID, err = ExtractUserID(c); err != nil {
		c.JSON(err.Code, err)
		return
	}

	if err := users.Delete(userID, sid, otp); err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func GetUserDevices(c *gin.Context) {
	var err *model.Err
	var userID uuid.UUID

	if userID, err = ExtractUserID(c); err != nil {
		c.JSON(err.Code, err)
		return
	}

	devices, err := users.GetDevices(userID)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusOK, devices)
}

func GetUserSuggestedEmails(c *gin.Context) {
	var err *model.Err
	var userID uuid.UUID

	if userID, err = ExtractUserID(c); err != nil {
		c.JSON(err.Code, err)
		return
	}

	suggestedEmails, err := users.GetSuggestedEmails(userID)
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.JSON(http.StatusOK, suggestedEmails)
}

func DeleteUserSuggestedEmail(c *gin.Context) {
	var err *model.Err
	var userID uuid.UUID

	suggestedEmailID, err := FormUuid(c, "suggestedEmailId")
	if err != nil {
		c.JSON(err.Code, err)
		return
	}

	if userID, err = ExtractUserID(c); err != nil {
		c.JSON(err.Code, err)
		return
	}

	if err := users.DeleteSuggestedEmail(userID, suggestedEmailID); err != nil {
		c.JSON(err.Code, err)
		return
	}

	c.Status(http.StatusNoContent)
}
