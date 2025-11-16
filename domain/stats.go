package domain

type PullRequestStats struct {
	UserID          string `json:"user_id"`
	UserName        string `json:"user_name"`
	TeamName        string `json:"team_name"`
	AssignedPRCount int    `json:"assigned_pr_count"`
	Open            int    `json:"open_pull_requests"`
	Merged          int    `json:"merged_pull_requests"`
}
