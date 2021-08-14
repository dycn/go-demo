package process

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"zdy_demo/common/rpc/thrift/gen-go/students"

	"github.com/apache/thrift/lib/go/thrift"
)

func shortSocketProcess(protocol thrift.TProtocolFactory, transport thrift.TTransportFactory, addr string) {
	socket, err := thrift.NewTSocketConf(addr, &thrift.TConfiguration{})
	if err != nil {
		log.Fatalf("建立连接失败:%v", err.Error())
	}
	defer func() {
		_ = socket.Close()
	}()
	conn, _ := transport.GetTransport(socket)

	client := students.NewClassMemberClientFactory(conn, protocol)
	ctx := context.TODO()

	for i := 0; i < 10000; i++ {
		err := socket.Open()
		if err != nil {
			log.Fatalf("连接打开失败:%v", err.Error())
		}
		_ = client.Ping(ctx, int64(i))
		_ = socket.Close()
	}
}

func longSocketProcess(protocol thrift.TProtocolFactory, transport thrift.TTransportFactory, addr string) {
	socket, err := thrift.NewTSocketConf(addr, &thrift.TConfiguration{})
	if err != nil {
		log.Fatalf("建立连接失败:%v", err.Error())
	}
	defer func() {
		_ = socket.Close()
	}()
	conn, _ := transport.GetTransport(socket)

	client := students.NewClassMemberClientFactory(conn, protocol)
	ctx := context.TODO()
	err = socket.Open()
	if err != nil {
		log.Fatalf("连接打开失败:%v", err.Error())
	}

	for i := 0; i < 10000; i++ {
		_ = client.Ping(ctx, int64(i))
	}
}

func goroutineSocketProcess(protocol thrift.TProtocolFactory, transport thrift.TTransportFactory, addr string) {
	ctx := context.TODO()
	var wg sync.WaitGroup

	config := &ThriftPoolConfig{
		MaxIdle: 100,
		Addr:    addr,
	}
	pool, err := NewThriftPool(config)
	if err != nil {
		log.Fatalln(err.Error())
	}

	wg.Add(10000)
	for i := 0; i < 10000; i++ {
		go func(i int) {
			socket, err := pool.Get()
			if err != nil || socket == nil {
				log.Fatalln("连接池取客户端异常")
			}
			conn, _ := transport.GetTransport(socket)
			client := students.NewClassMemberClientFactory(conn, protocol)
			err = socket.Open()
			if err != nil {
				log.Fatalf("连接打开失败:%v", err.Error())
			}

			err = client.Ping(ctx, int64(i))
			if err != nil {
				fmt.Println(err.Error())
			}
			_ = socket.Close()
			_ = pool.Put(socket)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func addStudent(protocol thrift.TProtocolFactory, transport thrift.TTransportFactory, addr string) {
	socket, err := thrift.NewTSocketConf(addr, &thrift.TConfiguration{})
	if err != nil {
		log.Fatalf("建立连接失败:%v", err.Error())
	}
	defer func() {
		_ = socket.Close()
	}()
	conn, _ := transport.GetTransport(socket)

	client := students.NewClassMemberClientFactory(conn, protocol)
	ctx := context.TODO()
	err = socket.Open()
	if err != nil {
		log.Fatalf("连接打开失败:%v", err.Error())
	}

	for i := 0; i < 10; i++ {
		stu := &students.Student{
			ID:   int32(i),
			Name: "testName" + strconv.Itoa(i),
			Sex:  students.Sex_MALE,
			Age:  int16(i),
		}
		_ = client.Add(ctx, stu)
	}
}

func getStudentList(protocol thrift.TProtocolFactory, transport thrift.TTransportFactory, addr string) {
	socket, err := thrift.NewTSocketConf(addr, &thrift.TConfiguration{})
	if err != nil {
		log.Fatalf("建立连接失败:%v", err.Error())
	}
	defer func() {
		_ = socket.Close()
	}()
	conn, _ := transport.GetTransport(socket)

	client := students.NewClassMemberClientFactory(conn, protocol)
	ctx := context.TODO()
	err = socket.Open()
	if err != nil {
		log.Fatalf("连接打开失败:%v", err.Error())
	}

	for i := 0; i < 10000; i++ {
		_, err := client.List(ctx)
		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Println(i)
	}
}

func count(protocol thrift.TProtocolFactory, transport thrift.TTransportFactory, addr string) {
	socket, err := thrift.NewTSocketConf(addr, &thrift.TConfiguration{})
	if err != nil {
		log.Fatalf("建立连接失败:%v", err.Error())
	}
	defer func() {
		_ = socket.Close()
	}()
	conn, _ := transport.GetTransport(socket)

	client := students.NewClassMemberClientFactory(conn, protocol)
	ctx := context.TODO()
	err = socket.Open()
	if err != nil {
		log.Fatalf("连接打开失败:%v", err.Error())
	}

	number, err := client.Count(ctx)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("有 %d 学生\n", number)
}
