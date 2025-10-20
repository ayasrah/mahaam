package handler

import (
	"mahaam-api/app/models"

	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/go-playground/validator/v10"
)

type SuggestedEmail = models.SuggestedEmail
type Device = models.Device
type Plan = models.Plan
type PlanIn = models.PlanIn
type Task = models.Task
type User = models.User
type CreatedUser = models.CreatedUser
type VerifiedUser = models.VerifiedUser
type Meta = models.Meta
type PlanType = models.PlanType

const PlanTypeMain = models.PlanTypeMain
const PlanTypeArchived = models.PlanTypeArchived

func parseRequestMeta(c *gin.Context) Meta {
	return Meta{
		UserID:   parseUserID(c),
		DeviceID: parseDeviceID(c),
	}
}

func parseUserID(c *gin.Context) uuid.UUID {
	userId, ok := c.Value("userId").(uuid.UUID)
	if !ok || userId == uuid.Nil {
		panic(models.UnauthorizedError("userId not found in context"))
	}
	return userId
}

func parseDeviceID(c *gin.Context) uuid.UUID {
	deviceId, ok := c.Value("deviceId").(uuid.UUID)
	if !ok || deviceId == uuid.Nil {
		panic(models.UnauthorizedError("deviceId not found in context"))
	}
	return deviceId
}

func parseTrafficID(c *gin.Context) uuid.UUID {
	trafficID, ok := c.Value("trafficID").(uuid.UUID)
	if !ok || trafficID == uuid.Nil {
		panic(models.UnauthorizedError("trafficID not found in context"))
	}
	return trafficID
}

func parseQueryParam(c *gin.Context, param string) string {
	value := c.Query(param)
	if strings.TrimSpace(value) == "" {
		panic(models.InputError(param + " is required"))
	}
	return value
}

func parsePathParam(c *gin.Context, param string) string {
	value := c.Param(param)
	if strings.TrimSpace(value) == "" {
		panic(models.InputError(param + " is required"))
	}
	return value
}

func parseFormParam(c *gin.Context, param string) string {
	value := c.PostForm(param)
	if strings.TrimSpace(value) == "" {
		panic(models.InputError(param + " is required"))
	}
	return value
}

func parseFormUuid(c *gin.Context, param string) uuid.UUID {
	valueStr := parseFormParam(c, param)
	value, err := uuid.Parse(valueStr)
	if err != nil {
		panic(models.InputError(param + " is not valid uuid"))
	}
	return value
}

func parseFormInt(c *gin.Context, param string) int {
	valueStr := parseFormParam(c, param)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		panic(models.InputError(param + " is not valid integer"))
	}
	return value
}

func parseFormBool(c *gin.Context, param string) bool {
	value := parseFormParam(c, param)
	val, err := strconv.ParseBool(value)
	if err != nil {
		panic(models.InputError(param + " is not valid boolean"))
	}
	return val
}

func parsePathUuid(c *gin.Context, param string) uuid.UUID {
	value := parsePathParam(c, param)
	id, err := uuid.Parse(value)
	if err != nil {
		panic(models.InputError(param + " is not valid uuid"))
	}
	return id
}

func parse(c *gin.Context, obj any) {
	if err := c.ShouldBind(obj); err != nil {
		panic(models.InputError("Invalid body." + err.Error()))
	}

	validate := validator.New()
	if err := validate.Struct(obj); err != nil {
		errors := err.(validator.ValidationErrors)
		msg := joinErrors(errors)
		panic(models.InputError(msg))
	}
}

func joinErrors(ves validator.ValidationErrors) string {
	messages := make([]string, 0)
	for _, e := range ves {
		messages = append(messages, fmt.Sprintf("'%s' failed validation (%s=%s).", e.Field(), e.Tag(), e.Param()))
	}

	return strings.Join(messages, " ")
}

func requiredParam(value any, param string) {
	if value == nil {
		panic(models.InputError(param + " is required"))
	}
}
