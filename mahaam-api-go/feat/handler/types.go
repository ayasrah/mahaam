package handler

import (
	"mahaam-api/feat/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UUID = uuid.UUID
type SuggestedEmail = models.SuggestedEmail
type Device = models.Device
type Plan = models.Plan
type PlanIn = models.PlanIn
type Task = models.Task
type User = models.User
type CreatedUser = models.CreatedUser
type VerifiedUser = models.VerifiedUser
type Ctx = *gin.Context
type Router = *gin.RouterGroup
type res = map[string]any
type Meta = models.Meta

const OK = http.StatusOK
const Created = http.StatusCreated
const NoContent = http.StatusNoContent
const BadRequest = http.StatusBadRequest
const Unauthorized = http.StatusUnauthorized
const Forbidden = http.StatusForbidden
const NotFound = http.StatusNotFound
const Conflict = http.StatusConflict
const ServerError = http.StatusInternalServerError

type HttpErr = models.HttpErr
