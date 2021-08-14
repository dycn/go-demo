package redis_timestamp_queue

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool

const RedisKey = "demo:timestamp:list"

func SetRedisConf(addr string, port int) {
	pool = &redis.Pool{
		MaxIdle:     100,             //最大数量 高频使用应该调高
		MaxActive:   100,             //最大连接数 高频使用应该调高
		IdleTimeout: 5 * time.Second, //空闲超时时间 高频使用应该调短
		Wait:        true,            //超过最大连接是否等待
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}
}

func AddOpTimestampQueue(data string) error {
	conn := pool.Get()
	defer func() {
		_ = conn.Close()
	}()
	_, err := conn.Do("LPUSH", RedisKey, data)
	return err
}

func GetOpTimestampQueue() (string, error) {
	conn := pool.Get()
	defer func() {
		_ = conn.Close()
	}()
	val, err := conn.Do("BRPOP", RedisKey, 0)

	if val == nil && err == nil {
		time.Sleep(time.Second * 2) //消费太快 超时 就休息2s再去取
		return "", err
	}
	return "", nil
	//
	//if err != nil {
	//	return "", err
	//}
	//val1, ok := val.([]interface{})
	//if !ok {
	//	//_, _ := val.(interface{})
	//}
	//if len(val1) != 2 {
	//	return "", errors.New("取出数据异常")
	//}
	//
	//recvKey, ok := val1[0].([]byte)
	//if !ok || string(recvKey) != RedisKey {
	//	return "", errors.New("redisKey 不一致")
	//}
	//
	//strTimestamp, ok := val1[1].([]byte)
	//if !ok {
	//	return "", errors.New("数据断言异常")
	//}
	//data := string(strTimestamp)
	//return data, nil
}

//func CountTimestampQueue() (int64, error) {
//	conn := pool.Get()
//	defer func() {
//		_ = conn.Close()
//	}()
//	val, err := conn.Do("LLen", RedisKey)
//	if err != nil {
//		return 0, err
//	}
//
//	num, ok := val.(int64)
//	if !ok {
//		return 0, errors.New("返回不是一个数字")
//	}
//
//	return num, nil
//}
//
//func GetPoolStats() (activeCount, idleCount int) {
//	return pool.Stats().ActiveCount, pool.Stats().IdleCount
//}
