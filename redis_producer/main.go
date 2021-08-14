package main

import (
	"zdy_demo/db/redis_timestamp_queue"
	"zdy_demo/redis_producer/process"
)

var (
	ADDR = "127.0.0.1"
	PORT = 6379
)

func main() {
	redis_timestamp_queue.SetRedisConf(ADDR, PORT)
	//process.ProduceOpTime()
	process.ProduceOpTimeMulti()
}
