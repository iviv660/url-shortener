package repository

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache miss")

type repoRedis struct {
	client *redis.Client
}

func NewRepoRedis(client *redis.Client) *repoRedis {
	return &repoRedis{client: client}
}

func (r *repoRedis) Set(ctx context.Context, short, original string, ttl time.Duration) error {
	return r.client.Set(ctx, short, original, ttl).Err()
}

func (r *repoRedis) Get(ctx context.Context, short string) (string, error) {
	val, err := r.client.Get(ctx, short).Result()
	if err == redis.Nil {
		return "", ErrCacheMiss
	}
	return val, err
}

func (r *repoRedis) Delete(ctx context.Context, short string) error {
	return r.client.Del(ctx, short).Err()
}
