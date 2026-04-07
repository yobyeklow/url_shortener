package handler

import (
	"net/http"
	"strings"
	"url_shortener/internal/dto"
	"url_shortener/internal/services"
	"url_shortener/internal/utils"
	"url_shortener/internal/validation"
	"url_shortener/pkg/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UrlsHandler struct {
	service    services.UrlServices
	jwtService auth.JWTService
}

func NewUrlHandler(service services.UrlServices) *UrlsHandler {
	return &UrlsHandler{
		service: service,
	}
}

func (uh *UrlsHandler) CreateShortUrl(ctx *gin.Context) {
	var input dto.UrlInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	authHeader := ctx.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	payload, err := uh.jwtService.DecryptAccesTokenPayload(token)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	randKey := utils.GenerateRandomKey(4)
	userUuid, err := uuid.Parse(payload.UserUUID)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	hash_value := utils.GenerateHashedValue(&input.IosDeepLink, &input.IosFallbackUrl, &input.AndroidDeepLink, &input.AndroidFallbackUrl, input.DefaultFallbackUrl)
	urlInput := input.MapCreateInputToModel(randKey, &hash_value, userUuid)

	urlData, existed, err := uh.service.CreateUrl(ctx, urlInput)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	shortKey := uh.service.MergeShortKey(urlData.RandomKey, int32(urlData.UrlID))

	if existed {
		utils.ResponseSuccess(ctx, http.StatusOK, "Short URL already exists", map[string]string{"short_key": shortKey})
	} else {
		utils.ResponseSuccess(ctx, http.StatusCreated, "Short URL created", map[string]string{"short_key": shortKey})
	}
}
