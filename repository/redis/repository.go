package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jwambugu/hexagonal-architecture-url-shortener/shortener"
	errs "github.com/pkg/errors"
	"strconv"
)

type redisRepository struct {
	client *redis.Client
}

func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

// Find fetches and returns shortener.Redirect using the code provided
func (r *redisRepository) Find(code string) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}

	key := r.generateKey(code)

	ctx := context.Background()
	data, err := r.client.HGetAll(ctx, key).Result()

	if err != nil {
		return nil, errs.Wrap(err, "repository.Redirect.Find")
	}

	if len(data) == 0 {
		return nil, errs.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Find")
	}

	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)

	if err != nil {
		return nil, errs.Wrap(err, "repository.Redirect.Find")
	}

	redirect.Code = data["code"]
	redirect.URL = data["url"]
	redirect.CreatedAt = createdAt

	return redirect, nil
}

// Store creates a new shortener.Redirect in the DB
func (r *redisRepository) Store(redirect *shortener.Redirect) error {
	key := r.generateKey(redirect.Code)

	data := map[string]interface{}{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}

	ctx := context.Background()
	_, err := r.client.HMSet(ctx, key, data).Result()

	if err != nil {
		return errs.Wrap(err, "repository.Redirect.Store")
	}

	return nil
}

// newRedisClient returns a client to connect to the db instance
func newRedisClient(url string) (*redis.Client, error) {
	opts, err := redis.ParseURL(url)

	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	ctx := context.Background()

	_, err = client.Ping(ctx).Result()

	if err != nil {
		return nil, err
	}

	return client, nil
}

// NewRedisRepository returns an interface to interact with the DB
func NewRedisRepository(url string) (shortener.RedirectRepository, error) {
	repo := &redisRepository{}

	client, err := newRedisClient(url)

	if err != nil {
		return nil, errs.Wrap(err, "repository.NewRedisRepository")
	}

	repo.client = client
	return repo, nil
}
