package database

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(addr string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	fmt.Println("Connected to Redis")
	return client
}
