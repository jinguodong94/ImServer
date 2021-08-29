package dao

import (
	"context"
	"fmt"
	"gindemo/conf"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

var Rdb *redis.Client

func InitRedis() {

	log.Println("初始化redis")

	Rdb = redis.NewClient(&redis.Options{
		Addr:     conf.Configs.RedisConfig.Address,
		Password: conf.Configs.RedisConfig.Pwd, // no password set
		DB:       0,                            // use default DB
		PoolSize: 50,                           // 连接池大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := Rdb.Ping(ctx).Result()

	if err != nil {
		panic(fmt.Sprintf("connect redis err %s", err))
	}
}

func CloseRedis() {
	Rdb.Close()
}
