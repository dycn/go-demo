package process

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"zdy_demo/db/redis_timestamp_queue"
)

//func ProduceOpTime() {
//	//time.Sleep(time.Millisecond * 15)
//	timestamp := time.Now().Local().UnixNano()
//	data := strconv.FormatInt(timestamp, 10)
//	err := redis_timestamp_queue.AddOpTimestampQueue(data)
//	if err != nil {
//		fmt.Printf("something bad happend, %s\n", err.Error())
//	}
//	count, err := redis_timestamp_queue.CountTimestampQueue()
//	if err == nil {
//		fmt.Printf("生产了一条消息, 此时一共有 %d 消息\n", count)
//	}
//}

func ProduceOpTimeMulti() {
	for {
		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go produce(&wg)
		}
		//go func() {
		//	fmt.Println(redis_timestamp_queue.GetPoolStats())
		//}()
		wg.Wait()
		time.Sleep(time.Millisecond * 4)
	}

}

func produce(wg *sync.WaitGroup) {
	timestamp := time.Now().Add(time.Hour).Local().UnixNano()
	data := strconv.FormatInt(timestamp, 10)
	err := redis_timestamp_queue.AddOpTimestampQueue(data)
	if err != nil {
		fmt.Printf("something bad happend, %s", err.Error())
	}
	wg.Done()
}
