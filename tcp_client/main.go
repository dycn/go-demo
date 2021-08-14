package main

import "zdy_demo/tcp_client/process"

func main() {
	//并发客户端 一万次请求
	//总耗时 1627479747841 - 1627479747298 = 543 ms
	//总耗时 1627539471521 - 1627539470104 = 1417 ms

	//process.Tcp()

	process.Udp()
}
