package redis

import (
	"context"
	"go-project/setting"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

// 初始化Redis客户端并返回RedisClient实例
func Init(cfg *setting.RedisConfig) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DefaultDb,
	})

	// 创建带有超时控制的上下文
	timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancelFunc()

	_, err := rdb.Ping(timeoutCtx).Result()
	if err != nil {
		panic("Redis初始化失败! " + err.Error())
	}
}

// Set设置Redis键值对，支持过期时间设置
func Set(ctx context.Context, key string, value interface{}, expiration ...time.Duration) error {
	exp := time.Duration(0)
	if len(expiration) > 0 {
		exp = expiration[0]
	}
	return rdb.Set(ctx, key, value, exp).Err()
}

// Get获取Redis键对应的值，如果键不存在返回redis.Nil错误
func Get(ctx context.Context, key string) (string, error) {
	val, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", err
	}
	return val, err
}

// MSet批量设置多个键值对
func MSet(ctx context.Context, values ...interface{}) error {
	return rdb.MSet(ctx, values...).Err()
}

// MGet批量获取多个键的值
func MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	return rdb.MGet(ctx, keys...).Result()
}

// Del删除一个或多个键
func Del(ctx context.Context, keys ...string) error {
	return rdb.Del(ctx, keys...).Err()
}

// Exists检查键是否存在
func Exists(ctx context.Context, key string) (bool, error) {
	n, err := rdb.Exists(ctx, key).Result()
	return n > 0, err
}

// SetNX设置键值对，如果键不存在
func SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return rdb.SetNX(ctx, key, value, expiration).Result()
}
