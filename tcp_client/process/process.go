package process

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
	"zdy_demo/common/tcp"
	"zdy_demo/utility"

	"gopkg.in/ini.v1"
)

func Tcp() {
	config := loadIni("tcp")
	if !checkTCPConfig(config) {
		os.Exit(1)
	}

	startTime := utility.GetCurrentMillisecond()
	timeout, err := strconv.Atoi(config["timeOut"])
	if err != nil {
		log.Fatalln(err.Error())
	}
	client(config["ip"], config["port"], timeout)
	//concurrentClient(config["ip"], config["port"], timeout)
	endTime := utility.GetCurrentMillisecond()
	fmt.Printf("总耗时 %d - %d = %d ms", endTime, startTime, endTime-startTime)
}

func client(ip, port string, timeout int) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", ip, port), time.Second*time.Duration(timeout))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	ctx, cancel := context.WithCancel(context.TODO())
	stopCh := make(chan struct{})
	//固定发心跳
	go heart(ctx, conn, timeout, stopCh)

	// 发消息
	go write(ctx, conn, stopCh)

	// 收消息
	go read(ctx, conn, stopCh)

	fmt.Println("请输入内容")
	<-stopCh
	cancel()
}

//concurrentClient 并发测试
func concurrentClient(ip, port string, timeout int) {
	config := &ConnPoolConfig{
		MaxIdle: 10,
		TimeOut: time.Duration(timeout) * time.Second,
		Addr:    fmt.Sprintf("%s:%s", ip, port),
		Factory: func() (Conner, error) {
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", ip, port))
			return Conner(conn), err
		},
	}

	pool, err := NewTCPConnPool(config)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(10000)

	for i := 0; i < 10000; i++ {
		go func(i int) {
			conn, err := pool.Get()
			if err != nil {
				log.Fatalln(err.Error())

			}
			tcpConn := (*conn).(net.Conn)
			OnClient(tcpConn, i)
			err = pool.Put(conn)
			if err != nil {
				log.Fatalln(err.Error())
			}
			reader := bufio.NewReader(tcpConn)
			for {
				lenByte, err := reader.Peek(4)
				if err == io.EOF {
					//fmt.Println("取不到")
					break
				}
				if err != nil {
					//fmt.Println(err.Error())
					break
				}
				//创建一个长度为内容长度对buffer
				buffer := bytes.NewBuffer(lenByte)
				var length int32
				err = binary.Read(buffer, binary.BigEndian, &length)
				if err != nil {
					fmt.Println(err.Error())
					break
				}
				if reader.Buffered() < int(length)+4 {
					fmt.Println("一次性读不完包")
					break
				}

				pack := make([]byte, int(length)+4)
				_, err = reader.Read(pack)
				if err == io.EOF {
					break
				}
				msg, err := tcp.UnPack(pack)
				if err != nil {
					fmt.Println(err.Error())
					break
				}

				if msg == "exit" {
					fmt.Println("退出会话")
					break
				}
				fmt.Println(msg)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	//for i := 0; i < 15; i++ {
	//	go func() {
	//		conn, err := pool.Get()
	//		if err != nil {
	//			fmt.Println(err.Error())
	//			os.Exit(1)
	//		}
	//		OnClient((*conn).(net.Conn), 10086)
	//		err = pool.Put(conn)
	//		if err != nil {
	//			fmt.Println(err.Error())
	//			os.Exit(1)
	//		}
	//	}()
	//}
	//time.Sleep(10 * time.Second)
}

func OnClient(conn net.Conn, clientId int) {
	for i := 0; i < 1; i++ {
		msg := fmt.Sprintf("hello.Im client=%d,seq=%d", clientId, i)
		content, err := tcp.Pack(msg)

		if err != nil {
			log.Fatalln("打包消息出错", err.Error())
		}
		//fmt.Println("发送消息:", string(content))
		n, err := conn.Write(content)
		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Println("发送消息:", msg, "写入了", n, "字节")
	}
}

func loadIni(sectionName string) map[string]string {
	ret := make(map[string]string)
	dir, _ := os.Getwd()
	file := dir + "/conf.d/config.ini"
	config, err := ini.Load(file)
	if err != nil {
		log.Fatalln("配置文件找不到", err.Error())
	}
	ret["ip"] = config.Section(sectionName).Key("ip").String()
	ret["port"] = config.Section(sectionName).Key("port").String()
	ret["timeOut"] = config.Section(sectionName).Key("timeout").String()
	return ret
}

func checkTCPConfig(config map[string]string) bool {
	_, ok := config["ip"]
	if !ok {
		fmt.Println("ip配置错误")
		return false
	}
	_, ok = config["port"]
	if !ok {
		fmt.Println("端口配置错误")
		return false
	}
	_, ok = config["timeOut"]
	if !ok {
		fmt.Println("超时配置错误")
		return false
	}
	return true
}

func Udp() {
	config := loadIni("udp")
	if !checkUDPConfig(config) {
		os.Exit(1)
	}
	port, err := strconv.Atoi(config["port"])
	if err != nil {
		log.Fatalln(err.Error())
	}
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		//IP:   net.IP([]byte{127, 0, 0, 1}),
		Port: port,
	})
	if err != nil {
		log.Fatalln(err.Error())
	}
	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func(i int) {
			//conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
			//	//IP:   net.IP([]byte{127, 0, 0, 1}),
			//	Port: port,
			//})
			//if err != nil {
			//	fmt.Println(err.Error())
			//	os.Exit(1)
			//}
			data := []byte("hello," + strconv.Itoa(i))
			conn.Write(data)
			msg := make([]byte, 128)
			n, err := conn.Read(msg)
			if err != nil {
				log.Fatalln(err.Error())
			}
			fmt.Println("收到服务端传来", string(msg), "长度是", n)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func checkUDPConfig(config map[string]string) bool {
	_, ok := config["ip"]
	if !ok {
		fmt.Println("ip配置错误")
		return false
	}
	_, ok = config["port"]
	if !ok {
		fmt.Println("端口配置错误")
		return false
	}
	return true
}
