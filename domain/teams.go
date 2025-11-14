package domain

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

// TODO: разобраться, зачем нужно, пока выглядит сомнительно, но взято из openapi.yml
type TeamMember struct {
	UserID   string `json:"user_id"`
	UserName string `json:"username"`
	IsActive bool   `json:"is_active"`
}
