package common

import (
	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

func init() {
	pool = &redis.Pool{
		MaxIdle:     2,
		MaxActive:   10, //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		IdleTimeout: 60, //连接关闭时间 300秒 （300秒不使用自动关闭）
		Dial: func() (redis.Conn, error) { //要连接的redis数据库
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
	}
}

func GetRedisClient() redis.Conn {
	return pool.Get()
}
