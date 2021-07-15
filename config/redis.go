package config

import (
	"context"
	"github.com/go-redis/redis/v8"
	errs "github.com/pkg/errors"
)

// RedisConfig stores the config for connecting to redis
type RedisConfig struct {
	URL string
}

// NewRedisConfig creates
func NewRedisConfig(url string) *RedisConfig {
	return &RedisConfig{
		URL: url,
	}
}

// newRedisClient creates a new redis client
func (c *RedisConfig) newRedisClient() (*redis.Client, error) {
	opts, err := redis.ParseURL(c.URL)

	if err != nil {
		return nil, errs.Wrap(err, "config.newRedisClient.ParseRedisURL")
	}

	client := redis.NewClient(opts)
	ctx := context.Background()

	_, err = client.Ping(ctx).Result()

	if err != nil {
		return nil, errs.Wrap(err, "config.newRedisClient.PingRedis")
	}

	return client, err
}

// RedisClient returns c client to the redis server
func (c *RedisConfig) RedisClient() (*redis.Client, error) {
	return c.newRedisClient()
}
