package services

import (
	"database/sql"
	"errors"
	"url_shortener/internal/database/sqlc"
	"url_shortener/internal/repository"
	"url_shortener/internal/utils"
	"url_shortener/pkg/cache"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type urlService struct {
	url_repo repository.UrlRepository
	cache    cache.RedisService
}

func NewUrlService(repo repository.UrlRepository, redisClient *redis.Client) UrlServices {
	return &urlService{
		url_repo: repo,
		cache:    cache.NewRedisCacheService(redisClient),
	}
}
func (us *urlService) CreateUrl(ctx *gin.Context, arg sqlc.CreateUrlParams) (sqlc.Url, bool, error) {
	context := ctx.Request.Context()
	urlData, err := us.url_repo.CreateUrl(context, arg)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			urlExisted, err := us.url_repo.FindUrlByHashed(context, *arg.HashedValueUrl)
			if err != nil {
				return sqlc.Url{}, false, utils.WrapError("Race condition: URL not found after conflict", utils.ErrCodeConflict, err)
			}
			return urlExisted, true, nil
		}

		return sqlc.Url{}, false, utils.WrapError("Failed to create url", utils.ErrCodeInternal, err)
	}
	return urlData, false, nil
}

func (us *urlService) MergeShortKey(randKey string, id int32) string {
	encoded := utils.Base62Encode(id)
	randPrefixLen := utils.GetEnvInt("RAND_PREFIX_LEN", 2)
	randomKeyLen := utils.GetEnvInt("RANDOM_KEY_LEN", 4)
	var prefix string
	if len(randKey) >= randPrefixLen {
		prefix = randKey[:randPrefixLen]
	} else {
		prefix = randKey
	}

	var suffix string
	if len(randKey) >= randomKeyLen {
		suffix = randKey[randPrefixLen:randomKeyLen]
	} else if len(randKey) > randPrefixLen {
		suffix = randKey[randPrefixLen:]
	} else {
		suffix = ""
	}

	return prefix + encoded + suffix
}
