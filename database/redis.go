package db

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"gochat/config"
	"gochat/logger"
)

var (
	redisquerycli *redis.Client
	redisupsertcli *redis.Client
)

func init() {
	redisquerycli = redis.NewClient(&redis.Options{
		Addr:     config.GetKey("redis_host"),
		Password: config.GetKey("redis_password"),
		DB:       0,
	})
	redisupsertcli = redis.NewClient(&redis.Options{
		Addr:     config.GetKey("redis_host"),
		Password: config.GetKey("redis_password"),
		DB:       0,
	})
	_, err := redisquerycli.Ping().Result()
	if err != nil {
		logger.Error("Connect to redis error", err)
		return
	}
	_, err = redisupsertcli.Ping().Result()
	if err != nil {
		_ = redisquerycli.Conn().Close()
		logger.Error("Connect to redis error", err)
		return
	}
}

func Upsert(addr, name string) error {
	var errstr string
	
	err := redisupsertcli.Set(addr, name, 0).Err()
	if err != nil {
		errstr = fmt.Sprintln("set name err", err)
		return errors.New(errstr)
	}

	return err
}
