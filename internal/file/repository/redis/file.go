package redis

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"time"

	"ttages/internal/file/entity"
)

type FileCache struct {
	client *redis.Client
	ttl    time.Duration
}

func NewFileCache(client *redis.Client, ttl time.Duration) *FileCache {
	return &FileCache{client: client, ttl: ttl}
}

func (c *FileCache) SetFiles(ctx context.Context, files []entity.File) error {
	data, err := json.Marshal(files)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, "files:all", data, c.ttl).Err()
}

func (c *FileCache) GetFiles(ctx context.Context) ([]entity.File, error) {
	data, err := c.client.Get(ctx, "files:all").Bytes()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var files []entity.File
	err = json.Unmarshal(data, &files)
	return files, err
}
