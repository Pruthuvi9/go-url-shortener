package usecase

import (
	"context"
	"time"

	"github.com/pruthuvifernando/url-shortner/internal/domain"
)

type URLUsecase struct {
	repo  URLRepository
	cache URLCache
}

func NewURLUsecase(repo URLRepository, cache URLCache) *URLUsecase {
	return &URLUsecase{repo: repo, cache: cache}
}

type ShortenRequest struct {
	LongURL        string
	CustomAlias    string
	ExpirationDate *time.Time
}

func (u *URLUsecase) Shorten(ctx context.Context, req ShortenRequest) (domain.URL, error) {
	// TODO: validate URL, generate/validate short code, persist, cache
	return domain.URL{}, nil
}

func (u *URLUsecase) Redirect(ctx context.Context, shortCode string) (string, error) {
	// TODO: cache lookup -> DB lookup -> expiry check
	return "", nil
}

func (u *URLUsecase) Delete(ctx context.Context, shortCode string) error {
	// TODO: delete from DB then invalidate cache
	return nil
}
