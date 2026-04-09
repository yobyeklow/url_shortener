package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
	"url_shortener/internal/database/sqlc"
	"url_shortener/internal/utils"
	"url_shortener/pkg/cache"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	cache cache.RedisService
}
type EncryptedPayload struct {
	UserUUID string `json:"user_uuid"`
	Email    string `json:"user_email"`
	Role     int32  `json:"user_role"`
}
type RefreshToken struct {
	Token     string    `json:"token"`
	User_uuid string    `json:"user_uuid"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
}

var (
	jwt_secret_key  = []byte(utils.GetEnv("JWT_SECRET_KEY", "01594057772770059127161404529220"))
	jwt_encrypt_key = []byte(utils.GetEnv("JWT_ENCRYPT_KEY", "01594057772770059127161404523216"))
)

const (
	AccessTokenTTL  = 40 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
)

func NewJWTService(cache cache.RedisService) TokenService {
	return &JWTService{
		cache: cache,
	}
}
func (js *JWTService) GenerateAccessToken(user sqlc.User) (string, error) {
	payload := EncryptedPayload{
		UserUUID: user.UserUuid.String(),
		Email:    user.UserEmail,
		Role:     user.UserRole,
	}

	raw_data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("Parse Raw Data: %w", err)
	}
	encrypted, err := utils.EncrypAES(raw_data, jwt_encrypt_key)
	if err != nil {
		return "", fmt.Errorf("Encrypt payload: %w", err)
	}
	claims := jwt.MapClaims{
		"data": encrypted,
		"jti":  uuid.NewString(),
		"exp":  jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
		"iat":  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwt_secret_key)
}
func (js *JWTService) ParseToken(tokenString string) (*jwt.Token, jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return jwt_secret_key, nil
	})
	if err != nil || !token.Valid {
		return nil, nil, utils.NewError("Invalid token", utils.ErrCodeUnauthorized)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, nil, utils.NewError("Invalid token claims", utils.ErrCodeUnauthorized)
	}
	return token, claims, nil
}
func (js *JWTService) DecryptAccesTokenPayload(tokenString string) (*EncryptedPayload, error) {
	_, claims, err := js.ParseToken(tokenString)
	if err != nil {
		return nil, utils.WrapError("Can't parse JWT token", utils.ErrCodeInternal, err)
	}
	encryptedData, ok := claims["data"].(string)
	if !ok {
		return nil, utils.NewError("Encoded data not found", utils.ErrCodeUnauthorized)
	}
	decryptBytes, err := utils.DecryptAES(encryptedData, jwt_encrypt_key)
	if err != nil {
		return nil, utils.WrapError("Can't decode data", utils.ErrCodeInternal, err)
	}

	var payload EncryptedPayload
	if err := json.Unmarshal(decryptBytes, &payload); err != nil {
		return nil, utils.WrapError("Invalid data format", utils.ErrCodeInternal, err)
	}
	return &payload, nil
}
func (js *JWTService) GenerateRefreshToken(user sqlc.User) (RefreshToken, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return RefreshToken{}, err
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)
	return RefreshToken{
		Token:     token,
		User_uuid: user.UserUuid.String(),
		ExpiresAt: time.Now().Add(RefreshTokenTTL),
		Revoked:   false,
	}, nil
}
func (js *JWTService) StoreRefreshTokenToRedis(token RefreshToken) error {
	cacheKey := "refresh_token:" + token.Token
	return js.cache.Set(cacheKey, token, RefreshTokenTTL)
}
func (js *JWTService) ValidJWTToken(token string) (RefreshToken, error) {
	cacheKey := "refresh_token:" + token
	var refreshToken RefreshToken
	err := js.cache.Get(cacheKey, &refreshToken)
	if err != nil || refreshToken.Revoked == true || refreshToken.ExpiresAt.Before(time.Now()) {
		return RefreshToken{}, utils.WrapError("Can't get refresh token", utils.ErrCodeInternal, err)
	}
	return refreshToken, nil
}
func (js *JWTService) RevokeToken(token string) error {
	cacheKey := "refresh_token:" + token
	var refreshToken RefreshToken
	err := js.cache.Get(cacheKey, &refreshToken)
	if err != nil {
		return utils.WrapError("Can't revoke token", utils.ErrCodeInternal, err)
	}
	refreshToken.Revoked = true
	return js.cache.Set(cacheKey, refreshToken, time.Until(refreshToken.ExpiresAt))
}
