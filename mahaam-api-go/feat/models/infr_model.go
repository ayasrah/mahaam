package models

import (
	"fmt"
)

type HttpErr struct {
	Code    int
	Message string
	Key     string
}

func (e *HttpErr) Error() string {
	return fmt.Sprintf("HttpErr: %v, %v, %v", e.Code, e.Message, e.Key)
}

func UnauthorizedErr(message string) *HttpErr {
	return &HttpErr{
		Code:    401,
		Message: message,
	}
}

func ForbiddenErr(message string) *HttpErr {
	return &HttpErr{
		Code:    403,
		Message: message,
	}
}

func NotFoundErr(message string) *HttpErr {
	return &HttpErr{
		Code:    404,
		Message: message,
	}
}

func InputErr(message string) *HttpErr {
	return &HttpErr{
		Code:    400,
		Message: message,
	}
}

func LogicErr(message, key string) *HttpErr {
	return &HttpErr{
		Code:    409,
		Message: message,
		Key:     key,
	}
}

type DBErr struct {
	Err error
}

func (e *DBErr) Error() string {
	return fmt.Sprintf("DBErr: %v", e.Err.Error())
}
