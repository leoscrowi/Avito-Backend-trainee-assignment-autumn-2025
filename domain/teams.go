package domain

import "github.com/google/uuid"

type Team struct {
	TeamName string `json:"team_name"`
	Members  []User `json:"members"`
}

// TODO: разобраться, зачем нужно, пока выглядит сомнительно, но взято из openapi.yml
type TeamMember struct {
	UserID   uuid.UUID `json:"user_id"`
	UserName string    `json:"username"`
	IsActive bool      `json:"is_active"`
}
