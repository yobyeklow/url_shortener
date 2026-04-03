package services

import (
	"url_shortener/internal/database/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserServices interface {
	CreateUser(ctx *gin.Context, userParms sqlc.CreateUserParams) (sqlc.User, error)
	UpdateUser(ctx *gin.Context, userParams sqlc.UpdateUserParams) (sqlc.User, error)
	SoftDeleteUser(ctx *gin.Context, userUuid uuid.UUID) (sqlc.User, error)
	GetUserByUUID(ctx *gin.Context, userUuid uuid.UUID) (sqlc.User, error)
	GetAllUser(ctx *gin.Context, search string, page int32, limit int32, orderBy string, sort string, deleted bool) ([]sqlc.User, int32, error)
	CleanSoftDelete(ctx *gin.Context, userUuid uuid.UUID) error
	RestoreUser(ctx *gin.Context, userUuid uuid.UUID) (sqlc.User, error)
}
