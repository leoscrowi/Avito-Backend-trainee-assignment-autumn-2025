package postgresql

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/leoscrowi/pr-assignment-service/domain"
)

type Repository struct {
	db *sqlx.DB
}

func NewStatsRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetPullRequestStats(ctx context.Context) ([]domain.PullRequestStats, error) {
	const op = "stats.Repository.GetPullRequestStats"

	fail := func(code domain.ErrorCode, message string, err error) ([]domain.PullRequestStats, error) {
		log.Printf("%s: %v\n", op, err)
		return nil, domain.NewError(code, message, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	query, args, err := sq.Select(
		"u.user_id",
		"u.username",
		"u.team_name",
		"COUNT(prr.pull_request_id) as assigned_review_count",
		"COUNT(CASE WHEN pr.status = 'OPEN' THEN 1 END) as open_pr_review_count",
		"COUNT(CASE WHEN pr.status = 'MERGED' THEN 1 END) as merged_pr_review_count",
	).
		From("users u").
		LeftJoin("pull_request_reviewers prr ON u.user_id = prr.reviewer_id").
		LeftJoin("pull_requests pr ON prr.pull_request_id = pr.pull_request_id").
		GroupBy("u.user_id", "u.username", "u.team_name").
		OrderBy("assigned_review_count DESC").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	rows, err := tx.QueryxContext(ctx, query, args...)
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)

	var result []domain.PullRequestStats
	for rows.Next() {
		var userID string
		var userName string
		var teamName string
		var assignedPRCount int
		var open int
		var merged int

		if err = rows.Scan(&userID, &userName, &teamName, &assignedPRCount, &open, &merged); err != nil {
			return fail(domain.INTERNAL, "internal server error", err)
		}

		pr := domain.PullRequestStats{
			UserID:          userID,
			UserName:        userName,
			TeamName:        teamName,
			AssignedPRCount: assignedPRCount,
			Open:            open,
			Merged:          merged,
		}
		result = append(result, pr)
	}

	if err = rows.Err(); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	if err = tx.Commit(); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	return result, nil
}
