package process

import (
	"fmt"
	"sync"
	"zdy_demo/db/redis_timestamp_queue"
)

//func OnHandleOpTime() {
//	timestampStr, err := redis_timestamp_queue.GetOpTimestampQueue()
//	if err == nil {
//		fmt.Printf("消费到一条信息 [%s]\n", timestampStr)
//	}
//}

var wg sync.WaitGroup

func OnHandleOpTimeMulti() {
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go consumer(&wg)
	}
	wg.Wait()
}

func consumer(wg *sync.WaitGroup) {
	_, err := redis_timestamp_queue.GetOpTimestampQueue()
	if err == nil {
		fmt.Println("消费到一条信息")
	} else {
		fmt.Println(err.Error())
		wg.Done()
		return
	}

	//count, err := redis_timestamp_queue.CountTimestampQueue()
	//if err == nil {
	//	fmt.Printf("消费后队列中还有 %d 条信息\n", count)
	//}

	wg.Done()
}
