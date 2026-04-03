package repository

import (
	"context"
	"fmt"
	"strings"
	"url_shortener/internal/database"
	"url_shortener/internal/database/sqlc"

	"github.com/google/uuid"
)

type SQLUserRepository struct {
	db sqlc.Querier
}

func NewSQLUserRepository(db sqlc.Querier) UserRepository {
	return &SQLUserRepository{
		db: db,
	}
}
func (ur *SQLUserRepository) CreateUser(ctx context.Context, userParams sqlc.CreateUserParams) (sqlc.User, error) {
	userData, err := ur.db.CreateUser(ctx, userParams)
	if err != nil {
		return sqlc.User{}, err
	}
	return userData, nil
}
func (ur *SQLUserRepository) UpdateUser(ctx context.Context, userParams sqlc.UpdateUserParams) (sqlc.User, error) {
	userData, err := ur.db.UpdateUser(ctx, userParams)
	if err != nil {
		return sqlc.User{}, err
	}
	return userData, nil
}
func (ur *SQLUserRepository) RestoreUser(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error) {
	userData, err := ur.db.RestoreUser(ctx, userUuid)
	if err != nil {
		return sqlc.User{}, err
	}
	return userData, nil
}
func (ur *SQLUserRepository) GetUserByUUID(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error) {
	userData, err := ur.db.GetUserByUUID(ctx, userUuid)
	if err != nil {
		return sqlc.User{}, err
	}
	return userData, nil
}
func (ur *SQLUserRepository) SoftDeleteUser(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error) {
	userData, err := ur.db.SoftDeleteUser(ctx, userUuid)
	if err != nil {
		return sqlc.User{}, err
	}
	return userData, nil
}
func (ur *SQLUserRepository) CleanSoftDelete(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error) {
	userData, err := ur.db.CleanSoftDelete(ctx, userUuid)
	if err != nil {
		return sqlc.User{}, err
	}
	return userData, nil
}
func (ur *SQLUserRepository) GetUserAll(ctx context.Context, search string, page int32, limit int32, orderBy string, sort string, offset int32, deleted bool) ([]sqlc.User, error) {
	query := `SELECT *
			FROM users
			WHERE ( $1::TEXT IS NULL
				OR $1::TEXT = ''
				OR user_email ILIKE '%' || $1 || '%' )`

	if deleted {
		query += " AND user_deleted_at IS NOT NULL"
	} else {
		query += " AND user_deleted_at IS NULL"
	}

	switch orderBy {
	case "user_id", "user_created_at":
		query += fmt.Sprintf(" ORDER BY %s %s", orderBy, strings.ToUpper(sort))
	default:
		query += " ORDER BY user_id ASC"
	}

	query += " LIMIT $2 OFFSET $3"

	rows, err := database.DBPool.Query(ctx, query, search, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []sqlc.User{}
	for rows.Next() {
		var i sqlc.User
		if err := rows.Scan(
			&i.UserID,
			&i.UserUuid,
			&i.UserEmail,
			&i.UserPassword,
			&i.UserStatus,
			&i.UserRole,
			&i.UserCreatedAt,
			&i.UserUpdatedAt,
			&i.UserDeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
func (ur *SQLUserRepository) CountUsers(ctx context.Context, search string, deleted bool) (int64, error) {
	var params sqlc.CountRecordsParams
	params.Search = search
	params.Deleted = &deleted
	total, err := ur.db.CountRecords(ctx, params)
	if err != nil {
		return 0, err
	}
	return total, nil
}
