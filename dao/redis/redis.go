package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var rdb *redis.Client

func Init() (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s", viper.GetString("redis.host")),
		Password: viper.GetString("redis.password"), // 密码
		DB:       viper.GetInt("redis.db"),          // 数据库
		PoolSize: viper.GetInt("redis.poolSize"),    // 连接池大小
	})
	_, err = rdb.Ping(context.TODO()).Result()
	return
}

func Close() {
	_ = rdb.Close()
}
