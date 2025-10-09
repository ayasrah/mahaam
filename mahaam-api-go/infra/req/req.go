package req

import (
	"fmt"
	"strconv"
	"strings"

	"mahaam-api/feat/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func Query(c *gin.Context, param string) string {
	value := c.Query(param)
	if strings.TrimSpace(value) == "" {
		panic(models.InputErr(param + " is required"))
	}
	return value
}

func Path(c *gin.Context, param string) string {
	value := c.Param(param)
	if strings.TrimSpace(value) == "" {
		panic(models.InputErr(param + " is required"))
	}
	return value
}

func Form(c *gin.Context, param string) string {
	value := c.PostForm(param)
	if strings.TrimSpace(value) == "" {
		panic(models.InputErr(param + " is required"))
	}
	return value
}

func FormUuid(c *gin.Context, param string) uuid.UUID {
	valueStr := Form(c, param)
	value, err := uuid.Parse(valueStr)
	if err != nil {
		panic(models.InputErr(param + " is not valid uuid"))
	}
	return value
}

func FormInt(c *gin.Context, param string) int {
	valueStr := Form(c, param)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		panic(models.InputErr(param + " is not valid integer"))
	}
	return value
}

func FormBool(c *gin.Context, param string) bool {
	value := Form(c, param)
	val, err := strconv.ParseBool(value)
	if err != nil {
		panic(models.InputErr(param + "is not valid boolean"))
	}
	return val
}

func PathUuid(c *gin.Context, param string) uuid.UUID {
	value := Path(c, param)
	id, err := uuid.Parse(value)
	if err != nil {
		panic(models.InputErr(param + " is not valid uuid"))
	}
	return id
}

func Parse(c *gin.Context, obj any) {
	if err := c.ShouldBind(obj); err != nil {
		panic(models.InputErr("Invalid body." + err.Error()))
	}

	validate := validator.New()
	if err := validate.Struct(obj); err != nil {
		errors := err.(validator.ValidationErrors)
		msg := joinErrors(errors)
		panic(models.InputErr(msg))
	}
}

func joinErrors(ves validator.ValidationErrors) string {
	messages := make([]string, 0)
	for _, e := range ves {
		messages = append(messages, fmt.Sprintf("'%s' failed validation (%s=%s).", e.Field(), e.Tag(), e.Param()))
	}

	return strings.Join(messages, " ")
}

func Required(value any, param string) {
	if value == nil {
		panic(models.InputErr(param + " is required"))
	}
}
