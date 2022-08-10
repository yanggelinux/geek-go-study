package redis

import (
	"github.com/go-redis/redis"
	"time"
)

type RedisHandler struct {
	client *redis.ClusterClient
}

func NewRedisHandler() *RedisHandler {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{"addrs"},
		Password: "password",
	})

	return &RedisHandler{client: client}
}

func (rh *RedisHandler) CloseClient() error {
	err := rh.client.Close()
	return err
}

func (rh *RedisHandler) Set(key string, val interface{}) error {
	err := rh.client.Set(key, val, time.Hour*24).Err()
	return err
}

func (rh *RedisHandler) SetHash(key, field string, val interface{}) error {
	err := rh.client.HSet(key, field, val).Err()
	return err
}

func (rh *RedisHandler) Del(key string) error {
	err := rh.client.Del(key).Err()
	return err
}
