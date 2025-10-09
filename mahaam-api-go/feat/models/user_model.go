package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID    uuid.UUID `json:"id,omitempty"`
	Email *string   `json:"email,omitempty"`
	Name  *string   `json:"name,omitempty"`
}

type User2 struct {
	User
	UserID *uuid.UUID `json:"userId"`
}

type Device struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"userId" db:"user_id"`
	Platform    string    `json:"platform" db:"platform"`
	Fingerprint string    `json:"fingerprint" db:"fingerprint"`
	Info        string    `json:"info" db:"info"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type SuggestedEmail struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"userId"`
	Email     *string    `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
}

type VerifiedUser struct {
	UserID       uuid.UUID `json:"userId"`
	DeviceID     uuid.UUID `json:"deviceId"`
	Jwt          string    `json:"jwt"`
	UserFullName *string   `json:"userFullName"`
	Email        *string   `json:"email"`
}

type CreatedUser struct {
	ID       uuid.UUID `json:"id"`
	DeviceID uuid.UUID `json:"deviceId"`
	Jwt      string    `json:"jwt"`
}

type Meta struct {
	UserID   uuid.UUID
	DeviceID uuid.UUID
}
