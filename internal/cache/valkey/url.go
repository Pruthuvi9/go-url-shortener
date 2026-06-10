package valkey

import (
	"context"
	"time"
)

// Client is a minimal interface over a Valkey/Redis client.
// Wire in your preferred client (e.g. github.com/valkey-io/valkey-go) here.
type Client interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}

type URLCache struct {
	client Client
}

func NewURLCache(client Client) *URLCache {
	return &URLCache{client: client}
}

func (c *URLCache) Get(ctx context.Context, shortCode string) (string, error) {
	// TODO: return client.Get, map cache-miss to ("", nil)
	return "", nil
}

func (c *URLCache) Set(ctx context.Context, shortCode string, longURL string, ttl time.Duration) error {
	// TODO: return client.Set
	return nil
}

func (c *URLCache) Delete(ctx context.Context, shortCode string) error {
	// TODO: return client.Del
	return nil
}
