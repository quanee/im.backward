package db

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"gochat.udp/config"
	"gochat.udp/logger"
	"strconv"
	"time"
)

var (
	redisquerycli *redis.Client
	redisupsertcli *redis.Client
)

func init() {
	redisquerycli = redis.NewClient(&redis.Options{
		Addr:     config.GetKey("redis_host")+":"+config.GetKey("redis_port"),
		Password: config.GetKey("redis_password"),
		DB:       0,
	})
	redisupsertcli = redis.NewClient(&redis.Options{
		Addr:     config.GetKey("redis_host")+":"+config.GetKey("redis_port"),
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

func IdtoUdpAddr(id int, udpaddr string) {
	err := redisupsertcli.HSet("idtoudpaddr", strconv.Itoa(id), udpaddr).Err()
	if err != nil {
		logger.Errorf("insert id[%d] to udpaddr[%s] err: %v", id, udpaddr, err)
	}
}

func KeepOnline(id int, udpaddr string) {
	err := redisupsertcli.Set(strconv.Itoa(id), udpaddr, time.Duration(config.GetIntKey("expire"))*time.Millisecond).Err()
	if err != nil {
		logger.Errorf("insert id[%d] to udpaddr[%s] err: %v", id, udpaddr, err)
	}
}

func IsOnline(id int) bool {
	err := redisquerycli.Exists(strconv.Itoa(id))
	if err != nil {
		logger.Errorf("insert id[%d] to udpaddr[%s] err: %v", err)
		return false
	}
	return true
}

func GetUdpAddrById(id int) (res string, success bool) {
	res, err := redisquerycli.HGet("idtoudpaddr", strconv.Itoa(id)).Result()
	success = true
	if err != nil {
		logger.Errorf("get udp addr by id[%d] err: %v", id, err)
		success = false
	}
	return
}

func AllUDPAddrSet(udpaddr string) {
	err := redisupsertcli.SAdd("allupdaddrset", udpaddr).Err()
	if err != nil {
		logger.Errorf("set upd addr[%s] err: %v", udpaddr, err)
	}
}

func AllUDPAddrExist(udpaddr string) {
	err := redisupsertcli.SIsMember("allupdaddrset", udpaddr).Err()
	if err != nil {
		logger.Errorf("set upd addr[%s] err: %v", udpaddr, err)
	}
}

func QueryAddr(key string) {

}
