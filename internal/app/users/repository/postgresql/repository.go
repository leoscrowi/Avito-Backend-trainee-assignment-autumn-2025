package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"log"

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

func (r *Repository) SetIsActive(ctx context.Context, userID string, isActive bool) error {
	const op = "users.Repository.SetIsActive"

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

	query, args, err := sq.Update(tableName).
		Set("is_active", isActive).
		Where(sq.Eq{"user_id": userID}).
		ToSql()

	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	if err = tx.Commit(); err != nil {
		return fail(domain.INTERNAL, "", err)
	}

	return nil
}

func (r *Repository) CreateOrUpdateUser(ctx context.Context, user *domain.User) (string, error) {
	const op = "users.Repository.CreateOrUpdateUser"

	fail := func(code domain.ErrorCode, message string, err error) (string, error) {
		log.Printf("%s: %v", op, err)
		return "", domain.NewError(code, message, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
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
		return fail(domain.INTERNAL, "internal server error", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	if err = tx.Commit(); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	return user.UserID, nil
}

func (r *Repository) FetchByTeamName(ctx context.Context, teamName string) ([]domain.TeamMember, error) {
	const op = "users.Repository.FetchByTeamName"

	fail := func(code domain.ErrorCode, message string, err error) ([]domain.TeamMember, error) {
		log.Printf("%s: %v", op, err)
		return nil, domain.NewError(code, message, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	query, args, err := sq.Select("user_id", "username", "is_active").From(tableName).Where(sq.Eq{"team_name": teamName}).ToSql()
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

	var result []domain.TeamMember
	for rows.Next() {
		var userID string
		var userName string
		var isActive bool

		if err = rows.Scan(&userID, &userName, &isActive); err != nil {
			return fail(domain.INTERNAL, "internal server error", err)
		}

		pr := domain.TeamMember{
			UserID:   userID,
			UserName: userName,
			IsActive: isActive,
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

func (r *Repository) GetActiveUsersIDByTeam(ctx context.Context, teamName string) ([]string, error) {
	const op = "users.Repository.GetActiveUsersIDByTeam"

	fail := func(code domain.ErrorCode, message string, err error) ([]string, error) {
		log.Printf("%s: %v", op, err)
		return nil, domain.NewError(code, message, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	query, args, err := sq.Select("user_id").From(tableName).
		Where(sq.Eq{"team_name": teamName, "is_active": true}).ToSql()
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
		var userID string
		if err = rows.Scan(&userID); err != nil {
			return fail(domain.INTERNAL, "internal server error", err)
		}

		result = append(result, userID)
	}

	if err = rows.Err(); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	if err = tx.Commit(); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	return result, nil
}

func (r *Repository) FetchByID(ctx context.Context, userID string) (domain.User, error) {
	const op = "users.Repository.FetchByID"

	fail := func(code domain.ErrorCode, message string, err error) (domain.User, error) {
		log.Printf("%s: %v", op, err)
		return domain.User{}, domain.NewError(code, message, err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}
	defer func(tx *sqlx.Tx) {
		_ = tx.Rollback()
	}(tx)

	query, args, err := sq.Select("user_id", "username", "team_name", "is_active").From(tableName).Where(sq.Eq{"user_id": userID}).ToSql()
	if err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	var user domain.User
	if err = tx.GetContext(ctx, &user, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fail(domain.NOT_FOUND, "resource not found", err)
		}
		return fail(domain.INTERNAL, "internal server error", err)
	}

	if err = tx.Commit(); err != nil {
		return fail(domain.INTERNAL, "internal server error", err)
	}

	return user, nil
}
