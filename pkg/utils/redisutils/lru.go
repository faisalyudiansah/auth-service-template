package redisutils

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/redis/go-redis/v9"
)

type redisUtilLRU struct {
	client *redis.Client
	cache  *expirable.LRU[string, string]
	ttl    time.Duration
}

func NewRedisUtilsLRU(rdb *redis.Client, capacity int, ttl time.Duration) *redisUtilLRU {
	redisUtil := &redisUtilLRU{client: rdb, ttl: ttl}
	cache := expirable.NewLRU(capacity, func(key, value string) {
		redisUtil.Delete(context.Background(), key)
	}, ttl)

	redisUtil.cache = cache
	return redisUtil
}

func (r *redisUtilLRU) Set(ctx context.Context, key string, value any, duration time.Duration) error {
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

	err := r.client.Set(ctx, key, data, duration).Err()
	if err != nil {
		return err
	}

	r.cache.Add(key, data)
	return nil
}

func (r *redisUtilLRU) SetJSON(ctx context.Context, key string, value any, duration time.Duration) error {
	return r.Set(ctx, key, value, duration)
}

func (r *redisUtilLRU) Get(ctx context.Context, key string) (string, error) {
	if val, found := r.cache.Get(key); found {
		return val, nil
	}

	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return val, nil
}

func (r *redisUtilLRU) GetWithScan(ctx context.Context, key string, dest any) error {
	err := r.client.Get(ctx, key).Scan(dest)
	if err == redis.Nil {
		return nil
	}
	return err
}

func (r *redisUtilLRU) GetWithScanJSON(ctx context.Context, key string, dest any) error {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil
	} else if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (r *redisUtilLRU) Delete(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		r.cache.Remove(key)
	}

	return r.client.Del(ctx, keys...).Err()
}
