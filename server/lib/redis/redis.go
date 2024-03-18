package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Redis struct {
	client             *redis.Client
	KeyExpiredDuration time.Duration
}

func NewRedis(addr string, username string, password string, db int, keyExpiredDuration time.Duration) (*Redis, error) {
	if addr == "" {
		return nil, errors.New("address cannot be empty or null")
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       db,
	})

	return &Redis{client, keyExpiredDuration}, nil
}

func (r *Redis) GetClient() *redis.Client {
	return r.client
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get key: %s", err)
	}

	return val, nil
}

func (r *Redis) Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error {
	err := r.client.Set(ctx, key, val, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key: %s", err)
	}

	return nil
}

func (r *Redis) Del(ctx context.Context, key string) (int64, error) {
	return r.client.Del(ctx, key).Result()
}
