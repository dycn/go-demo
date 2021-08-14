package main

import (
	"time"
	"zdy_demo/db/redis_timestamp_queue"
	"zdy_demo/redis_consumer/process"
)

var (
	ADDR = "127.0.0.1"
	PORT = 6379
)

func main() {
	stopCh := time.After(time.Minute * 5)
	redis_timestamp_queue.SetRedisConf(ADDR, PORT)
	for {
		select {
		case <-stopCh:
			return
		default:
			//process.OnHandleOpTime()
			//本地环境 CPU : 3 GHz 双核Intel Core i7  内存 : 16GB  默认redis配置
			//本例消费 单个消费者情况下 速度与10个生产者每10毫秒生产一次对速度相当 qps 10000左右
			//本例消费 10个消费者情况下 速度快于每5毫秒100个生产者生产一次的速度相当 qps 23000左右
			process.OnHandleOpTimeMulti()
		}
	}
}
