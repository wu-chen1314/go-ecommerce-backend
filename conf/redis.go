package conf

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// 全局 Redis 客户端
var RDB *redis.Client

// 上下文 (Redis v9 必须要有 context)
var Ctx = context.Background()

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     "host.docker.internal:6379", // Redis 地址
		Password: "",               // 密码 (没有就留空)
		DB:       0,                // 默认数据库 0
	})

	// 测试连接
	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		panic("Redis 连接失败: " + err.Error())
	}
	fmt.Println("Redis 连接成功！")
}