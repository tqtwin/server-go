package config

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func ConnectRedis() (*redis.Client, error) {
	RDB := redis.NewClient(&redis.Options{
		Addr: GetEnv("REDIS_ADDR"),
	})

	res, err := RDB.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	println(res)

	return RDB, nil
}
