package repository

import (
	"context"
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
