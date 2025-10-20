package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        uuid.UUID  `db:"id"`
	PlanID    uuid.UUID  `db:"plan_id"`
	Title     string     `db:"title"`
	Done      bool       `db:"done"`
	SortOrder int        `db:"sort_order"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
