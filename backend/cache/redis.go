package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache interface {
	Set(shortenURL string, longURL string) error
	Get(shortenURL string) (result string, err error)
}

type redisCache struct {
	host     string
	db       int
	password string
	expires  time.Duration
	ctx      context.Context
}

func NewRedisCache(host string, db int, password string, expires time.Duration) RedisCache {
	return &redisCache{
		host:     host,
		db:       db,
		password: password,
		expires:  expires,
		ctx:      context.Background(),
	}
}

func (r *redisCache) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     r.host,
		Password: r.password,
		DB:       r.db,
	})
}

func (r *redisCache) Set(shortenURL string, longURL string) error {
	err := r.getClient().Set(r.ctx, shortenURL, longURL, r.expires*time.Hour).Err()
	if err != nil {
		return errors.New("failed to set url into redis cache")
	}
	fmt.Printf("Successfully set shorten url: shorten url = %s, long url = %s\n", shortenURL, longURL)
	return nil
}

func (r *redisCache) Get(shortenURL string) (result string, err error) {
	result, err = r.getClient().Get(r.ctx, shortenURL).Result()
	if err != nil {
		return result, errors.New("failed to get shorten url from redis cache")
	}

	return result, nil
}
