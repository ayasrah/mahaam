package models

import (
	"fmt"
	"net/http"
)

type Err struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Key     string `json:"key"`
}

func (e *Err) Error() string {
	return fmt.Sprintf("HttpError: %v, %v, %v", e.Code, e.Message, e.Key)
}

func UnauthorizedError(message string) *Err {
	return &Err{
		Code:    http.StatusUnauthorized,
		Message: message,
		Key:     "UNAUTHORIZED",
	}
}

func ForbiddenError(message string) *Err {
	return &Err{
		Code:    http.StatusForbidden,
		Message: message,
		Key:     "FORBIDDEN",
	}
}

func NotFoundError(message string) *Err {
	return &Err{
		Code:    http.StatusNotFound,
		Message: message,
		Key:     "NOT_FOUND",
	}
}

func InputError(message string) *Err {
	return &Err{
		Code:    http.StatusBadRequest,
		Message: message,
		Key:     "INPUT_ERROR",
	}
}

func LogicError(message, key string) *Err {
	return &Err{
		Code:    http.StatusConflict,
		Message: message,
		Key:     key,
	}
}

func ServerError(message string) *Err {
	return &Err{
		Code:    http.StatusInternalServerError,
		Message: message,
		Key:     "SERVER_ERROR",
	}
}
