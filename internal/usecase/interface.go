package usecase

import (
	"context"
	"time"
)

type RepoPostgres interface {
	SaveURL(ctx context.Context, short, original string) error
	FindByShort(ctx context.Context, short string) (string, error)
	FindByOriginal(ctx context.Context, original string) (string, error)
}

type RepoRedis interface {
	Set(ctx context.Context, short, original string, ttl time.Duration) error
	Get(ctx context.Context, short string) (string, error) // возвращает original
	Delete(ctx context.Context, short string) error
}
