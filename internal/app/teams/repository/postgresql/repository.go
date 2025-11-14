package postgresql

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/leoscrowi/pr-assignment-service/domain"
)

const tableName = "teams"

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FetchByName(ctx context.Context, teamName string) (domain.Team, error) {
	const op = "teams.Repository.FetchByName"

	fail := func(err error) (domain.Team, error) { return domain.Team{}, err }

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	query, args, err := sq.Select("team_name").From(tableName).Where(sq.Eq{"team_name": teamName}).ToSql()
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

	var team domain.Team
	if err = tx.GetContext(ctx, &team, query, args...); err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	return team, nil
}
