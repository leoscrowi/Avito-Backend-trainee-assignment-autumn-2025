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
	const op = "users.Repository.SetIsActive"

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

	query, args, err := sq.Update(tableName).
		Set("is_active", isActive).
		Where(sq.Eq{"user_id": userID}).
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

	return nil
}

func (r *Repository) CreateOrUpdateUser(ctx context.Context, user *domain.User) (uuid.UUID, error) {
	const op = "users.Repository.CreateOrUpdateUser"

	if user == nil {
		return uuid.Nil, fmt.Errorf("%s: user is nil", op)
	}

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

func (r *Repository) FetchByTeamName(ctx context.Context, teamName string) ([]domain.TeamMember, error) {
	const op = "users.Repository.FetchByTeamName"

	fail := func(err error) ([]domain.TeamMember, error) {
		return []domain.TeamMember{}, fmt.Errorf("%s: %v", op, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	query, args, err := sq.Select("user_id", "username", "is_active").From(tableName).Where(sq.Eq{"team_name": teamName}).ToSql()
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

	var result []domain.TeamMember
	for rows.Next() {
		var userID uuid.UUID
		var userName string
		var isActive bool

		if err = rows.Scan(&userID, &userName, &isActive); err != nil {
			return fail(err)
		}

		pr := domain.TeamMember{
			UserID:   userID,
			UserName: userName,
			IsActive: isActive,
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
