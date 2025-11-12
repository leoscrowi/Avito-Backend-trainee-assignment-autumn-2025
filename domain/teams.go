package domain

import "github.com/google/uuid"

type Team struct {
	TeamName string `yml:"team_name"`
	Members  []User `yml:"members"`
}

// TODO: разобраться, зачем нужно, пока выглядит сомнительно, но взято из openapi.yml
type TeamMember struct {
	UserID   uuid.UUID `yml:"user_id"`
	UserName string    `yml:"username"`
	IsActive bool      `yml:"is_active"`
}
