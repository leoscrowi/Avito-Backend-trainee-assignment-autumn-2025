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

func NewPullRequestsRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreatePullRequest(ctx context.Context, pr *domain.PullRequest) error {
	const op = "pull_requests.Repository.CreatePullRequest"

	fail := func(code domain.ErrorCode, message string, err error) error {
		log.Printf("%s: %v\n", op, err)
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
		PlaceholderFormat(sq.Dollar).
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

func (r *Repository) FindPullRequestsIDByUserID(ctx context.Context, userID string) ([]string, error) {
	const op = "pull_requests.Repository.FindPullRequestsIDByUserID"

	fail := func(code domain.ErrorCode, message string, err error) ([]string, error) {
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
		"pull_request_id",
	).From(reviewersTableName).Where(sq.Eq{"reviewer_id": userID}).PlaceholderFormat(sq.Dollar).ToSql()
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

	var result []string
	for rows.Next() {
		var reviewerID string

		if err = rows.Scan(&reviewerID); err != nil {
			return fail(domain.INTERNAL, "internal server error", err)
		}

		result = append(result, reviewerID)
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
		log.Printf("%s: %v\n", op, err)
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
		PlaceholderFormat(sq.Dollar).
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
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

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

func (r *Repository) FetchByIDWithMergeAt(ctx context.Context, prID string) (domain.PullRequest, error) {
	const op = "pull_requests.Repository.FetchByIDWithMergedAt"

	fail := func(code domain.ErrorCode, message string, err error) (domain.PullRequest, error) {
		log.Printf("%s: %v\n", op, err)
		return domain.PullRequest{}, domain.NewError(code, message, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	query, args, err := sq.Select("pull_request_id", "pull_request_name", "author_id", "status", "COALESCE(merged_at, '0001-01-01'::timestamp) as merged_at").
		From(tableName).Where(sq.Eq{"pull_request_id": prID}).
		Limit(1).PlaceholderFormat(sq.Dollar).ToSql()
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

	var reviewers []string
	query, args, err = sq.Select("reviewer_id").
		From(reviewersTableName).
		Where(sq.Eq{"pull_request_id": prID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}
	if err = tx.SelectContext(ctx, &reviewers, query, args...); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	pr.AssignedReviewers = reviewers

	if err = tx.Commit(); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	return pr, nil
}

func (r *Repository) FetchByID(ctx context.Context, prID string) (domain.PullRequest, error) {
	const op = "pull_requests.Repository.FetchByID"

	fail := func(code domain.ErrorCode, message string, err error) (domain.PullRequest, error) {
		log.Printf("%s: %v\n", op, err)
		return domain.PullRequest{}, domain.NewError(code, message, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	query, args, err := sq.Select("pull_request_id", "pull_request_name", "author_id", "status").
		From(tableName).Where(sq.Eq{"pull_request_id": prID}).
		Limit(1).PlaceholderFormat(sq.Dollar).ToSql()
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

	var reviewers []string
	query, args, err = sq.Select("reviewer_id").
		From(reviewersTableName).
		Where(sq.Eq{"pull_request_id": prID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}
	if err = tx.SelectContext(ctx, &reviewers, query, args...); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	pr.AssignedReviewers = reviewers

	if err = tx.Commit(); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	return pr, nil
}

func (r *Repository) FetchShortByID(ctx context.Context, prID string) (domain.PullRequestShort, error) {
	const op = "pull_requests.Repository.FetchShortByID"

	fail := func(code domain.ErrorCode, message string, err error) (domain.PullRequestShort, error) {
		log.Printf("%s: %v\n", op, err)
		return domain.PullRequestShort{}, domain.NewError(code, message, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	query, args, err := sq.Select(
		"pull_request_id", "pull_request_name", "author_id", "status",
	).From(tableName).Where(sq.Eq{"pull_request_id": prID}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	var pr domain.PullRequestShort
	if err = tx.GetContext(ctx, &pr, query, args...); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	if err = tx.Commit(); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	return pr, nil
}
