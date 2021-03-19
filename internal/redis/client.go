package redis

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type Client struct {
	cl *redis.ClusterClient
}

type Member struct {
	Score float64
	Value interface{}
}

var newOnce, stopOnce sync.Once
var rdb *redis.ClusterClient

func NewClient() *Client {
	newOnce.Do(func() {
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: []string{
				viper.GetString("redis.address"),
			},
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		})
		if err := rdb.Ping(context.Background()).Err(); err != nil {
			log.Fatal("rdb ping fail, terminating", err)
		}
	})

	return &Client{
		cl: rdb,
	}
}

// Stop closes redis cluster client
// 두번 이상 호출 되면 이미 종료했다는 메세지를 return
func Stop() error {
	err := errors.New("redis cluster client is already stopped")

	stopOnce.Do(func() {
		err = rdb.Close()
	})

	return err
}
