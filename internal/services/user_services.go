package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
	"url_shortener/internal/database/sqlc"
	"url_shortener/internal/repository"
	"url_shortener/internal/utils"
	"url_shortener/pkg/cache"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	user_repo repository.UserRepository
	cache     cache.RedisService
}

func NewUserService(repo repository.UserRepository, redisClient *redis.Client) UserServices {
	return &userService{
		user_repo: repo,
		cache:     cache.NewRedisCacheService(redisClient),
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

	//Clear cache Redis
	if err := us.cache.Clear("users:*"); err != nil {
		log.Printf("Failed to clear cache: %v", err)
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

	//Clear cache Redis
	if err := us.cache.Clear("users:*"); err != nil {
		log.Printf("Failed to clear cache: %v", err)
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
	//Clear cache Redis
	if err := us.cache.Clear("users:*"); err != nil {
		log.Printf("Failed to clear cache: %v", err)
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
	//Clear cache Redis
	if err := us.cache.Clear("users:*"); err != nil {
		log.Printf("Failed to clear cache: %v", err)
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
	//Clear cache Redis
	if err := us.cache.Clear("users:*"); err != nil {
		log.Printf("Failed to clear cache: %v", err)
	}
	return userData, nil
}
func (us *userService) GetAllUser(ctx *gin.Context, search string, page int32, limit int32, orderBy string, sort string, deleted bool) ([]sqlc.User, int32, error) {
	context := ctx.Request.Context()
	//Get data from Redit
	cacheKey := generateCacheKey(search, page, limit, orderBy, sort, deleted)
	var cacheData struct {
		Users []sqlc.User `json:"users"`
		Total int32       `json:"total"`
	}
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
	//Get data to Redit
	cacheKey = generateCacheKey(search, page, limit, orderBy, sort, deleted)
	cacheData = struct {
		Users []sqlc.User `json:"users"`
		Total int32       `json:"total"`
	}{
		Users: usersData,
		Total: int32(total),
	}
	if err = us.cache.Set(cacheKey, cacheData, 5*time.Minute); err != nil {
		log.Printf("Failed to save data to Redis: %v", err)
	}
	return usersData, int32(total), nil
}
func (us *userService) GetUserByEmail(ctx *gin.Context, email string) (sqlc.User, error) {
	context := ctx.Request.Context()
	userData, err := us.user_repo.GetUserByEmail(context, email)
	if err != nil {
		return sqlc.User{}, utils.WrapError("Failed to fetch user by email", utils.ErrCodeInternal, err)
	}
	return userData, nil
}
func generateCacheKey(search string, page int32, limit int32, order_by string, sort string, deleted bool) string {
	search = strings.TrimSpace(search)
	if search == "" {
		search = "none"
	}
	order_by = strings.TrimSpace(order_by)
	if order_by == "" {
		order_by = "user_created_at"
	}
	sort = strings.ToLower(strings.TrimSpace(sort))
	if sort == "" {
		sort = "desc"
	}
	return fmt.Sprintf("users:%s:%d:%d:%s:%s:%t", search, page, limit, order_by, sort, deleted)
}
