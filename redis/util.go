package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

type CacheOptions struct {
	HotMaxKeepTime time.Duration
}

func ThroughCache[T any](ctx context.Context, key string, fn func() (T, error), expireTime time.Duration, opts *CacheOptions) (T, error) {
	var res T

	if expireTime == 0 {
		expireTime = 24 * time.Hour
	}

	if opts == nil {
		opts = &CacheOptions{}
	}

	val, err := Get(ctx, key)
	if err != nil {
		return res, err
	}

	if val != "" {
		if err = json.Unmarshal([]byte(val), &res); err != nil {
			return res, err
		}
	} else {
		res, err = fn()
		if err != nil {
			return res, err
		}
		cacheVal, err := json.Marshal(res)
		if err != nil {
			return res, err
		}
		if err := Set(ctx, key, string(cacheVal), expireTime); err != nil {
			return res, err
		}
	}

	if opts.HotMaxKeepTime > 0 {
		ttl, err := TTL(ctx, key)
		if err == nil && ttl > 0 && ttl < opts.HotMaxKeepTime/2 {
			randomSeconds := expireTime.Seconds() * rand.Float64() * 10
			newExpireTime := time.Duration(ttl) + time.Duration(randomSeconds)*time.Second
			_ = Expire(ctx, key, newExpireTime)
		}
	}

	return res, nil
}

func RedisLockMutex[T any](ctx context.Context, lockKey string, fn func() (T, error), expireTime time.Duration) (T, error) {
	var result T

	if expireTime == 0 {
		expireTime = 60 * time.Second
	}

	ok, err := SetNX(ctx, lockKey, "1", expireTime)
	if err != nil {
		return result, err
	}

	if !ok {
		fmt.Printf("lockKey: %s, failed to acquire lock, cancelled\n", lockKey)
		return result, nil
	}

	defer func() {
		_ = Del(ctx, lockKey)
	}()

	result, err = fn()
	if err != nil {
		return result, err
	}

	return result, nil
}
