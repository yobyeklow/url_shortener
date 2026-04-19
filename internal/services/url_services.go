package services

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"url_shortener/internal/database/sqlc"
	"url_shortener/internal/repository"
	"url_shortener/internal/utils"
	"url_shortener/pkg/cache"
	"url_shortener/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type urlService struct {
	url_repo repository.UrlRepository
	cache    cache.RedisService
}

var (
	randPrefixLen = utils.GetEnvInt("RAND_PREFIX_LEN", 2)
	randomKeyLen  = utils.GetEnvInt("RANDOM_KEY_LEN", 4)
)
var cacheData struct {
	Url sqlc.Url `json:"url"`
}

func NewUrlService(repo repository.UrlRepository, redisClient *redis.Client) UrlServices {
	return &urlService{
		url_repo: repo,
		cache:    cache.NewRedisCacheService(redisClient),
	}
}
func (us *urlService) CreateUrl(ctx *gin.Context, arg sqlc.CreateUrlParams) (sqlc.Url, string, bool, error) {
	var shortKey string
	context := ctx.Request.Context()
	urlData, err := us.url_repo.CreateUrl(context, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			urlExisted, err := us.url_repo.FindUrlByHashed(context, *arg.HashedValueUrl)
			if err != nil {
				return sqlc.Url{}, "", false, utils.WrapError("Race condition: URL not found after conflict", utils.ErrCodeConflict, err)
			}
			shortKey = mergeShortKey(urlExisted.RandomKey, int32(urlExisted.UrlID))
			return urlExisted, shortKey, true, nil
		}

		return sqlc.Url{}, "", false, utils.WrapError("Failed to create url", utils.ErrCodeInternal, err)
	}
	shortKey = mergeShortKey(urlData.RandomKey, int32(urlData.UrlID))

	return urlData, shortKey, false, nil
}
func (us *urlService) FindUrlById(ctx *gin.Context, id int32) (sqlc.Url, error) {
	context := ctx.Request.Context()
	urlData, err := us.url_repo.FindUrlById(context, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sqlc.Url{}, utils.NewError("URL not existed!", utils.ErrCodeNotFound)
		}
		return sqlc.Url{}, utils.WrapError("Failed to fetch URL", utils.ErrCodeInternal, err)
	}
	return urlData, nil
}
func mergeShortKey(randKey string, id int32) string {
	prefix := randKey[:randPrefixLen]
	suffix := randKey[randPrefixLen:randomKeyLen]

	encoded := utils.Base62Encode(id)

	return prefix + encoded + suffix
}
func splitShortKey(shortKey string) (randKey string, id int32, err error) {
	if len(shortKey) < randomKeyLen+1 {
		return "", 0, errors.New("shortKey too short")
	}

	// Split prefix, suffix, encoded
	prefix := shortKey[:randPrefixLen]
	suffix := shortKey[len(shortKey)-(randomKeyLen-randPrefixLen):]
	encoded := shortKey[randPrefixLen : len(shortKey)-(randomKeyLen-randPrefixLen)]

	// Decode
	decodedID, err := utils.DecodeBase62(encoded)
	if err != nil {
		return "", 0, fmt.Errorf("URL not found:%s", err)
	}

	randKey = prefix + suffix

	return randKey, decodedID, nil
}
func (us *urlService) DecryptShortKey(ctx *gin.Context, shortKey string) (sqlc.Url, error) {
	context := ctx.Request.Context()
	//Get data from Redis
	cacheKey := fmt.Sprintf("urls:%s", shortKey)
	if err := us.cache.Get(cacheKey, &cacheData); err == nil && cacheData.Url.IsActive {
		return cacheData.Url, nil
	}

	randKey, id, err := splitShortKey(shortKey)
	if err != nil {
		return sqlc.Url{}, utils.NewError("URL not found!", utils.ErrCodeNotFound)
	}

	urlData, err := us.url_repo.FindUrlById(context, id)
	if err != nil {
		return sqlc.Url{}, utils.NewError("URL not found!", utils.ErrCodeNotFound)
	}

	if urlData.RandomKey != randKey {
		return sqlc.Url{}, utils.NewError("URL not found!", utils.ErrCodeNotFound)
	}

	//Save URL to Redis
	cacheData = struct {
		Url sqlc.Url `json:"url"`
	}{
		Url: urlData,
	}
	if err = us.cache.Set(cacheKey, cacheData, 5*time.Minute); err != nil {

		logger.Log.Warn().Err(err).Msg("Failed to save data to Redis")
	}

	return urlData, nil
}
