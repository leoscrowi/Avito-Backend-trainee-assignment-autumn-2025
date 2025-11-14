package postgresql

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

const reviewersTableName = "pull_requests_reviewers"

func (r *Repository) GetReviewersID(ctx context.Context, prID string) ([]string, error) {
	const op = "pull_requests.Repository.GetReviewersID"

	fail := func(err error) ([]string, error) {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	query, args, err := sq.Select("reviewer_id").From(reviewersTableName).Where(sq.Eq{"pull_request_id": prID}).ToSql()
	if err != nil {
		return fail(err)
	}

	rows, err := tx.QueryxContext(ctx, query, args...)
	if err != nil {
		return fail(err)
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)

	var ids []string
	for rows.Next() {
		var id string
		if err = rows.Scan(&id); err != nil {
			return fail(err)
		}
		ids = append(ids, id)
	}
	if err = rows.Err(); err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	return ids, nil
}

func (r *Repository) DeleteReviewer(ctx context.Context, prID, reviewerID string) error {
	const op = "pull_requests.Repository.DeleteReviewer"

	fail := func(err error) error {
		return fmt.Errorf("%s: %v", op, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	query, args, err := sq.Delete(reviewersTableName).Where(sq.Eq{"pull_request_id": prID, "reviewer_id": reviewerID}).ToSql()
	if err != nil {
		return fail(err)
	}

	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	return nil
}

func (r *Repository) AddReviewer(ctx context.Context, prID, reviewerID string) error {
	const op = "pull_requests.Repository.AddReviewer"

	fail := func(err error) error {
		return fmt.Errorf("%s: %v", op, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	query, args, err := sq.Insert(reviewersTableName).Columns("pull_request_id", "reviewer_id").Values(prID, reviewerID).ToSql()
	if err != nil {
		return fail(err)
	}

	if _, err = tx.ExecContext(ctx, query, args...); err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	return nil
}
