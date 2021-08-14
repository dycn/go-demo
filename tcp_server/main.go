package main

import (
	"net/http"
	_ "net/http/pprof"
	"zdy_demo/tcp_server/process"
)

func main() {
	go func() {
		http.ListenAndServe("127.0.0.1:8090", nil)
	}()
	//process.TCPServer()

	process.UDPServer()
}
