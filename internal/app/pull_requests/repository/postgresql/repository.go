package postgresql

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/leoscrowi/pr-assignment-service/domain"
)

const tableName = "pull_requests"
const reviewersTableName = "pull_requests_reviewers"

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

	fail := func(err error) (domain.PullRequest, error) {
		return domain.PullRequest{}, fmt.Errorf("%s: %v", op, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(err)
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
			pr.PullRequestID,
			pr.PullRequestName,
			pr.AuthorID,
			pr.Status,
			pq.Array(pr.AssignedReviewers),
			pr.NeedMoreReviewers,
			pr.CreatedAt,
			pr.MergedAt,
		).
		ToSql()
	if err != nil {
		return fail(err)
	}

	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	return *pr, nil
}

func (r *Repository) FindPullRequestsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.PullRequestShort, error) {
	const op = "pull_requests.Repository.FindPullRequestsByUserID"

	fail := func(err error) ([]domain.PullRequestShort, error) {
		return []domain.PullRequestShort{}, fmt.Errorf("%s: %v", op, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	query, args, err := sq.Select(
		"pull_request_id",
		"pull_request_name",
		"author_id",
		"status",
	).From(tableName).Where(sq.Eq{"author_id": userID}).ToSql()
	if err != nil {
		return fail(err)
	}

	rows, err := tx.Queryx(query, args...)
	if err != nil {
		return fail(err)
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)

	var result []domain.PullRequestShort
	for rows.Next() {
		var pullRequestID uuid.UUID
		var pullRequestName string
		var authorID uuid.UUID
		var status domain.Status

		if err = rows.Scan(&pullRequestID, &pullRequestName, &authorID, &status); err != nil {
			return fail(err)
		}

		pr := domain.PullRequestShort{
			PullRequestID:   pullRequestID,
			PullRequestName: pullRequestName,
			AuthorID:        authorID,
			Status:          status,
		}
		result = append(result, pr)
	}

	if err = rows.Err(); err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	return result, nil
}

func (r *Repository) MergePullRequest(ctx context.Context, prID uuid.UUID) (domain.PullRequest, error) {
	const op = "pull_requests.Repository.MergePullRequest"

	fail := func(err error) (domain.PullRequest, error) {
		return domain.PullRequest{}, fmt.Errorf("%s: %v", op, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	query, args, err := sq.Update(tableName).
		Set("status", "MERGED").
		Set("merged_at", time.Now()).
		Where(sq.Eq{"pull_request_id": prID}).
		Suffix("RETURNING pull_request_id, pull_request_name, author_id, status, merged_at").
		ToSql()
	if err != nil {
		return fail(err)
	}

	var pr domain.PullRequest
	if err = tx.GetContext(ctx, &pr, query, args...); err != nil {
		return fail(err)
	}

	var reviewers []uuid.UUID
	query, args, err = sq.Select("reviewer_id").
		From(reviewersTableName).
		Where(sq.Eq{"pull_request_id": prID}).
		ToSql()
	if err = tx.SelectContext(ctx, &reviewers, query, args...); err != nil {
		return fail(err)
	}

	updated := domain.PullRequest{
		PullRequestID:     pr.PullRequestID,
		PullRequestName:   pr.PullRequestName,
		AuthorID:          pr.AuthorID,
		Status:            pr.Status,
		AssignedReviewers: reviewers,
		MergedAt:          pr.MergedAt,
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	return updated, nil
}
