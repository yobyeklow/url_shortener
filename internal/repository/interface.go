package repository

import (
	"context"
	"url_shortener/internal/database/sqlc"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, userParams sqlc.CreateUserParams) (sqlc.User, error)
	UpdateUser(ctx context.Context, userParams sqlc.UpdateUserParams) (sqlc.User, error)
	RestoreUser(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error)
	GetUserByUUID(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error)
	SoftDeleteUser(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error)
	CleanSoftDelete(ctx context.Context, userUuid uuid.UUID) (sqlc.User, error)
}
