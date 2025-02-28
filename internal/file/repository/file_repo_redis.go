package repository

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
}

func (r *RedisRepository) CacheFile(id string, data string) error {
}

func (r *RedisRepository) GetCachedFile(id string) (string, error) {
}
