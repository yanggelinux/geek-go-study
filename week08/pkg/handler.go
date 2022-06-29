package pkg

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type RedisHandler struct {
	client *redis.ClusterClient
}

func NewRedisHandler() *RedisHandler {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{"192.168.0.2:6400", "192.168.0.2:6410", "192.168.0.2:6420"},
		Password: "password",
	})

	return &RedisHandler{client: client}

}

func (rh *RedisHandler) CloseClient() error {
	err := rh.client.Close()
	return err
}

func (rh *RedisHandler) Set(key string, val interface{}) error {
	err := rh.client.Set(key, val, time.Second*60).Err()
	return err
}

func (rh *RedisHandler) Get(key string) interface{} {

	res, err := rh.client.Get(key).Result()
	if err != nil {
		fmt.Println("get key errr:", err)
	}
	return res
}

func (rh *RedisHandler) Del(key string) error {
	err := rh.client.Del(key).Err()
	return err
}
