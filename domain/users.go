package domain

import "github.com/google/uuid"

type User struct {
	UserID   uuid.UUID `yml:"user_id"`
	Username string    `yml:"username"`
	IsActive bool      `yml:"is_active"`
}
