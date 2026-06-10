package domain

import (
	"errors"
	"time"
)

type URL struct {
	ShortCode string
	LongURL   string
	CreatedAt time.Time
	ExpiresAt *time.Time
}

var (
	ErrNotFound = errors.New("url not found")
	ErrExpired  = errors.New("url has expired")
	ErrConflict = errors.New("short code already exists")
)
