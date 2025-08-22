package database

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

func ConnectRedis(ctx context.Context, urlRedis string) (*redis.Client, error) {
	opt, err := redis.ParseURL(urlRedis)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opt)

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	log.Println("Connected to Redis")
	return client, nil
}
