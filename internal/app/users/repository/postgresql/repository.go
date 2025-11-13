package postgresql

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/leoscrowi/pr-assignment-service/domain"

	sq "github.com/Masterminds/squirrel"
)

const tableName = "users"

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SetIsActive(ctx context.Context, userID uuid.UUID, isActive bool) error {
	// TODO: implement
	panic("implement me")
}

func (r *Repository) CreateOrUpdateUser(ctx context.Context, user domain.User) (uuid.UUID, error) {
	const op = "users.Repository.CreateOrUpdateUser"
	fail := func(err error) (uuid.UUID, error) {
		return uuid.Nil, fmt.Errorf("%s: %v", op, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	query, args, err := sq.Insert(tableName).
		Columns("user_id", "username", "team_name", "is_active").
		Values(user.UserID, user.Username, user.TeamName, user.IsActive).
		Suffix("ON CONFLICT (user_id) DO UPDATE SET " +
			"username = EXCLUDED.username, " +
			"team_name = EXCLUDED.team_name, " +
			"is_active = EXCLUDED.is_active").
		ToSql()

	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	return user.UserID, nil
}
