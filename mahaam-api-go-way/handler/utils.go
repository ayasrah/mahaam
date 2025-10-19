package handler

import (
	"mahaam-api/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"fmt"
	"strconv"
	"strings"
)

func ExtractMeta(c *gin.Context) (*model.Meta, error) {
	userId, err := ExtractUserID(c)
	if err != nil {
		return nil, err
	}
	deviceId, err := ExtractDeviceID(c)
	if err != nil {
		return nil, err
	}
	return &model.Meta{
		UserID:   userId,
		DeviceID: deviceId,
	}, nil
}

func ExtractUserID(c *gin.Context) (uuid.UUID, *model.Err) {
	userId, ok := c.Value("userId").(uuid.UUID)
	if !ok || userId == uuid.Nil {
		return uuid.Nil, model.UnauthorizedError("userId not found in context")
	}
	return userId, nil
}

func ExtractDeviceID(c *gin.Context) (uuid.UUID, *model.Err) {
	deviceId, ok := c.Value("deviceId").(uuid.UUID)
	if !ok || deviceId == uuid.Nil {
		return uuid.Nil, model.UnauthorizedError("deviceId not found in context")
	}
	return deviceId, nil
}

func ExtractTrafficID(c *gin.Context) (uuid.UUID, *model.Err) {
	trafficID, ok := c.Value("trafficID").(uuid.UUID)
	if !ok || trafficID == uuid.Nil {
		return uuid.Nil, model.UnauthorizedError("trafficID not found in context")
	}
	return trafficID, nil
}

func Query(c *gin.Context, param string) (string, *model.Err) {
	value := c.Query(param)
	if strings.TrimSpace(value) == "" {
		return "", model.InputError(param + " is required")
	}
	return value, nil
}

func Path(c *gin.Context, param string) (string, *model.Err) {
	value := c.Param(param)
	if strings.TrimSpace(value) == "" {
		return "", model.InputError(param + " is required")
	}
	return value, nil
}

func Form(c *gin.Context, param string) (string, *model.Err) {
	value := c.PostForm(param)
	if strings.TrimSpace(value) == "" {
		return "", model.InputError(param + " is required")
	}
	return value, nil
}

func FormUuid(c *gin.Context, param string) (uuid.UUID, *model.Err) {
	value, err := Form(c, param)
	if err != nil {
		return uuid.Nil, err
	}
	val, error := uuid.Parse(value)
	if error != nil {
		return uuid.Nil, model.InputError(param + " is not valid uuid")
	}
	return val, nil
}

func FormInt(c *gin.Context, param string) (int, *model.Err) {
	value, err := Form(c, param)
	if err != nil {
		return 0, err
	}
	val, error := strconv.Atoi(value)
	if error != nil {
		return 0, model.InputError(param + " is not valid integer")
	}
	return val, nil
}

func FormBool(c *gin.Context, param string) (bool, *model.Err) {
	value, err := Form(c, param)
	if err != nil {
		return false, err
	}
	val, error := strconv.ParseBool(value)
	if error != nil {
		return false, model.InputError(param + "is not valid boolean")
	}
	return val, nil
}

func PathUuid(c *gin.Context, param string) (uuid.UUID, *model.Err) {
	value, err := Path(c, param)
	if err != nil {
		return uuid.Nil, err
	}
	val, error := uuid.Parse(value)
	if error != nil {
		return uuid.Nil, model.InputError(param + " is not valid uuid")
	}
	return val, nil
}

func Parse(c *gin.Context, obj any) *model.Err {
	if err := c.ShouldBind(obj); err != nil {
		return model.InputError("Invalid body." + err.Error())
	}

	validate := validator.New()
	if err := validate.Struct(obj); err != nil {
		errors := err.(validator.ValidationErrors)
		msg := joinErrors(errors)
		return model.InputError(msg)
	}
	return nil
}

func joinErrors(ves validator.ValidationErrors) string {
	messages := make([]string, 0)
	for _, e := range ves {
		messages = append(messages, fmt.Sprintf("'%s' failed validation (%s=%s).", e.Field(), e.Tag(), e.Param()))
	}

	return strings.Join(messages, " ")
}

func Required(value any, param string) *model.Err {
	if value == nil {
		return model.InputError(param + " is required")
	}
	return nil
}
