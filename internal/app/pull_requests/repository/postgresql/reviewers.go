package postgresql

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/leoscrowi/pr-assignment-service/domain"
)

const reviewersTableName = "pull_requests_reviewers"

func (r *Repository) GetReviewersID(ctx context.Context, prID string) ([]string, error) {
	const op = "pull_requests.Repository.GetReviewersID"

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

	query, args, err := sq.Select("reviewer_id").From(reviewersTableName).Where(sq.Eq{"pull_request_id": prID}).PlaceholderFormat(sq.Dollar).ToSql()
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

	var ids []string
	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			return fail(domain.INTERNAL, "internal server error", err)
		}
		ids = append(ids, id)
	}
	if err = rows.Err(); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	if err = tx.Commit(); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	return ids, nil
}

func (r *Repository) DeleteReviewer(ctx context.Context, prID, reviewerID string) error {
	const op = "pull_requests.Repository.DeleteReviewer"

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

	query, args, err := sq.Delete(reviewersTableName).Where(sq.Eq{"pull_request_id": prID, "reviewer_id": reviewerID}).PlaceholderFormat(sq.Dollar).ToSql()
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

func (r *Repository) AddReviewer(ctx context.Context, prID, reviewerID string) error {
	const op = "pull_requests.Repository.AddReviewer"

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

	query, args, err := sq.Insert(reviewersTableName).Columns("pull_request_id", "reviewer_id").Values(prID, reviewerID).PlaceholderFormat(sq.Dollar).ToSql()
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
