package auth

import (
	"url_shortener/internal/database/sqlc"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateAccessToken(user sqlc.User) (string, error)
	ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error)
	DecryptAccesTokenPayload(tokenString string) (*EncryptedPayload, error)
	GenerateRefreshToken(user sqlc.User) (RefreshToken, error)
	StoreRefreshTokenToRedis(token RefreshToken) error
	ValidJWTToken(token string) (RefreshToken, error)
	RevokeToken(token string) error
}
