package handler

import (
	"net/http"
	"time"
	"url_shortener/internal/dto"
	"url_shortener/internal/services"
	"url_shortener/internal/utils"
	"url_shortener/internal/validation"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service services.AuthServices
}

func NewAuthHandler(service services.AuthServices) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (ah *AuthHandler) Login(ctx *gin.Context) {
	var input dto.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return

	}
	accessToken, refreshToken, expiredAt, err := ah.service.Login(ctx, input.Email, input.Password)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	resp := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiredAt:    time.Now().Add(time.Duration(expiredAt) * time.Second).Format("2006-01-02 15:04:05"),
	}

	utils.ResponseSuccess(ctx, http.StatusOK, "Login successfully!", resp)
}
func (ah *AuthHandler) Logout(ctx *gin.Context) {
	var input dto.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return

	}
	err := ah.service.Logout(ctx, input.RefreshToken)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}

	utils.ResponseSuccess(ctx, http.StatusOK, "Logout successfully!")
}
func (ah *AuthHandler) RefreshToken(ctx *gin.Context) {
	var input dto.RefreshTokenInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.ResponseWValidator(ctx, validation.HandleValidationErrors(err))
		return
	}
	accessToken, refreshToken, expired_at, err := ah.service.RefreshToken(ctx, input.RefreshToken)
	if err != nil {
		utils.ResponseError(ctx, err)
		return
	}
	resp := dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiredAt:    time.Now().Add(time.Duration(expired_at) * time.Second).Format("2006-01-02 15:04:05"),
	}
	utils.ResponseSuccess(ctx, http.StatusOK, "Refresh token generated successfully!", resp)
}
