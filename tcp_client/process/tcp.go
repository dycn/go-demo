package process

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"time"
	"zdy_demo/common/tcp"
)

func heart(ctx context.Context, conn net.Conn, timeout int, stopCh chan<- struct{}) {

	ticker := time.NewTicker(time.Second * time.Duration(timeout) / 2)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			msg, err := tcp.Pack("")
			if err != nil {
				//fmt.Println(err.Error())
				stopCh <- struct{}{}
				return
			}
			_, _ = conn.Write(msg)
		}
	}
}

func write(ctx context.Context, conn net.Conn, stopCh chan<- struct{}) {
	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("停止client写消息")
			return
		default:

		}
		str, err := reader.ReadBytes('\n') //从键盘读取内容，放在str
		if err != nil {
			//fmt.Println(err.Error())
			stopCh <- struct{}{}
			return
		}

		strS := string(str[:len(str)-1])
		//把输入的内容给服务器
		msg, err := tcp.Pack(strS)
		if err != nil {
			//fmt.Println(err.Error())
			stopCh <- struct{}{}
			return
		}
		_, _ = conn.Write(msg)
	}
}

func read(ctx context.Context, conn net.Conn, stopCh chan<- struct{}) {
	reader := bufio.NewReader(conn)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("停止client读消息")
			return
		default:

		}
		//取包头内容
		lenByte, err := reader.Peek(4)
		if err == io.EOF {
			continue
		}
		if err != nil {
			//fmt.Println(err.Error())
			stopCh <- struct{}{}
			return
		}
		length := tcp.Byte2int32(lenByte)

		if reader.Buffered() < int(length)+4 {
			fmt.Println("数据有问题")
			stopCh <- struct{}{}
			return
		}

		pack := make([]byte, int(length)+4)
		_, err = reader.Read(pack)
		if err == io.EOF {
			continue
		}

		msg, err := tcp.UnPack(pack)
		if err != nil {
			//fmt.Println(err.Error())
			stopCh <- struct{}{}
			return
		}

		if msg == "exit" {
			fmt.Println("退出会话")
			stopCh <- struct{}{}
			return
		}

		fmt.Println(msg)
	}
}
