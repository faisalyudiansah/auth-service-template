package redisutils

import (
	"context"
	"encoding/json"
	"time"

	"github.com/faisalyudiansah/auth-service-template/pkg/config"

	"github.com/redis/go-redis/v9"
)

type RedisUtil interface {
	Set(ctx context.Context, key string, value any, duration time.Duration) error
	SetJSON(ctx context.Context, key string, value any, duration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	GetWithScan(ctx context.Context, key string, dest any) error
	GetWithScanJSON(ctx context.Context, key string, dest any) error
	Delete(ctx context.Context, keys ...string) error
}

type redisUtil struct {
	cfg    *config.RedisConfig
	client *redis.Client
}

func NewRedisUtils(cfg *config.RedisConfig, rdb *redis.Client) *redisUtil {
	return &redisUtil{
		cfg:    cfg,
		client: rdb,
	}
}

func (r *redisUtil) Set(ctx context.Context, key string, value any, duration time.Duration) error {
	if duration == -1 {
		duration = time.Duration(r.cfg.DefaultExpiration) * time.Second
	}

	var data string
	switch v := value.(type) {
	case string:
		data = v
	default:
		bytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		data = string(bytes)
	}
	return r.client.Set(ctx, key, data, duration).Err()
}

func (r *redisUtil) SetJSON(ctx context.Context, key string, value any, duration time.Duration) error {
	if duration == -1 {
		duration = time.Duration(r.cfg.DefaultExpiration) * time.Second
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, string(data), duration).Err()
}

func (r *redisUtil) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (r *redisUtil) GetWithScan(ctx context.Context, key string, dest any) error {
	err := r.client.Get(ctx, key).Scan(dest)
	if err == redis.Nil {
		return nil
	}
	return err
}

func (r *redisUtil) GetWithScanJSON(ctx context.Context, key string, dest any) error {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (r *redisUtil) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}
