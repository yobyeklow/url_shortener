-- name: CreateUrl :one
INSERT INTO urls(
    user_uuid,random_key, ios_deep_link, ios_fallback_url,
    android_deep_link, android_fallback_url, default_fallback_url,
    hashed_value_url, webhook_url, opengraph_title, opengraph_description,
    opengraph_image, is_active, url_created_at, url_updated_at
)
VALUES($13,$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())
ON CONFLICT (hashed_value_url) WHERE url_deleted_at IS NULL
DO NOTHING
RETURNING *;
-- name: FindUrlByHashed :one
SELECT *
FROM urls
WHERE hashed_value_url = sqlc.arg(hashed_value_url)::TEXT AND url_deleted_at IS NULL
LIMIT 1;
