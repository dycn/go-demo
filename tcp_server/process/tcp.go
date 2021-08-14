package process

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"time"
	"zdy_demo/common/tcp"
)

func HandleTCPConnection(maxGoroutine <-chan struct{}, conn net.Conn, timeout int) {
	defer func() {
		_ = conn.Close()
	}()
	//fmt.Println("client ip:", conn.RemoteAddr().String())

	ctx, cancel := context.WithCancel(context.Background())
	msgCh := make(chan string)
	stopCh := make(chan struct{})
	heartCh := make(chan struct{})

	//回复消息
	go write(ctx, conn, msgCh, stopCh)

	//读消息
	go read(ctx, conn, msgCh, stopCh, heartCh, timeout)

	<-stopCh
	cancel()
	<-maxGoroutine
}

// write server写
func write(ctx context.Context, conn net.Conn, msgCh <-chan string, stopCh chan<- struct{}) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("停止server写信息")
			return
		default:
		}

		select {
		case <-ctx.Done():
			fmt.Println("停止server写信息")
			return
		case msg := <-msgCh:
			resp := "[此消息由server返回]收到:" + msg
			sendResp, err := tcp.Pack(resp)
			if err != nil {
				fmt.Println("打包错误", err.Error())
				stopCh <- struct{}{}
				return
			}
			conn.Write(sendResp)
		}

	}
}

// read server读
func read(ctx context.Context, conn net.Conn, msgCh chan<- string, stopCh chan<- struct{}, heartCh chan struct{}, timeout int) {
	reader := bufio.NewReader(conn)
	conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))

	for {
		select {
		case <-ctx.Done():
			fmt.Println("停止server读信息")
			return
		default:
			//心跳
			go heart(ctx, conn, heartCh, timeout)
		}

		//取包头长度
		lenByte, err := reader.Peek(4)
		if err == io.EOF {
			fmt.Println("读到结尾")
			stopCh <- struct{}{}
			break
		}
		if err != nil {
			fmt.Println("取包头长度错误", err.Error())
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
			//fmt.Println("已经读到最后了")
			continue
		}
		msg, err := tcp.UnPack(pack)
		if err != nil {
			fmt.Println("解包头长度有错误", err.Error())
			stopCh <- struct{}{}
			return
		}
		//长度为0 是心跳包
		if length == 0 {
			fmt.Println("收到心跳")
			heartCh <- struct{}{}
			continue
		}

		if msg == "exit" {
			fmt.Println("退出会话")
			stopCh <- struct{}{}
			return
		}
		fmt.Println(msg)
		msgCh <- msg
	}
}

// heart 心跳
func heart(ctx context.Context, conn net.Conn, heartCh <-chan struct{}, timeout int) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	select {
	case <-ctx.Done():
		return
	case <-heartCh:
		conn.SetDeadline(time.Now().Add(time.Second * time.Duration(timeout)))
		return
		//case <-time.After(time.Second * time.Duration(timeout)):
		//conn.Close()
		//stopCh <- struct{}{}
		//return
	}
}
