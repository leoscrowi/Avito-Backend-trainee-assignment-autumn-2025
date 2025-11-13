package postgresql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/leoscrowi/pr-assignment-service/domain"
	"github.com/lib/pq"
)

const tableName = "pull_requests"

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreatePullRequest(ctx context.Context, pr *domain.PullRequest) (domain.PullRequest, error) {
	const op = "pull_requests.Repository.CreatePullRequest"

	if pr == nil {
		return domain.PullRequest{}, fmt.Errorf("%s: pr is nil", op)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return domain.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	query, args, err := sq.Insert(tableName).
		Columns(
			"pull_request_id",
			"pull_request_name",
			"author_id",
			"status",
			"assigned_reviewers",
			"need_more_reviewers",
			"created_at",
			"updated_at",
		).
		Values(
			pr.PullRequestId,
			pr.PullRequestName,
			pr.AuthorId,
			pr.Status,
			pq.Array(pr.AssignedReviewers),
			pr.NeedMoreReviewers,
			pr.CreatedAt,
			pr.UpdatedAt,
		).
		ToSql()
	if err != nil {
		return domain.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return domain.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		return domain.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	return *pr, nil
}

func (r *Repository) UpdatePullRequest(ctx context.Context, pr *domain.PullRequest) (domain.PullRequest, error) {
	const op = "pull_requests.Repository.UpdatePullRequest"

	if pr == nil {
		return domain.PullRequest{}, fmt.Errorf("%s: pr is nil", op)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return domain.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	// TODO: проверить, будут ли затираться пустые значения, если что - пофиксить
	query, args, err := sq.Update(tableName).
		Set("pull_request_name", pr.PullRequestName).
		Set("author_id", pr.AuthorId).
		Set("status", pr.Status).
		Set("assigned_reviewers", pq.Array(pr.AssignedReviewers)).
		Set("need_more_reviewers", pr.NeedMoreReviewers).
		Set("created_at", pr.CreatedAt).
		Set("updated_at", pr.UpdatedAt).
		Where(sq.Eq{"pull_request_id": pr.PullRequestId}).
		ToSql()
	if err != nil {
		return domain.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return domain.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		return domain.PullRequest{}, fmt.Errorf("%s: %w", op, err)
	}

	return *pr, nil
}
