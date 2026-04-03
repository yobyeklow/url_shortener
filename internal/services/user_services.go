package services

import (
	"database/sql"
	"errors"
	"url_shortener/internal/database/sqlc"
	"url_shortener/internal/repository"
	"url_shortener/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	user_repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserServices {
	return &userService{
		user_repo: repo,
	}
}
func (us *userService) CreateUser(ctx *gin.Context, userParms sqlc.CreateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()
	userParms.UserEmail = utils.NormalizeString(userParms.UserEmail)

	//Generate hash password
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(userParms.UserPassword), bcrypt.DefaultCost)
	if err != nil {
		return sqlc.User{}, utils.WrapError("Failed to hash password", utils.ErrCodeInternal, err)
	}
	userParms.UserPassword = string(hashed_password)
	userData, err := us.user_repo.CreateUser(context, userParms)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return sqlc.User{}, utils.NewError("Email already existed", utils.ErrCodeConflict)
		}
		return sqlc.User{}, utils.WrapError("Failed to create new user", utils.ErrCodeInternal, err)
	}

	return userData, nil
}
func (us *userService) UpdateUser(ctx *gin.Context, userParams sqlc.UpdateUserParams) (sqlc.User, error) {
	context := ctx.Request.Context()

	if *userParams.UserPassword != "" && userParams.UserPassword != nil {
		hashed_pass, err := bcrypt.GenerateFromPassword([]byte(*userParams.UserPassword), bcrypt.DefaultCost)
		if err != nil {
			return sqlc.User{}, utils.WrapError("Failed to hash password", utils.ErrCodeInternal, err)
		}
		hash := string(hashed_pass)
		userParams.UserPassword = &hash
	}

	updatedUser, err := us.user_repo.UpdateUser(context, userParams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("User not existed!", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError("Failed to update user", utils.ErrCodeInternal, err)
	}

	return updatedUser, nil
}
func (us *userService) SoftDeleteUser(ctx *gin.Context, userUuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()

	userDeleted, err := us.user_repo.SoftDeleteUser(context, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("User not existed!", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError("Failed to delete user", utils.ErrCodeInternal, err)
	}
	return userDeleted, nil
}
func (us *userService) CleanSoftDelete(ctx *gin.Context, userUuid uuid.UUID) error {
	context := ctx.Request.Context()
	_, err := us.user_repo.CleanSoftDelete(context, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NewError("User not existed!", utils.ErrCodeNotFound)
		}
		return utils.WrapError("Failed to delete user", utils.ErrCodeInternal, err)
	}
	return nil
}
func (us *userService) GetUserByUUID(ctx *gin.Context, userUuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()
	userData, err := us.user_repo.GetUserByUUID(context, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("User not existed!", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError("Failed to fetch user by uuid", utils.ErrCodeInternal, err)
	}
	return userData, nil
}
func (us *userService) RestoreUser(ctx *gin.Context, userUuid uuid.UUID) (sqlc.User, error) {
	context := ctx.Request.Context()
	userData, err := us.user_repo.RestoreUser(context, userUuid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.User{}, utils.NewError("User not existed!", utils.ErrCodeNotFound)
		}
		return sqlc.User{}, utils.WrapError("Failed to restore user", utils.ErrCodeInternal, err)
	}
	return userData, nil
}
func (us *userService) GetAllUser(ctx *gin.Context, search string, page int32, limit int32, orderBy string, sort string, deleted bool) ([]sqlc.User, int32, error) {
	context := ctx.Request.Context()
	if sort == "" {
		sort = "desc"
	}
	if orderBy == "" {
		orderBy = "user_id"
	}
	if limit <= 0 {
		envLimit := utils.GetEnvInt("LIMIT_ITEM_ON_PER_PAGE", 10)
		limit = int32(envLimit)
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	usersData, err := us.user_repo.GetUserAll(context, search, page, limit, orderBy, sort, offset, deleted)
	if err != nil {
		return []sqlc.User{}, 0, utils.WrapError("Failed to fetch users", utils.ErrCodeInternal, err)
	}
	total, err := us.user_repo.CountUsers(context, search, deleted)
	if err != nil {
		return []sqlc.User{}, 0, utils.WrapError("Failed to count users", utils.ErrCodeInternal, err)
	}
	return usersData, int32(total), nil
}
