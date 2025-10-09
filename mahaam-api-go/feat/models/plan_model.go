package models

import (
	"time"

	"github.com/google/uuid"
)

type Plan struct {
	ID          uuid.UUID  `json:"id,omitempty"`
	Title       *string    `json:"title,omitempty"`
	Type        *string    `json:"type,omitempty"`
	SortOrder   int        `json:"sortOrder,omitempty" db:"sort_order"`
	Starts      *time.Time `json:"starts,omitempty"`
	Ends        *time.Time `json:"ends,omitempty"`
	DonePercent *string    `json:"donePercent,omitempty" db:"done_percent"`
	CreatedAt   *time.Time `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
	Members     []User     `json:"members,omitempty" db:"user_id"`
	IsShared    bool       `json:"isShared,omitempty" db:"is_shared"`
	User        User       `json:"user,omitempty" db:"user"`
}

type PlanIn struct {
	ID     uuid.UUID `json:"id"`
	Title  *string   `json:"title"`
	Starts *string   `json:"starts"`
	Ends   *string   `json:"ends"`
}

type PlanType string

const (
	PlanTypeMain     PlanType = "Main"
	PlanTypeArchived PlanType = "Archived"
)

var AllPlanTypes = []PlanType{
	PlanTypeMain,
	PlanTypeArchived,
}
