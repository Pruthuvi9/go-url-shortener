package usecase

import (
	"context"
	"time"

	"github.com/pruthuvifernando/url-shortner/internal/domain"
)

type URLRepository interface {
	Save(ctx context.Context, url domain.URL) error
	FindByCode(ctx context.Context, shortCode string) (domain.URL, error)
	Delete(ctx context.Context, shortCode string) error
}

type URLCache interface {
	Get(ctx context.Context, shortCode string) (string, error)
	Set(ctx context.Context, shortCode string, longURL string, ttl time.Duration) error
	Delete(ctx context.Context, shortCode string) error
}
