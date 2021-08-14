package process

import (
	"fmt"
	"log"
	"net"
)

func HandleUDPConnection(conn *net.UDPConn) {
	for {
		data := make([]byte, 128)
		n, remote, err := conn.ReadFromUDP(data)
		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Println("由", remote, "发送的", string(data), "数据长度是", n)
		conn.WriteToUDP([]byte("收到"+string(data)), remote)
	}

}
