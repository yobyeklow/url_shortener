package handler

import (
	"net/http"
	"net/url"
	"strings"
	"url_shortener/internal/database/sqlc"
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

	_, shortKey, existed, err := uh.service.CreateUrl(ctx, urlInput)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	if existed {
		utils.ResponseSuccess(ctx, http.StatusOK, "Short URL already exists", map[string]string{"short_key": shortKey})
	} else {
		utils.ResponseSuccess(ctx, http.StatusCreated, "Short URL created", map[string]string{"short_key": shortKey})
	}
}
func (uh *UrlsHandler) DecryptShortKey(ctx *gin.Context) {
	var input dto.ParamsShortKey
	if err := ctx.ShouldBindUri(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	urlData, err := uh.service.DecryptShortKey(ctx, input.ShortKey)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	user_agent := ctx.GetHeader("User-Agent")
	target := chooseTargetURL(urlData, user_agent)
	ctx.Redirect(http.StatusMovedPermanently, target)
}

func chooseTargetURL(data sqlc.Url, ua string) string {
	isIOS := strings.Contains(ua, "iPhone") || strings.Contains(ua, "iPad")
	isAndroid := strings.Contains(ua, "Android")

	if isIOS && data.IosDeepLink != "" {
		return appendFallback(data.IosDeepLink, data.IosFallbackUrl)
	}
	if isAndroid && data.AndroidDeepLink != "" {
		return appendFallback(data.AndroidDeepLink, data.AndroidFallbackUrl)
	}
	if isIOS && data.IosFallbackUrl != "" {
		return data.IosFallbackUrl
	}
	if isAndroid && data.AndroidFallbackUrl != "" {
		return data.AndroidFallbackUrl
	}
	return data.DefaultFallbackUrl
}
func appendFallback(base, fallback string) string {
	if base == "" {
		return fallback
	}
	if strings.Contains(base, "?") {
		return base + "&fallback=" + url.QueryEscape(fallback)
	}
	return base + "?fallback=" + url.QueryEscape(fallback)
}
