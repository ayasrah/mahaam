package repo

import (
	"mahaam-api/feat/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UUID = uuid.UUID
type Param = map[string]any
type SuggestedEmail = models.SuggestedEmail
type Device = models.Device
type Plan = models.Plan
type PlanIn = models.PlanIn
type Task = models.Task
type User = models.User
type HttpErr = models.HttpErr
type DBErr = models.DBErr
type Tx = *sqlx.Tx
