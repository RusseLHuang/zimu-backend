package utils

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var redisdb *redis.Client

// InitRedis is init redis client
func InitRedis() {
	redisdb = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis"),
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     10,
		PoolTimeout:  30 * time.Second,
	})

	_, err := redisdb.Ping().Result()
	if err != nil {
		panic(err)
	}
}

func RedisGet(key string) string {
	res, err := redisdb.Get(key).Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}

	return res
}

func RedisSet(key string, val interface{}) {
	err := redisdb.Set(key, val, 0).Err()
	if err != nil {
		panic(err)
	}
}
