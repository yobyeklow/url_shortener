package dto

import (
	"url_shortener/internal/database/sqlc"

	"github.com/google/uuid"
)

type UrlInput struct {
	IosDeepLink          string  `json:"ios_deep_link" binding:"omitempty"`
	IosFallbackUrl       string  `json:"ios_fallback_url" binding:"omitempty"`
	AndroidDeepLink      string  `json:"android_deep_link" binding:"omitempty"`
	AndroidFallbackUrl   string  `json:"android_fallback_url" binding:"omitempty"`
	DefaultFallbackUrl   string  `json:"default_fallback_url" binding:"required"`
	WebhookUrl           string  `json:"webhook_url" binding:"omitempty"`
	OpengraphTitle       *string `json:"opengraph_title" binding:"omitempty"`
	OpengraphDescription string  `json:"opengraph_description" binding:"omitempty"`
	OpengraphImage       string  `json:"opengraph_image" binding:"omitempty"`
}
type UrlCacheData struct {
	UrlID              int32  `json:"url_id"`
	RandomKey          string `json:"random_key"`
	IosDeepLink        string `json:"ios_deep_link"`
	IosFallbackUrl     string `json:"ios_fallback_url"`
	AndroidDeepLink    string `json:"android_deep_link"`
	AndroidFallbackUrl string `json:"android_fallback_link"`
	DefaultFallbackUrl string `json:"default_fallback_url"`
	WebhookUrl         string `json:"webhook_url"`
	OgTitle            string `json:"og_title"`
	OgDescription      string `json:"og_description"`
	OgImageUrl         string `json:"og_image_url"`
	IsActive           bool   `json:"is_active"`
}
type UrlDTO struct {
	DefaultFallbackUrl string `json:"default_fallback_url"`
}
type ParamsShortKey struct {
	ShortKey string `uri:"short_key" binding:"required,short_key"`
}

func (input *UrlInput) MapCreateInputToModel(randomKey string, hashed_value *string, userUUID uuid.UUID) sqlc.CreateUrlParams {
	return sqlc.CreateUrlParams{
		IosDeepLink:          input.IosDeepLink,
		IosFallbackUrl:       input.IosFallbackUrl,
		AndroidDeepLink:      input.AndroidDeepLink,
		AndroidFallbackUrl:   input.AndroidFallbackUrl,
		DefaultFallbackUrl:   input.DefaultFallbackUrl,
		WebhookUrl:           input.WebhookUrl,
		OpengraphTitle:       input.OpengraphTitle,
		OpengraphDescription: input.OpengraphDescription,
		OpengraphImage:       input.OpengraphImage,
		RandomKey:            randomKey,
		HashedValueUrl:       hashed_value,
		IsActive:             true,
		UserUuid:             userUUID,
	}
}
