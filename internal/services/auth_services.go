package services

import (
	"strings"
	"time"
	"url_shortener/internal/repository"
	"url_shortener/internal/utils"
	"url_shortener/pkg/auth"
	"url_shortener/pkg/cache"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authServices struct {
	repo         repository.UserRepository
	tokenService auth.TokenService
	cache        cache.RedisService
}

func NewAuthServices(repo repository.UserRepository, tokenService auth.TokenService, cache cache.RedisService) AuthServices {
	return &authServices{
		repo:         repo,
		tokenService: tokenService,
		cache:        cache,
	}
}
func (as *authServices) Login(ctx *gin.Context, email string, password string) (string, string, int, error) {
	email = utils.NormalizeString(email)
	userData, err := as.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", 0, utils.NewError("Invalid  email or password", utils.ErrCodeUnauthorized)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userData.UserPassword), []byte(password)); err != nil {
		return "", "", 0, utils.NewError("Invalid  email or password", utils.ErrCodeUnauthorized)
	}

	accessToken, err := as.tokenService.GenerateAccessToken(userData)
	if err != nil {
		return "", "", 0, utils.NewError("Failed to generate access token", utils.ErrCodeUnauthorized)
	}
	refreshToken, err := as.tokenService.GenerateRefreshToken(userData)
	if err != nil {
		return "", "", 0, utils.NewError("Failed to generate refresh token", utils.ErrCodeUnauthorized)
	}
	if err := as.tokenService.StoreRefreshTokenToRedis(refreshToken); err != nil {
		return "", "", 0, utils.NewError("Can't save refresh token", utils.ErrCodeUnauthorized)
	}
	return accessToken, refreshToken.Token, int(auth.AccessTokenTTL.Seconds()), nil
}
func (as *authServices) Logout(ctx *gin.Context, refreshTokenStr string) error {
	//Get accessToken from header
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return utils.NewError("Missing Authorization Header", utils.ErrCodeUnauthorized)
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")
	_, claims, err := as.tokenService.ParseToken(accessToken)
	if err != nil {
		return utils.NewError("Invalid Access Token", utils.ErrCodeUnauthorized)
	}

	if jti, ok := claims["jti"].(string); ok {
		expUnix, _ := claims["exp"].(float64)

		exp := time.Unix(int64(expUnix), 0)
		key := "blacklist:" + jti
		ttl := time.Until(exp)
		// Save access token to blacklist access_token
		as.cache.Set(key, "revoked", ttl)
	}

	//Validate Refresh Token
	_, err = as.tokenService.ValidJWTToken(refreshTokenStr)
	if err != nil {
		return utils.NewError("Refresh token is invalid or revoked", utils.ErrCodeUnauthorized)
	}
	//Revoke old refreshToken
	err = as.tokenService.RevokeToken(refreshTokenStr)
	if err != nil {
		return utils.NewError("Unable to revoke refresh token", utils.ErrCodeUnauthorized)
	}
	return nil
}
func (as *authServices) RefreshToken(ctx *gin.Context, refreshTokenStr string) (string, string, int, error) {
	context := ctx.Request.Context()
	// Validate  Token and get user uuid
	token, err := as.tokenService.ValidJWTToken(refreshTokenStr)
	if err != nil {
		return "", "", 0, utils.NewError("Refresh token is invalid or revoked", utils.ErrCodeUnauthorized)
	}

	user_uuid, _ := uuid.Parse(token.User_uuid)
	//Get user data from user uuid
	userData, err := as.repo.GetUserByUUID(context, user_uuid)
	if err != nil {
		return "", "", 0, utils.NewError("User not found", utils.ErrCodeUnauthorized)
	}
	//Create new  accessToken and new RefreshToken
	access_token, err := as.tokenService.GenerateAccessToken(userData)
	if err != nil {
		return "", "", 0, utils.NewError("Failed to generate access token", utils.ErrCodeUnauthorized)
	}

	refresh_token, err := as.tokenService.GenerateRefreshToken(userData)
	if err != nil {
		return "", "", 0, utils.NewError("Failed to generate refresh token", utils.ErrCodeUnauthorized)
	}
	//Revoke old refresh token
	if err := as.tokenService.RevokeToken(refreshTokenStr); err != nil {
		return "", "", 0, utils.NewError("Unable to revoke refresh token", utils.ErrCodeInternal)
	}
	// Save new refresh token to Redis
	if err := as.tokenService.StoreRefreshTokenToRedis(refresh_token); err != nil {
		return "", "", 0, utils.NewError("Can't save refresh token", utils.ErrCodeUnauthorized)
	}
	return access_token, refresh_token.Token, int(auth.AccessTokenTTL.Seconds()), nil
}
