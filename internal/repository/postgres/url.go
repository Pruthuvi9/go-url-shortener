package postgres

import (
	"context"
	"database/sql"

	"github.com/pruthuvifernando/url-shortner/internal/domain"
)

type URLRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{db: db}
}

func (r *URLRepository) Save(ctx context.Context, url domain.URL) error {
	// TODO: INSERT INTO urls (short_code, long_url, created_at, expires_at)
	return nil
}

func (r *URLRepository) FindByCode(ctx context.Context, shortCode string) (domain.URL, error) {
	// TODO: SELECT short_code, long_url, created_at, expires_at FROM urls WHERE short_code = $1
	return domain.URL{}, nil
}

func (r *URLRepository) Delete(ctx context.Context, shortCode string) error {
	// TODO: DELETE FROM urls WHERE short_code = $1
	return nil
}
