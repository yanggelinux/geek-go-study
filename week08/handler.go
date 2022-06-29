package test

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
		Addrs:    []string{"192.168.20.32:6400", "192.168.20.32:6410", "192.168.20.32:6420"},
		Password: "8D_af05a01eb94fa35cd563931e3d440@B",
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
