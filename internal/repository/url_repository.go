package repository

import (
	"context"
	"url_shortener/internal/database/sqlc"
)

type SQLUrlRepository struct {
	db sqlc.Querier
}

func NewSQLUrlRepository(db sqlc.Querier) UrlRepository {
	return &SQLUrlRepository{
		db: db,
	}
}
func (ur *SQLUrlRepository) CreateUrl(ctx context.Context, arg sqlc.CreateUrlParams) (sqlc.Url, error) {
	urlData, err := ur.db.CreateUrl(ctx, arg)
	if err != nil {
		return sqlc.Url{}, err
	}
	return urlData, nil
}
func (ur *SQLUrlRepository) FindUrlByHashed(ctx context.Context, hashedValueUrl string) (sqlc.Url, error) {
	urlData, err := ur.db.FindUrlByHashed(ctx, hashedValueUrl)
	if err != nil {
		return sqlc.Url{}, err
	}
	return urlData, nil
}
