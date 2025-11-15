package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/leoscrowi/pr-assignment-service/domain"
)

const tableName = "pull_requests"

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreatePullRequest(ctx context.Context, pr *domain.PullRequest) error {
	const op = "pull_requests.Repository.CreatePullRequest"

	fail := func(code domain.ErrorCode, message string, err error) error {
		log.Printf("%s: %v", op, err)
		return domain.NewError(code, message, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
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
			"need_more_reviewers",
			"created_at",
		).
		Values(
			pr.PullRequestID,
			pr.PullRequestName,
			pr.AuthorID,
			pr.Status,
			len(pr.AssignedReviewers) < 2,
			time.Now(),
		).
		ToSql()
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	if err = tx.Commit(); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	return nil
}

func (r *Repository) FindPullRequestsByUserID(ctx context.Context, userID string) ([]domain.PullRequestShort, error) {
	const op = "pull_requests.Repository.FindPullRequestsByUserID"

	fail := func(code domain.ErrorCode, message string, err error) ([]domain.PullRequestShort, error) {
		log.Printf("%s: %v", op, err)
		return []domain.PullRequestShort{}, domain.NewError(code, message, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
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
		return fail(domain.INTERNAL, "internal server error", err)
	}

	rows, err := tx.QueryxContext(ctx, query, args...)
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)

	var result []domain.PullRequestShort
	for rows.Next() {
		var pullRequestID string
		var pullRequestName string
		var authorID string
		var status domain.Status

		if err = rows.Scan(&pullRequestID, &pullRequestName, &authorID, &status); err != nil {
			return fail(domain.INTERNAL, "internal server error", err)
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
		return fail(domain.INTERNAL, "internal server error", err)
	}

	if err = tx.Commit(); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	return result, nil
}

func (r *Repository) MergePullRequest(ctx context.Context, prID string) (domain.PullRequest, error) {
	const op = "pull_requests.Repository.MergePullRequest"

	fail := func(code domain.ErrorCode, message string, err error) (domain.PullRequest, error) {
		log.Printf("%s: %v", op, err)
		return domain.PullRequest{}, domain.NewError(code, message, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
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
		return fail(domain.INTERNAL, "internal server error", err)
	}

	var pr domain.PullRequest
	if err = tx.GetContext(ctx, &pr, query, args...); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	var reviewers []string
	query, args, err = sq.Select("reviewer_id").
		From(reviewersTableName).
		Where(sq.Eq{"pull_request_id": prID}).
		ToSql()
	if err = tx.SelectContext(ctx, &reviewers, query, args...); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
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
		return fail(domain.INTERNAL, "internal server error", err)
	}

	return updated, nil
}

func (r *Repository) IsMerged(ctx context.Context, prID string) (bool, error) {
	const op = "pull_requests.Repository.IsMerged"

	fail := func(code domain.ErrorCode, message string, err error) (bool, error) {
		log.Printf("%s: %v", op, err)
		return false, domain.NewError(code, message, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	query, args, err := sq.Select("pull_request_id").From(tableName).Where(sq.Eq{"pull_request_id": prID}).Limit(1).ToSql()
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	var pr domain.PullRequest
	if err = tx.GetContext(ctx, &pr, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if commitErr := tx.Commit(); commitErr != nil {
				return fail(domain.INTERNAL, "internal server error", err)
			}
			return false, nil
		}
		return fail(domain.INTERNAL, "internal server error", err)
	}

	return pr.Status == domain.MERGED, nil
}

func (r *Repository) FetchByID(ctx context.Context, prID string) (domain.PullRequest, error) {
	const op = "pull_requests.Repository.CreatePullRequest"

	fail := func(code domain.ErrorCode, message string, err error) (domain.PullRequest, error) {
		log.Printf("%s: %v", op, err)
		return domain.PullRequest{}, domain.NewError(code, message, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	query, args, err := sq.Select("pull_request_id", "pull_request_name", "author_id", "status").From(tableName).Where(sq.Eq{"pull_request_id": prID}).Limit(1).ToSql()
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	var pr domain.PullRequest
	if err = tx.GetContext(ctx, &pr, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fail(domain.NOT_FOUND, "resource not found", err)
		}
		return fail(domain.INTERNAL, "internal server error", err)
	}

	if err = tx.Commit(); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	var reviewers []string
	query, args, err = sq.Select("reviewer_id").
		From(reviewersTableName).
		Where(sq.Eq{"pull_request_id": prID}).
		ToSql()
	if err = tx.SelectContext(ctx, &reviewers, query, args...); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	pr.AssignedReviewers = reviewers

	if err = tx.Commit(); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	return pr, nil
}
