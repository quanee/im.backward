package config

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"gochat/logger"
	"time"
)

var (
	config *clientv3.Client
	ctx    context.Context
)

func init() {
	var err error
	config, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	ctx = context.Background()
	if err != nil {
		logger.Info("Etcd init error, ", err)
	}
	logger.Info("Etcd init successfully!")
}

func GetKey(key string) string {
	resp, err := config.Get(ctx, key)
	if err != nil {
		logger.Error("etcd get key err", err)
	}
	var value string
	for _, v := range resp.Kvs {
		value = string(v.Value)
		logger.Infof("get etcd key=%v, value=%v\n", key, value)
	}
	logger.Info("return value", value)
	return value
}
